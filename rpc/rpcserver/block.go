package rpcserver

import (
	"encoding/hex"
	"fmt"

	"reflect"

	"github.com/mihongtech/appchain/rpc/rpcobject"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/util/log"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func getBestBlock(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	block := GetNodeAPI(s).GetBestBlock()
	b := getBlockObject(block)
	return b, nil
}

func getBlockByHeight(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c, ok := cmd.(*rpcobject.GetBlockByHeightCmd)
	if !ok {
		fmt.Println("Type error:", reflect.TypeOf(cmd))
		return nil, nil
	}

	height := c.Height

	if uint32(height) > GetNodeAPI(s).GetBestBlock().GetHeight() || height < 0 {
		log.Error("getblockbyheight ", "error", "height is out of range", "best", GetNodeAPI(s).GetBestBlock().GetHeight())
		return nil, nil
	}

	// get block
	block, err := GetNodeAPI(s).GetBlockByHeight(uint32(height))
	if err != nil {
		log.Error("getblockbyheight ", "error", err)
		return nil, err
	}

	b := getBlockObject(block)
	return b, nil
}

func getBlockByHash(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c, ok := cmd.(*rpcobject.GetBlockByHashCmd)
	if !ok {
		fmt.Println("Type error:", reflect.TypeOf(cmd))
		return nil, nil
	}

	hash, err := math.NewHashFromStr(c.Hash)
	if err != nil {
		return nil, err
	}

	block, err := GetNodeAPI(s).GetBlockByID(*hash)
	if err != nil {
		log.Error("getblockbyhash ", "error", err)
		return nil, err
	}

	b := getBlockObject(block)
	return b, nil
}

func getBlockObject(block *node_meta.Block) *rpcobject.BlockRSP {

	buffer, _ := block.EncodeToBytes()
	txids := make([]string, 0)
	for i := range block.TXs.Txs {
		txid := block.TXs.Txs[i].GetTxID()
		txids = append(txids, txid.String())
		log.Error("getBlockObject", "txid", txid.String())
	}

	return &rpcobject.BlockRSP{
		block.GetHeight(),
		block.GetBlockID().String(),
		&block.Header,
		txids,
		hex.EncodeToString(buffer),
	}
}
