package entities

import (
	"time"
)

// EntityNamePortfolio contains the name of the entity
const EntityNamePortfolio = "portfolio"

// User Portfolio
type Portfolio struct {
	PortfolioID            string    `bson:"_portfolioid,omitempty"`
	Userid                 string    `bson:"userid"`
	Name                   string    `bson:"name"`
	Extid                  string    `bson:"extid"`
	TaxCountryID           string    `bson:"tax_countryid"`
	FinancialYear          time.Time `bson:"financial_year"`
	PerformenceCalculation int       `bson:"performence_calculation"`
	Summary                bool      `bson:"summary"`
	PriceAlert             bool      `bson:"price_alert"`
	CompanyEventAlert      bool      `bson:"company_event_alert"`
	CreatedAt              time.Time `bson:"created_at"`
	UpdatedAt              time.Time `bson:"updated_at"`
}
