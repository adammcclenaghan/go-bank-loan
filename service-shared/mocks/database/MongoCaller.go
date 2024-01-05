// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	mongo "go.mongodb.org/mongo-driver/mongo"

	options "go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCaller is an autogenerated mock type for the MongoCaller type
type MongoCaller struct {
	mock.Mock
}

// DeleteOne provides a mock function with given fields: ctx, filter, opts
func (_m *MongoCaller) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *mongo.DeleteResult
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.DeleteOptions) *mongo.DeleteResult); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.DeleteResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.DeleteOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, filter, opts
func (_m *MongoCaller) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *mongo.Cursor
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.FindOptions) *mongo.Cursor); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.Cursor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.FindOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindOne provides a mock function with given fields: ctx, filter, opts
func (_m *MongoCaller) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *mongo.SingleResult
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.SingleResult)
		}
	}

	return r0
}

// InsertOne provides a mock function with given fields: ctx, document, opts
func (_m *MongoCaller) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, document)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *mongo.InsertOneResult
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.InsertOneOptions) *mongo.InsertOneResult); ok {
		r0 = rf(ctx, document, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.InsertOneResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.InsertOneOptions) error); ok {
		r1 = rf(ctx, document, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateByID provides a mock function with given fields: ctx, id, update, opts
func (_m *MongoCaller) UpdateByID(ctx context.Context, id interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, id, update)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *mongo.UpdateResult
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) *mongo.UpdateResult); ok {
		r0 = rf(ctx, id, update, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, id, update, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewMongoCaller interface {
	mock.TestingT
	Cleanup(func())
}

// NewMongoCaller creates a new instance of MongoCaller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMongoCaller(t mockConstructorTestingTNewMongoCaller) *MongoCaller {
	mock := &MongoCaller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
