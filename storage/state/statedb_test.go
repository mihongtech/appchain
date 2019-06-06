package state

import (
	"fmt"
	"testing"

	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/linkchain-core/common/btcec"
	"github.com/mihongtech/linkchain-core/common/lcdb"
	"github.com/mihongtech/linkchain-core/common/math"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func TestNew(t *testing.T) {
	db, _ := lcdb.NewMemDatabase()
	sdb, _ := New(math.Hash{}, db)

	hash := sdb.GetRootHash()
	fmt.Println("Root hash:", hash.String())
}

func TestStateDB(t *testing.T) {
	//new state db
	db, _ := lcdb.NewMemDatabase()

	root := writeStateDB(math.Hash{}, db)

	//load state rpcobject
	readStateDB(root, db)
}

var (
	obj1key string
	obj2key string
)

func writeStateDB(root math.Hash, db lcdb.Database) math.Hash {
	sdb, _ := New(root, db)

	hash := sdb.GetRootHash()
	fmt.Println("Root hash:", hash.String())

	//create state rpcobject
	a1 := getNewAccount(10)
	obj1 := sdb.NewObject(math.HashH(a1.GetAccountID().CloneBytes()), *a1)

	sdb.SetObject(obj1)
	obj1key = obj1.key.String()
	fmt.Println("rpcobject 1 key:", obj1key)

	a2 := getNewAccount(20)
	obj2 := sdb.NewObject(math.HashH(a2.GetAccountID().CloneBytes()), *a2)

	sdb.SetObject(obj2)
	obj2key = obj2.key.String()
	fmt.Println("rpcobject 2 key:", obj2key)

	//save all state
	sdb.Commit()
	fmt.Println("New Root hash:", sdb.GetRootHash())

	return sdb.GetRootHash()
}

func readStateDB(root math.Hash, db lcdb.Database) {
	sdb, _ := New(root, db)

	hash := sdb.GetRootHash()
	fmt.Println("Root hash:", hash.String())

	//create state rpcobject
	h1, _ := math.NewHashFromStr(obj1key)
	robj1 := sdb.GetObject(*h1)
	_ = robj1

	h2, _ := math.NewHashFromStr(obj2key)
	robj2 := sdb.GetObject(*h2)
	_ = robj2

	fmt.Println("Finish read rpcobject")

}

func getNewAccount(amount int64) *meta.Account {
	ex, _ := btcec.NewPrivateKey(btcec.S256())
	id := node_meta.NewAddress(ex.PubKey())

	utxos := make([]meta.UTXO, 0)

	account := meta.NewAccount(*id, config.NormalAccount, utxos)
	txid, _ := math.NewHashFromStr("5e6e12fc6cddbcdac39a9b265402960473fd2640a65ef32e558f89b47be40f64")
	ticket := meta.NewTicket(*txid, 0)

	u := meta.NewUTXO(ticket, 120, 150, *meta.NewAmount(amount))

	account.UTXOs = append(account.UTXOs, *u)

	return account
}
