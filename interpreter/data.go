package interpreter

import (
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/storage/state"
)

type Params interface {
	GetBlockSigner() meta.AccountID
	GetBlockHeader() *meta.BlockHeader
	GetStateDB() *state.StateDB
	GetChainReader() meta.ChainReader
}

type Result interface {
	GetTxFee() *meta.Amount
	GetReceipt() *core.Receipt
	WriteResult() error
}
