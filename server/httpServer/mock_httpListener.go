// Code generated by mockery v2.43.2. DO NOT EDIT.

package httpServer

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// mockHttpListener is an autogenerated mock type for the httpListener type
type mockHttpListener struct {
	mock.Mock
}

// ListenAndServe provides a mock function with given fields:
func (_m *mockHttpListener) ListenAndServe() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ListenAndServe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListenAndServeTLS provides a mock function with given fields: certFile, keyFile
func (_m *mockHttpListener) ListenAndServeTLS(certFile string, keyFile string) error {
	ret := _m.Called(certFile, keyFile)

	if len(ret) == 0 {
		panic("no return value specified for ListenAndServeTLS")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(certFile, keyFile)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields: ctx
func (_m *mockHttpListener) Shutdown(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Shutdown")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// newMockHttpListener creates a new instance of mockHttpListener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockHttpListener(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockHttpListener {
	mock := &mockHttpListener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}