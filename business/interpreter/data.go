package interpreter

import (
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/storage/state"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type Params interface {
	GetBlockSigner() node_meta.Address
	GetBlockHeader() *node_meta.BlockHeader
	GetStateDB() *state.StateDB
	GetChainReader() node_meta.ChainReader
}

type Result interface {
	GetTxFee() *meta.Amount
	GetReceipt() *core.Receipt
	WriteResult() error
}
