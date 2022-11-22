package cmd

import (
	"github.com/spf13/cobra"

	"github.com/go-gosh/gask/app/query"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

func mustGetFlag[T any](t T, err error) T {
	cobra.CheckErr(err)
	return t
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all milestones",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.PaginateMilestone(
			service.NewMilestone(query.Q),
			mustGetFlag(cmd.Flags().GetInt("page")),
			mustGetFlag(cmd.Flags().GetInt("limit")),
			mustGetFlag(cmd.Flags().GetBool("checkpoint")),
		)
	},
	SilenceUsage: true,
}

// listCheckpointCmd represents the list command
var listCheckpointCmd = &cobra.Command{
	Use: "checkpoint",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.PaginateCheckpoint(
			service.NewMilestone(query.Q),
			mustGetFlag(cmd.Flags().GetInt("page")),
			mustGetFlag(cmd.Flags().GetInt("limit")),
			service.CheckpointQuery{
				MilestoneId: mustGetFlag(cmd.Flags().GetUint("milestone")),
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listCheckpointCmd)
	listCmd.Flags().IntP("page", "p", 1, "page of data")
	listCheckpointCmd.Flags().IntP("page", "p", 1, "page of data")
	listCmd.Flags().IntP("limit", "l", 10, "limit per page")
	listCheckpointCmd.Flags().IntP("limit", "l", 10, "limit per page")
	listCmd.Flags().BoolP("checkpoint", "c", false, "show all checkpoints of milestone")
	listCheckpointCmd.Flags().UintP("milestone", "m", 0, "filter by milestone id")
}
