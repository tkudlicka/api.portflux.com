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

// TestNewDividendRepository_Ok checks that NewDividendRepository creates a new dividendRepository struct
func TestNewDividendRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewDividendRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateDividend_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateDividend_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newDividend := entities.Dividend{}

	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newDividend)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateDividend_OkInsertError checks that Create returns an error when the insert statement fails
func TestCreateDividend_OkInsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newDividend := entities.Dividend{}

	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO dividend").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newDividend)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetDividend_Ok checks that Get returns the expected response when a valid filter is received
func TestGetDividend_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expecteddividend := entities.Dividend{
		DividendID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"name": "test-name", "holdingid": "1"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid", "stockid", "dividend_per_share", "dividend_date", "created_at", "updated_at"}).
		AddRow(expecteddividend.DividendID, expecteddividend.StockID, expecteddividend.DividendPerShare, expecteddividend.DividendDate, expecteddividend.CreatedAt, expecteddividend.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Dividend))
	assert.Equal(t, expecteddividend, entity)
}

// TestGetDividend_SelectError checks that Get returns an error when the select query fails
func TestGetDividend_SelectError(t *testing.T) {
	// TestGetDividend_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM dividend").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetDividend_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetDividend_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid", "holdingid", "name", "symbol", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetdividendByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetdividendByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expecteddividend := entities.Dividend{
		DividendID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid", "stockid", "dividend_per_share", "dividend_date", "created_at", "updated_at"}).
		AddRow(expecteddividend.DividendID, expecteddividend.StockID, expecteddividend.DividendPerShare, expecteddividend.DividendDate, expecteddividend.CreatedAt, expecteddividend.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expecteddividend.DividendID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Dividend))
	assert.Equal(t, expecteddividend, entity)
}

// TestGetdividendByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetdividendByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM dividend").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetdividendByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetdividendByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid", "stockid", "dividend_per_share", "dividend_date", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateDividend_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateDividend_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE dividend").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Dividend{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateDividend_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateDividend_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE dividend").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Dividend{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateDividend_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateDividend_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE dividend").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Dividend{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteDividend_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteDividend_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM dividend").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteDividend_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteDividend_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM dividend").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteDividend_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteDividend_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM dividend").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyDividends_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyDividends_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newDividends := []interface{}{entities.Dividend{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newDividends)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyDividends_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyDividends_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newDividends := []interface{}{entities.Dividend{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO dividend").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newDividends)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyDividends_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyDividends_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newDividends := []interface{}{entities.Dividend{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newDividends)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyDividends_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyDividends_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &dividendRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newDividends := []interface{}{entities.Dividend{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO dividend").WillReturnRows(sqlmock.NewRows([]string{"dividendid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newDividends)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
