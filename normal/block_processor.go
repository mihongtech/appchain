package normal

import (
	"errors"
	"github.com/mihongtech/linkchain-core/common/math"

	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/interpreter"
	"github.com/mihongtech/appchain/storage/state"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func (n *Interpreter) ProcessBlockState(block *node_meta.Block, stateDb *state.StateDB, chain node_meta.ChainReader, validator interpreter.Validator) (error, []interpreter.Result) {
	//update mine account status
	actualReward, fee, results, root, err := n.processBlockState(block, stateDb, chain, validator)
	if err != nil {
		return err, nil
	}

	if err := validator.VerifyBlockState(block, *root, actualReward, fee, nil); err != nil {
		return err, nil
	}
	return nil, results
}

func (n *Interpreter) processBlockState(block *node_meta.Block, stateDb *state.StateDB, chain node_meta.ChainReader, validator interpreter.Validator) (*meta.Amount, *meta.Amount, []interpreter.Result, *math.Hash, error) {
	txs := block.GetTxs()

	coinBase := meta.NewAmount(0)
	txFee := meta.NewAmount(0)
	inputData := Input{&block.Header, stateDb, chain, block.TXs[0].To.Coins[0].Id}
	outputDatas := make([]interpreter.Result, 0)
	for index := range txs {
		if err := validator.VerifyTx(&txs[index], &inputData); err != nil {
			return nil, nil, nil, nil, errors.New(err.Error() + ",txid=" + txs[index].GetTxID().String())
		}
		err, outputData := n.ProcessTxState(&txs[index], &inputData)
		if err != nil {
			return nil, nil, nil, nil, errors.New(err.Error() + ",txid=" + txs[index].GetTxID().String())
		}
		outputDatas = append(outputDatas, outputData)
		if txs[index].GetType() != config.CoinBaseTx {
			txFee.Addition(*outputData.GetTxFee())
		} else {
			coinBase.Addition(*txs[index].GetToValue())
		}
	}

	root := stateDb.IntermediateRoot()
	return coinBase, txFee, outputDatas, &root, nil
}
