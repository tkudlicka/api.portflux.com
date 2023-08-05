package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// CryptoCurrencyResp crypto currency response struct
type CryptoCurrencyResp struct {
	CryptoCurrencyID string    `json:"-"`
	HoldingID        string    `json:"-"`
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateCryptoCurrencyReq crypto currency request struct
type CreateCryptoCurrencyReq struct {
	CryptoCurrencyID string    `json:"-"`
	HoldingID        string    `json:"-"`
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (req CreateCryptoCurrencyReq) Validate() error {
	var msgs []string

	if req.Symbol == "" {
		msgs = append(msgs, "symbol cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdateCryptoCurrencyReq update crypto currency request struct
type UpdateCryptoCurrencyReq struct {
	CryptoCurrencyID string    `json:"-"`
	HoldingID        string    `json:"-"`
	Symbol           string    `json:"symbol"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

func (req UpdateCryptoCurrencyReq) Validate() error {
	var msgs []string

	if req.Symbol == "" {
		msgs = append(msgs, "symbol cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
