package cmd

import (
	"github.com/go-gosh/gask/app"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run web server",
	Long:  `Run task management application web server`,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(app.Run())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
