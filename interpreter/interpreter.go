package interpreter

import "github.com/mihongtech/appchain/common/lcdb"

type Interpreter interface {
	Validator
	Processor
	CreateOffChain(db lcdb.Database) OffChain
}
