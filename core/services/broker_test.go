package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/testutils"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tkudlicka/portflux-api/config"
	"github.com/tkudlicka/portflux-api/core/entities"
	"github.com/tkudlicka/portflux-api/core/models"
	"github.com/tkudlicka/portflux-api/core/ports"
	"github.com/tkudlicka/portflux-api/test/mocks"
)

// TestNewBrokerService_Ok checks that NewBrokerService creates a new brokerService struct
func TestNewBrokerService_Ok(t *testing.T) {
	// Arrange
	cfg := config.Config{}
	brokerRepositoryMock := mocks.NewBrokerRepository(t)

	// Act
	service := NewBrokerService(cfg, brokerRepositoryMock)

	// Assert
	assert.NotEmpty(t, service)
}

// TestCreateBroker_Ok checks that Create returns the expected response when a valid request is received
func TestCreateBroker_Ok(t *testing.T) {
	// Arrange
	req := models.CreateBrokerReq{
		Name:        "test",
		Description: "test",
		Extid:       "1",
	}

	expectedResponse := models.CreationResp{
		InsertedID: "new-id",
	}

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Create), context.Background(), mock.AnythingOfType("entities.Broker")).Return(expectedResponse.InsertedID, nil).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	resp, err := service.Create(context.Background(), req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

// TestCreateBroker_CreateError checks that Create returns an error when the Create function from the repository fails
func TestCreateBroker_CreateError(t *testing.T) {
	// Arrange
	req := models.CreateBrokerReq{
		Name:        "test",
		Description: "test",
		Extid:       "1",
	}

	expectedError := "repository-error"

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Create), context.Background(), mock.AnythingOfType("entities.Broker")).Return("", fmt.Errorf(expectedError)).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	_, err := service.Create(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateBroker_InvalidRequest checks that Create returns an error when the received request is not valid
func TestCreateBroker_InvalidRequest(t *testing.T) {
	// Arrange
	req := models.CreateBrokerReq{
		Name:        "",
		Description: "",
	}

	expectedError := "external ID cannot be empty | name cannot be empty | description cannot be empty"

	service := &brokerService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.Create(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyBrokers_Ok checks that CreateMany does not return an error when a valid request is received
func TestCreateManyBrokers_Ok(t *testing.T) {
	// Arrange
	req := []models.CreateBrokerReq{
		{
			Extid:       "1",
			Name:        "test",
			Description: "test",
		},
	}

	expectedResponse := models.MultiCreationResp{
		InsertedIDs: []string{"new-id"},
	}

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.CreateMany), context.Background(), mock.AnythingOfType("[]interface {}")).Return(expectedResponse.InsertedIDs, nil).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	resp, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)
}

// TestCreateManyBroker_CreateManyError checks that CreateMany returns an error when the CreateMany function from the repository fails
func TestCreateManyBroker_CreateManyError(t *testing.T) {
	// Arrange
	req := []models.CreateBrokerReq{
		{
			Name:        "test",
			Description: "test",
			Extid:       "1",
		},
	}

	expectedError := "repository-error"

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.CreateMany), context.Background(), mock.AnythingOfType("[]interface {}")).Return([]string{}, fmt.Errorf(expectedError)).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	_, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateMany_InvalidRequest checks that CreateMany returns an error when one of the users in the received request is not valid
func TestCreateManyBrokers_InvalidRequest(t *testing.T) {
	// Arrange
	req := []models.CreateBrokerReq{
		{
			Name:        "",
			Description: "",
		},
	}

	expectedError := "external ID cannot be empty | name cannot be empty | description cannot be empty"

	service := &brokerService{
		config:     config.Config{},
		repository: nil,
	}

	// Act
	_, err := service.CreateMany(context.Background(), req)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestGetAll_Ok checks that GetAll returns the expected response when everything goes as expected
func TestGetAllBrokers_Ok(t *testing.T) {
	// Arrange
	var result []interface{}
	expectedBroker := entities.Broker{
		Name:        "test",
		Description: "test",
	}
	result = append(result, &expectedBroker)

	var nilPointer *int
	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Get), context.Background(), map[string]interface{}{}, nilPointer, nilPointer).Return(result, nil).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	resp, err := service.GetAll(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, models.BrokerResp(expectedBroker), resp[0])
}

