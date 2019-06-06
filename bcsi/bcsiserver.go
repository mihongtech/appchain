package bcsi

import (
	"github.com/mihongtech/linkchain-core/common/lcdb"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type BCSIServer struct {
	db lcdb.Database
}

func (s *BCSIServer) GetBlockState(id node_meta.BlockID) node_meta.TreeID {
	return id
}

func (s *BCSIServer) UpdateChain(head *node_meta.Block) error {
	return nil
}

func (s *BCSIServer) ProcessBlock(block *node_meta.Block) error {
	return nil
}

func (s *BCSIServer) Commit(id node_meta.BlockID) error {
	return nil
}

func (s *BCSIServer) CheckBlock(block *node_meta.Block) error {
	return nil
}

func (s *BCSIServer) CheckTx(transaction node_meta.Transaction) error {
	return nil
}

func (s *BCSIServer) FilterTx(txs []node_meta.Transaction) []node_meta.Transaction {
	return nil
}
