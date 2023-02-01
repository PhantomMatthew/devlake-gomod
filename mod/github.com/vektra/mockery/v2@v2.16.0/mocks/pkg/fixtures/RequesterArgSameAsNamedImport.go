// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	json "encoding/json"

	mock "github.com/stretchr/testify/mock"
)

// RequesterArgSameAsNamedImport is an autogenerated mock type for the RequesterArgSameAsNamedImport type
type RequesterArgSameAsNamedImport struct {
	mock.Mock
}

// Get provides a mock function with given fields: _a0
func (_m *RequesterArgSameAsNamedImport) Get(_a0 string) *json.RawMessage {
	ret := _m.Called(_a0)

	var r0 *json.RawMessage
	if rf, ok := ret.Get(0).(func(string) *json.RawMessage); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*json.RawMessage)
		}
	}

	return r0
}

type mockConstructorTestingTNewRequesterArgSameAsNamedImport interface {
	mock.TestingT
	Cleanup(func())
}

// NewRequesterArgSameAsNamedImport creates a new instance of RequesterArgSameAsNamedImport. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRequesterArgSameAsNamedImport(t mockConstructorTestingTNewRequesterArgSameAsNamedImport) *RequesterArgSameAsNamedImport {
	mock := &RequesterArgSameAsNamedImport{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}