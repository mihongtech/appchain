package meta

import (
	"github.com/mihongtech/linkchain-core/common/math"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type ChainInfo struct {
	ChainId    int
	BestHeight uint32
	BestHash   string
}

type ChainEvent struct {
	Block *node_meta.Block
	Hash  math.Hash
}

type ChainSideEvent struct {
	Block *node_meta.Block
}

type ChainHeadEvent struct{ Block *node_meta.Block }
