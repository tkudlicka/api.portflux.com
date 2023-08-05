package models

import (
	"testing"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
)

// TestValidateCreateHoldingReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateCreateHoldingReq_Ok(t *testing.T) {
	// Arrange
	req := CreateHoldingReq{
		PortfolioID:        "PortfolioID",
		BrokerID:           "BrokerID",
		Extid:              "Extid",
		Name:               "Name",
		Description:        "Description",
		Slug:               "Slug",
		TradeDate:          time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
		TradeType:          "TradeType",
		Quantity:           1,
		SharePrice:         1,
		ExchangeRate:       1,
		ExchangeCurrencyID: "czk",
		BrokerageUnitPrice: 10,
		BrokerageCurrency:  "czk",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateCreateHoldingReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateCreateHoldingReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := CreateHoldingReq{}
	expectedError := "portfolio id cannot be empty | broker id cannot be empty | external id cannot be empty | name cannot be empty | trade date cannot be empty | trade type cannot be empty | quantity cannot be zero | share price cannot be zero | exchange rate cannot be zero | exchange currency id cannot be empty | brokerage unit price cannot be zero | brokerage currency cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}

// TestValidateCreateHoldingReq_Ok checks that Validate does not return an error when a valid request is received
func TestValidateUpdateHoldingReq_Ok(t *testing.T) {
	// Arrange
	req := UpdateHoldingReq{
		PortfolioID:        "PortfolioID",
		BrokerID:           "BrokerID",
		Extid:              "Extid",
		Name:               "Name",
		Description:        "Description",
		Slug:               "Slug",
		TradeDate:          time.Date(2023, time.August, 3, 0, 0, 0, 0, time.UTC),
		TradeType:          "TradeType",
		Quantity:           1,
		SharePrice:         1,
		ExchangeRate:       1,
		ExchangeCurrencyID: "czk",
		BrokerageUnitPrice: 10,
		BrokerageCurrency:  "czk",
	}

	// Act
	err := req.Validate()

	// Assert
	assert.Nil(t, err)
}

// TestValidateUpdateHoldingReq_InvalidRequest checks that Validate returns an error when the received request is not valid
func TestValidateUpdateHoldingReq_InvalidRequest(t *testing.T) {
	// Arrange
	req := UpdateHoldingReq{}
	expectedError := "portfolio id cannot be empty | broker id cannot be empty | name cannot be empty | trade date cannot be empty | trade type cannot be empty | quantity cannot be zero | share price cannot be zero | exchange rate cannot be zero | exchange currency id cannot be empty | brokerage unit price cannot be zero | brokerage currency cannot be empty"

	// Act
	err := req.Validate()

	// Assert
	assert.NotEmpty(t, err)
	assert.IsType(t, wrappers.ValidationErr, err)
	assert.Equal(t, expectedError, err.Error())
}
