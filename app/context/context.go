package context

import (
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/interpreter"
	"github.com/mihongtech/linkchain-core/node"
)

type Context struct {
	NodeAPI        *node.CoreAPI
	WalletAPI      interpreter.Wallet
	InterpreterAPI interpreter.Interpreter
	Config         *config.LinkChainConfig
}
