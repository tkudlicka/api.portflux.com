package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sergicanet9/scv-go-tools/v3/wrappers"
	"github.com/tkudlicka/portflux-api/config"
	"github.com/tkudlicka/portflux-api/core/entities"
	"github.com/tkudlicka/portflux-api/core/models"
	"github.com/tkudlicka/portflux-api/core/ports"
)

// brokerService adapter of an broker service
type brokerService struct {
	config     config.Config
	repository ports.BrokerRepository
}

// NewBrokerService creates a new broker service
func NewBrokerService(cfg config.Config, repo ports.BrokerRepository) ports.BrokerService {
	return &brokerService{
		config:     cfg,
		repository: repo,
	}
}

// Create broker
func (s *brokerService) Create(ctx context.Context, broker models.CreateBrokerReq) (resp models.CreationResp, err error) {
	if err = broker.Validate(); err != nil {
		return
	}

	broker.Slug = broker.Name

	now := time.Now().UTC()
	broker.CreatedAt = now
	broker.UpdatedAt = now
	insertedID, err := s.repository.Create(ctx, entities.Broker(broker))
	if err != nil {
		return
	}

	resp = models.CreationResp{
		InsertedID: insertedID,
	}

	return
}

// CreateMany brokers
func (s *brokerService) CreateMany(ctx context.Context, brokers []models.CreateBrokerReq) (resp models.MultiCreationResp, err error) {
	var create []interface{}
	now := time.Now().UTC()

	for _, broker := range brokers {
		if err = broker.Validate(); err != nil {
			return
		}

		broker.CreatedAt = now
		broker.UpdatedAt = now

		create = append(create, entities.Broker(broker))
	}

	insertedIDs, err := s.repository.CreateMany(ctx, create)
	if err != nil {
		return
	}

	resp = models.MultiCreationResp{
		InsertedIDs: insertedIDs,
	}
	return
}

// GetAll brokers
func (s *brokerService) GetAll(ctx context.Context) (resp []models.BrokerResp, err error) {
	result, err := s.repository.Get(ctx, map[string]interface{}{}, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = nil
		}
		return
	}

	resp = make([]models.BrokerResp, len(result))
	for i, v := range result {
		resp[i] = models.BrokerResp(*(v.(*entities.Broker)))
	}

	return
}

// GetBySlug broker
func (s *brokerService) GetBySlug(ctx context.Context, slug string) (resp models.BrokerResp, err error) {
	filter := map[string]interface{}{"slug": slug}
	result, err := s.repository.Get(ctx, filter, nil, nil)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("slug %s not found", slug))
		}
		return
	}

	resp = models.BrokerResp(*(result[0].(*entities.Broker)))

	return
}

// GetByID broker
func (s *brokerService) GetByID(ctx context.Context, ID string) (resp models.BrokerResp, err error) {
	broker, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		if errors.Is(err, wrappers.NonExistentErr) {
			err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
		}
		return
	}

	resp = models.BrokerResp(*broker.(*entities.Broker))

	return
}

// Update user
func (s *brokerService) Update(ctx context.Context, ID string, broker models.UpdateBrokerReq) (err error) {
	dbBroker, err := s.GetByID(ctx, ID)
	if err != nil {
		return
	}

	if broker.Name != dbBroker.Name {
		dbBroker.Name = broker.Name
	}
	if broker.Description != dbBroker.Description {
		dbBroker.Description = broker.Description
	}

	dbBroker.UpdatedAt = time.Now().UTC()

	err = s.repository.Update(ctx, ID, entities.Broker(dbBroker))
	return err
}

// Delete broker
func (s *brokerService) Delete(ctx context.Context, ID string) (err error) {
	err = s.repository.Delete(ctx, ID)
	if errors.Is(err, wrappers.NonExistentErr) {
		err = wrappers.NewNonExistentErr(fmt.Errorf("ID %s not found", ID))
	}

	return
}
