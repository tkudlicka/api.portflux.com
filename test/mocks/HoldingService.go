// Code generated by mockery v2.32.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	models "github.com/tkudlicka/portflux-api/core/models"
)

// HoldingService is an autogenerated mock type for the HoldingService type
type HoldingService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, holding
func (_m *HoldingService) Create(ctx context.Context, holding models.CreateHoldingReq) (models.CreationResp, error) {
	ret := _m.Called(ctx, holding)

	var r0 models.CreationResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateHoldingReq) (models.CreationResp, error)); ok {
		return rf(ctx, holding)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.CreateHoldingReq) models.CreationResp); ok {
		r0 = rf(ctx, holding)
	} else {
		r0 = ret.Get(0).(models.CreationResp)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.CreateHoldingReq) error); ok {
		r1 = rf(ctx, holding)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateMany provides a mock function with given fields: ctx, holdings
func (_m *HoldingService) CreateMany(ctx context.Context, holdings []models.CreateHoldingReq) (models.MultiCreationResp, error) {
	ret := _m.Called(ctx, holdings)

	var r0 models.MultiCreationResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []models.CreateHoldingReq) (models.MultiCreationResp, error)); ok {
		return rf(ctx, holdings)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []models.CreateHoldingReq) models.MultiCreationResp); ok {
		r0 = rf(ctx, holdings)
	} else {
		r0 = ret.Get(0).(models.MultiCreationResp)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []models.CreateHoldingReq) error); ok {
		r1 = rf(ctx, holdings)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, ID
func (_m *HoldingService) Delete(ctx context.Context, ID string) error {
	ret := _m.Called(ctx, ID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: ctx
func (_m *HoldingService) GetAll(ctx context.Context) ([]models.HoldingResp, error) {
	ret := _m.Called(ctx)

	var r0 []models.HoldingResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.HoldingResp, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.HoldingResp); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.HoldingResp)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, ID
func (_m *HoldingService) GetByID(ctx context.Context, ID string) (models.HoldingResp, error) {
	ret := _m.Called(ctx, ID)

	var r0 models.HoldingResp
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.HoldingResp, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.HoldingResp); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Get(0).(models.HoldingResp)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, ID, holding
func (_m *HoldingService) Update(ctx context.Context, ID string, holding models.UpdateHoldingReq) error {
	ret := _m.Called(ctx, ID, holding)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.UpdateHoldingReq) error); ok {
		r0 = rf(ctx, ID, holding)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewHoldingService creates a new instance of HoldingService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHoldingService(t interface {
	mock.TestingT
	Cleanup(func())
}) *HoldingService {
	mock := &HoldingService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
