package cmd

import (
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/milestone"
	"github.com/go-gosh/gask/ui/cli"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use: "check",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := global.GetDatabase()
		if err != nil {
			return err
		}
		return cli.CheckMilestone(milestone.New(db))
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
