package rpcserver

import (
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/rpc/rpcobject"
)

func getBlockChainInfo(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	info := GetNodeAPI(s).GetBlockChainInfo()
	return &rpcobject.ChainRSP{
		Chains: info.(*meta.ChainInfo),
	}, nil
}
