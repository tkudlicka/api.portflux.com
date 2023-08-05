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

// transactionRepository adapter of an transaction repository for postgres
type transactionRepository struct {
	infrastructure.PostgresRepository
}

// NewTransactionRepository creates a transaction repository for postgres
func NewTransactionRepository(db *sql.DB) ports.TransactionRepository {
	return &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction interface{}) (string, error) {
	q := `
	INSERT INTO transaction (stockid, cryptocurrencyid, quantity,transaction_price,transaction_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7)
        RETURNING transactionid;
    `

	c := transaction.(entities.Transaction)
	row := r.DB.QueryRowContext(
		ctx, q, c.StockID, c.CryptocurrencyID, c.Quantity, c.TransactionPrice, c.TransactionDate, c.CreatedAt, c.UpdatedAt,
	)

	err := row.Scan(&c.TransactionID)
	if err != nil {
		return "", err
	}

	return c.TransactionID, nil
}

func (r *transactionRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT transactionid, stockid, cryptocurrencyid, quantity,transaction_price,transaction_date, created_at, updated_at
	    FROM transaction %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []interface{}
	for rows.Next() {
		var c entities.Transaction
		err := rows.Scan(&c.TransactionID, &c.StockID, &c.CryptocurrencyID, &c.Quantity, &c.TransactionPrice, &c.TransactionDate, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &c)
	}

	if len(transactions) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return transactions, nil
}

func (r *transactionRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT transactionid, stockid, cryptocurrencyid, quantity,transaction_price,transaction_date, created_at, updated_at
        FROM transaction WHERE transactionid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var c entities.Transaction
	err := row.Scan(&c.TransactionID, &c.StockID, &c.CryptocurrencyID, &c.Quantity, &c.TransactionPrice, &c.TransactionDate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *transactionRepository) GetBySymbol(ctx context.Context, Symbol string) (interface{}, error) {
	q := `
    SELECT transactionid, stockid, cryptocurrencyid, quantity,transaction_price,transaction_date, created_at, updated_at
	FROM transaction WHERE symbol = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Symbol)

	var c entities.Transaction
	err := row.Scan(&c.TransactionID, &c.StockID, &c.CryptocurrencyID, &c.Quantity, &c.TransactionPrice, &c.TransactionDate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *transactionRepository) GetByCode(ctx context.Context, Code string) (interface{}, error) {
	q := `
    SELECT transactionid, stockid, cryptocurrencyid, quantity,transaction_price,transaction_date, created_at, updated_at
        FROM transaction WHERE code = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Code)

	var c entities.Transaction
	err := row.Scan(&c.TransactionID, &c.StockID, &c.CryptocurrencyID, &c.Quantity, &c.TransactionPrice, &c.TransactionDate, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &c, nil
}

func (r *transactionRepository) Update(ctx context.Context, ID string, transaction interface{}) error {
	q := `
	UPDATE transaction set stockid=$1, cryptocurrencyid=$2, quantity=$3,transaction_price=$4,transaction_date=$5, updated_at=$6
	    WHERE transactionid=$5;
	`

	b := transaction.(entities.Transaction)
	result, err := r.DB.ExecContext(
		ctx, q, b.StockID, b.CryptocurrencyID, b.Quantity, b.TransactionPrice, b.TransactionDate, b.UpdatedAt, ID,
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

func (r *transactionRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM transaction WHERE transactionid=$1;`

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

func (r *transactionRepository) CreateMany(ctx context.Context, transaction_models []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range transaction_models {
		c := entity.(entities.Transaction)

		q := `
		INSERT INTO transaction (stockid, cryptocurrencyid, quantity,transaction_price,transaction_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6,$7)
			RETURNING transactionid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, c.StockID, c.CryptocurrencyID, c.Quantity, c.TransactionPrice, c.TransactionDate, c.CreatedAt, c.UpdatedAt,
		)
		err := row.Scan(&c.TransactionID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, c.TransactionID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
