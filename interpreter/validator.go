package interpreter

import (
	"github.com/mihongtech/appchain/common/math"
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/node/consensus"
)

type BlockValidator interface {
	VerifyBlockState(block *meta.Block, root math.Hash, actualReward *meta.Amount, fee *meta.Amount, headerData []byte) error
	ValidateBlockBody(txValidator TransactionValidator, chain core.Chain, block *meta.Block) error
	ValidateBlockHeader(engine consensus.Engine, chain core.Chain, block *meta.Block) error
}

type TransactionValidator interface {
	CheckTx(tx *meta.Transaction) error
	VerifyTx(tx *meta.Transaction, data Params) error
}

type Validator interface {
	TransactionValidator
	BlockValidator
}
