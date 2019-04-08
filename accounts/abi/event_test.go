package abi

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/linkchain/common"
	"github.com/linkchain/common/math"
	"github.com/linkchain/core/meta"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jsonEventTransfer = []byte(`{
  "anonymous": false,
  "inputs": [
    {
      "indexed": true, "name": "from", "type": "address"
    }, {
      "indexed": true, "name": "to", "type": "address"
    }, {
      "indexed": false, "name": "value", "type": "uint256"
  }],
  "name": "Transfer",
  "type": "event"
}`)

var jsonEventPledge = []byte(`{
  "anonymous": false,
  "inputs": [{
      "indexed": false, "name": "who", "type": "address"
    }, {
      "indexed": false, "name": "wad", "type": "uint128"
    }, {
      "indexed": false, "name": "currency", "type": "bytes3"
  }],
  "name": "Pledge",
  "type": "event"
}`)

// 1000000
var transferData1 = "00000000000000000000000000000000000000000000000000000000000f4240"

// "0x00Ce0d46d924CC8437c806721496599FC3FFA268", 2218516807680, "usd"
var pledgeData1 = "00000000000000000000000000ce0d46d924cc8437c806721496599fc3ffa2680000000000000000000000000000000000000000000000000000020489e800007573640000000000000000000000000000000000000000000000000000000000"

func TestEventId(t *testing.T) {
	var table = []struct {
		definition   string
		expectations map[string]math.Hash
	}{
		{
			definition: `[
			{ "type" : "event", "name" : "balance", "inputs": [{ "name" : "in", "type": "uint256" }] },
			{ "type" : "event", "name" : "check", "inputs": [{ "name" : "t", "type": "address" }, { "name": "b", "type": "uint256" }] }
			]`,
			expectations: map[string]math.Hash{
				"balance": math.HashH([]byte("balance(uint256)")),
				"check":   math.HashH([]byte("check(address,uint256)")),
			},
		},
	}

	for _, test := range table {
		abi, err := JSON(strings.NewReader(test.definition))
		if err != nil {
			t.Fatal(err)
		}

		for name, event := range abi.Events {
			if event.Id() != test.expectations[name] {
				t.Errorf("expected id to be %x, got %x", test.expectations[name], event.Id())
			}
		}
	}
}

// TestEventMultiValueWithArrayUnpack verifies that array fields will be counted after parsing array.
func TestEventMultiValueWithArrayUnpack(t *testing.T) {
	definition := `[{"name": "test", "type": "event", "inputs": [{"indexed": false, "name":"value1", "type":"uint8[2]"},{"indexed": false, "name":"value2", "type":"uint8"}]}]`
	type testStruct struct {
		Value1 [2]uint8
		Value2 uint8
	}
	abi, err := JSON(strings.NewReader(definition))
	require.NoError(t, err)
	var b bytes.Buffer
	var i uint8 = 1
	for ; i <= 3; i++ {
		b.Write(packNum(reflect.ValueOf(i)))
	}
	var rst testStruct
	require.NoError(t, abi.Unpack(&rst, "test", b.Bytes()))
	require.Equal(t, [2]uint8{1, 2}, rst.Value1)
	require.Equal(t, uint8(3), rst.Value2)
}

func TestEventTupleUnpack(t *testing.T) {

	type EventTransfer struct {
		Value *big.Int
	}

	type EventPledge struct {
		Who      meta.AccountID
		Wad      *big.Int
		Currency [3]byte
	}

	type BadEventPledge struct {
		Who      string
		Wad      int
		Currency [3]byte
	}

	bigint := new(big.Int)
	bigintExpected := big.NewInt(1000000)
	bigintExpected2 := big.NewInt(2218516807680)
	addr, _ := meta.HexToAccountID("0x00Ce0d46d924CC8437c806721496599FC3FFA268")
	var testCases = []struct {
		data     string
		dest     interface{}
		expected interface{}
		jsonLog  []byte
		error    string
		name     string
	}{{
		transferData1,
		&EventTransfer{},
		&EventTransfer{Value: bigintExpected},
		jsonEventTransfer,
		"",
		"Can unpack ERC20 Transfer event into structure",
	}, {
		transferData1,
		&[]interface{}{&bigint},
		&[]interface{}{&bigintExpected},
		jsonEventTransfer,
		"",
		"Can unpack ERC20 Transfer event into slice",
	}, {
		pledgeData1,
		&EventPledge{},
		&EventPledge{
			addr,
			bigintExpected2,
			[3]byte{'u', 's', 'd'}},
		jsonEventPledge,
		"",
		"Can unpack Pledge event into structure",
	}, {
		pledgeData1,
		&[]interface{}{&meta.AccountID{}, &bigint, &[3]byte{}},
		&[]interface{}{
			&addr,
			&bigintExpected2,
			&[3]byte{'u', 's', 'd'}},
		jsonEventPledge,
		"",
		"Can unpack Pledge event into slice",
	}, {
		pledgeData1,
		&[3]interface{}{&meta.AccountID{}, &bigint, &[3]byte{}},
		&[3]interface{}{
			&addr,
			&bigintExpected2,
			&[3]byte{'u', 's', 'd'}},
		jsonEventPledge,
		"",
		"Can unpack Pledge event into an array",
	}, {
		pledgeData1,
		&[]interface{}{new(int), 0, 0},
		&[]interface{}{},
		jsonEventPledge,
		"abi: cannot unmarshal meta.AccountID in to int",
		"Can not unpack Pledge event into slice with wrong types",
	}, {
		pledgeData1,
		&BadEventPledge{},
		&BadEventPledge{},
		jsonEventPledge,
		"abi: cannot unmarshal meta.AccountID in to string",
		"Can not unpack Pledge event into struct with wrong filed types",
	}, {
		pledgeData1,
		&[]interface{}{meta.AccountID{}, new(big.Int)},
		&[]interface{}{},
		jsonEventPledge,
		"abi: insufficient number of elements in the list/array for unpack, want 3, got 2",
		"Can not unpack Pledge event into too short slice",
	}, {
		pledgeData1,
		new(map[string]interface{}),
		&[]interface{}{},
		jsonEventPledge,
		"abi: cannot unmarshal tuple into map[string]interface {}",
		"Can not unpack Pledge event into map",
	}}

	for _, tc := range testCases {
		assert := assert.New(t)
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := unpackTestEventData(tc.dest, tc.data, tc.jsonLog, assert)
			if tc.error == "" {
				assert.Nil(err, "Should be able to unpack event data.")
				assert.Equal(tc.expected, tc.dest, tc.name)
			} else {
				assert.EqualError(err, tc.error)
			}
		})
	}
}

