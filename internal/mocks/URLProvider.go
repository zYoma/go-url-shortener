// Code generated by mockery v2.39.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// URLProvider is an autogenerated mock type for the URLProvider type
type URLProvider struct {
	mock.Mock
}

// GetURL provides a mock function with given fields: ctx, shortURL
func (_m *URLProvider) GetURL(ctx context.Context, shortURL string) (string, error) {
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

// Init provides a mock function with given fields:
func (_m *URLProvider) Init() error {
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
func (_m *URLProvider) Ping(ctx context.Context) error {
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

// SaveURL provides a mock function with given fields: ctx, fullURL, shortURL
func (_m *URLProvider) SaveURL(ctx context.Context, fullURL string, shortURL string) error {
	ret := _m.Called(ctx, fullURL, shortURL)

	if len(ret) == 0 {
		panic("no return value specified for SaveURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, fullURL, shortURL)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewURLProvider creates a new instance of URLProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewURLProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *URLProvider {
	mock := &URLProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
