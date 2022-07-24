// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Logging is an autogenerated mock type for the Logging type
type Logging struct {
	mock.Mock
}

// Error provides a mock function with given fields: args
func (_m *Logging) Error(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Errorf provides a mock function with given fields: format, args
func (_m *Logging) Errorf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Print provides a mock function with given fields: args
func (_m *Logging) Print(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Printf provides a mock function with given fields: format, args
func (_m *Logging) Printf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Println provides a mock function with given fields: args
func (_m *Logging) Println(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: args
func (_m *Logging) Warn(args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

// Warnf provides a mock function with given fields: format, args
func (_m *Logging) Warnf(format string, args ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, format)
	_ca = append(_ca, args...)
	_m.Called(_ca...)
}

type mockConstructorTestingTNewLogging interface {
	mock.TestingT
	Cleanup(func())
}

// NewLogging creates a new instance of Logging. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLogging(t mockConstructorTestingTNewLogging) *Logging {
	mock := &Logging{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
