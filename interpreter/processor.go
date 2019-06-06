package interpreter

import (
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/storage/state"
)

type Processor interface {
	ProcessTxState(tx *meta.Transaction, param Params) (error, Result)
	ProcessBlockState(block *meta.Block, stateDb *state.StateDB, chain core.Chain, validator Validator) (error, []Result)
}
