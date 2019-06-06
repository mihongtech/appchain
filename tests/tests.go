package tests

import (
	"encoding/hex"
	"encoding/json"
	"github.com/mihongtech/linkchain-core/common/lcdb"
	"github.com/mihongtech/linkchain-core/common/math"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/mihongtech/appchain/core/meta"
	"github.com/mihongtech/appchain/protobuf"
	"github.com/mihongtech/appchain/storage/state"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"

	"path"
	"path/filepath"
	"runtime"
	"testing"
)

//Get block by testData
func ConvertTestBlock(hexData string) *node_meta.Block {
	buffer, _ := hex.DecodeString(hexData)
	pb := &protobuf.Block{}
	proto.Unmarshal(buffer, pb)

	block := node_meta.Block{}
	block.Deserialize(pb)

	return &block
}

//Get stateDb by testData
func ConvertTestStateDb(db lcdb.Database, data ...StateDB) *state.StateDB {
	stateDb, _ := state.New(math.Hash{}, db)
	for _, d := range data {
		buffer, _ := hex.DecodeString(d.Account)
		pa := &protobuf.Account{}

		proto.Unmarshal(buffer, pa)

		a := meta.Account{}
		a.Deserialize(pa)

		aObj := stateDb.NewObject(meta.GetAccountHash(a.Id), a)
		stateDb.SetObject(aObj)
	}
	stateDb.Commit()

	return stateDb
}

//Get tx by testData
func ConvertTestTransaction(hexData string) *meta.Transaction {
	buffer, _ := hex.DecodeString(hexData)
	tx := &protobuf.Transaction{}

	proto.Unmarshal(buffer, tx)

	newTx := meta.Transaction{}
	newTx.Deserialize(tx)

	return &newTx
}

//Get accountId by testData
func ConvertTestAccountId(data string) *node_meta.Address {
	buffer, _ := hex.DecodeString(data)
	id := node_meta.BytesToAddress(buffer)
	return &id
}

type Transaction struct {
	TxId    string `json:"TxId"`
	Height  int    `json:"Height"`
	Signer  string `json:"Signer"`
	HexData string `json:"HexData"`
}

type StateDB struct {
	AccountId string `json:"AccountId"`
	Account   string `json:"Account"`
	Amount    int    `json:"Amount"`
}

type Block struct {
	BlockId string        `json:"BlockId"`
	Height  int           `json:"Height"`
	Miner   string        `json:"Miner"`
	HexData string        `json:"HexData"`
	StateDB []StateDB     `json:"stateDB"`
	Tx      []Transaction `json:"Tx"`
}

type Chain struct {
	Db    lcdb.Database
	Chain []Block `json:"chain"`
}

func (c *Chain) GetBlock(height int) *node_meta.Block {
	return ConvertTestBlock(c.Chain[height].HexData)
}

func (c *Chain) GetStateDB(height int) *state.StateDB {
	return ConvertTestStateDb(c.Db, c.Chain[height].StateDB...)
}

func GetChain(db lcdb.Database, t *testing.T) *Chain {
	_, filename, _, _ := runtime.Caller(0)
	testData := filepath.Join(path.Dir(filename), "testdata/testdata.json")
	t.Log("get file", "file path", testData)
	file, err := os.Open(testData)
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d\n open test file error:%s",
			filepath.Base(file), line, err)
		return nil
	}
	defer file.Close()

	chain := new(Chain)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&chain); err != nil {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d\n json decode error:%s",
			filepath.Base(file), line, err)
		return nil
	}
	chain.Db = db
	return chain
}
