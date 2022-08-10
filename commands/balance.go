package commands

import (
	"context"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/mdedys/fuse/lightning"
	"github.com/mdedys/fuse/lnd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getBalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Get wallet balance",
	Long:  `Fetches the balance of the wallet`,
	Run: func(cmd *cobra.Command, args []string) {

		ctx := context.Background()

		var root Config
		err := viper.Unmarshal(&root)
		cobra.CheckErr(err)

		config := root.Configs[root.Active]

		lndClient, err := lnd.NewClient(
			config.Address,
			config.Network,
			config.Credentials.MacPath,
			config.Credentials.TlsPath,
			btcutil.Amount(1000),
		)

		cobra.CheckErr(err)

		client := lightning.New(lndClient)

		balance, err := client.WalletBalance(ctx)
		cobra.CheckErr(err)

		fmt.Println(balance.ToUnit(btcutil.AmountSatoshi))
	},
}

func init() {
	getCmd.AddCommand(getBalanceCmd)
}
