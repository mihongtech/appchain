package poa

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mihongtech/linkchain/common/btcec"
	"github.com/mihongtech/linkchain/common/lcdb"
	"github.com/mihongtech/linkchain/common/math"
	"github.com/mihongtech/linkchain/common/util/log"
	"github.com/mihongtech/linkchain/core/meta"
	"github.com/mihongtech/linkchain/node/config"
	"sync"
)

// SignerFn is a signer callback function to request a hash to be signed by a
// backing account.
type SignerFn func(meta.Account, []byte) ([]byte, error)

// Poa is the proof-of-authority consensus engine proposed
type Poa struct {
	chainConfig *config.ChainConfig // Consensus engine configuration parameters
	db          lcdb.Database       // Database to store and retrieve snapshot checkpoints

	proposals map[math.Hash]bool // Current list of proposals we are pushing

	signer math.Hash    // address of the signing key
	signFn SignerFn     // Signer function to authorize hashes with
	lock   sync.RWMutex // Protects the signer fields
}

// New creates a proof-of-authority consensus engine with the initial
// signers set to the ones provided by the user.
func NewPoa(chainConfig *config.ChainConfig, db lcdb.Database) *Poa {
	// Set any missing consensus parameters to their defaults
	conf := *chainConfig

	return &Poa{
		chainConfig: &conf,
		db:          db,
		proposals:   make(map[math.Hash]bool),
	}
}

func (p *Poa) Author(header *meta.BlockHeader) ([]byte, error) {
	pub, _, err := btcec.RecoverCompact(btcec.S256(), header.Sign.Code, (*header.GetBlockID())[:])
	if err != nil {
		return nil, err
	}

	id := meta.NewAccountId(pub)

	return id.CloneBytes(), nil
}

func (p *Poa) VerifyBlock(chain meta.ChainReader, block *meta.Block) error {
	prevBlock, err := chain.GetBlockByID(*block.GetPrevBlockID())

	if err != nil {
		log.Error("BlockManage", "checkBlock", err)
		return err
	}

	if prevBlock.GetHeight()+1 != block.GetHeight() {
		log.Error("BlockManage", "checkBlock", "current block height is error")
		return errors.New("Check block height failed")
	}

	croot := block.CalculateTxTreeRoot()
	if !block.GetMerkleRoot().IsEqual(&croot) {
		log.Error("POA checkBlock", "check merkle root", false)
		return errors.New("Check block merkle root failed")
	}

	//check txs have the same tx
	txs := block.GetTxs()
	txCount := len(txs)
	for i := 0; i < txCount; i++ {
		for j := i + 1; j < txCount; j++ {
			if txs[i].GetTxID().IsEqual(txs[j].GetTxID()) {
				return errors.New("the block have two same tx")
			}
		}
	}
	return p.verifySeal(block)
}

func (p *Poa) verifySeal(block *meta.Block) error {
	signerIndex := block.GetHeight() % uint32(len(config.SignMiners))
	miner, err := hex.DecodeString(config.SignMiners[signerIndex])
	if err != nil {
		return err
	}
	pubkey, _, err := btcec.RecoverCompact(btcec.S256(), block.Header.Sign.Code, block.GetBlockID().CloneBytes())
	if err != nil {
		return err
	}

	accountID := meta.NewAccountId(pubkey)
	if accountID.IsEqual(meta.BytesToAccountID(miner)) {
		return nil
	}

	return errors.New(fmt.Sprintf("Verify seal failed %s\n, want %s\n", accountID.String(), meta.BytesToAccountID(miner).String()))
}