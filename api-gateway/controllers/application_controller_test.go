package controllers

import (
	mocks "api-gateway/mocks/repositorys"
	"api-gateway/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"service-shared/database"
	sharedmocks "service-shared/mocks/database"
	sharedmodels "service-shared/shared-models"
	"strings"
	"testing"
)

const (
	applicationID = "TestID"
	dbID          = "dbID"
)

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestGetApplicationRequiresApplicationID(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)

	// Setup router
	router := SetUpRouter()
	router.GET("/api/application", controller.GetApplication)

	req, _ := http.NewRequest("GET", "/api/application", nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)
	responseData, _ := ioutil.ReadAll(respRecorder.Body)

	actualResponse := string(responseData)
	expectedResponse := "The application_id parameter is required"

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
	assert2.True(t, strings.Contains(actualResponse, expectedResponse))
}

func TestGetApplicationDbInternalError(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	repository.On("GetApplication", applicationID).Return(nil, database.InternalError)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/application", controller.GetApplication)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/application?application_id=%s", applicationID), nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
}

func TestGetApplicationDoesNotExist(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	repository.On("GetApplication", applicationID).Return(nil, errors.New(""))

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/application", controller.GetApplication)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/application?application_id=%s", applicationID), nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusNotFound, respRecorder.Code)
}

func TestGetApplicationExists(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	dbEntry := &sharedmodels.ApplicationEntry{
		ID:        primitive.ObjectID{},
		Status:    sharedmodels.Pending,
		FirstName: "First",
		LastName:  "Last",
	}

	repository.On("GetApplication", applicationID).Return(dbEntry, nil)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/application", controller.GetApplication)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/application?application_id=%s", applicationID), nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)
	responseData, _ := ioutil.ReadAll(respRecorder.Body)

	var actualClientView models.ClientApplicationView
	json.Unmarshal(responseData, &actualClientView)

	expectedClientView := dbEntryToClientView(dbEntry)

	assert.Equal(t, http.StatusOK, respRecorder.Code)
	assert.Equal(t, expectedClientView, actualClientView)
}

func TestGetApplicationWithStatusRequiresStatus(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/applications-with-status", controller.GetApplicationsWithStatus)

	req, _ := http.NewRequest("GET", "/api/applications-with-status", nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
}

func TestGetApplicationWithStatusInvalidStatus(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/applications-with-status", controller.GetApplicationsWithStatus)

	req, _ := http.NewRequest("GET", "/api/applications-with-status?status=invalidstatustype", nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
}

func TestGetApplicationWithStatusInternalDbError(t *testing.T) {
	status := sharedmodels.Pending
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	repository.On("GetApplicationsWithStatus", status).Return(nil, database.InternalError)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/applications-with-status", controller.GetApplicationsWithStatus)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/applications-with-status?status=%s", status), nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
}

func TestGetApplicationWithStatusSuccess(t *testing.T) {
	status := sharedmodels.Pending
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	var entries []sharedmodels.ApplicationEntry
	dbEntry := sharedmodels.ApplicationEntry{
		ID:        primitive.ObjectID{},
		Status:    sharedmodels.Pending,
		FirstName: "First",
		LastName:  "Last",
	}
	entries = append(entries, dbEntry)

	repository.On("GetApplicationsWithStatus", status).Return(entries, nil)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.GET("/api/applications-with-status", controller.GetApplicationsWithStatus)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/applications-with-status?status=%s", status), nil)
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)
	responseData, _ := ioutil.ReadAll(respRecorder.Body)

	var actualClientView models.GetAppsWithStatusResponse
	json.Unmarshal(responseData, &actualClientView)

	expectedClientView := models.GetAppsWithStatusResponse{ApplicationsWithStatus: dbEntryToClientResp(entries)}

	assert.Equal(t, http.StatusOK, respRecorder.Code)
	assert.Equal(t, expectedClientView, actualClientView)

}

func TestCreateApplicationInvalidBodyJSON(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.POST("/api/application", controller.CreateApplication)

	requestBody := "{invalidjson,"

	req, _ := http.NewRequest("POST", "/api/application", bytes.NewBuffer([]byte(requestBody)))
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusBadRequest, respRecorder.Code)
}

func TestCreateApplicationInternalDbError(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	repository.On("CreateApplication", mock.Anything, mock.Anything).Return("", database.InternalError)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.POST("/api/application", controller.CreateApplication)

	validRequest := &models.CreateApplicationRequest{
		FirstName: "First",
		LastName:  "Last",
	}
	jsonReqBody, _ := json.Marshal(validRequest)

	req, _ := http.NewRequest("POST", "/api/application", bytes.NewBuffer(jsonReqBody))
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
}

func TestCreateApplicationQueueError(t *testing.T) {
	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	repository.On("CreateApplication", mock.Anything, mock.Anything).Return(dbID, nil)
	repository.On("RemoveApplication", mock.Anything).Return(nil)
	messageQueue.On("PublishLoanRequest", mock.Anything).Return(errors.New(""))

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.POST("/api/application", controller.CreateApplication)

	validRequest := &models.CreateApplicationRequest{
		FirstName: "First",
		LastName:  "Last",
	}
	jsonReqBody, _ := json.Marshal(validRequest)

	req, _ := http.NewRequest("POST", "/api/application", bytes.NewBuffer(jsonReqBody))
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)

	repository.AssertCalled(t, "RemoveApplication", dbID)
	assert.Equal(t, http.StatusInternalServerError, respRecorder.Code)
}

func TestCreateApplication(t *testing.T) {
	// Expected response values
	firstName := "First"
	lastName := "Last"

	// Create mocks
	repository := new(sharedmocks.Repository)
	messageQueue := new(mocks.MessageQueue)

	repository.On("CreateApplication", mock.Anything, mock.Anything).Return(dbID, nil)
	messageQueue.On("PublishLoanRequest", mock.Anything).Return(nil)

	// Create real controller
	controller := NewLoanAppController(repository, messageQueue)
	// Setup router
	router := SetUpRouter()
	router.POST("/api/application", controller.CreateApplication)

	validRequest := &models.CreateApplicationRequest{
		FirstName: firstName,
		LastName:  lastName,
	}
	jsonReqBody, _ := json.Marshal(validRequest)

	req, _ := http.NewRequest("POST", "/api/application", bytes.NewBuffer(jsonReqBody))
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)
	responseData, _ := ioutil.ReadAll(respRecorder.Body)

	var actualClientView models.ClientApplicationView
	json.Unmarshal(responseData, &actualClientView)

	expectedClientView := models.ClientApplicationView{
		ApplicationID: dbID,
		Status:        sharedmodels.Pending,
		FirstName:     firstName,
		LastName:      lastName,
	}

	assert.Equal(t, http.StatusCreated, respRecorder.Code)
	assert.Equal(t, expectedClientView, actualClientView)
}
