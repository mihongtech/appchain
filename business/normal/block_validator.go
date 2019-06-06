package normal

import (
	"errors"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/util/log"

	"github.com/mihongtech/appchain/business/interpreter"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/core/meta"

	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func (n *Interpreter) ValidateBlockBody(header *node_meta.BlockHeader, txs []meta.Transaction, txValidator interpreter.TransactionValidator, chain node_meta.ChainReader) error {
	//check TXs only one coinBase
	for i := range txs {
		if i != 0 && txs[i].Type == config.CoinBaseTx {
			return errors.New("the block must be only one coinBase tx")
		} else if i == 0 && txs[i].Type != config.CoinBaseTx {
			return errors.New("the first tx of block must be coinBase tx")
		}
	}

	//check txs have the same tx
	txCount := len(txs)
	for i := 0; i < txCount; i++ {
		for j := i + 1; j < txCount; j++ {
			if txs[i].GetTxID().IsEqual(txs[j].GetTxID()) {
				return errors.New("the block have two same tx")
			}
		}
	}

	//check tx body
	for i := range txs {
		if err := txValidator.CheckTx(&txs[i]); err != nil {
			return err
		}
	}
	return nil
}

func (n *Interpreter) VerifyBlockState(header *node_meta.BlockHeader, txs []meta.Transaction, root math.Hash, actualReward *meta.Amount, fee *meta.Amount, headerData []byte) error {
	log.Debug("VerifyBlockState", "actualReward", actualReward.GetInt64(), "fee", fee.GetInt64())
	//Check block reward
	if actualReward.Subtraction(*meta.NewAmount(config.DefaultBlockReward)).GetInt64() != fee.GetInt64() && len(txs) > 0 {
		return errors.New("coin base tx reward is error")
	}

	return nil
}
