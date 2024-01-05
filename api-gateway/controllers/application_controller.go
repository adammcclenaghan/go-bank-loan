//Package controllers contains controllers for the API gateway. They are responsible for
//handling HTTP requests.
package controllers

import (
	"api-gateway/models"
	"api-gateway/repositorys"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"service-shared/database"
	sharedmodels "service-shared/shared-models"
	"strings"
)

/*
LoanAppController is an HTTP controller which manages receiving
HTTP requests to the API. It hands of responsibilities for DB
interaction or message queue interaction to database.Repository
and repositorys.PublishQueue respectively.
*/
type LoanAppController struct {
	repository   database.Repository
	messageQueue repositorys.PublishQueue
}

//NewLoanAppController returns a LoanAppController struct
func NewLoanAppController(repository database.Repository, messageQueue repositorys.PublishQueue) *LoanAppController {
	return &LoanAppController{repository: repository, messageQueue: messageQueue}
}

// The following error structs are used to build nicer API documentation via swagger

// HTTPBadRequestError is returned on a bad request
type HTTPBadRequestError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"error" example:"status bad request"`
}

// HTTPInternalServerError is returned on an internal server error
type HTTPInternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"error" example:"status internal server error"`
}

// HTTPNotFoundError is returned for an HTTP not found response
type HTTPNotFoundError struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"error" example:"status not found"`
}

//GetApplication godoc
//@Summary Gets a loan application
//@Tags applications
//@Description Gets a loan application based on a provided application ID
//@Produce json
//@Param application_id query string true "Loan Application ID"
//@Success 200 {object} models.ClientApplicationView "Application retrieved"
//@Failure 400 {object} HTTPBadRequestError "When an application ID is not provided"
//@Failure 404 {object} HTTPNotFoundError "When an application ID is not found"
//@Failure 500 {object} HTTPInternalServerError "When an internal server error occurs"
//@Router /api/application/ [get]
func (controller LoanAppController) GetApplication(ginCtx *gin.Context) {
	applicationID := ginCtx.Query("application_id")
	if len(applicationID) == 0 {
		newBadRequest(ginCtx, http.StatusBadRequest, errors.New("The application_id parameter is required"))
		return
	}

	statusResponse, err := controller.repository.GetApplication(applicationID)
	if err != nil {
		if errors.Is(err, database.InternalError) {
			newInternalError(ginCtx, http.StatusInternalServerError, err)
			return
		}

		newNotFoundError(ginCtx, http.StatusNotFound, err)
		return
	}

	clientResponse := dbEntryToClientView(statusResponse)
	ginCtx.IndentedJSON(http.StatusOK, clientResponse)
}

//GetApplicationsWithStatus godoc
//@Summary Gets all loans with status
//@Tags applications
//@Description Gets all loans based on a provided status
//@Produce json
//@Param status query string true "Status [pending, completed, rejected]"
//@Success 200 {object} models.GetAppsWithStatusResponse "Applications retrieved"
//@Failure 400 {object} HTTPBadRequestError "When the status parameter is not provided or is not a valid value"
//@Failure 500 {object} HTTPInternalServerError "When an internal server error occurs"
//@Router /api/applications-with-status [get]
func (controller LoanAppController) GetApplicationsWithStatus(ginCtx *gin.Context) {
	status := strings.ToLower(ginCtx.Query("status"))
	if !sharedmodels.Status(status).IsValid() {
		newBadRequest(ginCtx, http.StatusBadRequest, errors.New(fmt.Sprintf("The status parameter is required and must be one of [%s %s %s]",
			sharedmodels.Pending, sharedmodels.Completed, sharedmodels.Rejected)))
		return
	}

	applications, err := controller.repository.GetApplicationsWithStatus(sharedmodels.Status(status))
	if err != nil {
		newInternalError(ginCtx, http.StatusInternalServerError, err)
		return
	}

	clientResponse := models.GetAppsWithStatusResponse{ApplicationsWithStatus: dbEntryToClientResp(applications)}
	ginCtx.IndentedJSON(http.StatusOK, clientResponse)
}

//CreateApplication godoc
//@Summary Create a loan application
//@Tags applications
//@Description Creates a new loan application
//@Accept json
//@Param application body models.CreateApplicationRequest true "Create loan application"
//@Produce json
//@Success 201 {object} models.CreateApplicationResponse "Loan application created"
//@Failure 400 {object} HTTPBadRequestError "When the request body is malformed"
//@Failure 500 {object} HTTPInternalServerError "When an internal server error occurs"
//@Router /api/application [post]
func (controller LoanAppController) CreateApplication(ginCtx *gin.Context) {
	var createRequest models.CreateApplicationRequest
	if err := ginCtx.ShouldBindJSON(&createRequest); err != nil {
		newBadRequest(ginCtx, http.StatusBadRequest, err)
		return
	}

	// Add to the DB
	applicationID, err := controller.repository.CreateApplication(createRequest.FirstName, createRequest.LastName)
	if err != nil {
		newInternalError(ginCtx, http.StatusInternalServerError, err)
		return
	}

	// Push message onto the message queue
	loanApplication := sharedmodels.CreateLoanMessage{
		ApplicationID: applicationID,
		FirstName:     createRequest.FirstName,
		LastName:      createRequest.LastName,
	}

	queueErr := controller.messageQueue.PublishLoanRequest(loanApplication)
	if queueErr != nil {
		// We failed to publish the loan application to the create queue ...
		// Unfortunately, according to RabbitMQ docs, this does not indicate whether
		// the MQ server has received this msg.

		// Remove the entry from the DB and respond with an internal error
		controller.repository.RemoveApplication(applicationID)
		fmt.Println(fmt.Sprintf("Encountered an error publishing to the queue %s", queueErr))
		newInternalError(ginCtx, http.StatusInternalServerError, queueErr)
		return
	}

	fmt.Printf("Sent message %#v to queue\n", loanApplication)
	clientResponse := models.CreateApplicationResponse{
		ApplicationID: applicationID,
		Status:        sharedmodels.Pending,
		FirstName:     createRequest.FirstName,
		LastName:      createRequest.LastName,
	}
	ginCtx.IndentedJSON(http.StatusCreated, clientResponse)
}

// Helper funcs for converting db entry type to a client friendly view
func dbEntryToClientResp(dbEntries []sharedmodels.ApplicationEntry) []models.ClientApplicationView {
	// Purposefully init to empty so that clients don't get 'nil' in JSON response
	clientView := []models.ClientApplicationView{}

	for _, entry := range dbEntries {
		clientView = append(clientView, dbEntryToClientView(&entry))
	}

	return clientView
}

func dbEntryToClientView(dbEntry *sharedmodels.ApplicationEntry) models.ClientApplicationView {
	return models.ClientApplicationView{
		ApplicationID: dbEntry.ID.Hex(),
		Status:        dbEntry.Status,
		FirstName:     dbEntry.FirstName,
		LastName:      dbEntry.LastName,
	}
}

func newBadRequest(ctx *gin.Context, status int, err error) {
	er := HTTPBadRequestError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

func newInternalError(ctx *gin.Context, status int, err error) {
	er := HTTPInternalServerError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}

func newNotFoundError(ctx *gin.Context, status int, err error) {
	er := HTTPNotFoundError{
		Code:    status,
		Message: err.Error(),
	}
	ctx.JSON(status, er)
}
