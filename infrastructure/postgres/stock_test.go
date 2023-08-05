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

// TestNewStockRepository_Ok checks that NewStockRepository creates a new stockRepository struct
func TestNewStockRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewStockRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateStock_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateStock_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newStock := entities.Stock{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO stock").WillReturnRows(sqlmock.NewRows([]string{"stockid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newStock)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateStock_InsertError checks that Create returns an error when the insert statement fails
func TestCreateStock_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newStock := entities.Stock{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO stock").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newStock)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetStock_Ok checks that Get returns the expected response when a valid filter is received
func TestGetStock_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedStock := entities.Stock{
		StockID:     "f8352727-231e-4de1-8257-c235a0af5c4a",
		CompanyName: "test-name",
		Slug:        "1",
	}
	filter := map[string]interface{}{"company_name": "test-name", "slug": "1"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM stock").WillReturnRows(sqlmock.NewRows([]string{"stockid", "holdingid", "extid", "ticker_symbol", "company_name", "slug", "created_at", "updated_at"}).
		AddRow(expectedStock.StockID, expectedStock.HoldingID, expectedStock.Extid, expectedStock.TickerSymbol, expectedStock.CompanyName, expectedStock.Slug, expectedStock.CreatedAt, expectedStock.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Stock))
	assert.Equal(t, expectedStock, entity)
}

// TestGetStock_SelectError checks that Get returns an error when the select query fails
func TestGetStock_SelectError(t *testing.T) {
	// TestGetStock_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM stock").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetStock_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetStock_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM stock").WillReturnRows(sqlmock.NewRows([]string{"stockid", "holdingid", "extid", "ticker_symbol", "company_name", "slug", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetStockByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetStockByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedStock := entities.Stock{
		StockID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM stock").WillReturnRows(sqlmock.NewRows([]string{"stockid", "holdingid", "extid", "ticker_symbol", "company_name", "slug", "created_at", "updated_at"}).
		AddRow(expectedStock.StockID, expectedStock.HoldingID, expectedStock.Extid, expectedStock.TickerSymbol, expectedStock.CompanyName, expectedStock.Slug, expectedStock.CreatedAt, expectedStock.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedStock.StockID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Stock))
	assert.Equal(t, expectedStock, entity)
}

// TestGetStockByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetStockByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM stock").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetStockByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetStockByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM stock").WillReturnRows(sqlmock.NewRows([]string{"stockid", "holdingid", "extid", "ticker_symbol", "company_name", "slug", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateStock_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateStock_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE stock").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Stock{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateStock_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateStock_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE stock").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Stock{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateStock_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateStock_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE stock").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Stock{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteStock_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteStock_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM stock").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteStock_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteStock_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM stock").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteStock_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteStock_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM stock").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyStocks_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyStocks_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newStocks := []interface{}{entities.Stock{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO stock").WillReturnRows(sqlmock.NewRows([]string{"stockid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newStocks)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyStocks_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyStocks_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newStocks := []interface{}{entities.Stock{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO stock").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newStocks)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyStocks_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyStocks_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newStocks := []interface{}{entities.Stock{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newStocks)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyStocks_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyStocks_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &stockRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newStocks := []interface{}{entities.Stock{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO stock").WillReturnRows(sqlmock.NewRows([]string{"stockid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newStocks)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
