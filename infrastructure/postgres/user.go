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

// userRepository adapter of an user repository for postgres
type userRepository struct {
	infrastructure.PostgresRepository
}

// NewUserRepository creates a user repository for postgres
func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
}

func (r *userRepository) Create(ctx context.Context, user interface{}) (string, error) {
	q := `
	INSERT INTO "user" (firstname, lastname, email, password_hash, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING userid;
    `

	u := user.(entities.User)
	row := r.DB.QueryRowContext(
		ctx, q, u.Firstname, u.Lastname, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
	)

	err := row.Scan(&u.UserID)
	if err != nil {
		return "", err
	}

	return u.UserID, nil
}

func (r *userRepository) Get(ctx context.Context, filter map[string]interface{}, skip, take *int) ([]interface{}, error) {
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
	SELECT userid, firstname, lastname, email, password_hash, created_at, updated_at
	    FROM "user" %s;
	`, where)

	rows, err := r.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []interface{}
	for rows.Next() {
		var u entities.User
		err = rows.Scan(&u.UserID, &u.Firstname, &u.Lastname, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	if len(users) < 1 {
		return nil, wrappers.NewNonExistentErr(sql.ErrNoRows)
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, ID string) (interface{}, error) {
	q := `
    SELECT userid, firstname, lastname, email, password_hash, created_at, updated_at
        FROM "user" WHERE userid = $1;
    `

	row := r.DB.QueryRowContext(ctx, q, ID)

	var u entities.User
	err := row.Scan(&u.UserID, &u.Firstname, &u.Lastname, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = wrappers.NewNonExistentErr(err)
		}
		return nil, err
	}

	return &u, nil
}

func (r *userRepository) Update(ctx context.Context, ID string, user interface{}) error {
	q := `
	UPDATE "user" set firstname=$1, lastname=$2, email=$3, password_hash=$4, updated_at=$5
	    WHERE userid=$6;
	`

	u := user.(entities.User)
	result, err := r.DB.ExecContext(
		ctx, q, u.Firstname, u.Lastname, u.Email, u.PasswordHash, u.UpdatedAt, ID,
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

func (r *userRepository) Delete(ctx context.Context, ID string) error {
	q := `DELETE FROM "user" WHERE userid=$1;`

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

func (r *userRepository) CreateMany(ctx context.Context, users []interface{}) ([]string, error) {
	// `tx` is an instance of `*sql.Tx` through which we can execute our queries
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, entity := range users {
		u := entity.(entities.User)

		q := `
		INSERT INTO "user" (firstname, lastname, email, password_hash, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING userid;`

		// Here, the query is executed on the transaction instance, and not applied to the database yet
		row := tx.QueryRowContext(
			ctx, q, u.Firstname, u.Lastname, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt,
		)
		err := row.Scan(&u.UserID)
		if err != nil {
			// Incase we find any error in the query execution, rollback the transaction
			tx.Rollback()
			return nil, err
		}
		result = append(result, u.UserID)
	}

	// Finally, if no errors are recieved from the queries, commit the transaction
	// this applies the above changes to our database
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}
