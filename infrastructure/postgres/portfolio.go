package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/tkudlicka/portflux-api/core/entities"
	"github.com/tkudlicka/portflux-api/core/ports"
)

// portfolioRepository adapter of an portfolio repository for postgres
type portfolioRepository struct {
	infrastructure.PostgresRepository
}

// NewPortfolioRepository creates a portfolio repository for postgres
func NewPortfolioRepository(db *sql.DB) ports.PortfolioRepository {
	return &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *portfolioRepository) Create(ctx context.Context, portfolio interface{}) (string, error) {
	q := `
	INSERT INTO portfolio (userid,name,extid,tax_countryid,financial_year,performence_calculation,summary,price_alert,company_event_alert,created_at,updated_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11)
        RETURNING portfolioid;
    `

	c := portfolio.(entities.Portfolio)
	row := r.DB.QueryRowContext(
		ctx, q, c.Userid, c.Name, c.Extid, c.TaxCountryID, c.FinancialYear, c.PerformenceCalculation, c.Summary, c.PriceAlert, c.CompanyEventAlert, c.CreatedAt, c.UpdatedAt,
	)

	err := row.Scan(&c.PortfolioID)
	if err != nil {
		return "", err
	}

	return c.PortfolioID, nil
}

func (r *portfolioRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
	var where string
	for k, v := range filter {
		if where == "" {
			where = "WHERE"
		} else {
			where = fmt.Sprintf("%s AND", where)
		}
		where = fmt.Sprintf("%s %s = '%v'", where, k, v)
	}
	if skip != nil {
		where = fmt.Sprintf("%s OFFSET %d", where, *skip)
	}
	if take != nil {
		where = fmt.Sprintf("%s LIMIT %d", where, *take)
	}

	q := fmt.Sprintf(`
		SELECT portfolioid,userid,name,extid,tax_countryid,financial_year,performence_calculation,summary,price_alert,company_event_alert,created_at,updated_at
	    FROM portfolio %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var portfolios []interface{}
	for rows.Next() {
		var c entities.Portfolio
		err := rows.Scan(&c.PortfolioID, &c.Userid, &c.Name, &c.Extid, &c.TaxCountryID, &c.FinancialYear, &c.PerformenceCalculation, &c.Summary, &c.PriceAlert, &c.CompanyEventAlert, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		portfolios = append(portfolios, &c)
	}

	if len(portfolios) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return portfolios, nil
}

func (r *portfolioRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    	SELECT portfolioid,userid,name,extid,tax_countryid,financial_year,performence_calculation,summary,price_alert,company_event_alert,created_at,updated_at
        FROM portfolio WHERE portfolioid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.Portfolio
	err := row.Scan(&c.PortfolioID, &c.Userid, &c.Name, &c.Extid, &c.TaxCountryID, &c.FinancialYear, &c.PerformenceCalculation, &c.Summary, &c.PriceAlert, &c.CompanyEventAlert, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *portfolioRepository) GetBySymbol(ctx context.Context, Symbol string) (interface{}, error) {
	q := `
	SELECT portfolioid,userid,name,extid,tax_countryid,financial_year,performence_calculation,summary,price_alert,company_event_alert,created_at,updated_at
	FROM portfolio WHERE symbol = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Symbol)

	var c entities.Portfolio
	err := row.Scan(&c.PortfolioID, &c.Userid, &c.Name, &c.Extid, &c.TaxCountryID, &c.FinancialYear, &c.PerformenceCalculation, &c.Summary, &c.PriceAlert, &c.CompanyEventAlert, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *portfolioRepository) GetByCode(ctx context.Context, Code string) (interface{}, error) {
	q := `
		SELECT portfolioid,userid,name,extid,tax_countryid,financial_year,performence_calculation,summary,price_alert,company_event_alert,created_at,updated_at
        FROM portfolio WHERE code = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Code)

	var c entities.Portfolio
	err := row.Scan(&c.PortfolioID, &c.Userid, &c.Name, &c.Extid, &c.TaxCountryID, &c.FinancialYear, &c.PerformenceCalculation, &c.Summary, &c.PriceAlert, &c.CompanyEventAlert, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *portfolioRepository) Update(ctx context.Context, ID string, portfolio interface{}) error {
	q := `
	UPDATE portfolio set name=$1,tax_countryid=$2,financial_year=$3,performence_calculation=$4,summary=$5,price_alert=$6,company_event_alert=$7,updated_at=$8
	    WHERE portfolioid=$9;
	`

	b := portfolio.(entities.Portfolio)
	result, err := r.DB.ExecContext(
		ctx, q, b.Name, b.TaxCountryID, b.FinancialYear, b.PerformenceCalculation, b.Summary, b.PriceAlert, b.CompanyEventAlert, b.UpdatedAt, ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("RowsAffected Error", err)
	}
	if rows < 1 {
		return wrappers.NewNonExistentErr(sql.ErrNoRows)
	}
	return nil
}

func (r *portfolioRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM portfolio WHERE portfolioid=$1;`

	result, err := r.DB.ExecContext(ctx, q, ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows < 1 {
		return wrappers.NewNonExistentErr(sql.ErrNoRows)
	}
	return nil
}

func (r *portfolioRepository) CreateMany(ctx context.Context, portfolio_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range portfolio_models {
		c := entity.(entities.Portfolio)

		q := `
		INSERT INTO portfolio (userid,name,extid,tax_countryid,financial_year,performence_calculation,summary,price_alert,company_event_alert,created_at,updated_at)
        	VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11)
			RETURNING portfolioid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, c.Userid, c.Name, c.Extid, c.TaxCountryID, c.FinancialYear, c.PerformenceCalculation, c.Summary, c.PriceAlert, c.CompanyEventAlert, c.CreatedAt, c.UpdatedAt,
		)
		err := row.Scan(&c.PortfolioID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, c.PortfolioID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
