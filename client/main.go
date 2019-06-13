package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/mihongtech/appchain/client/cmd"
	"github.com/mihongtech/linkchain-core/common/util/log"
)

func main() {
	logLevel := flag.Int("loglevel", 3, "log level")
	rpcIp := flag.String("rpcip", "127.0.0.1", "linkchain rpc ip")
	rpcPort := flag.Int("rpcport", 8082, "linkchain rpc port")
	flag.Parse()

	//init log
	log.Root().SetHandler(
		log.LvlFilterHandler(log.Lvl(*logLevel),
			log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	log.Info("rpcserver client is running")

	cmd.StartCmd(*rpcIp + ":" + strconv.Itoa(*rpcPort))

}
