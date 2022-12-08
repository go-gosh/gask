package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-gosh/gask/app/conf"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Init config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		conf.GetConfig()
		return viper.SafeWriteConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
