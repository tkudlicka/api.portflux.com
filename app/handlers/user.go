package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/sergicanet9/scv-go-tools/v3/api/middlewares"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
	"github.com/tkudlicka/portflux-api/config"
	"github.com/tkudlicka/portflux-api/core/models"
	"github.com/tkudlicka/portflux-api/core/ports"
)

// SetUserRoutes creates user routes
func SetUserRoutes(ctx context.Context, cfg config.Config, r *mux.Router, s ports.UserService) {
	r.Handle("/v1/user/login", loginUser(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/v1/user", createUser(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/v1/user/many", createManyUsers(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/v1/user", middlewares.JWT(getAllUsers(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/v1/user/email/{email}", middlewares.JWT(getUserByEmail(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/v1/user/{id}", middlewares.JWT(getUserByID(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/v1/user/{id}", middlewares.JWT(updateUser(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodPatch)
	r.Handle("/v1/user/{id}", middlewares.JWT(deleteUser(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{"admin": true})).Methods(http.MethodDelete)
	//r.Handle("/v1/claims", middlewares.JWT(getUserClaims(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
}

// @Summary Login user
// @Description Logs in an user
// @Tags Users
// @Param login body models.LoginUserReq true "Login request"
// @Success 200 {object} models.LoginUserResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user/login [post]
func loginUser(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var credentials models.LoginUserReq
		err = json.Unmarshal(body, &credentials)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		response, err := s.Login(ctx, credentials)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusOK, response)
	})
}

// @Summary Create user
// @Description Creates a new user
// @Tags Users
// @Param user body models.CreateUserReq true "New user to be created"
// @Success 201 {object} models.CreationResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user [post]
func createUser(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var user models.CreateUserReq
		err = json.Unmarshal(body, &user)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		result, err := s.Create(ctx, user)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusCreated, result)
	})
}

// @Summary Create many users
// @Description Creates many users atomically
// @Tags Users
// @Param users body []models.CreateUserReq true "New users to be created"
// @Success 201 {object} models.MultiCreationResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user/many [post]
func createManyUsers(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var users []models.CreateUserReq
		err = json.Unmarshal(body, &users)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		result, err := s.CreateMany(ctx, users)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusCreated, result)
	})
}

// @Summary Get all users
// @Description Gets all the users
// @Tags Users
// @Security Bearer
// @Success 200 {array} models.UserResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user [get]
func getAllUsers(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		users, err := s.GetAll(ctx)
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, users)
	})
}

// @Summary Get user by email
// @Description Gets a user by email
// @Tags Users
// @Security Bearer
// @Param email path string true "Email"
// @Success 200 {object} models.UserResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user/email/{email} [get]
func getUserByEmail(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		user, err := s.GetByEmail(ctx, params["email"])
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, user)
	})
}

// @Summary Get user by ID
// @Description Gets a user by ID
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 {object} models.UserResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user/{id} [get]
func getUserByID(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		user, err := s.GetByID(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, user)
	})
}

// @Summary Update user
// @Description Updates a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Param User body models.UpdateUserReq true "User"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user/{id} [patch]
func updateUser(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var params = mux.Vars(r)
		var user models.UpdateUserReq
		err = json.Unmarshal(body, &user)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		err = s.Update(ctx, params["id"], user)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusOK, nil)
	})
}

// @Summary Delete user
// @Description Delete a user
// @Tags Users
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/user/{id} [delete]
func deleteUser(ctx context.Context, cfg config.Config, s ports.UserService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		err := s.Delete(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, nil)
	})
}
