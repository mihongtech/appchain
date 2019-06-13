package rpcserver

import (
	"github.com/mihongtech/appchain/rpc/rpcobject"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func getBlockChainInfo(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	chainId := GetNodeAPI(s).GetChainID()
	best := GetNodeAPI(s).GetBestBlock()

	chainInfo := node_meta.ChainInfo{
		ChainId:    int(chainId.Int64()),
		BestHeight: best.GetHeight(),
		BestHash:   best.GetBlockID().String(),
	}
	return &rpcobject.ChainRSP{
		Chains: &chainInfo,
	}, nil
}
