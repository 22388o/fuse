package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fusecli",
	Short: "Lightning Node CLI",
	Long: `Fuse is a CLI library to help make interacting
with lightning nodes much easier. Aimed at local development, it can be 
used with any lightning node. The idea being that one cli can be used
to interact with any type of lightning node (lnd, c-lightning, etc)`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}
