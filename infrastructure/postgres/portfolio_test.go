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

// TestNewPortfolioRepository_Ok checks that NewPortfolioRepository creates a new portfolioRepository struct
func TestNewPortfolioRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewPortfolioRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreatePortfolio_Ok checks that Create returns the expected response when a valid entity is received
func TestCreatePortfolio_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newPortfolio := entities.Portfolio{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newPortfolio)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreatePortfolio_InsertError checks that Create returns an error when the insert statement fails
func TestCreatePortfolio_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newPortfolio := entities.Portfolio{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO portfolio").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newPortfolio)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetPortfolio_Ok checks that Get returns the expected response when a valid filter is received
func TestGetPortfolio_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedPortfolio := entities.Portfolio{
		PortfolioID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"name": "test-name", "holdingid": "1"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid", "userid", "name", "extid", "tax_countryid", "financial_year", "performence_calculation", "summary", "price_alert", "company_event_alert", "created_at", "updated_at"}).
		AddRow(expectedPortfolio.PortfolioID, expectedPortfolio.Userid, expectedPortfolio.Name, expectedPortfolio.Extid, expectedPortfolio.TaxCountryID, expectedPortfolio.FinancialYear, expectedPortfolio.PerformenceCalculation, expectedPortfolio.Summary, expectedPortfolio.PriceAlert, expectedPortfolio.CompanyEventAlert, expectedPortfolio.CreatedAt, expectedPortfolio.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Portfolio))
	assert.Equal(t, expectedPortfolio, entity)
}

// TestGetPortfolio_SelectError checks that Get returns an error when the select query fails
func TestGetPortfolio_SelectError(t *testing.T) {
	// TestGetPortfolio_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM portfolio").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetPortfolio_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetPortfolio_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid", "userid", "name", "extid", "tax_countryid", "financial_year", "performence_calculation", "summary", "price_alert", "company_event_alert", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetPortfolioByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetPortfolioByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedPortfolio := entities.Portfolio{
		PortfolioID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid", "userid", "name", "extid", "tax_countryid", "financial_year", "performence_calculation", "summary", "price_alert", "company_event_alert", "created_at", "updated_at"}).
		AddRow(expectedPortfolio.PortfolioID, expectedPortfolio.Userid, expectedPortfolio.Name, expectedPortfolio.Extid, expectedPortfolio.TaxCountryID, expectedPortfolio.FinancialYear, expectedPortfolio.PerformenceCalculation, expectedPortfolio.Summary, expectedPortfolio.PriceAlert, expectedPortfolio.CompanyEventAlert, expectedPortfolio.CreatedAt, expectedPortfolio.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedPortfolio.PortfolioID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Portfolio))
	assert.Equal(t, expectedPortfolio, entity)
}

// TestGetPortfolioByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetPortfolioByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM portfolio").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetPortfolioByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetPortfolioByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid", "userid", "name", "extid", "tax_countryid", "financial_year", "performence_calculation", "summary", "price_alert", "company_event_alert", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdatePortfolio_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdatePortfolio_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE portfolio").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Portfolio{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdatePortfolio_UpdateError checks that Update returns an error when the update statement fails
func TestUpdatePortfolio_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE portfolio").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Portfolio{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdatePortfolio_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdatePortfolio_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE portfolio").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Portfolio{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeletePortfolio_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeletePortfolio_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM portfolio").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDeltePortfolio_DeleteError checks that Delete returns an error when the delete statement fails
func TestDeltePortfolio_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM portfolio").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeletePortfolio_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeletePortfolio_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM portfolio").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyPortfolios_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyPortfolios_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newPortfolios := []interface{}{entities.Portfolio{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newPortfolios)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyPortfolios_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyPortfolios_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newPortfolios := []interface{}{entities.Portfolio{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO portfolio").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newPortfolios)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyPortfolios_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyPortfolios_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newPortfolios := []interface{}{entities.Portfolio{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newPortfolios)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyPortfolios_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyPortfolios_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &portfolioRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newPortfolios := []interface{}{entities.Portfolio{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO portfolio").WillReturnRows(sqlmock.NewRows([]string{"portfolioid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newPortfolios)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
