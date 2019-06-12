package rpcserver

import (
	"github.com/mihongtech/appchain/rpc/rpcobject"
)

func getBlockChainInfo(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	//info := GetNodeAPI(s).GetChainID()
	return &rpcobject.ChainRSP{
		//Chains: info.(*node_meta.ChainInfo),
	}, nil
}
