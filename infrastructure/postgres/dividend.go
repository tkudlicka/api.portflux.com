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

// dividendRepository adapter of an dividend repository for postgres
type dividendRepository struct {
	infrastructure.PostgresRepository
}

// NewDividendRepository creates a dividend repository for postgres
func NewDividendRepository(db *sql.DB) ports.DividendRepository {
	return &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *dividendRepository) Create(ctx context.Context, dividend interface{}) (string, error) {
	q := `
	INSERT INTO dividend (stockid, dividend_per_share, dividend_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING dividendid;
    `

	d := dividend.(entities.Dividend)
	row := r.DB.QueryRowContext(
		ctx, q, d.StockID, d.DividendPerShare, d.DividendDate, d.CreatedAt, d.UpdatedAt,
	)

	err := row.Scan(&d.DividendID)
	if err != nil {
		return "", err
	}

	return d.DividendID, nil
}

func (r *dividendRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT dividendid, stockid, dividend_per_share, dividend_date, created_at, updated_at
	    FROM dividend %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var dividends []interface{}
	for rows.Next() {
		var c entities.Dividend
		err := rows.Scan(&c.DividendID, &c.StockID, &c.DividendPerShare, &c.DividendDate, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		dividends = append(dividends, &c)
	}

	if len(dividends) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return dividends, nil
}

func (r *dividendRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT dividendid, stockid, dividend_per_share, dividend_date, created_at, updated_at
        FROM dividend WHERE dividendid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.Dividend
	err := row.Scan(&c.DividendID, &c.StockID, &c.DividendPerShare, &c.DividendDate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *dividendRepository) GetBySymbol(ctx context.Context, Symbol string) (interface{}, error) {
	q := `
    SELECT dividendid, stockid, dividend_per_share, dividend_date, created_at, updated_at
	FROM dividend WHERE symbol = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Symbol)

	var c entities.Dividend
	err := row.Scan(&c.DividendID, &c.StockID, &c.DividendPerShare, &c.DividendDate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *dividendRepository) GetByCode(ctx context.Context, Code string) (interface{}, error) {
	q := `
    SELECT dividendid, stockid, dividend_per_share, dividend_date, created_at, updated_at
        FROM dividend WHERE code = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Code)

	var c entities.Dividend
	err := row.Scan(&c.DividendID, &c.StockID, &c.DividendPerShare, &c.DividendDate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *dividendRepository) Update(ctx context.Context, ID string, dividend interface{}) error {
	q := `
	UPDATE dividend set dividend_per_share=$1, dividend_date=$2, updated_at=$3
	    WHERE dividendid=$4;
	`

	b := dividend.(entities.Dividend)
	result, err := r.DB.ExecContext(
		ctx, q, b.DividendPerShare, b.DividendDate, b.UpdatedAt, ID,
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

func (r *dividendRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM dividend WHERE dividendid=$1;`

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

func (r *dividendRepository) CreateMany(ctx context.Context, dividend_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range dividend_models {
		d := entity.(entities.Dividend)

		q := `
		INSERT INTO dividend (stockid, dividend_per_share, dividend_date, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING dividendid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, d.StockID, d.DividendPerShare, d.DividendDate, d.CreatedAt, d.UpdatedAt,
		)
		err := row.Scan(&d.DividendID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, d.DividendID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
