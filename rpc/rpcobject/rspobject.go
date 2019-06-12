package rpcobject

import (
	"github.com/mihongtech/appchain/core/meta"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

//block
type BlockRSP struct {
	Height uint32                 `json:"height"`
	Hash   string                 `json:"hash"`
	Header *node_meta.BlockHeader `json:"header"`
	TXIDs  []string               `json:"txids"`
	Hex    string                 `json:"hex"`
}

//chain
type ChainRSP struct {
	Chains *node_meta.ChainInfo `json:"chains"`
}

//wallet
type WalletAccountRSP struct {
	ID     string `json:"id"`
	Type   uint32 `json:"type"`
	Amount int64  `json:"amount"`
}

type WalletInfoRSP struct {
	Accounts []*WalletAccountRSP `json:"accounts"`
}

//account
type TxRSP struct {
	TxID          string `json:"txid"`
	Index         uint32 `json:"index"`
	Value         int64  `json:"value"`
	LocatedHeight uint32 `json:"locatedHeight"`
	EffectHeight  uint32 `json:"effectHeight"`
}

type AccountRSP struct {
	ID          string   `json:"id"`
	Type        uint32   `json:"type"`
	Amount      int64    `json:"amount"`
	UTXO        []*TxRSP `json:"utxo"`
	StorageRoot string   `json:"storageRoot"`
	CodeHash    string   `json:"codeHash"`
	Code        string   `json:"code"`
}

//sendmoney
type TransactionWithIDRSP struct {
	ID string            `json:"id"`
	Tx *meta.Transaction `json:"tx"`
}
