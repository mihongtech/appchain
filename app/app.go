package app

import (
	"time"

	"github.com/mihongtech/appchain/app/context"
	"github.com/mihongtech/appchain/common/util/log"
	"github.com/mihongtech/appchain/config"
	"github.com/mihongtech/appchain/contract"
	"github.com/mihongtech/appchain/interpreter"

	"github.com/mihongtech/appchain/normal"
	"github.com/mihongtech/appchain/rpc/rpcserver"
	"github.com/mihongtech/appchain/wallet"
	"github.com/mihongtech/linkchain-core/node"
	node_config "github.com/mihongtech/linkchain-core/node/config"
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

	//create interpreterAPI and Excutor by config choice different function
	appContext.InterpreterAPI = chooseInterpreterAPI(globalConfig.InterpreterAPI)

	//create service
	nodecfg := node.Config{BaseConfig: node_config.BaseConfig{
		DataDir:            globalConfig.DataDir,
		GenesisPath:        globalConfig.GenesisPath,
		ListenAddress:      globalConfig.ListenAddress,
		NoDiscovery:        globalConfig.NoDiscovery,
		BootstrapNodes:     globalConfig.BootstrapNodes,
		InterpreterAPIType: globalConfig.InterpreterAPI,
		RpcAddr:            globalConfig.RpcAddr,
	},
		BcsiAPI: nil,
	}
	nodeSvc = node.NewNode(nodecfg.BaseConfig)

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
	case "contract":
		return &contract.Interpreter{}
	default:
		return &normal.Interpreter{}
	}
}
