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

// currencyRepository adapter of an currency repository for postgres
type currencyRepository struct {
	infrastructure.PostgresRepository
}

// NewCurrencyRepository creates a currency repository for postgres
func NewCurrencyRepository(db *sql.DB) ports.CurrencyRepository {
	return &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *currencyRepository) Create(ctx context.Context, currency interface{}) (string, error) {
	q := `
	INSERT INTO currency (name, code, symbol, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING currencyid;
    `

	c := currency.(entities.Currency)
	row := r.DB.QueryRowContext(
		ctx, q, c.Name, c.Code, c.Symbol, c.CreatedAt, c.UpdatedAt,
	)

	err := row.Scan(&c.CurrencyID)
	if err != nil {
		return "", err
	}

	return c.CurrencyID, nil
}

func (r *currencyRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT currencyid, code, name, symbol, created_at, updated_at
	    FROM currency %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var currencys []interface{}
	for rows.Next() {
		var c entities.Currency
		err := rows.Scan(&c.CurrencyID, &c.Code, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		currencys = append(currencys, &c)
	}

	if len(currencys) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return currencys, nil
}

func (r *currencyRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT currencyid, code, name, symbol, created_at, updated_at
        FROM currency WHERE currencyid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.Currency
	err := row.Scan(&c.CurrencyID, &c.Code, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *currencyRepository) GetBySymbol(ctx context.Context, Symbol string) (interface{}, error) {
	q := `
    SELECT currencyid, code, name, symbol, created_at, updated_at
        FROM currency WHERE symbol = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Symbol)

	var c entities.Currency
	err := row.Scan(&c.CurrencyID, &c.Code, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *currencyRepository) GetByCode(ctx context.Context, Code string) (interface{}, error) {
	q := `
    SELECT currencyid, code, name, symbol, created_at, updated_at
        FROM currency WHERE code = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Code)

	var c entities.Currency
	err := row.Scan(&c.CurrencyID, &c.Code, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *currencyRepository) Update(ctx context.Context, ID string, currency interface{}) error {
	q := `
	UPDATE currency set code=$1, name=$2, symbol=$3, updated_at=$4
	    WHERE currencyid=$5;
	`

	b := currency.(entities.Currency)
	result, err := r.DB.ExecContext(
		ctx, q, b.Code, b.Name, b.Symbol, b.UpdatedAt, ID,
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

func (r *currencyRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM currency WHERE currencyid=$1;`

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

func (r *currencyRepository) CreateMany(ctx context.Context, currency_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range currency_models {
		c := entity.(entities.Currency)

		q := `
		INSERT INTO currency (name, code, symbol, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING currencyid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, c.Name, c.Code, c.Symbol, c.CreatedAt, c.UpdatedAt,
		)
		err := row.Scan(&c.CurrencyID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, c.CurrencyID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
