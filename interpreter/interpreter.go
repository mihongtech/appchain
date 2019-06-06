package interpreter

import "github.com/mihongtech/linkchain-core/common/lcdb"

type Interpreter interface {
	Validator
	Processor
	CreateOffChain(db lcdb.Database) OffChain
}
