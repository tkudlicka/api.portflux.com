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

// holdingRepository adapter of an holding repository for postgres
type holdingRepository struct {
	infrastructure.PostgresRepository
}

// NewHoldingRepository creates a holding repository for postgres
func NewHoldingRepository(db *sql.DB) ports.HoldingRepository {
	return &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *holdingRepository) Create(ctx context.Context, holding interface{}) (string, error) {
	q := `
	INSERT INTO holding (portfolioid,brokerid,extid,name,description,slug,trade_date,trade_type,quantity,share_price,exchange_rate,exhange_currencyid,brokerage_unit_price,brokerage_currency, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
        RETURNING holdingid;
    `

	c := holding.(entities.Holding)
	row := r.DB.QueryRowContext(
		ctx, q, c.PortfolioID, c.BrokerID, c.Extid, c.Name, c.Description, c.Slug, c.TradeDate, c.TradeType, c.Quantity, c.SharePrice, c.ExchangeRate, c.ExchangeCurrencyID, c.BrokerageUnitPrice, c.BrokerageCurrency, c.CreatedAt, c.UpdatedAt,
	)

	err := row.Scan(&c.HoldingID)
	if err != nil {
		return "", err
	}

	return c.HoldingID, nil
}

func (r *holdingRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT holdingid,portfolioidbrokerid,extid,name,description,slug,trade_date,trade_type,quantity,share_price,exchange_rate,exchange_currencyid,brokerage_unit_price,brokerage_currency,created_at,updated_at
	    FROM holding %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var holdings []interface{}
	for rows.Next() {
		var c entities.Holding
		err := rows.Scan(&c.HoldingID, &c.PortfolioID, &c.BrokerID, &c.Extid, &c.Name, &c.Description, &c.Slug, &c.TradeDate, &c.TradeType, &c.Quantity, &c.SharePrice, &c.ExchangeRate, &c.ExchangeCurrencyID, &c.BrokerageUnitPrice, &c.BrokerageCurrency, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		holdings = append(holdings, &c)
	}

	if len(holdings) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return holdings, nil
}

func (r *holdingRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT holdingid,portfolioidbrokerid,extid,name,description,slug,trade_date,trade_type,quantity,share_price,exchange_rate,exchange_currencyid,brokerage_unit_price,brokerage_currency,created_at,updated_at
        FROM holding WHERE holdingid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.Holding
	err := row.Scan(&c.HoldingID, &c.PortfolioID, &c.BrokerID, &c.Extid, &c.Name, &c.Description, &c.Slug, &c.TradeDate, &c.TradeType, &c.Quantity, &c.SharePrice, &c.ExchangeRate, &c.ExchangeCurrencyID, &c.BrokerageUnitPrice, &c.BrokerageCurrency, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *holdingRepository) Update(ctx context.Context, ID string, holding interface{}) error {
	q := `
	UPDATE holding set brokerid=$1, name=$2, description=$3, slug=$4,trade_date=$5,trade_type=$6,quantity=$7,share_price=$8,exhcange_rate=$9,exchange_currencyid=$10,brokerage_unit_price=$11,brokerage_currency=$12,,updated_at=$13
	    WHERE holdingid=$14;
	`

	c := holding.(entities.Holding)
	result, err := r.DB.ExecContext(
		ctx, q, c.BrokerID, c.Name, c.Description, c.Slug, c.TradeDate, c.TradeType, c.Quantity, c.SharePrice, c.ExchangeRate, c.ExchangeCurrencyID, c.BrokerageUnitPrice, c.BrokerageCurrency, c.UpdatedAt, ID,
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

func (r *holdingRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM holding WHERE holdingid=$1;`

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

func (r *holdingRepository) CreateMany(ctx context.Context, holding_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range holding_models {
		c := entity.(entities.Holding)

		q := `
		INSERT INTO holding (portfolioid,brokerid,extid,name,description,slug,trade_date,trade_type,quantity,share_price,exchange_rate,exhange_currencyid,brokerage_unit_price,brokerage_currency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
			RETURNING holdingid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, c.PortfolioID, c.BrokerID, c.Extid, c.Name, c.Description, c.Slug, c.TradeDate, c.TradeType, c.Quantity, c.SharePrice, c.ExchangeRate, c.ExchangeCurrencyID, c.BrokerageUnitPrice, c.BrokerageCurrency, c.CreatedAt, c.UpdatedAt,
		)
		err := row.Scan(&c.HoldingID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, c.HoldingID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
