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

// SetBrokerRoutes creates broker routes
func SetBrokerRoutes(ctx context.Context, cfg config.Config, r *mux.Router, s ports.BrokerService) {
	r.Handle("/v1/broker", createBroker(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/v1/broker/many", createManyBrokers(ctx, cfg, s)).Methods(http.MethodPost)
	r.Handle("/v1/broker", middlewares.JWT(getAllBrokers(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/v1/broker/slug/{slug}", middlewares.JWT(getBrokerBySlug(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/v1/broker/{id}", middlewares.JWT(getBrokerByID(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
	r.Handle("/v1/broker/{id}", middlewares.JWT(updateBroker(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodPatch)
	r.Handle("/v1/broker/{id}", middlewares.JWT(deleteBroker(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{"admin": true})).Methods(http.MethodDelete)
	//r.Handle("/v1/claims", middlewares.JWT(getBrokerClaims(ctx, cfg, s), cfg.JWTSecret, jwt.MapClaims{})).Methods(http.MethodGet)
}

// @Summary Create broker
// @Description Creates a new broker
// @Tags Brokers
// @Param broker body models.CreateBrokerReq true "New broker to be created"
// @Success 201 {object} models.CreationResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers [post]
func createBroker(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var broker models.CreateBrokerReq
		err = json.Unmarshal(body, &broker)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		result, err := s.Create(ctx, broker)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusCreated, result)
	})
}

// @Summary Create many brokers
// @Description Creates many brokers atomically
// @Tags Brokers
// @Param brokers body []models.CreateBrokerReq true "New brokers to be created"
// @Success 201 {object} models.MultiCreationResp "OK"
// @Failure 400 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers/many [post]
func createManyBrokers(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var brokers []models.CreateBrokerReq
		err = json.Unmarshal(body, &brokers)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		result, err := s.CreateMany(ctx, brokers)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusCreated, result)
	})
}

// @Summary Get all brokers
// @Description Gets all the brokers
// @Tags Brokers
// @Security Bearer
// @Success 200 {array} models.BrokerResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers [get]
func getAllBrokers(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		brokers, err := s.GetAll(ctx)
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, brokers)
	})
}

// @Summary Get broker by slug
// @Description Gets a broker by slug
// @Tags Brokers
// @Security Bearer
// @Param slug path string true "Slug"
// @Success 200 {object} models.BrokerResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers/slug/{slug} [get]
func getBrokerBySlug(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		broker, err := s.GetBySlug(ctx, params["slug"])
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, broker)
	})
}

// @Summary Get broker by ID
// @Description Gets a broker by ID
// @Tags Brokers
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 {object} models.BrokerResp "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers/{id} [get]
func getBrokerByID(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		var params = mux.Vars(r)
		broker, err := s.GetByID(ctx, params["id"])
		if err != nil {
			utils.ResponseError(w, r, nil, err)
			return
		}
		utils.ResponseJSON(w, r, nil, http.StatusOK, broker)
	})
}

// @Summary Update broker
// @Description Updates a broker
// @Tags Brokers
// @Security Bearer
// @Param id path string true "ID"
// @Param Broker body models.UpdateBrokerReq true "Broker"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers/{id} [patch]
func updateBroker(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
	return middlewares.Recover(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, cfg.Timeout.Duration)
		defer cancel()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		var params = mux.Vars(r)
		var broker models.UpdateBrokerReq
		err = json.Unmarshal(body, &broker)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}

		err = s.Update(ctx, params["id"], broker)
		if err != nil {
			utils.ResponseError(w, r, body, err)
			return
		}
		utils.ResponseJSON(w, r, body, http.StatusOK, nil)
	})
}

// @Summary Delete broker
// @Description Delete a broker
// @Tags Brokers
// @Security Bearer
// @Param id path string true "ID"
// @Success 200 "OK"
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 408 {object} object
// @Failure 500 {object} object
// @Router /v1/brokers/{id} [delete]
func deleteBroker(ctx context.Context, cfg config.Config, s ports.BrokerService) http.Handler {
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
