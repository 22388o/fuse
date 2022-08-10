package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var btcNode string

var btcdCmd = &cobra.Command{
	Use:   "btcd",
	Short: "Bitcoin node utility commands",
	Long: `Network root command encompasses all command used to interact 
with localhost lightning environment. The local environment is run through 
docker and contains a bitcoin and lightning node
`,
}

var blocks int

var mineCmd = &cobra.Command{
	Use:   "mine <address>",
	Short: "Mine some blocks",
	Long:  `Mine blocks immediately to a specified address`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(_ *cobra.Command, args []string) {

		str := fmt.Sprintf("docker exec -t %s /bin/bash -c \"bitcoin-cli -chain=regtest -rpcuser=regtest -rpcpassword=regtest -rpcwait generatetoaddress %v %v\"", btcNode, blocks, args[0])
		cmd := exec.Command("/bin/sh", "-c", str)

		output, err := cmd.CombinedOutput()
		cobra.CheckErr(err)

		fmt.Fprint(os.Stdout, string(output))
	},
}

var lightningCmd = &cobra.Command{
	Use:   "lightning",
	Short: "Lightning node utlity commands",
	Long: `Set of commands to interact with local lightning node
directly. Mostly used to setup the environment. (Funding, Addresses, etc)
`,
}

var newAddressCmd = &cobra.Command{
	Use:   "newaddress",
	Short: "Generates a new address",
	Run: func(_ *cobra.Command, __ []string) {
		str := "docker exec -t fuse_lnd /bin/bash -c \"lncli --tlscertpath /root/.lnd/tls.cert --macaroonpath /root/.lnd/data/chain/bitcoin/regtest/admin.macaroon newaddress np2wkh\""
		cmd := exec.Command("/bin/bash", "-c", str)
		output, err := cmd.CombinedOutput()
		cobra.CheckErr(err)
		fmt.Println(string(output))
	},
}

func init() {

	// Lightning Commands
	lightningCmd.AddCommand(newAddressCmd)
	rootCmd.AddCommand(lightningCmd)

	// Btcd Commands
	mineCmd.Flags().IntVarP(&blocks, "blocks", "b", 10, "number of blocks to mine")
	btcdCmd.PersistentFlags().StringVarP(&btcNode, "node", "n", "bitcoind", "Name of docker container running btcd")
	btcdCmd.AddCommand(mineCmd)
	rootCmd.AddCommand(btcdCmd)
}
