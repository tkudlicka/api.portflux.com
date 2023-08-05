package models

import (
	"testing"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreatePortfolioReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreatePortfolioReq_Ok(t *testing.T) {
	// Arrange
	req := CreatePortfolioReq{
		Userid:                 "1",
		Name:                   "name",
		Extid:                  "1",
		TaxCountryID:           "1",
		FinancialYear:          time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
		PerformenceCalculation: 1,
		Summary:                false,
		PriceAlert:             false,
		CompanyEventAlert:      false,
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreatePortfolioReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreatePortfolioReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreatePortfolioReq{}
	expectedError := "user id cannot be empty | name cannot be empty | external id cannot be empty | tax country id cannot be empty | financial year cannot be empty | performance calculation cannot be zero"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreatePortfolioReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdatePortfolioReq_Ok(t *testing.T) {
	// Arrange
	req := UpdatePortfolioReq{
		Userid:                 "1",
		Name:                   "name",
		Extid:                  "1",
		TaxCountryID:           "1",
		FinancialYear:          time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
		PerformenceCalculation: 1,
		Summary:                false,
		PriceAlert:             false,
		CompanyEventAlert:      false,
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdatePortfolioReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdatePortfolioReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdatePortfolioReq{}
	expectedError := "user id cannot be empty | name cannot be empty | external id cannot be empty | tax country id cannot be empty | financial year cannot be empty | performance calculation cannot be zero"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
