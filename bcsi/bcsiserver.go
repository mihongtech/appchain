package bcsi

import (
	"errors"
	"github.com/bliblicode/library/log"

	"sync/atomic"

	"github.com/hashicorp/golang-lru"
	"github.com/mihongtech/appchain/business/interpreter"
	"github.com/mihongtech/appchain/storage/state"
	"github.com/mihongtech/linkchain-core/common/lcdb"
	"github.com/mihongtech/linkchain-core/common/math"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
	"github.com/mihongtech/linkchain-core/node/bcsi"
	"github.com/mihongtech/linkchain-core/node/chain"
)

const (
	statusCacheLimit = 256
	ErrorProcessed   = "the block had processed"
)

type BCSIServer struct {
	db          lcdb.Database
	interpreter interpreter.Interpreter

	chain chain.ChainReader

	cacheState  map[node_meta.BlockID]*state.StateDB
	statusCache *lru.Cache // Cache for status of block

	currentBlock atomic.Value
}

func NewBCSIServer(db lcdb.Database, interpreter interpreter.Interpreter) bcsi.BCSI {
	statusCache, _ := lru.New(statusCacheLimit)
	cacheState := make(map[node_meta.BlockID]*state.StateDB)
	return &BCSIServer{db: db, interpreter: interpreter, cacheState: cacheState, statusCache: statusCache}
}

func (s *BCSIServer) GetBlockState(id node_meta.BlockID) (node_meta.TreeID, error) {
	//TODO need use map[blockid]status
	block, err := s.chain.GetBlockByID(id)
	if err != nil {
		return math.Hash{}, nil
	}
	stateDB, err := state.New(*block.GetStatus(), s.db)
	if err != nil {
		return math.Hash{}, err
	}
	return stateDB.GetRootHash(), err
}

func (s *BCSIServer) UpdateChain(head *node_meta.Block) error {
	if s.currentBlock.Load() == nil {
		log.Info("BCSIServer", "UpdateChain", "init chain", "best block", head.GetBlockID().String())
	} else {
		log.Info("BCSIServer", "UpdateChain", "update chain", "best block", head.GetBlockID().String())
	}
	s.currentBlock.Store(head)
	return nil
}

func (s *BCSIServer) ProcessBlock(block *node_meta.Block) error {
	//check app have status
	if _, ok := s.statusCache.Get(*block.GetBlockID()); ok {
		return errors.New(ErrorProcessed)
	}
	stateDB, err := state.New(*block.GetStatus(), s.db)
	if err != nil {
		return err
	}
	//TODO need encode node-tx to app-tx
	err, _ = s.interpreter.ProcessBlockState(&block.Header, nil, stateDB, s.chain, s.interpreter)
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
	return s.interpreter.CheckTx(nil)
}

func (s *BCSIServer) FilterTx(txs []node_meta.Transaction) []node_meta.Transaction {
	return txs
}
