package models

import (
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateBrokerReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateBrokerReq_Ok(t *testing.T) {
	// Arrange
	req := CreateBrokerReq{
		Name:        "test@test.com",
		Extid:       "test",
		Description: "test",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateBrokerReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateBrokerReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateBrokerReq{}
	expectedError := "external ID cannot be empty | name cannot be empty | description cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateBrokerReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateBrokerReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateBrokerReq{
		Name:        "test@test.com",
		Description: "test",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateBrokerReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateBrokerReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateBrokerReq{}
	expectedError := "name cannot be empty | description cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
