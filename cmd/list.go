package cmd

import (
	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/query"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all milestones",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.PaginateMilestone(
			cmd.Context(),
			service.NewMilestoneV2(tk.Must(global.GetDatabase())),
			&service.MilestoneQuery{
				Pager: repo.Pager{
					Page:     tk.Must(cmd.Flags().GetInt("page")),
					PageSize: tk.Must(cmd.Flags().GetInt("limit")),
				},
			},
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
			tk.Must(cmd.Flags().GetInt("page")),
			tk.Must(cmd.Flags().GetInt("limit")),
			service.CheckpointQuery{
				MilestoneId: tk.Must(cmd.Flags().GetUint("milestone")),
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
