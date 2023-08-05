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

// TestNewHoldingRepository_Ok checks that NewHoldingRepository creates a new holdingRepository struct
func TestNewHoldingRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewHoldingRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateHolding_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateHolding_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newHolding := entities.Holding{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newHolding)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateHolding_InsertError checks that Create returns an error when the insert statement fails
func TestCreateHolding_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newHolding := entities.Holding{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO holding").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newHolding)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetHolding_Ok checks that Get returns the expected response when a valid filter is received
func TestGetHolding_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedHolding := entities.Holding{
		HoldingID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"name": "test-name", "holdingid": "1"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid", "portfolioid", "brokerid", "extid", "name", "description", "slug", "trade_date", "trade_type", "quantity", "share_price", "exchange_rate", "exchange_currencyid", "brokerage_unit_price", "brokerage_currency", "created_at", "updated_at"}).
		AddRow(expectedHolding.HoldingID, expectedHolding.PortfolioID, expectedHolding.BrokerID, expectedHolding.Extid, expectedHolding.Name, expectedHolding.Description, expectedHolding.Slug, expectedHolding.TradeDate, expectedHolding.TradeType, expectedHolding.Quantity, expectedHolding.SharePrice, expectedHolding.ExchangeRate, expectedHolding.ExchangeCurrencyID, expectedHolding.BrokerageUnitPrice, expectedHolding.BrokerageCurrency, expectedHolding.CreatedAt, expectedHolding.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Holding))
	assert.Equal(t, expectedHolding, entity)
}

// TestGetHolding_SelectError checks that Get returns an error when the select query fails
func TestGetHolding_SelectError(t *testing.T) {
	// TestGetHolding_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM holding").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetHolding_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetHolding_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid", "portfolioid", "brokerid", "extid", "name", "description", "slug", "trade_date", "trade_type", "quantity", "share_price", "exchange_rate", "exchange_currencyid", "brokerage_unit_price", "brokerage_currency", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetHoldingByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetHoldingByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedHolding := entities.Holding{
		HoldingID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid", "portfolioid", "brokerid", "extid", "name", "description", "slug", "trade_date", "trade_type", "quantity", "share_price", "exchange_rate", "exchange_currencyid", "brokerage_unit_price", "brokerage_currency", "created_at", "updated_at"}).
		AddRow(expectedHolding.HoldingID, expectedHolding.PortfolioID, expectedHolding.BrokerID, expectedHolding.Extid, expectedHolding.Name, expectedHolding.Description, expectedHolding.Slug, expectedHolding.TradeDate, expectedHolding.TradeType, expectedHolding.Quantity, expectedHolding.SharePrice, expectedHolding.ExchangeRate, expectedHolding.ExchangeCurrencyID, expectedHolding.BrokerageUnitPrice, expectedHolding.BrokerageCurrency, expectedHolding.CreatedAt, expectedHolding.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedHolding.HoldingID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Holding))
	assert.Equal(t, expectedHolding, entity)
}

// TestGetHoldingByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetHoldingByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM holding").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetHoldingByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetHoldingByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid", "portfolioid", "brokerid", "extid", "name", "description", "slug", "trade_date", "trade_type", "quantity", "share_price", "exchange_rate", "exchange_currencyid", "brokerage_unit_price", "brokerage_currency", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateHolding_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateHolding_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE holding").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Holding{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateHolding_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateHolding_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE holding").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Holding{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateHolding_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateHolding_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE holding").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Holding{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteHolding_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteHolding_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM holding").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteHolding_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteHolding_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM holding").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteHolding_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteHolding_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM holding").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyHoldings_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyHoldings_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newHoldings := []interface{}{entities.Holding{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newHoldings)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyHoldings_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyHoldings_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newHoldings := []interface{}{entities.Holding{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO holding").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newHoldings)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyHoldings_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyHoldings_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newHoldings := []interface{}{entities.Holding{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newHoldings)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyHoldings_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyHoldings_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &holdingRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newHoldings := []interface{}{entities.Holding{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO holding").WillReturnRows(sqlmock.NewRows([]string{"holdingid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newHoldings)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
