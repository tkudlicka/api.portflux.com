package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/tkudlicka/portflux-api/config"
	"github.com/tkudlicka/portflux-api/core/entities"
	"github.com/tkudlicka/portflux-api/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestLoginUser_Ok checks that Login endpoint returns the expected response when everything goes as expected
func TestLoginUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		testUser.Email = "testlogin@test.com"
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		body := models.CreateUserReq{
			Email:        "testlogin@test.com",
			PasswordHash: "test",
		}
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/user/login", cfg.Port)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		var response models.LoginUserResp
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.Equal(t, testUser.UserID, response.User.UserID)
		assert.NotEmpty(t, response.Token)
	})
}

// TestCreateUser checks that CreateUser endpoint returns the expected response when everything goes as expected
func TestCreateUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()

		// Act
		body := models.CreateUserReq(testUser)
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/user", cfg.Port)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusCreated, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		var response models.CreationResp
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.NotEmpty(t, response.InsertedID)
		createdUser, err := findUser(response.InsertedID, cfg)
		if err != nil {
			t.Fatalf("unexpected error while finding the created user: %s", err)
		}
		assert.Equal(t, testUser.Firstname, createdUser.Firstname)
		assert.Equal(t, testUser.Lastname, createdUser.Lastname)
		assert.Equal(t, testUser.Email, createdUser.Email)
	})
}

// TestCreateManyUsers_Ok checks that CreateManyUsers endpoint returns the expected response when everything goes as expected
func TestCreateManyUsers_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		users := []entities.User{getNewTestUser(), getNewTestUser()}

		// Act
		body := []models.CreateUserReq{
			models.CreateUserReq(users[0]),
			models.CreateUserReq(users[1]),
		}
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/user/many", cfg.Port)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusCreated, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		var response models.MultiCreationResp
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.Equal(t, 2, len(response.InsertedIDs))
		for i, id := range response.InsertedIDs {
			createdUser1, err := findUser(id, cfg)
			if err != nil {
				t.Fatalf("unexpected error while finding the created user: %s", err)
			}
			assert.Equal(t, users[i].Firstname, createdUser1.Firstname)
			assert.Equal(t, users[i].Lastname, createdUser1.Lastname)
			assert.Equal(t, users[i].Email, createdUser1.Email)
		}
	})
}

// TestGetAllUsers_Ok checks that GetAllUsers endpoint returns the expected response when everything goes as expected
func TestGetAllUsers_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)

		// Act
		url := fmt.Sprintf("http://:%d/v1/user", cfg.Port)

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", nonExpiryToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		var response []models.UserResp
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.NotEmpty(t, response)
	})
}

// TestGetUserByEmail_Ok checks that GetUserByEmail endpoint returns the expected response when everything goes as expected
func TestGetUserByEmail_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("http://:%d/v1/user/email/%s", cfg.Port, testUser.Email)

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", nonExpiryToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		var response models.UserResp
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.Equal(t, testUser.Email, response.Email)
	})
}

// TestGetUserByID_Ok checks that GetUserByID endpoint returns the expected response when everything goes as expected
func TestGetUserByID_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		url := fmt.Sprintf("http://:%d/v1/user/%s", cfg.Port, testUser.UserID)

		req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", nonExpiryToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		var response models.UserResp
		if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatalf("unexpected error parsing the response while calling %s: %s", resp.Request.URL, err)
		}
		assert.Equal(t, testUser.UserID, response.UserID)
	})
}

// TestUpdateUser_Ok checks that UpdateUser endpoint returns the expected response when everything goes as expected
func TestUpdateUser_Ok(t *testing.T) {
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		testUser.Firstname = "modified"
		testUser.Lastname = "modified"
		body := models.CreateUserReq(testUser)
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/user/%s", cfg.Port, testUser.UserID)

		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", nonExpiryToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		updatedUser, err := findUser(testUser.UserID, cfg)
		if err != nil {
			t.Fatalf("unexpected error while finding the created user: %s", err)
		}
		assert.Equal(t, testUser.Firstname, updatedUser.Firstname)
		assert.Equal(t, testUser.Lastname, updatedUser.Lastname)
		assert.Equal(t, testUser.Email, updatedUser.Email)
	})
}

// TestDeleteUser_Ok checks that DeleteUser endpoint returns the expected response when everything goes as expected
func TestDeleteUser_Ok(t *testing.T) {
	notFoundError := map[string]error{"postgres": sql.ErrNoRows}
	Databases(t, func(t *testing.T, database string) {
		// Arrange
		cfg := New(t, database)
		testUser := getNewTestUser()
		err := insertUser(&testUser, cfg)
		if err != nil {
			t.Fatal(err)
		}

		// Act
		body := models.CreateUserReq(testUser)
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}

		url := fmt.Sprintf("http://:%d/v1/user/%s", cfg.Port, testUser.UserID)

		req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", contentType)
		req.Header.Set("Authorization", nonExpiryToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		// Assert
		if want, got := http.StatusOK, resp.StatusCode; want != got {
			t.Fatalf("unexpected http status code while calling %s: want=%d but got=%d", resp.Request.URL, want, got)
		}
		_, err = findUser(testUser.UserID, cfg)
		assert.Equal(t, notFoundError[database], err)
	})
}

// HELP FUNCTIONS

func getNewTestUser() entities.User {
	return entities.User{
		Firstname:    "test",
		Lastname:     "test",
		Email:        fmt.Sprintf("test%d@test.com", rand.Int()),
		PasswordHash: "$2a$10$Q71DDcyvQhzt2K1EbRp1cOh4ToUh9de9ETsixwXGOVeRorTh8tjN2", // test hashed
	}
}

func insertUser(u *entities.User, cfg config.Config) error {
	switch cfg.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.DSN)
		if err != nil {
			return err
		}

		result, err := db.Collection(entities.EntityNameUser).InsertOne(context.Background(), u)
		u.UserID = result.InsertedID.(primitive.ObjectID).Hex()
		return err

	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(context.Background(), cfg.DSN)
		if err != nil {
			return err
		}

		q := `
		INSERT INTO "user" (firstname, lastname, email, password_hash, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING userid, firstname, lastname, email, password_hash, created_at, updated_at;
		`

		row := db.QueryRowContext(
			context.Background(), q, u.Firstname, u.Lastname, u.Email, u.PasswordHash, time.Now().UTC(), time.Now().UTC(),
		)

		err = row.Scan(&u.UserID, &u.Firstname, &u.Lastname, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		return err

	default:
		return fmt.Errorf("database flag %s not valid", cfg.Database)
	}
}

func findUser(ID string, cfg config.Config) (entities.User, error) {
	switch cfg.Database {
	case "mongo":
		db, err := infrastructure.ConnectMongoDB(context.Background(), cfg.DSN)
		if err != nil {
			return entities.User{}, err
		}

		objectID, err := primitive.ObjectIDFromHex(ID)
		if err != nil {
			return entities.User{}, err
		}

		var u entities.User
		err = db.Collection(entities.EntityNameUser).FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&u)
		return u, err

	case "postgres":
		db, err := infrastructure.ConnectPostgresDB(context.Background(), cfg.DSN)
		if err != nil {
			return entities.User{}, err
		}

		q := `
		SELECT userid, firstname, lastname, email, password_hash, created_at, updated_at
			FROM "user" WHERE userid = $1;
		`

		row := db.QueryRowContext(context.Background(), q, ID)

		var u entities.User
		err = row.Scan(&u.UserID, &u.Firstname, &u.Lastname, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
		return u, err

	default:
		return entities.User{}, fmt.Errorf("database flag %s not valid", cfg.Database)
	}
}
