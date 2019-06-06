package normal

import (
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/interpreter"
	"github.com/mihongtech/appchain/storage/state"
	"github.com/mihongtech/linkchain-core/common/lcdb"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type Interpreter struct {
}

func (i *Interpreter) CreateOffChain(db lcdb.Database) interpreter.OffChain {
	return &OffChainState{}
}

type Input struct {
	Header      *node_meta.BlockHeader
	StateDB     *state.StateDB
	ChainReader node_meta.ChainReader
	BlockSigner node_meta.Address
}

func (i *Input) GetBlockSigner() node_meta.Address {
	return i.BlockSigner
}

func (i *Input) GetBlockHeader() *node_meta.BlockHeader {
	return i.Header
}

func (i *Input) GetStateDB() *state.StateDB {
	return i.StateDB
}

func (i *Input) GetChainReader() node_meta.ChainReader {
	return i.ChainReader
}

type Output struct {
	TxFee   *meta.Amount
	Receipt *core.Receipt
}

func (o *Output) GetTxFee() *meta.Amount {
	return o.TxFee
}

func (o *Output) GetReceipt() *core.Receipt {
	return o.Receipt
}

func (o *Output) WriteResult() error {
	return nil
}

func IsNormal(txType uint32) bool {
	return txType == config.CoinBaseTx || txType == config.NormalTx
}

func GetReceiptsByResult(results []interpreter.Result) []*core.Receipt {
	receipts := make([]*core.Receipt, 0)
	for i := range results {
		if results[i].GetReceipt() != nil {
			receipts = append(receipts, results[i].GetReceipt())
		}
	}
	return receipts
}