func unpackTestEventData(dest interface{}, hexData string, jsonEvent []byte, assert *assert.Assertions) error {
	data, err := hex.DecodeString(hexData)
	assert.NoError(err, "Hex data should be a correct hex-string")
	var e Event
	assert.NoError(json.Unmarshal(jsonEvent, &e), "Should be able to unmarshal event ABI")
	a := ABI{Events: map[string]Event{"e": e}}
	return a.Unpack(dest, "e", data)
}

/*
Taken from
https://github.com/ethereum/go-ethereum/pull/15568
*/

type testResult struct {
	Values [2]*big.Int
	Value1 *big.Int
	Value2 *big.Int
}

type testCase struct {
	definition string
	want       testResult
}

func (tc testCase) encoded(intType, arrayType Type) []byte {
	var b bytes.Buffer
	if tc.want.Value1 != nil {
		val, _ := intType.pack(reflect.ValueOf(tc.want.Value1))
		b.Write(val)
	}

	if !reflect.DeepEqual(tc.want.Values, [2]*big.Int{nil, nil}) {
		val, _ := arrayType.pack(reflect.ValueOf(tc.want.Values))
		b.Write(val)
	}
	if tc.want.Value2 != nil {
		val, _ := intType.pack(reflect.ValueOf(tc.want.Value2))
		b.Write(val)
	}
	return b.Bytes()
}

// TestEventUnpackIndexed verifies that indexed field will be skipped by event decoder.
func TestEventUnpackIndexed(t *testing.T) {
	definition := `[{"name": "test", "type": "event", "inputs": [{"indexed": true, "name":"value1", "type":"uint8"},{"indexed": false, "name":"value2", "type":"uint8"}]}]`
	type testStruct struct {
		Value1 uint8
		Value2 uint8
	}
	abi, err := JSON(strings.NewReader(definition))
	require.NoError(t, err)
	var b bytes.Buffer
	b.Write(packNum(reflect.ValueOf(uint8(8))))
	var rst testStruct
	require.NoError(t, abi.Unpack(&rst, "test", b.Bytes()))
	require.Equal(t, uint8(0), rst.Value1)
	require.Equal(t, uint8(8), rst.Value2)
}

// TestEventIndexedWithArrayUnpack verifies that decoder will not overlow when static array is indexed input.
func TestEventIndexedWithArrayUnpack(t *testing.T) {
	definition := `[{"name": "test", "type": "event", "inputs": [{"indexed": true, "name":"value1", "type":"uint8[2]"},{"indexed": false, "name":"value2", "type":"string"}]}]`
	type testStruct struct {
		Value1 [2]uint8
		Value2 string
	}
	abi, err := JSON(strings.NewReader(definition))
	require.NoError(t, err)
	var b bytes.Buffer
	stringOut := "abc"
	// number of fields that will be encoded * 32
	b.Write(packNum(reflect.ValueOf(32)))
	b.Write(packNum(reflect.ValueOf(len(stringOut))))
	b.Write(common.RightPadBytes([]byte(stringOut), 32))

	var rst testStruct
	require.NoError(t, abi.Unpack(&rst, "test", b.Bytes()))
	require.Equal(t, [2]uint8{0, 0}, rst.Value1)
	require.Equal(t, stringOut, rst.Value2)
}