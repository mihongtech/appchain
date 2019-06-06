package explorer

import (
	"github.com/mihongtech/appchain/client/explorer/rpc/rpcserver"
)

var shutdownRequestChannel = make(chan struct{})

func StartExplore() {
	s, err := rpcserver.NewRPCServer(&rpcserver.Config{})

	if err != nil {
		return
	}

	s.Start()

	go func() {
		<-s.RequestedProcessShutdown()
		shutdownRequestChannel <- struct{}{}
	}()
}