// TestGetAll_NoResourcesFound checks that GetAll does not return an error when the repository does not return an user
func TestGetAllBrokers_NoResourcesFound(t *testing.T) {
	// Arrange
	var nilPointer *int
	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Get), context.Background(), map[string]interface{}{}, nilPointer, nilPointer).Return(nil, wrappers.NonExistentErr).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	resp, err := service.GetAll(context.Background())

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, 0, len(resp))
}

// TestGetBrokerByID_Ok checks that GetByID returns the expected response when a valid ID is received
func TestGetBrokerByID_Ok(t *testing.T) {
	// Arrange
	expectedBroker := entities.Broker{
		BrokerID:    "test-id",
		Name:        "test",
		Description: "test",
	}

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.GetByID), context.Background(), expectedBroker.BrokerID).Return(&expectedBroker, nil).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	resp, err := service.GetByID(context.Background(), expectedBroker.BrokerID)

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, models.BrokerResp(expectedBroker), resp)
}

// TestGetByID_Ok checks that GetByID returns tan error when the provided ID does not exist
func TestGetBrokerByID_NotFound(t *testing.T) {
	// Arrange
	nonExistentID := "non-existent-id"
	expectedError := fmt.Sprintf("ID %s not found", nonExistentID)

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.GetByID), context.Background(), nonExistentID).Return(nil, wrappers.NonExistentErr).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	_, err := service.GetByID(context.Background(), nonExistentID)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateBroker_Ok checks that Update does not return an error when everything goes as expected

func TestUpdateBroker_Ok(t *testing.T) {
	// Arrange
	testParam := "test"
	testDesc := "test"

	req := models.UpdateBrokerReq{
		BrokerID:    "test-id",
		Name:        testParam,
		Description: testDesc,
	}

	existingBroker := entities.Broker{
		Slug: "test",
	}

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.GetByID), context.Background(), req.BrokerID).Return(&existingBroker, nil).Once()
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Update), context.Background(), req.BrokerID, mock.AnythingOfType("entities.Broker")).Return(nil).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	err := service.Update(context.Background(), req.BrokerID, req)

	// Assert
	assert.Nil(t, err)
}

// TestUpdateBroker_NotFound checks that Update returns an error when the provided ID does not exist
func TestUpdateBroker_NotFound(t *testing.T) {
	// Arrange
	nonExistentID := "non-existent-id"
	expectedError := fmt.Sprintf("ID %s not found", nonExistentID)

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.GetByID), context.Background(), nonExistentID).Return(nil, wrappers.NonExistentErr).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	err := service.Update(context.Background(), nonExistentID, models.UpdateBrokerReq{})

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteBroker_Ok checks that Delete does not return an error when everything goes as expected
func TestDeleteBroker_Ok(t *testing.T) {
	// Arrange
	testID := "test-id"
	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Delete), context.Background(), testID).Return(nil).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), testID)

	// Assert
	assert.Nil(t, err)
}

// TestDeleteBroker_NotFound checks that Delete returns an error when the provided ID does not exist
func TestDeleteBroker_NotFound(t *testing.T) {
	// Arrange
	nonExistentID := "non-existent-id"
	expectedError := fmt.Sprintf("ID %s not found", nonExistentID)

	brokerRepositoryMock := mocks.NewBrokerRepository(t)
	brokerRepositoryMock.On(testutils.FunctionName(t, ports.BrokerRepository.Delete), context.Background(), nonExistentID).Return(wrappers.NonExistentErr).Once()

	service := &brokerService{
		config:     config.Config{},
		repository: brokerRepositoryMock,
	}

	// Act
	err := service.Delete(context.Background(), nonExistentID)

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.NonExistentErr, err)
	assert.Equal(t, expectedError, err.Error())
}
