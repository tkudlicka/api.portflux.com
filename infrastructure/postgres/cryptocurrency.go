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

// cryptoCurrencyRepository adapter of an cryptoCurrency repository for postgres
type cryptoCurrencyRepository struct {
	infrastructure.PostgresRepository
}

// NewCryptoCurrencyRepository creates a cryptoCurrency repository for postgres
func NewCryptoCurrencyRepository(db *sql.DB) ports.CryptoCurrencyRepository {
	return &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *cryptoCurrencyRepository) Create(ctx context.Context, cryptoCurrency interface{}) (string, error) {
	q := `
	INSERT INTO crypto_currency (name, holdingid, symbol, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING crypto_currencyid;
    `

	c := cryptoCurrency.(entities.CryptoCurrency)
	row := r.DB.QueryRowContext(
		ctx, q, c.Name, c.HoldingID, c.Symbol, c.CreatedAt, c.UpdatedAt,
	)

	err := row.Scan(&c.CryptoCurrencyID)
	if err != nil {
		return "", err
	}

	return c.CryptoCurrencyID, nil
}

func (r *cryptoCurrencyRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT crypto_currencyid, holdingid, name, symbol, created_at, updated_at
	    FROM crypto_currency %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var cryptoCurrencys []interface{}
	for rows.Next() {
		var c entities.CryptoCurrency
		err := rows.Scan(&c.CryptoCurrencyID, &c.HoldingID, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cryptoCurrencys = append(cryptoCurrencys, &c)
	}

	if len(cryptoCurrencys) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return cryptoCurrencys, nil
}

func (r *cryptoCurrencyRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT crypto_currencyid, holdingid, name, symbol, created_at, updated_at
        FROM crypto_currency WHERE crypto_currencyid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.CryptoCurrency
	err := row.Scan(&c.CryptoCurrencyID, &c.HoldingID, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *cryptoCurrencyRepository) GetBySymbol(ctx context.Context, Symbol string) (interface{}, error) {
	q := `
    SELECT crypto_currencyid, holdingid, name, symbol, created_at, updated_at
	FROM crypto_currency WHERE symbol = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Symbol)

	var c entities.CryptoCurrency
	err := row.Scan(&c.CryptoCurrencyID, &c.CryptoCurrencyID, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *cryptoCurrencyRepository) GetByCode(ctx context.Context, Code string) (interface{}, error) {
	q := `
    SELECT crypto_currencyid, holdingid, name, symbol, created_at, updated_at
        FROM crypto_currency WHERE code = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Code)

	var c entities.CryptoCurrency
	err := row.Scan(&c.CryptoCurrencyID, &c.HoldingID, &c.Name, &c.Symbol, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *cryptoCurrencyRepository) Update(ctx context.Context, ID string, cryptoCurrency interface{}) error {
	q := `
	UPDATE crypto_currency set holdingid=$1, name=$2, symbol=$3, updated_at=$4
	    WHERE crypto_currencyid=$5;
	`

	b := cryptoCurrency.(entities.CryptoCurrency)
	result, err := r.DB.ExecContext(
		ctx, q, b.HoldingID, b.Name, b.Symbol, b.UpdatedAt, ID,
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

func (r *cryptoCurrencyRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM crypto_currency WHERE crypto_currencyid=$1;`

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

func (r *cryptoCurrencyRepository) CreateMany(ctx context.Context, cryptoCurrency_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range cryptoCurrency_models {
		c := entity.(entities.CryptoCurrency)

		q := `
		INSERT INTO crypto_currency (holdingid, name, symbol, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING crypto_currencyid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, c.Name, c.HoldingID, c.Symbol, c.CreatedAt, c.UpdatedAt,
		)
		err := row.Scan(&c.CryptoCurrencyID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, c.CryptoCurrencyID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
