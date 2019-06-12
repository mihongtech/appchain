package bcsi

import (
	"errors"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/helper"
	"sync/atomic"

	"github.com/mihongtech/appchain/business/interpreter"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/storage/state"
	"github.com/mihongtech/linkchain-core/common/lcdb"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/util/log"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
	"github.com/mihongtech/linkchain-core/node/chain"

	"github.com/hashicorp/golang-lru"
)

const (
	statusCacheLimit = 256
	ErrorProcessed   = "the block had processed"
)

type BCSIServer struct {
	Db          lcdb.Database
	interpreter interpreter.Interpreter

	chain chain.ChainReader

	cacheState  map[node_meta.BlockID]*state.StateDB
	statusCache *lru.Cache // Cache for status of block

	CurrentBlock atomic.Value
}

func NewBCSIServer(db lcdb.Database, interpreter interpreter.Interpreter) *BCSIServer {
	statusCache, _ := lru.New(statusCacheLimit)
	cacheState := make(map[node_meta.BlockID]*state.StateDB)
	return &BCSIServer{Db: db, interpreter: interpreter, cacheState: cacheState, statusCache: statusCache}
}

func (s *BCSIServer) Setup(i interface{}) bool {
	s.chain = i.(chain.ChainReader)
	return true
}

func (s *BCSIServer) Start() bool {
	return true
}

func (s *BCSIServer) Stop() {

}

func (s *BCSIServer) GetBlockState(id node_meta.BlockID) (node_meta.TreeID, error) {
	//TODO need use map[blockid]status
	block, err := s.chain.GetBlockByID(id)
	if err != nil {
		return math.Hash{}, nil
	}
	stateDB, err := state.New(*block.GetStatus(), s.Db)
	if err != nil {
		return math.Hash{}, err
	}
	return stateDB.GetRootHash(), err
}

func (s *BCSIServer) UpdateChain(head *node_meta.Block) error {
	if s.CurrentBlock.Load() == nil {
		log.Info("BCSIServer", "UpdateChain", "init chain", "best block", head.GetBlockID().String())
	} else {
		log.Info("BCSIServer", "UpdateChain", "update chain", "best block", head.GetBlockID().String())
	}
	s.CurrentBlock.Store(head)
	return nil
}

func (s *BCSIServer) ProcessBlock(block *node_meta.Block) error {
	//check app have status
	if _, ok := s.statusCache.Get(*block.GetBlockID()); ok {
		return errors.New(ErrorProcessed)
	}
	stateDB, err := state.New(*block.GetStatus(), s.Db)
	if err != nil {
		return err
	}
	err, _ = s.interpreter.ProcessBlockState(&block.Header, AppTransactionsConvert(&block.TXs), stateDB, s.chain, s.interpreter)
	if err != nil {
		return err
	}
	s.cacheState[*block.GetBlockID()] = stateDB
	return nil
}

func (s *BCSIServer) Commit(id node_meta.BlockID) error {
	stateDB := s.cacheState[id]
	status, err := stateDB.Commit()
	if err != nil {
		return err
	}
	s.statusCache.Add(id, status)
	delete(s.cacheState, id)
	return nil
}

func (s *BCSIServer) CheckBlock(block *node_meta.Block) error {
	return nil
}

func (s *BCSIServer) CheckTx(transaction node_meta.Transaction) error {
	tx, err := meta.ConvertToAppTX(transaction)
	if err != nil {
		return err
	}
	return s.interpreter.CheckTx(&tx)
}

func (s *BCSIServer) FilterTx(txs []node_meta.Transaction) []node_meta.Transaction {
	height := 1
	block := s.CurrentBlock.Load().(*node_meta.Block)
	if block != nil {
		height = int(block.GetHeight())
	}
	signer, _ := node_meta.NewAddressFromStr(config.FirstPubMiner)
	coinbase := helper.CreateCoinBaseTx(*signer, meta.NewAmount(config.DefaultBlockReward), uint32(height+1))
	nodeTx, _ := meta.ConvertToNodeTX(*coinbase)
	txs = append(txs, nodeTx)
	return txs
}
