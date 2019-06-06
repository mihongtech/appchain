package interpreter

import (
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/linkchain-core/common/math"
	node_core "github.com/mihongtech/linkchain-core/core"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type OffChain interface {
	node_core.Service
	UpdateMainChain(ev meta.ChainEvent)
	UpdateSideChain(ev meta.ChainSideEvent)
}

type Wallet interface {
	node_core.Service
	SignMessage(accountId node_meta.Address, hash []byte) (math.ISignature, error)
	SignTransaction(tx meta.Transaction) (*meta.Transaction, error)
	ImportAccount(privateKeyStr string) (*node_meta.Address, error)
	ExportAccount(id node_meta.Address) (string, error)
	GetAccount(key string) (*meta.Account, error)
	GetAllWAccount() []meta.Account
	AddAccount(account meta.Account)
	NewAccount() (*node_meta.Address, error)
}
