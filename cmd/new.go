package cmd

import (
	"github.com/spf13/cobra"

	"github.com/go-gosh/gask/app/query"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new milestone",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.NewMilestone(service.NewMilestone(query.Q))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
