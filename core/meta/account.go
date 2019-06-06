package meta

import (
	"encoding/json"
	"errors"
	"github.com/mihongtech/linkchain-core/common/serialize"
	"github.com/mihongtech/linkchain-core/common/util/log"
	"sort"

	"github.com/mihongtech/appchain/protobuf"
	"github.com/mihongtech/linkchain-core/common/math"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"

	"github.com/golang/protobuf/proto"
)

type UTXO struct {
	Ticket
	LocatedHeight uint32 `json:"locatedHeight"`
	EffectHeight  uint32 `json:"effectHeight"`
	Value         Amount `json:"value"`
}

func NewUTXO(tickets *Ticket, locatedHeight uint32, effectHeight uint32, value Amount) *UTXO {
	return &UTXO{*tickets, locatedHeight, effectHeight, value}
}

func (u *UTXO) String() string {
	data, err := json.Marshal(u)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

type Account struct {
	Id          node_meta.Address `json:"accountId"`
	AccountType uint32            `json:"accountType"`
	UTXOs       []UTXO            `json:"AccountUXTO"`
	StorageRoot node_meta.TreeID  `json:"storageRoot"`
	CodeHash    math.Hash         `json:"codeHash"`
}

func NewAccount(id node_meta.Address, accountType uint32, utxos []UTXO) *Account {
	return &Account{Id: id, AccountType: accountType, UTXOs: utxos, StorageRoot: math.Hash{}, CodeHash: math.Hash{}}
}

func (a Account) GetAccountID() *node_meta.Address {
	return &a.Id
}

func (a Account) GetAmount() *Amount {
	sum := NewAmount(0)
	for _, u := range a.UTXOs {
		sum.Addition(u.Value)
	}
	return sum
}

func (a *Account) GetFromCoinValue(fromCoin *FromCoin) (*Amount, error) {
	sum := NewAmount(0)
	if a.GetAccountID().IsEqual(fromCoin.GetId()) {
		tickets := fromCoin.GetTickets()
		for _, t := range tickets {
			u, err := a.getUTXOByTicket(t)
			if err != nil {
				return nil, err
			}
			sum.Addition(u.Value)
		}
		return sum, nil
	} else {
		return nil, errors.New("fromCoin's AccountId is error")
	}
}

//check fromCoin ticket effect.
//the current block height must be > effectHeight
func (a *Account) IsFromEffect(fromCoin *FromCoin, height uint32) bool {
	if a.GetAccountID().IsEqual(fromCoin.GetId()) {
		tickets := fromCoin.GetTickets()
		for _, t := range tickets {
			u, err := a.getUTXOByTicket(t)
			if err != nil {
				return false
			}
			if height < u.EffectHeight {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (a *Account) CheckFromCoin(fromCoin *FromCoin) bool {
	if a.GetAccountID().IsEqual(fromCoin.GetId()) {
		tickets := fromCoin.GetTickets()
		for _, t := range tickets {
			if a.Contains(t) {
				return true
			}
		}
	}
	return false
}

func (a *Account) RemoveUTXOByFromCoin(fromCoin *FromCoin) error {
	if a.GetAccountID().IsEqual(fromCoin.GetId()) {
		tickets := fromCoin.GetTickets()
		for _, t := range tickets {
			removeOk := false
			for index, u := range a.UTXOs {
				if t.GetTxid().IsEqual(&u.Txid) && t.GetIndex() == u.GetIndex() {
					a.UTXOs = append(a.UTXOs[:index], a.UTXOs[index+1:]...)
					removeOk = true
					break
				}
			}
			if !removeOk {
				return errors.New("removeUTXO():The ticket is not exist in Account")

			}
		}
		return nil
	} else {
		return errors.New("fromCoin's AccountId is error")
	}
}

func (a *Account) getUTXOByTicket(ticket Ticket) (*UTXO, error) {
	for _, t := range a.UTXOs {
		if ticket.GetTxid().IsEqual(&t.Txid) && t.Index == ticket.GetIndex() {
			return &t, nil
		}
	}
	return nil, errors.New("the ticket is not exist in Account")
}

func (a *Account) Contains(ticket Ticket) bool {
	for _, t := range a.UTXOs {
		if ticket.GetTxid().IsEqual(&t.Txid) && t.Index == ticket.GetIndex() {
			return true
		}
	}
	return false
}

func (a *Account) GetUTXO(ticket Ticket) *UTXO {
	for _, t := range a.UTXOs {
		if ticket.GetTxid().IsEqual(&t.Txid) && t.Index == ticket.GetIndex() {
			return &t
		}
	}
	return nil
}

//Serialize/Deserialize
func (a *Account) Serialize() serialize.SerializeStream {
	us := make([]*protobuf.UTXO, 0)
	for index := range a.UTXOs {
		t := NewTicket(a.UTXOs[index].Txid, a.UTXOs[index].Index)
		u := &protobuf.UTXO{
			Id:            t.Serialize().(*protobuf.Ticket),
			LocatedHeight: proto.Uint32(a.UTXOs[index].LocatedHeight),
			EffectHeight:  proto.Uint32(a.UTXOs[index].EffectHeight),
			Value:         proto.NewBuffer(a.UTXOs[index].Value.GetBytes()).Bytes(),
		}
		us = append(us, u)
	}
	s := &protobuf.Account{
		Id:    a.Id.Serialize().(*protobuf.AccountID),
		Type:  proto.Uint32(a.AccountType),
		Utxos: us,
	}

	if !a.CodeHash.IsEmpty() {
		s.CodeHash = a.CodeHash.Serialize().(*protobuf.Hash)
	}

	if !a.StorageRoot.IsEmpty() {
		s.StorageRoot = a.StorageRoot.Serialize().(*protobuf.Hash)
	}

	return s

}

func (a *Account) Deserialize(s serialize.SerializeStream) error {
	data := s.(*protobuf.Account)
	if err := a.Id.Deserialize(data.Id); err != nil {
		return err
	}

	if data.StorageRoot != nil {
		if err := a.StorageRoot.Deserialize(data.StorageRoot); err != nil {
			return err
		}
	}

	if data.CodeHash != nil {
		if err := a.CodeHash.Deserialize(data.CodeHash); err != nil {
			return err
		}
	}

	a.AccountType = *data.Type

	a.UTXOs = a.UTXOs[:0] // UTXOs clear

	for _, u := range data.Utxos {
		newUtxo := UTXO{}
		newTicket := &Ticket{}
		if err := newTicket.Deserialize(u.Id); err != nil {
			return err
		}
		newUtxo.Txid = newTicket.Txid
		newUtxo.Index = newTicket.Index
		newUtxo.Value = *NewAmount(0)
		newUtxo.Value.SetBytes(u.Value)
		newUtxo.LocatedHeight = *u.LocatedHeight
		newUtxo.EffectHeight = *u.EffectHeight

		a.UTXOs = append(a.UTXOs, newUtxo)
	}
	return nil
}

func (a *Account) String() string {
	data, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (a *Account) MakeFromCoin(value *Amount, blockHeight uint32) (*FromCoin, *Amount, error) {
	if a.GetAmount().GetInt64() < value.GetInt64() {
		log.Error("MakeFromCoin failed", "a.GetAmount().GetInt64()", a.GetAmount().GetInt64(), "value.GetInt64()", value.GetInt64())
		return nil, nil, errors.New("Account MakeFromCoin() amount is too large")
	}
	tempUTXOs := make([]UTXO, 0)
	tempUTXOs = append(tempUTXOs, a.UTXOs...)
	sort.Slice(tempUTXOs, func(i, j int) bool {
		if tempUTXOs[i].Txid.Big().Cmp(tempUTXOs[j].Txid.Big()) == 0 {
			return tempUTXOs[i].Index > tempUTXOs[j].Index
		} else if tempUTXOs[i].Txid.Big().Cmp(tempUTXOs[j].Txid.Big()) < 0 {
			return true
		} else {
			return false
		}
	})
	tickets := make([]Ticket, 0)
	fc := NewFromCoin(a.Id, tickets)
	fromAmount := NewAmount(0)
	for _, v := range tempUTXOs {
		if blockHeight < v.EffectHeight {
			continue
		}
		fromAmount.Addition(v.Value)
		t := NewTicket(v.Txid, v.Index)
		fc.AddTicket(t)
	}
	if len(fc.Ticket) == 0 || fromAmount.GetInt64() < value.GetInt64() {
		log.Error("MakeFromCoin failed", "len(fc.Ticket)", len(fc.Ticket), "a.GetAmount().GetInt64()", a.GetAmount().GetInt64(), "value.GetInt64()", value.GetInt64())
		return nil, nil, errors.New("Account MakeFromCoin() can not cover value.the value is too large")
	}
	return fc, fromAmount, nil
}

func GetAccountHash(id node_meta.Address) math.Hash {
	return math.HashH(id.CloneBytes())
}
