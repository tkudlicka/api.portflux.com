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

// TestNewCurrencyRepository_Ok checks that NewCurrencyRepository creates a new currencyRepository struct
func TestNewCurrencyRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewCurrencyRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateCurrency_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCurrency := entities.Currency{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newCurrency)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateCurrency_InsertError checks that Create returns an error when the insert statement fails
func TestCreateCurrency_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCurrency := entities.Currency{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newCurrency)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetCurrency_Ok checks that Get returns the expected response when a valid filter is received
func TestGetCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedCurrency := entities.Currency{
		CurrencyID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"email": "test-email", "name": "test-name"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid", "code", "name", "symbol", "created_at", "updated_at"}).
		AddRow(expectedCurrency.CurrencyID, expectedCurrency.Code, expectedCurrency.Name, expectedCurrency.Symbol, expectedCurrency.CreatedAt, expectedCurrency.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Currency))
	assert.Equal(t, expectedCurrency, entity)
}

// TestGetCurrency_SelectError checks that Get returns an error when the select query fails
func TestGetCurrency_SelectError(t *testing.T) {
	// TestGetCurrency_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetCurrency_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetCurrency_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid", "code", "name", "symbol", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetCurrencyByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetCurrencyByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedCurrency := entities.Currency{
		CurrencyID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid", "code", "name", "symbol", "created_at", "updated_at"}).
		AddRow(expectedCurrency.CurrencyID, expectedCurrency.Code, expectedCurrency.Name, expectedCurrency.Symbol, expectedCurrency.CreatedAt, expectedCurrency.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedCurrency.CurrencyID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Currency))
	assert.Equal(t, expectedCurrency, entity)
}

// TestGetCurrencyByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetCurrencyByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetCurrencyByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetCurrencyByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid", "code", "name", "symbol", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateCurrency_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE currency").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Currency{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateCurrency_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateCurrency_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Currency{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateCurrency_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateCurrency_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE currency").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Currency{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteCurrency_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteCurrency_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM currency").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteCurrency_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteCurrency_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteCurrency_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteCurrency_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM currency").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyCurrencys_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyCurrencys_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCurrencys := []interface{}{entities.Currency{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newCurrencys)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyCurrencys_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyCurrencys_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCurrencys := []interface{}{entities.Currency{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO currency").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newCurrencys)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyCurrencys_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyCurrencys_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCurrencys := []interface{}{entities.Currency{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newCurrencys)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyCurrencys_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyCurrencys_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &currencyRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newCurrencys := []interface{}{entities.Currency{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO currency").WillReturnRows(sqlmock.NewRows([]string{"currencyid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newCurrencys)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
