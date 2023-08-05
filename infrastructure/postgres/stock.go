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

// stockRepository adapter of an stock repository for postgres
type stockRepository struct {
	infrastructure.PostgresRepository
}

// NewStockRepository creates a stock repository for postgres
func NewStockRepository(db *sql.DB) ports.StockRepository {
	return &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *stockRepository) Create(ctx context.Context, stock interface{}) (string, error) {
	q := `
	INSERT INTO stock (holdingid,extid, ticker_symbol,company_name,slug, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7)
        RETURNING stockid;
    `

	c := stock.(entities.Stock)
	row := r.DB.QueryRowContext(
		ctx, q, c.HoldingID, c.Extid, c.TickerSymbol, c.CompanyName, c.Slug, c.CreatedAt, c.UpdatedAt,
	)

	err := row.Scan(&c.StockID)
	if err != nil {
		return "", err
	}

	return c.StockID, nil
}

func (r *stockRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT stockid, holdingid,extid, ticker_symbol,company_name,slug, created_at, updated_at
	    FROM stock %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stocks []interface{}
	for rows.Next() {
		var c entities.Stock
		err := rows.Scan(&c.StockID, &c.HoldingID, &c.Extid, &c.TickerSymbol, &c.CompanyName, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, &c)
	}

	if len(stocks) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return stocks, nil
}

func (r *stockRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT stockid, holdingid,extid, ticker_symbol,company_name,slug, created_at, updated_at
        FROM stock WHERE stockid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.Stock
	err := row.Scan(&c.StockID, &c.HoldingID, &c.Extid, &c.TickerSymbol, &c.CompanyName, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *stockRepository) GetBySymbol(ctx context.Context, Symbol string) (interface{}, error) {
	q := `
    SELECT stockid, holdingid,extid, ticker_symbol,company_name,slug, created_at, updated_at
	FROM stock WHERE symbol = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Symbol)

	var c entities.Stock
	err := row.Scan(&c.StockID, &c.HoldingID, &c.Extid, &c.TickerSymbol, &c.CompanyName, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *stockRepository) GetByCode(ctx context.Context, Code string) (interface{}, error) {
	q := `
    SELECT stockid, holdingid,extid, ticker_symbol,company_name,slug, created_at, updated_at
        FROM stock WHERE code = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Code)

	var c entities.Stock
	err := row.Scan(&c.StockID, &c.HoldingID, &c.Extid, &c.TickerSymbol, &c.CompanyName, &c.Slug, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *stockRepository) Update(ctx context.Context, ID string, stock interface{}) error {
	q := `
	UPDATE stock set holdingid=$1, ticker_symbol=$2,company_name=$3,slug=$4, updated_at=$5
	    WHERE stockid=$6;
	`

	b := stock.(entities.Stock)
	result, err := r.DB.ExecContext(
		ctx, q, b.HoldingID, b.HoldingID, b.TickerSymbol, b.CompanyName, b.Slug, b.UpdatedAt, ID,
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

func (r *stockRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM stock WHERE stockid=$1;`

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

func (r *stockRepository) CreateMany(ctx context.Context, stock_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range stock_models {
		c := entity.(entities.Stock)

		q := `
		INSERT INTO stock (holdingid,extid, ticker_symbol,company_name,slug, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7)
			RETURNING stockid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, c.HoldingID, c.Extid, c.TickerSymbol, c.CompanyName, c.Slug, c.CreatedAt, c.UpdatedAt,
		)
		err := row.Scan(&c.StockID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, c.StockID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
