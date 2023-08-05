package models

import (
	"testing"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateDividendReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateDividendReq_Ok(t *testing.T) {
	// Arrange
	req := CreateDividendReq{
		StockID:          "5",
		DividendPerShare: 5,
		DividendDate:     time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateDividendReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateDividendReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateDividendReq{}
	expectedError := "stock id cannot be empty | dividend per share cannot be zero | dividend date cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateDividendReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateDividendReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateDividendReq{
		StockID:          "5",
		DividendPerShare: 5,
		DividendDate:     time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateDividendReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateDividendReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateDividendReq{}
	expectedError := "stock id cannot be empty | dividend per share cannot be zero | dividend date cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
