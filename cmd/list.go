package cmd

import (
	"github.com/spf13/cobra"

	"github.com/go-gosh/gask/app/query"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all milestones",
	RunE: func(cmd *cobra.Command, args []string) error {
		page, err := cmd.Flags().GetInt("page")
		if err != nil {
			return err
		}
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return err
		}
		withCheck, err := cmd.Flags().GetBool("checkpoint")
		if err != nil {
			return err
		}
		return cli.PaginateMilestone(service.NewMilestone(query.Q), page, limit, withCheck)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntP("page", "p", 1, "page of data")
	listCmd.Flags().IntP("limit", "l", 10, "limit per page")
	listCmd.Flags().BoolP("checkpoint", "c", false, "show all checkpoints of milestone")
}
