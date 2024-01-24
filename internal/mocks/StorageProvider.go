// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	models "github.com/zYoma/go-url-shortener/internal/models"
)

// StorageProvider is an autogenerated mock type for the StorageProvider type
type StorageProvider struct {
	mock.Mock
}

// BulkSaveURL provides a mock function with given fields: ctx, data, userID
func (_m *StorageProvider) BulkSaveURL(ctx context.Context, data []models.InsertData, userID string) error {
	ret := _m.Called(ctx, data, userID)

	if len(ret) == 0 {
		panic("no return value specified for BulkSaveURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []models.InsertData, string) error); ok {
		r0 = rf(ctx, data, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetShortURL provides a mock function with given fields: ctx, shortURL
func (_m *StorageProvider) GetShortURL(ctx context.Context, shortURL string) (string, error) {
	ret := _m.Called(ctx, shortURL)

	if len(ret) == 0 {
		panic("no return value specified for GetShortURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, shortURL)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, shortURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shortURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetURL provides a mock function with given fields: ctx, shortURL
func (_m *StorageProvider) GetURL(ctx context.Context, shortURL string) (string, error) {
	ret := _m.Called(ctx, shortURL)

	if len(ret) == 0 {
		panic("no return value specified for GetURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, shortURL)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, shortURL)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, shortURL)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserURLs provides a mock function with given fields: ctx, baseURL, userID
func (_m *StorageProvider) GetUserURLs(ctx context.Context, baseURL string, userID string) ([]models.UserURLS, error) {
	ret := _m.Called(ctx, baseURL, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserURLs")
	}

	var r0 []models.UserURLS
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) ([]models.UserURLS, error)); ok {
		return rf(ctx, baseURL, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []models.UserURLS); ok {
		r0 = rf(ctx, baseURL, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.UserURLS)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, baseURL, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields:
func (_m *StorageProvider) Init() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ping provides a mock function with given fields: ctx
func (_m *StorageProvider) Ping(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Ping")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveURL provides a mock function with given fields: ctx, fullURL, shortURL, userID
func (_m *StorageProvider) SaveURL(ctx context.Context, fullURL string, shortURL string, userID string) error {
	ret := _m.Called(ctx, fullURL, shortURL, userID)

	if len(ret) == 0 {
		panic("no return value specified for SaveURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, fullURL, shortURL, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorageProvider creates a new instance of StorageProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorageProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *StorageProvider {
	mock := &StorageProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
