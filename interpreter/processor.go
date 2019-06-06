package interpreter

import (
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/storage/state"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type Processor interface {
	ProcessTxState(tx *meta.Transaction, param Params) (error, Result)
	ProcessBlockState(block *node_meta.Block, stateDb *state.StateDB, validator Validator) (error, []Result)
}
