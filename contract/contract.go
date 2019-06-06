package contract

import (
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/contract/vm"
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/normal"
)

type Input struct {
	normal.Input
	Chain   ChainContext
	Config  *config.ChainConfig
	VmCfg   vm.Config
	UsedGas *uint64
	Gp      *core.GasPool
}

type Interpreter struct {
	normal.Interpreter
}

type Output struct {
	normal.Output
	ResultTx *meta.Transaction
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
