package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// TransactionResp transaction response struct
type TransactionResp struct {
	TransactionID    string    `json:"transactionid"`
	StockID          string    `json:"stockid"`
	CryptocurrencyID string    `json:"cryptocurrencyid"`
	Quantity         string    `json:"quantity"`
	TransactionPrice string    `json:"transaction_price"`
	TransactionDate  time.Time `json:"transaction_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateTransactionReq broker request struct
type CreateTransactionReq struct {
	TransactionID    string    `json:"-"`
	StockID          string    `json:"stockid"`
	CryptocurrencyID string    `json:"cryptocurrencyid"`
	Quantity         int16     `json:"quantity"`
	TransactionPrice string    `json:"transaction_price"`
	TransactionDate  time.Time `json:"transaction_date"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (req CreateTransactionReq) Validate() error {
	var msgs []string

	if req.StockID == "" && req.CryptocurrencyID == "" {
		msgs = append(msgs, "either stock id or cryptocurrency id must be provided")
	}
	if req.Quantity == 0 {
		msgs = append(msgs, "quantity cannot be zero")
	}
	if req.TransactionPrice == "" {
		msgs = append(msgs, "transaction price cannot be empty")
	}
	if req.TransactionDate == (time.Time{}) {
		msgs = append(msgs, "transaction date cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateTransactionReq update transaction request struct
type UpdateTransactionReq struct {
	TransactionID    string    `json:"-"`
	StockID          string    `json:"stockid"`
	CryptocurrencyID string    `json:"cryptocurrencyid"`
	Quantity         int16     `json:"quantity"`
	TransactionPrice string    `json:"transaction_price"`
	TransactionDate  time.Time `json:"transaction_date"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (req UpdateTransactionReq) Validate() error {
	var msgs []string

	if req.StockID == "" && req.CryptocurrencyID == "" {
		msgs = append(msgs, "either stock id or cryptocurrency id must be provided")
	}
	if req.Quantity == 0 {
		msgs = append(msgs, "quantity cannot be zero")
	}
	if req.TransactionPrice == "" {
		msgs = append(msgs, "transaction price cannot be empty")
	}
	if req.TransactionDate == (time.Time{}) {
		msgs = append(msgs, "transaction date cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
