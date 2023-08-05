package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
)

// Portfoliosp broker response struct
type PortfolioResp struct {
	PortfolioID            string    `json:"_portfolioid,omitempty"`
	Userid                 string    `json:"userid"`
	Name                   string    `json:"name"`
	Extid                  string    `json:"extid"`
	TaxCountryID           string    `json:"v"`
	FinancialYear          time.Time `json:"financial_year"`
	PerformenceCalculation int       `json:"performence_calculation"`
	Summary                bool      `json:"summary"`
	PriceAlert             bool      `json:"price_alert"`
	CompanyEventAlert      bool      `json:"company_event_alert"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// CreatePortfolioReq broker request struct
type CreatePortfolioReq struct {
	PortfolioID            string    `json:"-"`
	Userid                 string    `json:"userid"`
	Name                   string    `json:"name"`
	Extid                  string    `json:"extid"`
	TaxCountryID           string    `json:"v"`
	FinancialYear          time.Time `json:"financial_year"`
	PerformenceCalculation int       `json:"performence_calculation"`
	Summary                bool      `json:"summary"`
	PriceAlert             bool      `json:"price_alert"`
	CompanyEventAlert      bool      `json:"company_event_alert"`
	CreatedAt              time.Time `json:"-"`
	UpdatedAt              time.Time `json:"-"`
}

func (req CreatePortfolioReq) Validate() error {
	var msgs []string

	if req.Userid == "" {
		msgs = append(msgs, "user id cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}
	if req.Extid == "" {
		msgs = append(msgs, "external id cannot be empty")
	}
	if req.TaxCountryID == "" {
		msgs = append(msgs, "tax country id cannot be empty")
	}
	if req.FinancialYear == (time.Time{}) {
		msgs = append(msgs, "financial year cannot be empty")
	}
	if req.PerformenceCalculation == 0 {
		msgs = append(msgs, "performance calculation cannot be zero")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}

// UpdatePortfolioReq update broker request struct
type UpdatePortfolioReq struct {
	Userid                 string    `json:"userid"`
	Name                   string    `json:"name"`
	Extid                  string    `json:"extid"`
	TaxCountryID           string    `json:"v"`
	FinancialYear          time.Time `json:"financial_year"`
	PerformenceCalculation int       `json:"performence_calculation"`
	Summary                bool      `json:"summary"`
	PriceAlert             bool      `json:"price_alert"`
	CompanyEventAlert      bool      `json:"company_event_alert"`
}

func (req UpdatePortfolioReq) Validate() error {
	var msgs []string

	if req.Userid == "" {
		msgs = append(msgs, "user id cannot be empty")
	}
	if req.Name == "" {
		msgs = append(msgs, "name cannot be empty")
	}
	if req.Extid == "" {
		msgs = append(msgs, "external id cannot be empty")
	}
	if req.TaxCountryID == "" {
		msgs = append(msgs, "tax country id cannot be empty")
	}
	if req.FinancialYear == (time.Time{}) {
		msgs = append(msgs, "financial year cannot be empty")
	}
	if req.PerformenceCalculation == 0 {
		msgs = append(msgs, "performance calculation cannot be zero")
	}

	if len(msgs) > 0 {
		return wrappers.NewValidationErr(fmt.Errorf(strings.Join(msgs, " | ")))
	}

	return nil
}
