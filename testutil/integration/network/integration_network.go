// Copyright Tharsis Labs Ltd.(Evmos)
// SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/evmos/blob/main/LICENSE)
package network

import (
	"encoding/json"
	"math"
	"math/big"

	"github.com/evmos/evmos/v15/app"
	"github.com/evmos/evmos/v15/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var _ EvmosNetwork = (*IntegrationNetwork)(nil)

// IntegrationNetwork is the implementation of the Network interface for integration tests.
type IntegrationNetwork struct {
	BaseNetwork
	app *app.Evmos
}

// New configures and initializes a new integration Network instance with
// the given configuration options. If no configuration options are provided
// it uses the default configuration.
//
// It panics if an error occurs.
func New(opts ...ConfigOption) *IntegrationNetwork {
	cfg := DefaultConfig()
	// Modify the default config with the given options
	for _, opt := range opts {
		opt(&cfg)
	}

	ctx := sdktypes.Context{}
	network := &IntegrationNetwork{
		BaseNetwork: BaseNetwork{
			ctx:        ctx,
			cfg:        cfg,
			validators: []stakingtypes.Validator{},
		},
	}

	err := network.configureAndInitChain()
	if err != nil {
		panic(err)
	}
	return network
}

var (
	// bondedAmt is the amount of tokens that each validator will have initially bonded
	bondedAmt = sdktypes.TokensFromConsensusPower(1, types.PowerReduction)
	// PrefundedAccountInitialBalance is the amount of tokens that each prefunded account has at genesis
	PrefundedAccountInitialBalance = sdktypes.NewInt(int64(math.Pow10(18) * 4))
)

// configureAndInitChain initializes the network with the given configuration.
// It creates the genesis state and starts the network.
func (n *IntegrationNetwork) configureAndInitChain() error {
	// Create funded accounts based on the config and
	// create genesis accounts
	coin := sdktypes.NewCoin(n.cfg.denom, PrefundedAccountInitialBalance)
	genAccounts := createGenesisAccounts(n.cfg.preFundedAccounts)
	fundedAccountBalances := createBalances(n.cfg.preFundedAccounts, coin)

	// Create validator set with the amount of validators specified in the config
	// with the default power of 1.
	valSet := createValidatorSet(n.cfg.amountOfValidators)
	totalBonded := bondedAmt.Mul(sdktypes.NewInt(int64(n.cfg.amountOfValidators)))

	// Build staking type validators and delegations
	validators, err := createStakingValidators(valSet.Validators, bondedAmt)
	if err != nil {
		return err
	}

	fundedAccountBalances = addBondedModuleAccountToFundedBalances(fundedAccountBalances, sdktypes.NewCoin(n.cfg.denom, totalBonded))

	delegations := createDelegations(valSet.Validators, genAccounts[0].GetAddress())

	// Create a new EvmosApp with the following params
	evmosApp := createEvmosApp(n.cfg.chainID)

	// Configure Genesis state
	genesisState := app.NewDefaultGenesisState()

	genesisState = setAuthGenesisState(evmosApp, genesisState, genAccounts)

	stakingParams := StakingCustomGenesisState{
		denom:       n.cfg.denom,
		validators:  validators,
		delegations: delegations,
	}
	genesisState = setStakingGenesisState(evmosApp, genesisState, stakingParams)

	genesisState = setInflationGenesisState(evmosApp, genesisState)

	totalSupply := calculateTotalSupply(fundedAccountBalances)
	bankParams := BankCustomGenesisState{
		totalSupply: totalSupply,
		balances:    fundedAccountBalances,
	}
	genesisState = setBankGenesisState(evmosApp, genesisState, bankParams)

	// Init chain
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		return err
	}

	evmosApp.InitChain(
		abcitypes.RequestInitChain{
			ChainId:         n.cfg.chainID,
			Validators:      []abcitypes.ValidatorUpdate{},
			ConsensusParams: app.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)
	// Commit genesis changes
	evmosApp.Commit()

	header := tmproto.Header{
		ChainID:            n.cfg.chainID,
		Height:             evmosApp.LastBlockHeight() + 1,
		AppHash:            evmosApp.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
		ProposerAddress:    valSet.Proposer.Address,
	}
	evmosApp.BeginBlock(abcitypes.RequestBeginBlock{Header: header})

	// Set networks global parameters
	n.app = evmosApp
	// TODO - this might not be the best way to initilize the context
	n.ctx = evmosApp.BaseApp.NewContext(false, header)
	n.validators = validators
	return nil
}

// GetEIP155ChainID returns the network EIp-155 chainID number
func (n *IntegrationNetwork) GetEIP155ChainID() *big.Int {
	return n.cfg.eip155ChainID
}

// BroadcastTxSync broadcasts the given txBytes to the network and returns the response.
// TODO - this should be change to gRPC
func (n *IntegrationNetwork) BroadcastTxSync(txBytes []byte) (abcitypes.ResponseDeliverTx, error) {
	req := abcitypes.RequestDeliverTx{Tx: txBytes}
	return n.app.BaseApp.DeliverTx(req), nil
}

// Simulate simulates the given txBytes to the network and returns the simulated response.
// TODO - this should be change to gRPC
func (n *IntegrationNetwork) Simulate(txBytes []byte) (*txtypes.SimulateResponse, error) {
	gas, result, err := n.app.BaseApp.Simulate(txBytes)
	if err != nil {
		return nil, err
	}
	return &txtypes.SimulateResponse{
		GasInfo: &gas,
		Result:  result,
	}, nil
}
