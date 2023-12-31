package models

import (
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateCryptoCurrencyReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateCryptoCurrencyReq_Ok(t *testing.T) {
	// Arrange
	req := CreateCryptoCurrencyReq{
		Name:   "Bitcoin",
		Symbol: "bc",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateCryptoCurrencyReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateCryptoCurrencyReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateCryptoCurrencyReq{}
	expectedError := "symbol cannot be empty | name cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateCryptoCurrencyReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateCryptoCurrencyReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateCryptoCurrencyReq{
		Name:   "Bitcoin",
		Symbol: "bc",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateCryptoCurrencyReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateCryptoCurrencyReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateCryptoCurrencyReq{}
	expectedError := "symbol cannot be empty | name cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
