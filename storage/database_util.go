package storage

import (
	_ "bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/golang/protobuf/proto"
	"github.com/linkchain/common/lcdb"
	"github.com/linkchain/common/math"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/config"
	"github.com/linkchain/meta/block"
)

// DatabaseReader wraps the Get method of a backing data store.
type DatabaseReader interface {
	Get(key []byte) (value []byte, err error)
}

// DatabaseDeleter wraps the Delete method of a backing data store.
type DatabaseDeleter interface {
	Delete(key []byte) error
}

var (
	headBlockKey = []byte("LastBlock")
	headFastKey  = []byte("LastFast")
	trieSyncKey  = []byte("TrieSync")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`).
	blockPrefix         = []byte("h") // blockPrefix + num (uint64 big endian) + hash -> header
	tdSuffix            = []byte("t") // blockPrefix + num (uint64 big endian) + hash + tdSuffix -> td
	numSuffix           = []byte("n") // blockPrefix + num (uint64 big endian) + numSuffix -> hash
	blockHashPrefix     = []byte("H") // blockHashPrefix + hash -> num (uint64 big endian)
	bodyPrefix          = []byte("b") // bodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts
	lookupPrefix        = []byte("l") // lookupPrefix + hash -> transaction/receipt lookup metadata
	bloomBitsPrefix     = []byte("B") // bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits

	preimagePrefix = "secure-key-"              // preimagePrefix + hash -> preimage
	configPrefix   = []byte("ethereum-config-") // config prefix for the db

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress

	// used by old db, now only used for conversion
	oldReceiptsPrefix = []byte("receipts-")
	oldTxMetaSuffix   = []byte{0x01}

	ErrChainConfigNotFound = errors.New("ChainConfig not found") // general config not found error
)

// TxLookupEntry is a positional metadata to help looking up the data content of
// a transaction or receipt given only its hash.
type TxLookupEntry struct {
	BlockHash  math.Hash
	BlockIndex uint64
	Index      uint64
}

// encodeBlockNumber encodes a block number as big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

// GetCanonicalHash retrieves a hash assigned to a canonical block number.
func GetCanonicalHash(db DatabaseReader, number uint64) math.Hash {
	data, _ := db.Get(append(append(blockPrefix, encodeBlockNumber(number)...), numSuffix...))
	if len(data) == 0 {
		return math.Hash{}
	}
	return math.BytesToHash(data)
}

// missingNumber is returned by GetBlockNumber if no header with the
// given block hash has been stored in the database
const MissingNumber = uint64(0xffffffffffffffff)

// GetBlockNumber returns the block number assigned to a block hash
// if the corresponding header is present in the database
func GetBlockNumber(db DatabaseReader, hash math.Hash) uint64 {
	data, _ := db.Get(append(blockHashPrefix, hash.Bytes()...))
	if len(data) != 8 {
		return MissingNumber
	}
	return binary.BigEndian.Uint64(data)
}

// GetHeadHeaderHash retrieves the hash of the current canonical head block's
// header. The difference between this and GetHeadBlockHash is that whereas the
// last block hash is only updated upon a full block import, the last header
// hash is updated already at header import, allowing head tracking for the
// light synchronization mechanism.
func GetHeadBlockHash(db DatabaseReader) math.Hash {
	data, _ := db.Get(headBlockKey)
	if len(data) == 0 {
		return math.Hash{}
	}
	return math.BytesToHash(data)
}

// GetHeadFastBlockHash retrieves the hash of the current canonical head block during
// fast synchronization. The difference between this and GetHeadBlockHash is that
// whereas the last block hash is only updated upon a full block import, the last
// fast hash is updated when importing pre-processed blocks.
func GetHeadFastBlockHash(db DatabaseReader) math.Hash {
	data, _ := db.Get(headFastKey)
	if len(data) == 0 {
		return math.Hash{}
	}
	return math.BytesToHash(data)
}

// GetTrieSyncProgress retrieves the number of tries nodes fast synced to allow
// reportinc correct numbers across restarts.
func GetTrieSyncProgress(db DatabaseReader) uint64 {
	data, _ := db.Get(trieSyncKey)
	if len(data) == 0 {
		return 0
	}
	return new(big.Int).SetBytes(data).Uint64()
}

// GetHeaderRLP retrieves a block header in its raw RLP database encoding, or nil
// if the header's not found.
// TODO: implement
//func GetHeaderRLP(db DatabaseReader, hash math.Hash, number uint64) rlp.RawValue {
//	data, _ := db.Get(headerKey(hash, number))
//	return data
//}

// GetHeader retrieves the block header corresponding to the hash, nil if none
// found.
//func GetHeader(db DatabaseReader, hash math.Hash, number uint64) *types.Header {
//	data := GetHeaderRLP(db, hash, number)
//	if len(data) == 0 {
//		return nil
//	}
//	header := new(types.Header)
//	if err := rlp.Decode(bytes.NewReader(data), header); err != nil {
//		log.Error("Invalid block header RLP", "hash", hash, "err", err)
//		return nil
//	}
//	return header
//}

// GetBodyRLP retrieves the block body (transactions and uncles) in RLP encoding.
//func GetBodyRLP(db DatabaseReader, hash math.Hash, number uint64) rlp.RawValue {
//	data, _ := db.Get(blockBodyKey(hash, number))
//	return data
//}

func blockKey(hash math.Hash, number uint64) []byte {
	return append(append(blockPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

func blockBodyKey(hash math.Hash, number uint64) []byte {
	return append(append(bodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

// GetBody retrieves the block body (transactons, uncles) corresponding to the
// hash, nil if none found.
//func GetBody(db DatabaseReader, hash math.Hash, number uint64) *types.Body {
//	data := GetBodyRLP(db, hash, number)
//	if len(data) == 0 {
//		return nil
//	}
//	body := new(types.Body)
//	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
//		log.Error("Invalid block body RLP", "hash", hash, "err", err)
//		return nil
//	}
//	return body
//}

// GetTd retrieves a block's total difficulty corresponding to the hash, nil if
// none found.
//func GetTd(db DatabaseReader, hash math.Hash, number uint64) *big.Int {
//	data, _ := db.Get(append(append(append(headerPrefix, encodeBlockNumber(number)...), hash[:]...), tdSuffix...))
//	if len(data) == 0 {
//		return nil
//	}
//	td := new(big.Int)
//	if err := rlp.Decode(bytes.NewReader(data), td); err != nil {
//		log.Error("Invalid block total difficulty RLP", "hash", hash, "err", err)
//		return nil
//	}
//	return td
//}

// GetBlock retrieves an entire block corresponding to the hash, assembling it
// back from the stored header and body. If either the header or body could not
// be retrieved nil is returned.
//
// Note, due to concurrent download of header and block body the header and thus
// canonical hash can be stored in the database but the body data not (yet).
//func GetBlock(db DatabaseReader, hash math.Hash, number uint64) *meta.Block {
//	// Retrieve the block header and body contents
//	header := GetHeader(db, hash, number)
//	if header == nil {
//		return nil
//	}
//	body := GetBody(db, hash, number)
//	if body == nil {
//		return nil
//	}
//	// Reassemble the block and return
//	return types.NewBlockWithHeader(header).WithBody(body.Transactions, body.Uncles)
//}

// GetBlockReceipts retrieves the receipts generated by the transactions included
// in a block given by its hash.
//func GetBlockReceipts(db DatabaseReader, hash math.Hash, number uint64) types.Receipts {
//	data, _ := db.Get(append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash[:]...))
//	if len(data) == 0 {
//		return nil
//	}
//	storageReceipts := []*types.ReceiptForStorage{}
//	if err := rlp.DecodeBytes(data, &storageReceipts); err != nil {
//		log.Error("Invalid receipt array RLP", "hash", hash, "err", err)
//		return nil
//	}
//	receipts := make(types.Receipts, len(storageReceipts))
//	for i, receipt := range storageReceipts {
//		receipts[i] = (*types.Receipt)(receipt)
//	}
//	return receipts
//}

// GetTxLookupEntry retrieves the positional metadata associated with a transaction
// hash to allow retrieving the transaction or receipt by hash.
//func GetTxLookupEntry(db DatabaseReader, hash math.Hash) (math.Hash, uint64, uint64) {
//	// Load the positional metadata from disk and bail if it fails
//	data, _ := db.Get(append(lookupPrefix, hash.Bytes()...))
//	if len(data) == 0 {
//		return math.Hash{}, 0, 0
//	}
//	// Parse and return the contents of the lookup entry
//	var entry TxLookupEntry
//	if err := rlp.DecodeBytes(data, &entry); err != nil {
//		log.Error("Invalid lookup entry RLP", "hash", hash, "err", err)
//		return math.Hash{}, 0, 0
//	}
//	return entry.BlockHash, entry.BlockIndex, entry.Index
//}

// GetTransaction retrieves a specific transaction from the database, along with
// its added positional metadata.
//func GetTransaction(db DatabaseReader, hash math.Hash) (*types.Transaction, math.Hash, uint64, uint64) {
//	// Retrieve the lookup metadata and resolve the transaction from the body
//	blockHash, blockNumber, txIndex := GetTxLookupEntry(db, hash)
//
//	if blockHash != (math.Hash{}) {
//		body := GetBody(db, blockHash, blockNumber)
//		if body == nil || len(body.Transactions) <= int(txIndex) {
//			log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash, "index", txIndex)
//			return nil, math.Hash{}, 0, 0
//		}
//		return body.Transactions[txIndex], blockHash, blockNumber, txIndex
//	}
//	// Old transaction representation, load the transaction and it's metadata separately
//	data, _ := db.Get(hash.Bytes())
//	if len(data) == 0 {
//		return nil, math.Hash{}, 0, 0
//	}
//	var tx types.Transaction
//	if err := rlp.DecodeBytes(data, &tx); err != nil {
//		return nil, math.Hash{}, 0, 0
//	}
//	// Retrieve the blockchain positional metadata
//	data, _ = db.Get(append(hash.Bytes(), oldTxMetaSuffix...))
//	if len(data) == 0 {
//		return nil, math.Hash{}, 0, 0
//	}
//	var entry TxLookupEntry
//	if err := rlp.DecodeBytes(data, &entry); err != nil {
//		return nil, math.Hash{}, 0, 0
//	}
//	return &tx, entry.BlockHash, entry.BlockIndex, entry.Index
//}

// GetReceipt retrieves a specific transaction receipt from the database, along with
// its added positional metadata.
//func GetReceipt(db DatabaseReader, hash math.Hash) (*types.Receipt, math.Hash, uint64, uint64) {
//	// Retrieve the lookup metadata and resolve the receipt from the receipts
//	blockHash, blockNumber, receiptIndex := GetTxLookupEntry(db, hash)
//
//	if blockHash != (math.Hash{}) {
//		receipts := GetBlockReceipts(db, blockHash, blockNumber)
//		if len(receipts) <= int(receiptIndex) {
//			log.Error("Receipt refereced missing", "number", blockNumber, "hash", blockHash, "index", receiptIndex)
//			return nil, math.Hash{}, 0, 0
//		}
//		return receipts[receiptIndex], blockHash, blockNumber, receiptIndex
//	}
//	// Old receipt representation, load the receipt and set an unknown metadata
//	data, _ := db.Get(append(oldReceiptsPrefix, hash[:]...))
//	if len(data) == 0 {
//		return nil, math.Hash{}, 0, 0
//	}
//	var receipt types.ReceiptForStorage
//	err := rlp.DecodeBytes(data, &receipt)
//	if err != nil {
//		log.Error("Invalid receipt RLP", "hash", hash, "err", err)
//	}
//	return (*types.Receipt)(&receipt), math.Hash{}, 0, 0
//}

// GetBloomBits retrieves the compressed bloom bit vector belonging to the given
// section and bit index from the.
func GetBloomBits(db DatabaseReader, bit uint, section uint64, head math.Hash) ([]byte, error) {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), head.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return db.Get(key)
}

// WriteCanonicalHash stores the canonical hash for the given block number.
func WriteCanonicalHash(db lcdb.Putter, hash math.Hash, number uint64) error {
	key := append(append(blockPrefix, encodeBlockNumber(number)...), numSuffix...)
	if err := db.Put(key, hash.Bytes()); err != nil {
		log.Crit("Failed to store number to hash mapping", "err", err)
	}
	return nil
}

// WriteHeadBlockHash stores the head block's hash.
func WriteHeadBlockHash(db lcdb.Putter, hash math.Hash) error {
	if err := db.Put(headBlockKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last block's hash", "err", err)
	}
	return nil
}

// WriteHeadFastBlockHash stores the fast head block's hash.
func WriteHeadFastBlockHash(db lcdb.Putter, hash math.Hash) error {
	if err := db.Put(headFastKey, hash.Bytes()); err != nil {
		log.Crit("Failed to store last fast block's hash", "err", err)
	}
	return nil
}

// WriteTrieSyncProgress stores the fast sync trie process counter to support
// retrieving it across restarts.
func WriteTrieSyncProgress(db lcdb.Putter, count uint64) error {
	if err := db.Put(trieSyncKey, new(big.Int).SetUint64(count).Bytes()); err != nil {
		log.Crit("Failed to store fast sync trie progress", "err", err)
	}
	return nil
}

// WriteBlock serializes a block into the database, header and body separately.
func WriteBlock(db lcdb.Putter, block block.IBlock) error {

	data := block.Serialize()
	bytesData, err := proto.Marshal(data)
	if err != nil {
		return err
	}

	hash := block.GetBlockID().(*math.Hash).Bytes()
	num := block.GetHeight()
	encNum := encodeBlockNumber(uint64(num))
	key := append(blockHashPrefix, hash...)
	if err := db.Put(key, encNum); err != nil {
		log.Crit("Failed to store hash to number mapping", "err", err)
	}
	key = append(append(blockPrefix, encNum...), hash...)

	if err := db.Put(key, bytesData); err != nil {
		log.Crit("Failed to store block", "err", err)
	}
	return nil
}

// WriteTxLookupEntries stores a positional metadata for every transaction from
// a block, enabling hash based transaction and receipt lookups.
//func WriteTxLookupEntries(db lcdb.Putter, block *types.Block) error {
//	// Iterate over each transaction and encode its metadata
//	for i, tx := range block.Transactions() {
//		entry := TxLookupEntry{
//			BlockHash:  block.Hash(),
//			BlockIndex: block.NumberU64(),
//			Index:      uint64(i),
//		}
//		data, err := rlp.EncodeToBytes(entry)
//		if err != nil {
//			return err
//		}
//		if err := db.Put(append(lookupPrefix, tx.Hash().Bytes()...), data); err != nil {
//			return err
//		}
//	}
//	return nil
//}

// WriteBloomBits writes the compressed bloom bits vector belonging to the given
// section and bit index.
func WriteBloomBits(db lcdb.Putter, bit uint, section uint64, head math.Hash, bits []byte) {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), head.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	if err := db.Put(key, bits); err != nil {
		log.Crit("Failed to store bloom bits", "err", err)
	}
}

// DeleteCanonicalHash removes the number to hash canonical mapping.
func DeleteCanonicalHash(db DatabaseDeleter, number uint64) {
	db.Delete(append(append(blockPrefix, encodeBlockNumber(number)...), numSuffix...))
}

// DeleteHeader removes all block header data associated with a hash.
func DeleteHeader(db DatabaseDeleter, hash math.Hash, number uint64) {
	db.Delete(append(blockHashPrefix, hash.Bytes()...))
	db.Delete(append(append(blockPrefix, encodeBlockNumber(number)...), hash.Bytes()...))
}

// DeleteBody removes all block body data associated with a hash.
func DeleteBody(db DatabaseDeleter, hash math.Hash, number uint64) {
	db.Delete(append(append(bodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...))
}

// DeleteTd removes all block total difficulty data associated with a hash.
func DeleteTd(db DatabaseDeleter, hash math.Hash, number uint64) {
	db.Delete(append(append(append(blockPrefix, encodeBlockNumber(number)...), hash.Bytes()...), tdSuffix...))
}

// DeleteBlock removes all block data associated with a hash.
func DeleteBlock(db DatabaseDeleter, hash math.Hash, number uint64) {
	DeleteBlockReceipts(db, hash, number)
	DeleteHeader(db, hash, number)
	DeleteBody(db, hash, number)
	DeleteTd(db, hash, number)
}

// DeleteBlockReceipts removes all receipt data associated with a block hash.
func DeleteBlockReceipts(db DatabaseDeleter, hash math.Hash, number uint64) {
	db.Delete(append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash.Bytes()...))
}

// DeleteTxLookupEntry removes all transaction data associated with a hash.
func DeleteTxLookupEntry(db DatabaseDeleter, hash math.Hash) {
	db.Delete(append(lookupPrefix, hash.Bytes()...))
}

// PreimageTable returns a Database instance with the key prefix for preimage entries.
func PreimageTable(db lcdb.Database) lcdb.Database {
	return lcdb.NewTable(db, preimagePrefix)
}

// WritePreimages writes the provided set of preimages to the database. `number` is the
// current block number, and is used for debug messages only.
func WritePreimages(db lcdb.Database, number uint64, preimages map[math.Hash][]byte) error {
	table := PreimageTable(db)
	batch := table.NewBatch()
	hitCount := 0
	for hash, preimage := range preimages {
		if _, err := table.Get(hash.Bytes()); err != nil {
			batch.Put(hash.Bytes(), preimage)
			hitCount++
		}
	}

	if hitCount > 0 {
		if err := batch.Write(); err != nil {
			return fmt.Errorf("preimage write fail for block %d: %v", number, err)
		}
	}
	return nil
}

// GetBlockChainVersion reads the version number from db.
//func GetBlockChainVersion(db DatabaseReader) int {
//	var vsn uint
//	enc, _ := db.Get([]byte("BlockchainVersion"))
//	rlp.DecodeBytes(enc, &vsn)
//	return int(vsn)
//}

// WriteBlockChainVersion writes vsn as the version number to db.
//func WriteBlockChainVersion(db lcdb.Putter, vsn int) {
//	enc, _ := rlp.EncodeToBytes(uint(vsn))
//	db.Put([]byte("BlockchainVersion"), enc)
//}

// WriteChainConfig writes the chain config settings to the database.
func WriteChainConfig(db lcdb.Putter, hash *math.Hash, cfg *config.ChainConfig) error {
	// short circuit and ignore if nil config. GetChainConfig
	// will return a default.
	if cfg == nil {
		return nil
	}

	jsonChainConfig, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return db.Put(append(configPrefix, hash.Bytes()[:]...), jsonChainConfig)
}

// GetChainConfig will fetch the network settings based on the given hash.
func GetChainConfig(db DatabaseReader, hash math.Hash) (*config.ChainConfig, error) {
	jsonChainConfig, _ := db.Get(append(configPrefix, hash[:]...))
	if len(jsonChainConfig) == 0 {
		return nil, ErrChainConfigNotFound
	}

	var chainConfig config.ChainConfig
	if err := json.Unmarshal(jsonChainConfig, &chainConfig); err != nil {
		return nil, err
	}

	return &chainConfig, nil
}

// FindCommonAncestor returns the last common ancestor of two block headers
//func FindCommonAncestor(db DatabaseReader, a, b *types.Header) *types.Header {
//	for bn := b.Number.Uint64(); a.Number.Uint64() > bn; {
//		a = GetHeader(db, a.ParentHash, a.Number.Uint64()-1)
//		if a == nil {
//			return nil
//		}
//	}
//	for an := a.Number.Uint64(); an < b.Number.Uint64(); {
//		b = GetHeader(db, b.ParentHash, b.Number.Uint64()-1)
//		if b == nil {
//			return nil
//		}
//	}
//	for a.Hash() != b.Hash() {
//		a = GetHeader(db, a.ParentHash, a.Number.Uint64()-1)
//		if a == nil {
//			return nil
//		}
//		b = GetHeader(db, b.ParentHash, b.Number.Uint64()-1)
//		if b == nil {
//			return nil
//		}
//	}
//	return a
//}
