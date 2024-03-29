package context

import (
	"github.com/mihongtech/appchain/bcsi"
	"github.com/mihongtech/appchain/business/interpreter"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/linkchain-core/node"
)

type Context struct {
	NodeAPI        *node.CoreAPI
	WalletAPI      interpreter.Wallet
	InterpreterAPI interpreter.Interpreter
	BCSIAPI        *bcsi.BCSIServer
	Config         *config.LinkChainConfig
}
