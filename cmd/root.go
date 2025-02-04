package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sniper-bot",
	Short: "sniper-bot",
	// Uncomment the following line if your bare application
	// has an action associated with it:
}
var chainType string
var quickMode bool

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&chainType, "chain", "c", "bsc", "chain type")
	rootCmd.PersistentFlags().BoolVarP(&quickMode, "quick", "q", false, "quick mode")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetConfigFile("config.yml")
	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			if args := os.Args; len(args) >= 2 && args[1] == "init" {
				return
			}
		}
		log.Fatal("init config before sniper")
	}
}
