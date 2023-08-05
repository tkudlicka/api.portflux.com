package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// HoldingResp holding response struct
type HoldingResp struct {
	HoldingID          string    `json:"holdingid"`
	PortfolioID        string    `json:"portfolioid"`
	BrokerID           string    `json:"brokerid"`
	Extid              string    `json:"extid"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Slug               string    `json:"slug"`
	TradeDate          time.Time `json:"trade_date"`
	TradeType          string    `json:"trade_type"`
	Quantity           int       `json:"quantity"`
	SharePrice         int       `json:"share_price"`
	ExchangeRate       int       `json:"exchange_rate"`
	ExchangeCurrencyID string    `json:"exchange_currencyid"`
	BrokerageUnitPrice int       `json:"brokerage_unit_price"`
	BrokerageCurrency  string    `json:"brokerage_currency"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CreateHoldingReq holding request struct
type CreateHoldingReq struct {
	HoldingID          string    `json:"-"`
	PortfolioID        string    `json:"portfolioid"`
	BrokerID           string    `json:"brokerid"`
	Extid              string    `json:"extid"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Slug               string    `json:"-"`
	TradeDate          time.Time `json:"trade_date"`
	TradeType          string    `json:"trade_type"`
	Quantity           int       `json:"quantity"`
	SharePrice         int       `json:"share_price"`
	ExchangeRate       int       `json:"exchange_rate"`
	ExchangeCurrencyID string    `json:"exchange_currencyid"`
	BrokerageUnitPrice int       `json:"brokerage_unit_price"`
	BrokerageCurrency  string    `json:"brokerage_currency"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

func (req CreateHoldingReq) Validate() error {
	var msgs []string

	if req.PortfolioID == "" {
		msgs = append(msgs, "portfolio id cannot be empty")
	}
	if req.BrokerID == "" {
		msgs = append(msgs, "broker id cannot be empty")
	}
	if req.Extid == "" {
		msgs = append(msgs, "external id cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}
	if req.TradeDate == (time.Time{}) {
		msgs = append(msgs, "trade date cannot be empty")
	}
	if req.TradeType == "" {
		msgs = append(msgs, "trade type cannot be empty")
	}
	if req.Quantity == 0 {
		msgs = append(msgs, "quantity cannot be zero")
	}
	if req.SharePrice == 0 {
		msgs = append(msgs, "share price cannot be zero")
	}
	if req.ExchangeRate == 0 {
		msgs = append(msgs, "exchange rate cannot be zero")
	}
	if req.ExchangeCurrencyID == "" {
		msgs = append(msgs, "exchange currency id cannot be empty")
	}
	if req.BrokerageUnitPrice == 0 {
		msgs = append(msgs, "brokerage unit price cannot be zero")
	}
	if req.BrokerageCurrency == "" {
		msgs = append(msgs, "brokerage currency cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateHoldingReq update holding request struct
type UpdateHoldingReq struct {
	HoldingID          string    `json:"-"`
	PortfolioID        string    `json:"portfolioid"`
	BrokerID           string    `json:"brokerid"`
	Extid              string    `json:"-"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Slug               string    `json:"-"`
	TradeDate          time.Time `json:"trade_date"`
	TradeType          string    `json:"trade_type"`
	Quantity           int       `json:"quantity"`
	SharePrice         int       `json:"share_price"`
	ExchangeRate       int       `json:"exchange_rate"`
	ExchangeCurrencyID string    `json:"exchange_currencyid"`
	BrokerageUnitPrice int       `json:"brokerage_unit_price"`
	BrokerageCurrency  string    `json:"brokerage_currency"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

func (req UpdateHoldingReq) Validate() error {
	var msgs []string

	if req.PortfolioID == "" {
		msgs = append(msgs, "portfolio id cannot be empty")
	}
	if req.BrokerID == "" {
		msgs = append(msgs, "broker id cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}
	if req.TradeDate == (time.Time{}) {
		msgs = append(msgs, "trade date cannot be empty")
	}
	if req.TradeType == "" {
		msgs = append(msgs, "trade type cannot be empty")
	}
	if req.Quantity == 0 {
		msgs = append(msgs, "quantity cannot be zero")
	}
	if req.SharePrice == 0 {
		msgs = append(msgs, "share price cannot be zero")
	}
	if req.ExchangeRate == 0 {
		msgs = append(msgs, "exchange rate cannot be zero")
	}
	if req.ExchangeCurrencyID == "" {
		msgs = append(msgs, "exchange currency id cannot be empty")
	}
	if req.BrokerageUnitPrice == 0 {
		msgs = append(msgs, "brokerage unit price cannot be zero")
	}
	if req.BrokerageCurrency == "" {
		msgs = append(msgs, "brokerage currency cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
