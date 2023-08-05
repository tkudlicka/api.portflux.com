package models

import (
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateCurrencyReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateCurrencyReq_Ok(t *testing.T) {
	// Arrange
	req := CreateCurrencyReq{
		Name:   "Bitcoin",
		Code:   "bc",
		Symbol: "bc",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateCurrencyReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateCurrencyReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateCurrencyReq{}
	expectedError := "symbol cannot be empty | code cannot be empty | name cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateCurrencyReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateCurrencyReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateCurrencyReq{
		Name:   "Bitcoin",
		Code:   "bc",
		Symbol: "bc",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateCurrencyReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateCurrencyReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateCurrencyReq{}
	expectedError := "symbol cannot be empty | code cannot be empty | name cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
