package node

import (
	"github.com/linkchain/common"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/consensus"
	"github.com/linkchain/function/wallet"
	"github.com/linkchain/p2p"
)

var (
	//service collection
	svcList = []common.IService{
		&consensus.Service{},
		&wallet.Wallet{},
		&p2p.Service{},
	}
)

func Init() {
	log.Info("Node init...")

	svcList[0].Init(nil)		//consensus init
	svcList[1].Init(GetConsensusService().GetAccountManager()) //wallet init
	svcList[2].Init(GetConsensusService()) //p2p init
}

func Run() {
	log.Info("Node is running...")

	//start all service
	for _, v := range svcList {
		v.Start()
	}

	/*block :=svcList[1].(*consensus.Service).GetBlockManager().NewBlock()
	svcList[1].(*consensus.Service).GetBlockManager().ProcessBlock(block)*/
}

//get service
func GetConsensusService() *consensus.Service {
	return svcList[0].(*consensus.Service)
}

func GetP2pService() *p2p.Service {
	return svcList[2].(*p2p.Service)
}

func GetWallet() *wallet.Wallet {
	return svcList[1].(*wallet.Wallet)
}
