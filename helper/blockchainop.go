package helper

import (
	"encoding/hex"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/linkchain-core/common"
	"github.com/mihongtech/linkchain-core/common/btcec"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
	"sort"
)

/*
	Account
*/

func CreateAccountIdByAddress(addr string) (*node_meta.Address, error) {
	buffer, err := hex.DecodeString(addr)
	if err != nil {
		return nil, err
	}

	id := node_meta.BytesToAddress(buffer)
	return &id, nil
}

func CreateAccountIdByPubKey(pubKey string) (*node_meta.Address, error) {
	pkBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}
	pk, err := btcec.ParsePubKey(pkBytes, btcec.S256())
	if err != nil {
		return nil, err
	}
	return node_meta.NewAddress(pk), nil
}

func CreateAccountIdByPrivKey(privKey string) (*node_meta.Address, error) {
	priv, err := hex.DecodeString(privKey)
	if err != nil {
		return nil, err
	}
	_, pk := btcec.PrivKeyFromBytes(btcec.S256(), priv)
	if err != nil {
		return nil, err
	}
	return node_meta.NewAddress(pk), nil
}

func CreateTemplateAccount(id node_meta.Address) *meta.Account {
	u := make([]meta.UTXO, 0)
	a := meta.NewAccount(id, config.NormalAccount, u)
	return a
}

func CreateNormalAccount(key *btcec.PrivateKey) (*meta.Account, error) {
	privateStr := hex.EncodeToString(key.Serialize())
	id, err := CreateAccountIdByPrivKey(privateStr)
	if err != nil {
		return nil, err
	}

	a := CreateTemplateAccount(*id)
	return a, nil
}

/*

	Transaction
*/

func CreateToCoin(to node_meta.Address, amount *meta.Amount) *meta.ToCoin {
	return meta.NewToCoin(to, amount)
}

func CreateFromCoin(from node_meta.Address, ticket ...meta.Ticket) *meta.FromCoin {
	tickets := make([]meta.Ticket, 0)
	fc := meta.NewFromCoin(from, tickets)
	for _, c := range ticket {
		fc.AddTicket(&c)
	}
	return fc
}

func CreateTempleteTx(version uint32, txtype uint32) *meta.Transaction {
	return meta.NewEmptyTransaction(version, txtype)
}

func CreateTransaction(fromCoin meta.FromCoin, toCoin meta.ToCoin) *meta.Transaction {
	transaction := CreateTempleteTx(config.DefaultTransactionVersion, config.NormalTx)
	transaction.AddFromCoin(fromCoin)
	transaction.AddToCoin(toCoin)
	return transaction
}

func CreateCoinBaseTx(to node_meta.Address, amount *meta.Amount, height uint32) *meta.Transaction {
	toCoin := meta.NewToCoin(to, amount)
	transaction := meta.NewEmptyTransaction(config.DefaultTransactionVersion, config.CoinBaseTx)
	transaction.AddToCoin(*toCoin)
	transaction.Data = common.UInt32ToBytes(height)
	return transaction
}

func SortTransaction(tx *meta.Transaction) {
	//sort from
	sort.Slice(tx.From.Coins, func(i, j int) bool {
		if tx.From.Coins[i].Id.Big().Cmp(tx.From.Coins[j].Id.Big()) == 0 {
			return false
		} else if tx.From.Coins[i].Id.Big().Cmp(tx.From.Coins[j].Id.Big()) < 0 {
			return true
		} else {
			return false
		}
	})

	//sort from ticket
	for k := range tx.From.Coins {
		sort.Slice(tx.From.Coins[k].Ticket, func(i, j int) bool {
			if tx.From.Coins[k].Ticket[i].Txid.Big().Cmp(tx.From.Coins[k].Ticket[j].Txid.Big()) == 0 {
				return tx.From.Coins[k].Ticket[i].Index > tx.From.Coins[k].Ticket[j].Index
			} else if tx.From.Coins[k].Ticket[i].Txid.Big().Cmp(tx.From.Coins[k].Ticket[j].Txid.Big()) < 0 {
				return true
			} else {
				return false
			}
		})
	}

	//sort to
	sort.Slice(tx.To.Coins, func(i, j int) bool {
		if tx.To.Coins[i].Id.Big().Cmp(tx.To.Coins[j].Id.Big()) == 0 {
			return tx.To.Coins[i].Value.GetInt64() < tx.To.Coins[j].Value.GetInt64()
		} else if tx.To.Coins[i].Id.Big().Cmp(tx.To.Coins[j].Id.Big()) < 0 {
			return true
		} else {
			return false
		}
	})
}
