package wallet

import (
	"encoding/hex"
	"errors"
	"github.com/mihongtech/appchain/bcsi"
	"github.com/mihongtech/appchain/storage/state"
	"path/filepath"

	"github.com/mihongtech/appchain/app/context"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/helper"
	"github.com/mihongtech/linkchain-core/accounts"
	"github.com/mihongtech/linkchain-core/accounts/keystore"
	"github.com/mihongtech/linkchain-core/common/btcec"
	"github.com/mihongtech/linkchain-core/common/math"
	"github.com/mihongtech/linkchain-core/common/util/log"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

type Wallet struct {
	keystore *keystore.KeyStore
	password string
	Name     string
	DataDir  string
	accounts map[string]meta.Account
	bcsiAPI  *bcsi.BCSIServer
}

func NewWallet() *Wallet {
	name := "wallet"
	password := "password"
	return &Wallet{accounts: make(map[string]meta.Account), Name: name, password: password}
}

func (w *Wallet) Setup(i interface{}) bool {
	globalConfig := i.(*context.Context).Config

	w.DataDir = globalConfig.DataDir
	path := w.instanceDir(w.DataDir)
	w.keystore = keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	w.bcsiAPI = i.(*context.Context).BCSIAPI
	return true
}

func (w *Wallet) Start() bool {
	ksAccounts := w.keystore.Accounts()
	for i := range ksAccounts {
		account := helper.CreateTemplateAccount(ksAccounts[i].Address)
		w.accounts[account.Id.String()] = *account
	}
	w.reScanAllAccount()

	return true
}

func (w *Wallet) Stop() {
	log.Info("Stop wallet...")
}

func (w *Wallet) reScanAllAccount() {
	newWas := make([]meta.Account, 0)
	for key := range w.accounts {
		wa := w.accounts[key]
		newWa, err := w.queryAccount(wa.Id)
		if err != nil {
			continue
		}

		newWas = append(newWas, newWa)
	}
	for _, wa := range newWas {
		w.updateWalletAccount(wa)
	}
}

func (w *Wallet) updateWalletAccount(account meta.Account) error {
	a, ok := w.accounts[account.GetAccountID().String()]
	if !ok {
		return errors.New("GetAccountID can not find account")
	}

	a = account
	w.AddAccount(a)
	return nil
}

func (w *Wallet) NewAccount() (*node_meta.Address, error) {
	ksAccount, err := w.keystore.NewAccount(w.password)
	if err != nil {
		log.Error("wallet", "newAccount", err)
		return nil, err
	}
	account := helper.CreateTemplateAccount(ksAccount.Address)
	w.AddAccount(*account)
	return &ksAccount.Address, nil
}

func (w *Wallet) AddAccount(account meta.Account) {
	w.accounts[account.Id.String()] = account
}

func (w *Wallet) GetAllWAccount() []meta.Account {
	w.reScanAllAccount()
	var WAs []meta.Account
	for a := range w.accounts {
		WAs = append(WAs, w.accounts[a])
	}
	return WAs
}

func (w *Wallet) GetAccount(key string) (*meta.Account, error) {
	wa, ok := w.accounts[key]
	if ok {
		return &wa, nil
	} else {
		id, err := node_meta.HexToAddress(key)
		if err != nil {
			return nil, err
		}
		newWa, err := w.queryAccount(id)
		return &newWa, err
	}
}

func (w *Wallet) queryAccount(id node_meta.Address) (meta.Account, error) {
	currentBlockData := w.bcsiAPI.CurrentBlock.Load()
	if currentBlockData == nil {
		return meta.Account{}, errors.New("best block have not store")
	}
	bestBlock := currentBlockData.(node_meta.Block)
	root, err := w.bcsiAPI.GetBlockState(*bestBlock.GetBlockID())
	if err != nil {
		return meta.Account{}, err
	}
	stateDB, err := state.New(root, w.bcsiAPI.Db)
	if err != nil {
		return meta.Account{}, err
	}
	stateObject := stateDB.GetObject(meta.GetAccountHash(id))
	if stateObject == nil {
		return meta.Account{}, errors.New("can not find IAccount")
	}
	return *stateObject.GetAccount(), nil
}

func (w *Wallet) SignTransaction(tx meta.Transaction) (*meta.Transaction, error) {
	for _, fc := range tx.GetFromCoins() {
		sign, err := w.SignMessage(fc.Id, tx.GetTxID().CloneBytes())
		if err != nil {
			return nil, err
		}
		tx.AddSignature(sign)
	}
	return &tx, nil
}

func (w *Wallet) SignMessage(accountId node_meta.Address, hash []byte) (math.ISignature, error) {
	_, ok := w.accounts[accountId.String()]
	if !ok {
		return nil, errors.New("SignMessage can not find account id")
	}

	ksAccount, err := w.keystore.Find(accounts.Account{Address: accountId})
	if err != nil {
		return nil, err
	}
	sign, err := w.keystore.SignHashWithPassphrase(ksAccount, w.password, hash)
	return node_meta.NewSignature(sign), nil
}

func (w *Wallet) importKey(privkeyStr string) (*node_meta.Address, error) {
	privkeyBuff, err := hex.DecodeString(privkeyStr)
	if err != nil {
		return nil, err
	}
	privkey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privkeyBuff)
	ksAccount, err := w.keystore.ImportECDSA(privkey, w.password)
	if err != nil {
		return nil, err
	}
	account := helper.CreateTemplateAccount(ksAccount.Address)
	w.AddAccount(*account)
	return &ksAccount.Address, err
}

func (w *Wallet) ImportAccount(privateKeyStr string) (*node_meta.Address, error) {
	a, err := w.importKey(privateKeyStr)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (w *Wallet) ExportAccount(id node_meta.Address) (string, error) {
	_, ok := w.accounts[id.String()]
	if !ok {
		return "", errors.New("export can not find account id")
	}
	ksAccount, err := w.keystore.Find(accounts.Account{Address: id})
	if err != nil {
		return "", err
	}

	return w.keystore.ExportECDSA(ksAccount, w.password)
}

func (w *Wallet) instanceDir(path string) string {
	if path == "" {
		return ""
	}
	return filepath.Join(path, w.Name)
}
