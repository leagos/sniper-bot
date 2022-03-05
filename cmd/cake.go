package cmd

import (
	"github.com/leagos/sniper-bot/runner"
	"github.com/spf13/cobra"
)

// cakeCmd represents the unicake command
var cakeCmd = &cobra.Command{
	Use:   "cake",
	Short: "sniper on pancake",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runner.NewEthRunner().SniperUniCake(chainType, quickMode)
	},
}

func init() {
	rootCmd.AddCommand(cakeCmd)
}
