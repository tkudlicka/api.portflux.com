package models

import (
	"testing"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateTransactionReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateTransactionReq_Ok(t *testing.T) {
	// Arrange
	req := CreateTransactionReq{
		StockID:          "StockID",
		CryptocurrencyID: "CryptocurrencyID",
		Quantity:         10,
		TransactionPrice: "TransactionPrice",
		TransactionDate:  time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateTransactionReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateTransactionReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateTransactionReq{}
	expectedError := "either stock id or cryptocurrency id must be provided | quantity cannot be zero | transaction price cannot be empty | transaction date cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateTransactionReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateTransactionReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateTransactionReq{
		StockID:          "StockID",
		CryptocurrencyID: "CryptocurrencyID",
		Quantity:         10,
		TransactionPrice: "TransactionPrice",
		TransactionDate:  time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateTransactionReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateTransactionReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateTransactionReq{}
	expectedError := "either stock id or cryptocurrency id must be provided | quantity cannot be zero | transaction price cannot be empty | transaction date cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
