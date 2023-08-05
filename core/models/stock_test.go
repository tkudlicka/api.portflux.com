package models

import (
	"testing"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateStockReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateStockReq_Ok(t *testing.T) {
	// Arrange
	req := CreateStockReq{
		HoldingID:    "HoldingID",
		Extid:        "Extid",
		TickerSymbol: "TickerSymbol",
		CompanyName:  "CompanyName",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateStockReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateStockReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateStockReq{}
	expectedError := "holding id cannot be empty | external id cannot be empty | ticker symbol cannot be empty | company name cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateStockReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateStockReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateStockReq{
		HoldingID:    "HoldingID",
		Extid:        "Extid",
		TickerSymbol: "TickerSymbol",
		CompanyName:  "CompanyName",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateStockReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateStockReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateStockReq{}
	expectedError := "holding id cannot be empty | external id cannot be empty | ticker symbol cannot be empty | company name cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
