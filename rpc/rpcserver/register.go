package rpcserver

import (
	"github.com/mihongtech/linkchain-core/node"
	"github.com/mihongtech/linkchain-core/node/net"
	"reflect"

	"github.com/mihongtech/appchain/rpc/rpcobject"
	"github.com/mihongtech/appchain/wallet"
)

type commandHandler func(*Server, interface{}, <-chan struct{}) (interface{}, error)

//handler pool
var handlerPool = map[string]commandHandler{
	"getBlockChainInfo": getBlockChainInfo,

	"addPeer":    addPeer,
	"listPeer":   listPeer,
	"selfPeer":   selfPeer,
	"removePeer": removePeer,

	"getBestBlock":     getBestBlock,
	"getBlockByHeight": getBlockByHeight,
	"getBlockByHash":   getBlockByHash,

	//wallet
	"exportAccount": exportAccount,
	"importAccount": importAccount,

	"getWalletInfo":  getWalletInfo,
	"getAccountInfo": getAccountInfo,
	"newAcount":      newAcount,

	"sendMoneyTransaction": sendMoneyTransaction,

	//transaction
	"getTxByHash": getTxByHash,

	//shutdown
	"shutdown": shutdown,
}

var cmdPool = map[string]reflect.Type{
	"version":    reflect.TypeOf((*rpcobject.VersionCmd)(nil)),
	"addPeer":    reflect.TypeOf((*rpcobject.PeerCmd)(nil)),
	"removePeer": reflect.TypeOf((*rpcobject.PeerCmd)(nil)),

	"getBlockByHeight": reflect.TypeOf((*rpcobject.GetBlockByHeightCmd)(nil)),
	"getBlockByHash":   reflect.TypeOf((*rpcobject.GetBlockByHashCmd)(nil)),

	"getAccountInfo": reflect.TypeOf((*rpcobject.SingleCmd)(nil)),

	"sendMoneyTransaction": reflect.TypeOf((*rpcobject.SendToTxCmd)(nil)),

	"getTxByHash": reflect.TypeOf((*rpcobject.GetTransactionByHashCmd)(nil)),

	"importAccount": reflect.TypeOf((*rpcobject.ImportAccountCmd)(nil)),
	"exportAccount": reflect.TypeOf((*rpcobject.ExportAccountCmd)(nil)),
}

func GetWalletAPI(s *Server) *wallet.Wallet {
	return s.appContext.WalletAPI.(*wallet.Wallet)
}
func GetNodeAPI(s *Server) *node.CoreAPI {
	return s.appContext.NodeAPI
}

func GetP2PAPI(s *Server) net.P2PNet {
	return s.appContext.NodeAPI
}
