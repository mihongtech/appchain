package core

import (
	"github.com/mihongtech/appchain/common/math"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/core/meta"
)

type Chain interface {
	meta.ChainReader

	// GetHeader returns the hash corresponding to their hash.
	GetHeader(math.Hash, uint64) *meta.BlockHeader

	Config() *config.ChainConfig
}
