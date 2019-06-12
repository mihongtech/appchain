package app

import (
	"time"

	"github.com/mihongtech/appchain/app/context"
	"github.com/mihongtech/appchain/bcsi"
	"github.com/mihongtech/appchain/business/interpreter"
	"github.com/mihongtech/appchain/business/normal"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/rpc/rpcserver"
	"github.com/mihongtech/appchain/wallet"
	"github.com/mihongtech/linkchain-core/common/util/log"
	"github.com/mihongtech/linkchain-core/node"
	node_config "github.com/mihongtech/linkchain-core/node/config"
	"github.com/mihongtech/linkchain-core/storage"
)

var (
	appContext context.Context
	nodeSvc    *node.Node
	walletSvc  *wallet.Wallet
)

func Setup(globalConfig *config.LinkChainConfig) bool {
	log.Info("App setup...")

	//prepare config
	appContext.Config = globalConfig

	//create storage
	s := storage.NewStrorage(appContext.Config.DataDir + "app")
	if s == nil {
		log.Error("init storage failed")
		return false
	}

	//create bcsi service
	appContext.BCSIAPI = bcsi.NewBCSIServer(s.GetDB(), chooseInterpreterAPI(appContext.Config.InterpreterAPI))

	//create core service
	nodecfg := node.Config{BaseConfig: node_config.BaseConfig{
		DataDir:            globalConfig.DataDir,
		GenesisPath:        globalConfig.GenesisPath,
		ListenAddress:      globalConfig.ListenAddress,
		NoDiscovery:        globalConfig.NoDiscovery,
		BootstrapNodes:     globalConfig.BootstrapNodes,
		InterpreterAPIType: globalConfig.InterpreterAPI,
		RpcAddr:            globalConfig.RpcAddr,
	},
		BcsiAPI: appContext.BCSIAPI,
	}
	nodeSvc = node.NewNode(nodecfg.BaseConfig)

	//create wallet
	walletSvc = wallet.NewWallet()

	//node init
	if !nodeSvc.Setup(&nodecfg) {
		return false
	}
	//consensus api init
	appContext.NodeAPI = node.NewPublicCoreAPI(nodeSvc)

	//wallet init
	if !walletSvc.Setup(&appContext) {
		return false
	}
	//wallet api init
	appContext.WalletAPI = walletSvc

	return true
}

func Run() {
	//start all service
	nodeSvc.Start()
	walletSvc.Start()

	//start rpc
	startRPC()

	//here waiting for the interruption
	log.Info("App is running...")

	// listen the exit signal
	interrupt := interruptListener()
	<-interrupt
}

func Stop() {
	log.Info("Stopping app...")
	walletSvc.Stop()
	nodeSvc.Stop()
	log.Info("App exit")
}

func GetAppContext() *context.Context {
	return &appContext
}

func GetNodeAPI() *node.CoreAPI {
	return appContext.NodeAPI
}

func startRPC() {
	//init rpc servce
	s, err := rpcserver.NewRPCServer(&rpcserver.Config{
		StartupTime: time.Now().Unix(),
		Addr:        appContext.Config.RpcAddr,
	}, &appContext)
	if err != nil {
		return
	}

	s.Start()

	go func() {
		<-s.RequestedProcessShutdown()
		shutdownRequestChannel <- struct{}{}
	}()
}

func chooseInterpreterAPI(interpreter string) interpreter.Interpreter {
	log.Info("App", "interpreter", interpreter)
	switch interpreter {
	case "normal":
		return &normal.Interpreter{}
	default:
		return &normal.Interpreter{}
	}
}
