package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// DividendResp dividend response struct
type DividendResp struct {
	DividendID       string    `json:"-"`
	StockID          string    `json:"stockid"`
	DividendPerShare int       `json:"dividend_per_share"`
	DividendDate     time.Time `json:"dividend_date"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

// CreateDividendReq dividend request struct
type CreateDividendReq struct {
	DividendID       string    `json:"-"`
	StockID          string    `json:"stockid"`
	DividendPerShare int       `json:"dividend_per_share"`
	DividendDate     time.Time `json:"dividend_date"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (req CreateDividendReq) Validate() error {
	var msgs []string

	if req.StockID == "" {
		msgs = append(msgs, "stock id cannot be empty")
	}
	if req.DividendPerShare == 0 {
		msgs = append(msgs, "dividend per share cannot be zero")
	}

	if req.DividendDate == (time.Time{}) {
		msgs = append(msgs, "dividend date cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateDividendReq update dividend request struct
type UpdateDividendReq struct {
	DividendID       string    `json:"-"`
	StockID          string    `json:"stockid"`
	DividendPerShare int       `json:"dividend_per_share"`
	DividendDate     time.Time `json:"dividend_date"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (req UpdateDividendReq) Validate() error {
	var msgs []string

	if req.StockID == "" {
		msgs = append(msgs, "stock id cannot be empty")
	}
	if req.DividendPerShare == 0 {
		msgs = append(msgs, "dividend per share cannot be zero")
	}

	if req.DividendDate == (time.Time{}) {
		msgs = append(msgs, "dividend date cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
