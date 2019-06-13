package lcclient

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/mihongtech/appchain/client/httpclient"
	"github.com/mihongtech/appchain/core"
	"github.com/mihongtech/appchain/rpc/rpcjson"
	"github.com/mihongtech/appchain/rpc/rpcobject"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/util/log"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

var httpConfig = &httpclient.Config{
	RPCUser:     "lc",
	RPCPassword: "lc",
	RPCServer:   "localhost:8082",
}

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	c string
}

// NewClient creates a client that uses the given RPC client.
func NewClient(c string) *Client {
	return &Client{c}
}

// Blockchain Access

// BlockByHash returns the given full block.
//
// Note that loading full blocks requires two requests. Use HeaderByHash
// if you don't need all transactions or uncle headers.
func (ec *Client) BlockByHash(ctx context.Context, hash math.Hash) (*node_meta.Block, error) {
	// return ec.getBlock(ctx, "eth_getBlockByHash", hash, true)
	return nil, nil
}

// BlockByNumber returns a block from the current canonical chain. If number is nil, the
// latest known block is returned.
//
// Note that loading full blocks requires two requests. Use HeaderByNumber
// if you don't need all transactions or uncle headers.
func (ec *Client) BlockByNumber(ctx context.Context, number *big.Int) (*node_meta.Block, error) {
	// return ec.getBlock(ctx, "eth_getBlockByNumber", toBlockNumArg(number), true)
	return nil, nil
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *Client) CodeAt(ctx context.Context, account node_meta.Address, blockNumber *big.Int) ([]byte, error) {
	method := "getCode"
	//call
	if blockNumber == nil {
		blockNumber = big.NewInt(-1)
	}

	data, err := rpc(method, &rpcobject.GetCodeCmd{account.String(), blockNumber.Int64()})

	return data, err
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *Client) PendingCodeAt(ctx context.Context, account node_meta.Address) ([]byte, error) {
	//	var result hexutil.Bytes
	//	err := ec.c.CallContext(ctx, &result, "eth_getCode", account, "pending")
	//	return result, err
	return nil, nil
}

func (ec *Client) TransactionReceipt(ctx context.Context, txHash math.Hash) (*core.Receipt, error) {
	method := "transactionReceipt"
	//call
	data, err := rpc(method, &rpcobject.GetTransactionReceiptCmd{txHash.String()})
	var receipt core.Receipt
	if err = json.Unmarshal(data, &receipt); err != nil {
		log.Error("Unmarshal json failed", "data", data)
		return nil, err
	}
	return &receipt, nil
}

//rpc call
func rpc(method string, cmd interface{}) ([]byte, error) {
	//param
	s, _ := rpcjson.MarshalCmd(1, method, cmd)
	//log.Info(method, "req", string(s))

	//response
	rawRet, err := httpclient.SendPostRequest(s, httpConfig)
	if err != nil {
		log.Error(method, "error", err)
		return nil, err
	}

	//log.Info(method, "rsp", string(rawRet))

	return rawRet, nil
}
