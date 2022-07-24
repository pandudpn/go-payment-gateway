// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// RequestInterface is an autogenerated mock type for the RequestInterface type
type RequestInterface struct {
	mock.Mock
}

// DoRequest provides a mock function with given fields: ctx
func (_m *RequestInterface) DoRequest(ctx context.Context) ([]byte, int, error) {
	ret := _m.Called(ctx)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(context.Context) []byte); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context) int); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context) error); ok {
		r2 = rf(ctx)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// SetBasicAuth provides a mock function with given fields: username, password
func (_m *RequestInterface) SetBasicAuth(username string, password string) {
	_m.Called(username, password)
}

// SetClient provides a mock function with given fields: client
func (_m *RequestInterface) SetClient(client *http.Client) {
	_m.Called(client)
}

// SetHeader provides a mock function with given fields: header
func (_m *RequestInterface) SetHeader(header map[string]string) {
	_m.Called(header)
}

type mockConstructorTestingTNewRequestInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewRequestInterface creates a new instance of RequestInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRequestInterface(t mockConstructorTestingTNewRequestInterface) *RequestInterface {
	mock := &RequestInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
