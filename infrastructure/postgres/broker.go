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

// brokerRepository adapter of an broker repository for postgres
type brokerRepository struct {
	infrastructure.PostgresRepository
}

// NewBrokerRepository creates a broker repository for postgres
func NewBrokerRepository(db *sql.DB) ports.BrokerRepository {
	return &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *brokerRepository) Create(ctx context.Context, broker interface{}) (string, error) {
	q := `
	INSERT INTO broker (extid, name, description, slug, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING brokerid;
    `

	b := broker.(entities.Broker)
	row := r.DB.QueryRowContext(
		ctx, q, b.Extid, b.Name, b.Description, b.Slug, b.CreatedAt, b.UpdatedAt,
	)

	err := row.Scan(&b.BrokerID)
	if err != nil {
		return "", err
	}

	return b.BrokerID, nil
}

func (r *brokerRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT brokerid, extid, name, description, slug, created_at, updated_at
	    FROM broker %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var brokers []interface{}
	for rows.Next() {
		var b entities.Broker
		err = rows.Scan(&b.BrokerID, &b.Extid, &b.Name, &b.Description, &b.Slug, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		brokers = append(brokers, &b)
	}

	if len(brokers) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return brokers, nil
}

func (r *brokerRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT brokerid, extid, name, description, slug, created_at, updated_at
        FROM broker WHERE brokerid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var b entities.Broker
	err := row.Scan(&b.BrokerID, &b.Extid, &b.Name, &b.Description, &b.Slug, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &b, nil
}

func (r *brokerRepository) GetBySlug(ctx context.Context, Slug string) (interface{}, error) {
	q := `
    SELECT brokerid, extid, name, description, slug, created_at, updated_at
        FROM broker WHERE slug = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, Slug)

	var b entities.Broker
	err := row.Scan(&b.BrokerID, &b.Extid, &b.Name, &b.Description, &b.Slug, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &b, nil
}

func (r *brokerRepository) Update(ctx context.Context, ID string, broker interface{}) error {
	q := `
	UPDATE broker set name=$1, description=$2, slug=$3, updated_at=$4
	    WHERE brokerid=$5;
	`

	b := broker.(entities.Broker)
	result, err := r.DB.ExecContext(
		ctx, q, b.Name, b.Description, b.Slug, b.UpdatedAt, ID,
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

func (r *brokerRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM broker WHERE brokerid=$1;`

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

func (r *brokerRepository) CreateMany(ctx context.Context, brokers []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range brokers {
		b := entity.(entities.Broker)

		q := `
		INSERT INTO broker (extid, name, description, slug, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING brokerid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, b.Extid, b.Name, b.Description, b.Slug, b.CreatedAt, b.UpdatedAt,
		)
		err := row.Scan(&b.BrokerID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, b.BrokerID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
