package interpreter

import (
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/linkchain-core/common/math"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type BlockValidator interface {
	VerifyBlockState(header *node_meta.BlockHeader, txs []meta.Transaction, root math.Hash, actualReward *meta.Amount, fee *meta.Amount, headerData []byte) error
	ValidateBlockBody(header *node_meta.BlockHeader, txs []meta.Transaction, txValidator TransactionValidator, chain node_meta.ChainReader) error
}

type TransactionValidator interface {
	CheckTx(tx *meta.Transaction) error
	VerifyTx(tx *meta.Transaction, data Params) error
}

type Validator interface {
	TransactionValidator
	BlockValidator
}
