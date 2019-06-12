package rpcserver

import (
	"errors"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/util/log"
	"reflect"

	"github.com/mihongtech/appchain/rpc/rpcobject"
)

func getTxByHash(s *Server, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c, ok := cmd.(*rpcobject.GetTransactionByHashCmd)
	if !ok {
		log.Error("getTxByHash ", "Type error:", reflect.TypeOf(cmd))
		return nil, nil
	}

	hash, err := math.NewHashFromStr(c.Hash)
	if err != nil {
		return nil, err
	}

	transaction, _, _, _ := GetNodeAPI(s).GetTXByID(*hash)
	if transaction == nil {
		log.Error("getTxByHash ", "error", err)
		return nil, errors.New("getTxByHash failed")
	}
	appTx, err := meta.ConvertToAppTX(*transaction)
	return &rpcobject.TransactionWithIDRSP{transaction.GetTxID().GetString(), &appTx}, err
}
