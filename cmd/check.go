package cmd

import (
	"github.com/go-gosh/gask/app/query"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Create a checkpoint for milestone",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.CheckMilestone(service.NewMilestone(query.Q))
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
