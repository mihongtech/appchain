package block

import (
	"github.com/linkchain/common/serialize"
	"github.com/linkchain/meta"
	"github.com/linkchain/meta/tx"
)

type NewMinedBlockEvent struct{ Block IBlock }

type IBlock interface {
	//block content
	SetTx([]tx.ITx) error

	GetTxs() []tx.ITx

	GetHeight() uint32

	GetBlockID() meta.DataID

	GetPrevBlockID() meta.DataID
	//verifiy
	Verify() error

	IsGensis() bool

	//serialize
	serialize.ISerialize
}
