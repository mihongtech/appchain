package rpcserver

import (
	"github.com/mihongtech/appchain/rpc/rpcobject"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func getBlockChainInfo(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	info := GetNodeAPI(s).GetBlockChainInfo()
	return &rpcobject.ChainRSP{
		Chains: info.(*node_meta.ChainInfo),
	}, nil
}
