package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// CurrencyResp currency response struct
type CurrencyResp struct {
	CurrencyID string    `json:"-"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Symbol     string    `json:"symbol"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateCurrencyReq currency request struct
type CreateCurrencyReq struct {
	CurrencyID string    `json:"-"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Symbol     string    `json:"symbol"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

func (req CreateCurrencyReq) Validate() error {
	var msgs []string

	if req.Symbol == "" {
		msgs = append(msgs, "symbol cannot be empty")
	}
	if req.Code == "" {
		msgs = append(msgs, "code cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateCurrencyReq update currency request struct
type UpdateCurrencyReq struct {
	CurrencyID string    `json:"-"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	Symbol     string    `json:"symbol"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

func (req UpdateCurrencyReq) Validate() error {
	var msgs []string

	if req.Symbol == "" {
		msgs = append(msgs, "symbol cannot be empty")
	}
	if req.Code == "" {
		msgs = append(msgs, "code cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
