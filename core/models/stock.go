package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// StockResp broker response struct
type StockResp struct {
	StockID      string    `json:"stockid"`
	HoldingID    string    `json:"holdingid"`
	Extid        string    `json:"extid"`
	TickerSymbol string    `json:"ticker_symbol"`
	CompanyName  string    `json:"company_name"`
	Slug         string    `json:"slug"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateStockReq broker request struct
type CreateStockReq struct {
	StockID      string    `json:"-"`
	HoldingID    string    `json:"holdingid"`
	Extid        string    `json:"extid"`
	TickerSymbol string    `json:"ticker_symbol"`
	CompanyName  string    `json:"company_name"`
	Slug         string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func (req CreateStockReq) Validate() error {
	var msgs []string

	if req.HoldingID == "" {
		msgs = append(msgs, "holding id cannot be empty")
	}
	if req.Extid == "" {
		msgs = append(msgs, "external id cannot be empty")
	}
	if req.TickerSymbol == "" {
		msgs = append(msgs, "ticker symbol cannot be empty")
	}
	if req.CompanyName == "" {
		msgs = append(msgs, "company name cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateStockReq update broker request struct
type UpdateStockReq struct {
	StockID      string    `json:"-"`
	HoldingID    string    `json:"holdingid"`
	Extid        string    `json:"extid"`
	TickerSymbol string    `json:"ticker_symbol"`
	CompanyName  string    `json:"company_name"`
	Slug         string    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

func (req UpdateStockReq) Validate() error {
	var msgs []string

	if req.HoldingID == "" {
		msgs = append(msgs, "holding id cannot be empty")
	}
	if req.Extid == "" {
		msgs = append(msgs, "external id cannot be empty")
	}
	if req.TickerSymbol == "" {
		msgs = append(msgs, "ticker symbol cannot be empty")
	}
	if req.CompanyName == "" {
		msgs = append(msgs, "company name cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
