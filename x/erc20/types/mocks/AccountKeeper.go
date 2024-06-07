// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// AccountKeeper is an autogenerated mock type for the AccountKeeper type
type AccountKeeper struct {
	mock.Mock
}

// GetAccount provides a mock function with given fields: _a0, _a1
func (_m *AccountKeeper) GetAccount(_a0 types.Context, _a1 types.AccAddress) authtypes.AccountI {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetAccount")
	}

	var r0 authtypes.AccountI
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) authtypes.AccountI); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(authtypes.AccountI)
		}
	}

	return r0
}

// GetSequence provides a mock function with given fields: _a0, _a1
func (_m *AccountKeeper) GetSequence(_a0 types.Context, _a1 types.AccAddress) (uint64, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetSequence")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) (uint64, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(types.Context, types.AccAddress) uint64); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(types.Context, types.AccAddress) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccountKeeper creates a new instance of AccountKeeper. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountKeeper(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountKeeper {
	mock := &AccountKeeper{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
