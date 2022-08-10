package commands

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get {resource}",
	Short: "Get a specific resource",
	Long:  `Get detailed information about specific resource`,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
