package consensus

import (
	"github.com/linkchain/common/util/event"
	"github.com/linkchain/consensus/manager"
	"github.com/linkchain/poa"
)

var (
	service poa.Service
)

type Service struct {
}

func (s *Service) Init(i interface{}) bool {
	//log.Info("consensus service init...");
	service = poa.Service{}
	service.Init(i)
	return true
}

func (s *Service) Start() bool {
	//log.Info("consensus service start...");
	service.Start()
	return true
}

func (s *Service) Stop() {
	//log.Info("consensus service stop...");
	service.Stop()
}

func (s *Service) GetBlockManager() manager.BlockManager {
	return service.GetManager().BlockManager
}

func (s *Service) GetTXManager() manager.TransactionManager {
	return service.GetManager().TransactionManager
}

func (s *Service) GetAccountManager() manager.AccountManager {
	return service.GetManager().AccountManager
}

func (s *Service) GetChainManager() manager.ChainManager {
	return service.GetManager().ChainManager
}

func (s *Service) GetBlockEvent() *event.TypeMux {
	return service.GetManager().NewBlockEvent
}

func (s *Service) GetTxEvent() *event.Feed {
	return service.GetManager().NewTxEvent
}