// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	shared_models "service-shared/shared-models"
)

// PublishQueue is an autogenerated mock type for the PublishQueue type
type PublishQueue struct {
	mock.Mock
}

// PublishLoanRequest provides a mock function with given fields: loanApplication
func (_m *PublishQueue) PublishLoanRequest(loanApplication shared_models.CreateLoanMessage) error {
	ret := _m.Called(loanApplication)

	var r0 error
	if rf, ok := ret.Get(0).(func(shared_models.CreateLoanMessage) error); ok {
		r0 = rf(loanApplication)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPublishQueue interface {
	mock.TestingT
	Cleanup(func())
}

// NewPublishQueue creates a new instance of PublishQueue. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPublishQueue(t mockConstructorTestingTNewPublishQueue) *PublishQueue {
	mock := &PublishQueue{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
