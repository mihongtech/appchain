package cmd

import (
	"fmt"
	"strconv"

	"github.com/mihongtech/appchain/common/util/log"
	"github.com/mihongtech/appchain/rpc/rpcobject"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(blockCmd)
	blockCmd.AddCommand(heightCmd, hashCmd)
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "block command",
	Long:  "This is all block command for handling block",
}

var heightCmd = &cobra.Command{
	Use:     "height",
	Short:   "block <height>",
	Long:    "This is get block body command by height",
	Example: "block height 0",
	Run: func(cmd *cobra.Command, args []string) {
		example := []string{"example", "block height 0"}
		if len(args) != 1 {
			log.Error("getblockbyheight", "error", "please input height", example[0], example[1])
			return
		}

		height, err := strconv.Atoi(args[0])
		if err != nil {
			log.Error("getblockbyheight ", "error", err, example[0], example[1])
			return
		}

		method := "getBlockByHeight"

		//call
		out, err := rpc(method, &rpcobject.GetBlockByHeightCmd{height})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(out)
	},
}

var hashCmd = &cobra.Command{
	Use:     "hash",
	Short:   "block <hash>",
	Long:    "This is get block body command by blockhash",
	Example: "block hash 98acd27a58c79eaab05ea4abd0daa8e63021df3bf2e65fcb38e2474fb706c3fe",
	Run: func(cmd *cobra.Command, args []string) {
		example := []string{"example", "block hash 98acd27a58c79eaab05ea4abd0daa8e63021df3bf2e65fcb38e2474fb706c3fe"}
		if len(args) != 1 {
			log.Error("getblockbyhash", "error", "please input hash", example[0], example[1])
			return
		}

		hash := args[0]
		method := "getBlockByHash"

		//call
		out, err := rpc(method, &rpcobject.GetBlockByHashCmd{hash})
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(out)
	},
}
