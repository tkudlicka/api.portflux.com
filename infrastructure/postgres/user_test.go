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

// TestNewUserRepository_Ok checks that NewUserRepository creates a new userRepository struct
func TestNewUserRepository_Ok(t *testing.T) {
	// Arrange
	_, db := mocks.NewSqlDB(t)
	defer db.Close()

	// Act
	repo := NewUserRepository(db)

	// Assert
	assert.NotEmpty(t, repo)
}

// TestCreateUser_Ok checks that Create returns the expected response when a valid entity is received
func TestCreateUser_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUser := entities.User{}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectQuery("INSERT INTO \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid"}).AddRow(expectedID))

	// Act
	id, err := repo.Create(context.Background(), newUser)

	// Assert
	assert.Equal(t, expectedID, id)
	assert.Nil(t, err)
}

// TestCreateUser_InsertError checks that Create returns an error when the insert statement fails
func TestCreateUser_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUser := entities.User{}
	expectedError := "insert error"
	mock.ExpectQuery("INSERT INTO \"user\"").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Create(context.Background(), newUser)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetUser_Ok checks that Get returns the expected response when a valid filter is received
func TestGetUser_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedUser := entities.User{
		UserID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	filter := map[string]interface{}{"email": "test-email", "name": "test-name"}
	skip := 1
	take := 1
	mock.ExpectQuery("SELECT (.+) FROM \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid", "firstname", "lastname", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(expectedUser.UserID, expectedUser.Firstname, expectedUser.Lastname, expectedUser.Email, expectedUser.PasswordHash, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	result, err := repo.Get(context.Background(), filter, &skip, &take)

	// Assert
	assert.Nil(t, err)
	assert.True(t, len(result) == 1)

	entity := *(result[0].(*entities.User))
	assert.Equal(t, expectedUser, entity)
}

// TestGetUser_SelectError checks that Get returns an error when the select query fails
func TestGetUser_SelectError(t *testing.T) {
	// TestGetUser_SelectError
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM \"user\"").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetUser_NoResourcesFound checks that Get returns an error when no resources are found
func TestGetUser_NoResourcesFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid", "firstname", "lastname", "email", "password_hash", "created_at", "updated_at"}))

	// Act
	_, err := repo.Get(context.Background(), map[string]interface{}{}, nil, nil)

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestGetUserByID_Ok checks that GetByID returns the expected response when the received ID has a valid format
func TestGetUserByID_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedUser := entities.User{
		UserID: "f8352727-231e-4de1-8257-c235a0af5c4a",
	}
	mock.ExpectQuery("SELECT (.+) FROM \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid", "firstname", "lastname", "email", "password_hash", "created_at", "updated_at"}).
		AddRow(expectedUser.UserID, expectedUser.Firstname, expectedUser.Lastname, expectedUser.Email, expectedUser.PasswordHash, expectedUser.CreatedAt, expectedUser.UpdatedAt))

	// Act
	result, err := repo.GetByID(context.Background(), expectedUser.UserID)

	// Assert
	assert.Nil(t, err)

	entity := *(result.(*entities.User))
	assert.Equal(t, expectedUser, entity)
}

// TestGetUserByID_SelectError checks that GetByID returns an error when the select query fails
func TestGetUserByID_SelectError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	expectedError := "select error"
	mock.ExpectQuery("SELECT (.+) FROM \"user\"").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestGetUserByID_ResourceNotFound checks that GetByID returns an error when the resource is not found
func TestGetUserByID_ResourceNotFound(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid", "firstname", "lastname", "email", "password_hash", "created_at", "updated_at"}))

	// Act
	_, err := repo.GetByID(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestUpdateUser_Ok checks that Update does not return an error when the received ID has a valid format
func TestUpdateUser_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE \"user\"").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Update(context.Background(), "", entities.User{})

	// Assert
	assert.Nil(t, err)
}

// TestUpdateUser_UpdateError checks that Update returns an error when the update statement fails
func TestUpdateUser_UpdateError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "update error"
	mock.ExpectExec("UPDATE \"user\"").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Update(context.Background(), "", entities.User{})

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestUpdateUser_NotUpdatedError checks that Update returns an error when the update statement does not update any document
func TestUpdateUser_NotUpdatedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("UPDATE \"user\"").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Update(context.Background(), "", entities.User{})

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestDeleteUser_Ok checks that Delete does not return an error when the received ID has a valid format
func TestDeleteUser_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM \"user\"").WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Nil(t, err)
}

// TestDelteUser_DeleteError checks that Delete returns an error when the delete statement fails
func TestDelteUser_DeleteError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	expectedError := "delete error"
	mock.ExpectExec("DELETE FROM \"user\"").WillReturnError(fmt.Errorf(expectedError))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestDeleteUser_NotDeletedError checks that Delete returns an error when the delete statement does not delete any document
func TestDeleteUser_NotDeletedError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	mock.ExpectExec("DELETE FROM \"user\"").WillReturnResult(sqlmock.NewResult(1, 0))

	// Act
	err := repo.Delete(context.Background(), "")

	// Assert
	assert.Equal(t, wrappers.NewNonExistentErr(sql.ErrNoRows), err)
}

// TestCreateManyUsers_Ok checks that CreateMany does not return an error when everything goes as expected
func TestCreateManyUsers_Ok(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid"}).AddRow(expectedID))
	mock.ExpectCommit()

	// Act
	ids, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.True(t, len(ids) == 1)
	assert.Equal(t, expectedID, ids[0])
	assert.Nil(t, err)
}

// TestCreateManyUsers_InsertError checks that CreateMany returns an error when the insert statement fails
func TestCreateManyUsers_InsertError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedError := "insert error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"user\"").WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyUsers_BeginError checks that CreateMany returns an error when the begin statement fails
func TestCreateManyUsers_BeginError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedError := "begin error"
	mock.ExpectBegin().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}

// TestCreateManyUsers_CommitError checks that CreateMany returns an error when the commit statement fails
func TestCreateManyUsers_CommitError(t *testing.T) {
	// Arrange
	mock, db := mocks.NewSqlDB(t)
	defer db.Close()

	repo := &userRepository{
		infrastructure.PostgresRepository{
			DB: db,
		},
	}

	newUsers := []interface{}{entities.User{}}
	expectedID := "f8352727-231e-4de1-8257-c235a0af5c4a"
	expectedError := "commit error"
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"user\"").WillReturnRows(sqlmock.NewRows([]string{"userid"}).AddRow(expectedID))
	mock.ExpectCommit().WillReturnError(fmt.Errorf(expectedError))

	// Act
	_, err := repo.CreateMany(context.Background(), newUsers)

	// Assert
	assert.Equal(t, expectedError, err.Error())
}
