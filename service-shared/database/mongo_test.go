package database

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mocks "service-shared/mocks/database"
	shared_models "service-shared/shared-models"
	"testing"
)

const (
	firstName          = "First"
	lastName           = "Last"
	validApplicationID = "62ceaefa5338ed06fe445e18"
)

func TestCreateApplicationInternalError(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)
	mongo.On("InsertOne", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	repo := NewMongoRepository(mongo)
	_, err := repo.CreateApplication(firstName, lastName)

	assert.Equal(t, InternalError, err)
}

func TestCreateApplicationDuplicateKeyErr(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)
	dupeKeyError := getDuplicatekeyError()

	mongo.On("InsertOne", mock.Anything, mock.Anything).Return(nil, dupeKeyError).Once()
	mongo.On("InsertOne", mock.Anything, mock.Anything).Return(getInsertOneResult(), nil)

	repo := NewMongoRepository(mongo)
	resp, err := repo.CreateApplication(firstName, lastName)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	mongo.AssertNumberOfCalls(t, "InsertOne", 2)
}

func TestCreateApplicationSuccess(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)

	mongo.On("InsertOne", mock.Anything, mock.Anything).Return(getInsertOneResult(), nil)

	repo := NewMongoRepository(mongo)
	resp, err := repo.CreateApplication(firstName, lastName)

	assert.Nil(t, err)
	assert.NotNil(t, resp)
	mongo.AssertNumberOfCalls(t, "InsertOne", 1)
}

func TestGetApplicationInvalidID(t *testing.T) {
	repo := NewMongoRepository(new(mocks.MongoCaller))
	_, err := repo.GetApplication("an-invalid-id")

	assert.NotNil(t, err)
}

func TestGetApplicationInternalError(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)
	mongo.On("FindOne", mock.Anything, mock.Anything).Return(getSingleResult())

	repo := NewMongoRepository(mongo)

	_, err := repo.GetApplication(validApplicationID)

	assert.Equal(t, InternalError, err)
}

func TestGetApplicationsWithStatusInternalError(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)
	mongo.On("Find", mock.Anything, mock.Anything).Return(nil, errors.New(""))

	repo := NewMongoRepository(mongo)

	_, err := repo.GetApplicationsWithStatus(shared_models.Pending)

	assert.Equal(t, InternalError, err)
}

func TestUpdateApplicationStatusInternalError(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)
	mongo.On("UpdateByID", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(""))

	repo := NewMongoRepository(mongo)

	err := repo.UpdateApplicationStatus(validApplicationID, shared_models.Pending)

	assert.Equal(t, InternalError, err)
}

func TestRemoveApplicationInternalError(t *testing.T) {
	// Setup
	mongo := new(mocks.MongoCaller)
	mongo.On("DeleteOne", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New(""))

	repo := NewMongoRepository(mongo)

	err := repo.RemoveApplication(validApplicationID)

	assert.Equal(t, InternalError, err)
}

func getInsertOneResult() *mongo.InsertOneResult {
	return &mongo.InsertOneResult{InsertedID: primitive.ObjectID{}}
}

func getSingleResult() *mongo.SingleResult {
	var err error
	var i interface{}
	// This causes mongo to return a doc with an error
	return mongo.NewSingleResultFromDocument(i, err, nil)
}

func getDuplicatekeyError() error {
	return mongo.CommandError{
		Code:    11000,
		Message: "",
		Labels:  nil,
		Name:    "",
		Wrapped: nil,
		Raw:     nil,
	}
}
