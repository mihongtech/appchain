package interpreter

import (
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/storage/state"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type Processor interface {
	ProcessTxState(tx *meta.Transaction, param Params) (error, Result)
	ProcessBlockState(header *node_meta.BlockHeader, txs []meta.Transaction, stateDb *state.StateDB, chain node_meta.ChainReader, validator Validator) (error, []Result)
}
