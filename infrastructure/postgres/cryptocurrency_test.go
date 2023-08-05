package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/mocks"
	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/stretchr/testify/assert"
	"github.com/tkudlicka/portflux-api/core/entities"
)

// TestNewCryptoCurrencyRepository_Ok checks that NewCryptoCurrencyRepository creates a new cryptoCurrencyRepository struct
func TestNewCryptoCurrencyRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewCryptoCurrencyRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateCryptoCurrency_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateCryptoCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCryptoCurrency := entities.CryptoCurrency{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newCryptoCurrency)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateCryptoCurrency_InsertError checks that Create returns an error when the insert statement fails
func TestCreateCryptoCurrency_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCryptoCurrency := entities.CryptoCurrency{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO crypto_currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newCryptoCurrency)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetCryptoCurrency_Ok checks that Get returns the expected response when a valid filter is received
func TestGetCryptoCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedCryptoCurrency := entities.CryptoCurrency{
		CryptoCurrencyID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"name": "test-name", "holdingid": "1"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid", "holdingid", "name", "symbol", "created_at", "updated_at"}).
		AddRow(expectedCryptoCurrency.CryptoCurrencyID, expectedCryptoCurrency.HoldingID, expectedCryptoCurrency.Name, expectedCryptoCurrency.Symbol, expectedCryptoCurrency.CreatedAt, expectedCryptoCurrency.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.CryptoCurrency))
	assert.Equal(t, expectedCryptoCurrency, entity)
}

// TestGetCryptoCurrency_SelectError checks that Get returns an error when the select query fails
func TestGetCryptoCurrency_SelectError(t *testing.T) {
	// TestGetCryptoCurrency_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM crypto_currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetCryptoCurrency_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetCryptoCurrency_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid", "holdingid", "name", "symbol", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetCryptoCurrencyByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetCryptoCurrencyByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedCryptoCurrency := entities.CryptoCurrency{
		CryptoCurrencyID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid", "holdingid", "name", "symbol", "created_at", "updated_at"}).
		AddRow(expectedCryptoCurrency.CryptoCurrencyID, expectedCryptoCurrency.HoldingID, expectedCryptoCurrency.Name, expectedCryptoCurrency.Symbol, expectedCryptoCurrency.CreatedAt, expectedCryptoCurrency.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedCryptoCurrency.CryptoCurrencyID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.CryptoCurrency))
	assert.Equal(t, expectedCryptoCurrency, entity)
}

// TestGetCryptoCurrencyByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetCryptoCurrencyByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM crypto_currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetCryptoCurrencyByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetCryptoCurrencyByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid", "holdingid", "name", "symbol", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateCryptoCurrency_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateCryptoCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE crypto_currency").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.CryptoCurrency{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateCryptoCurrency_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateCryptoCurrency_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE crypto_currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.CryptoCurrency{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateCryptoCurrency_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateCryptoCurrency_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE crypto_currency").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.CryptoCurrency{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteCryptoCurrency_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteCryptoCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM crypto_currency").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteCryptoCurrency_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteCryptoCurrency_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM crypto_currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteCryptoCurrency_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteCryptoCurrency_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM crypto_currency").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyCryptoCurrencys_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyCryptoCurrencys_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCryptoCurrencys := []interface{}{entities.CryptoCurrency{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newCryptoCurrencys)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyCryptoCurrencys_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyCryptoCurrencys_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCryptoCurrencys := []interface{}{entities.CryptoCurrency{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO crypto_currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newCryptoCurrencys)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyCryptoCurrencys_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyCryptoCurrencys_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCryptoCurrencys := []interface{}{entities.CryptoCurrency{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newCryptoCurrencys)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyCryptoCurrencys_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyCryptoCurrencys_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &cryptoCurrencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCryptoCurrencys := []interface{}{entities.CryptoCurrency{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO crypto_currency").WillReturnRows(sqlmock.NewRows([]string{"crypto_currencyid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newCryptoCurrencys)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
