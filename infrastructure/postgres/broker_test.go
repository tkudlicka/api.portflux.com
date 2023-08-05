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

// TestNewBrokerRepository_Ok checks that NewBrokerRepository creates a new brokerRepository struct
func TestNewBrokerRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewBrokerRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateBroker_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateBroker_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newBroker := entities.Broker{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newBroker)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateBroker_InsertError checks that Create returns an error when the insert statement fails
func TestCreateBroker_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newBroker := entities.Broker{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO broker").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newBroker)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetBroker_Ok checks that Get returns the expected response when a valid filter is received
func TestGetBroker_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedBroker := entities.Broker{
		BrokerID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"slug": "test-email", "name": "test-name"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid", "name", "description", "extid", "slug", "created_at", "updated_at"}).
		AddRow(expectedBroker.BrokerID, expectedBroker.Name, expectedBroker.Description, expectedBroker.Extid, expectedBroker.Slug, expectedBroker.CreatedAt, expectedBroker.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.Broker))
	assert.Equal(t, expectedBroker, entity)
}

// TestGetBroker_SelectError checks that Get returns an error when the select query fails
func TestGetBroker_SelectError(t *testing.T) {
	// TestGetBroker_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM broker").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetBroker_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetBroker_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid", "name", "description", "extid", "slug", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetBrokerByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetBrokerByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedBroker := entities.Broker{
		BrokerID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid", "name", "description", "extid", "slug", "created_at", "updated_at"}).
		AddRow(expectedBroker.BrokerID, expectedBroker.Name, expectedBroker.Description, expectedBroker.Extid, expectedBroker.Slug, expectedBroker.CreatedAt, expectedBroker.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedBroker.BrokerID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.Broker))
	assert.Equal(t, expectedBroker, entity)
}

// TestGetBrokerByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetBrokerByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM broker").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetBrokerByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetBrokerByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid", "firstname", "lastname", "email", "password_hash", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateBroker_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateBroker_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE broker").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.Broker{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateBroker_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateBroker_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE broker").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.Broker{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateBroker_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateBroker_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE broker").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.Broker{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteBroker_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteBroker_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM broker").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteBroker_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteBroker_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM broker").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteBroker_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteBroker_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM broker").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyBrokers_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyBrokers_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newBrokers := []interface{}{entities.Broker{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newBrokers)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyBrokers_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyBrokers_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newBrokers := []interface{}{entities.Broker{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO broker").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newBrokers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyBrokers_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyBrokers_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newBrokers := []interface{}{entities.Broker{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newBrokers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyBrokers_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyBrokers_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &brokerRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newBrokers := []interface{}{entities.Broker{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO broker").WillReturnRows(sqlmock.NewRows([]string{"brokerid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newBrokers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
