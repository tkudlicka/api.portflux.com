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

// TestNewTransactionRepository_Ok checks that NewTransactionRepository creates a new transactionRepository struct
func TestNewTransactionRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewTransactionRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateTransaction_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateTransaction_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newTransaction := entities.Transaction{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newTransaction)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateTransaction_InsertError checks that Create returns an error when the insert statement fails
func TestCreateTransaction_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newTransaction := entities.Transaction{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO transaction").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newTransaction)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetTransaction_Ok checks that Get returns the expected response when a valid filter is received
func TestGetTransaction_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedTransaction := entities.Transaction{
		TransactionID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"stockid": "test-STOCK", "cryptocurrencyid": "1"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid", "stockid", "cryptocurrencyid", "quantity", "transaction_price", "transaction_date", "created_at", "updated_at"}).
		AddRow(expectedTransaction.TransactionID, expectedTransaction.StockID, expectedTransaction.CryptocurrencyID, expectedTransaction.Quantity, expectedTransaction.TransactionPrice, expectedTransaction.TransactionDate, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Transaction))
	assert.Equal(t, expectedTransaction, entity)
}

// TestGetTransaction_SelectError checks that Get returns an error when the select query fails
func TestGetTransaction_SelectError(t *testing.T) {
	// TestGetTransaction_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM transaction").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetTransaction_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetTransaction_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid", "holdingid", "name", "symbol", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetTransactionByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetTransactionByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedTransaction := entities.Transaction{
		TransactionID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid", "stockid", "cryptocurrencyid", "quantity", "transaction_price", "transaction_date", "created_at", "updated_at"}).
		AddRow(expectedTransaction.TransactionID, expectedTransaction.StockID, expectedTransaction.CryptocurrencyID, expectedTransaction.Quantity, expectedTransaction.TransactionPrice, expectedTransaction.TransactionDate, expectedTransaction.CreatedAt, expectedTransaction.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedTransaction.TransactionID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Transaction))
	assert.Equal(t, expectedTransaction, entity)
}

// TestGetTransactionByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetTransactionByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM transaction").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetTransactionByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetTransactionByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid", "stockid", "cryptocurrencyid", "quantity", "transaction_price", "transaction_date", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateTransaction_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateTransaction_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE transaction").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Transaction{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateTransaction_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateTransaction_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE transaction").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Transaction{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateTransaction_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateTransaction_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE transaction").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Transaction{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteTransaction_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteTransaction_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM transaction").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteTransaction_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteTransaction_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM transaction").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteTransaction_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteTransaction_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM transaction").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyTransactions_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyTransactions_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newTransactions := []interface{}{entities.Transaction{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newTransactions)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyTransactions_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyTransactions_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newTransactions := []interface{}{entities.Transaction{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO transaction").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newTransactions)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyTransactions_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyTransactions_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newTransactions := []interface{}{entities.Transaction{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newTransactions)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyTransactions_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyTransactions_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &transactionRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newTransactions := []interface{}{entities.Transaction{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO transaction").WillReturnRows(sqlmock.NewRows([]string{"transactionid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newTransactions)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
