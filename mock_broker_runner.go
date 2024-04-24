// Code generated by mockery v2.20.2. DO NOT EDIT.

package trengin

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockBrokerRunner is an autogenerated mock type for the BrokerRunner type
type MockBrokerRunner struct {
	mock.Mock
}

// ChangeConditionalOrder provides a mock function with given fields: ctx, action
func (_m *MockBrokerRunner) ChangeConditionalOrder(ctx context.Context, action ChangeConditionalOrderAction) (Position, error) {
	ret := _m.Called(ctx, action)

	var r0 Position
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ChangeConditionalOrderAction) (Position, error)); ok {
		return rf(ctx, action)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ChangeConditionalOrderAction) Position); ok {
		r0 = rf(ctx, action)
	} else {
		r0 = ret.Get(0).(Position)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ChangeConditionalOrderAction) error); ok {
		r1 = rf(ctx, action)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ClosePosition provides a mock function with given fields: ctx, action
func (_m *MockBrokerRunner) ClosePosition(ctx context.Context, action ClosePositionAction) (Position, error) {
	ret := _m.Called(ctx, action)

	var r0 Position
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ClosePositionAction) (Position, error)); ok {
		return rf(ctx, action)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ClosePositionAction) Position); ok {
		r0 = rf(ctx, action)
	} else {
		r0 = ret.Get(0).(Position)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ClosePositionAction) error); ok {
		r1 = rf(ctx, action)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OpenPosition provides a mock function with given fields: ctx, action
func (_m *MockBrokerRunner) OpenPosition(ctx context.Context, action OpenPositionAction) (Position, PositionClosed, error) {
	ret := _m.Called(ctx, action)

	var r0 Position
	var r1 PositionClosed
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, OpenPositionAction) (Position, PositionClosed, error)); ok {
		return rf(ctx, action)
	}
	if rf, ok := ret.Get(0).(func(context.Context, OpenPositionAction) Position); ok {
		r0 = rf(ctx, action)
	} else {
		r0 = ret.Get(0).(Position)
	}

	if rf, ok := ret.Get(1).(func(context.Context, OpenPositionAction) PositionClosed); ok {
		r1 = rf(ctx, action)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(PositionClosed)
		}
	}

	if rf, ok := ret.Get(2).(func(context.Context, OpenPositionAction) error); ok {
		r2 = rf(ctx, action)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Run provides a mock function with given fields: ctx
func (_m *MockBrokerRunner) Run(ctx context.Context) error {
	ret := _m.Called(ctx)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockBrokerRunner interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockBrokerRunner creates a new instance of MockBrokerRunner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockBrokerRunner(t mockConstructorTestingTNewMockBrokerRunner) *MockBrokerRunner {
	mock := &MockBrokerRunner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}