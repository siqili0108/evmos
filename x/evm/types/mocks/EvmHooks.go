// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	core "github.com/ethereum/go-ethereum/core"
	coretypes "github.com/ethereum/go-ethereum/core/types"

	mock "github.com/stretchr/testify/mock"

	types "github.com/cosmos/cosmos-sdk/types"
)

// EvmHooks is an autogenerated mock type for the EvmHooks type
type EvmHooks struct {
	mock.Mock
}

// PostTxProcessing provides a mock function with given fields: ctx, msg, receipt
func (_m *EvmHooks) PostTxProcessing(ctx types.Context, msg core.Message, receipt *coretypes.Receipt) error {
	ret := _m.Called(ctx, msg, receipt)

	if len(ret) == 0 {
		panic("no return value specified for PostTxProcessing")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(types.Context, core.Message, *coretypes.Receipt) error); ok {
		r0 = rf(ctx, msg, receipt)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewEvmHooks creates a new instance of EvmHooks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEvmHooks(t interface {
	mock.TestingT
	Cleanup(func())
},
) *EvmHooks {
	mock := &EvmHooks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
