package interpreter

import (
	"github.com/mihongtech/appchain/common/math"
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/core/meta"
)

type OffChain interface {
	core.Service
	UpdateMainChain(ev meta.ChainEvent)
	UpdateSideChain(ev meta.ChainSideEvent)
}

type Wallet interface {
	core.Service
	SignMessage(accountId meta.AccountID, hash []byte) (math.ISignature, error)
	SignTransaction(tx meta.Transaction) (*meta.Transaction, error)
	ImportAccount(privateKeyStr string) (*meta.AccountID, error)
	ExportAccount(id meta.AccountID) (string, error)
	GetAccount(key string) (*meta.Account, error)
	GetAllWAccount() []meta.Account
	AddAccount(account meta.Account)
	NewAccount() (*meta.AccountID, error)
}
