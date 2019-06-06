package interpreter

import (
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/linkchain-core/common/math"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type BlockValidator interface {
	VerifyBlockState(block *node_meta.Block, root math.Hash, actualReward *meta.Amount, fee *meta.Amount, headerData []byte) error
	ValidateBlockBody(txValidator TransactionValidator, block *node_meta.Block) error
}

type TransactionValidator interface {
	CheckTx(tx *meta.Transaction) error
	VerifyTx(tx *meta.Transaction, data Params) error
}

type Validator interface {
	TransactionValidator
	BlockValidator
}
