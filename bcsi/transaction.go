package bcsi

import (
	"github.com/mihongtech/appchain/core/meta"
	node_meta "github.com/mihongtech/linkchain-core/core/meta"
)

func AppTransactionsConvert(txs *node_meta.Transactions) []meta.Transaction {
	resultTxs := make([]meta.Transaction, 0)
	for i := range txs.Txs {
		if resultTx, err := meta.ConvertToAppTX(txs.Txs[i]); err == nil {
			resultTxs = append(resultTxs, resultTx)
		}
	}
	return resultTxs
}
