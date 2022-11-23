package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/service"
)

var (
	// updateCmd represents the update command
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resources.",
	}
	// updateMilestoneCmd represents the update milestone command
	updateMilestoneCmd = &cobra.Command{
		Use:   "milestone <id>",
		Short: "Update milestone by id",
		Args:  checkArgsId("milestone"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := uint(tk.Must(strconv.Atoi(args[0])))
			svc := service.NewMilestoneV2(tk.Must(global.GetDatabase()))
			entity := tk.Must(svc.OneById(cmd.Context(), id))
			fmt.Printf("Will update milestone:\nid:%v\ttitle:%v\n\n", entity.ID, entity.Title)
			return svc.UpdateById(cmd.Context(), id, &service.MilestoneUpdate{
				CleanDeadline: tk.Must(cmd.Flags().GetBool("forever")),
				Point:         tk.ZeroNilPointer(tk.Must(cmd.Flags().GetInt("point"))),
				Title:         tk.ZeroNilPointer(tk.Must(cmd.Flags().GetString("title"))),
				StartedAt:     tk.Must(tk.ParseTimePointer(tk.Must(cmd.Flags().GetString("start")))),
				Deadline:      tk.Must(tk.ParseTimePointer(tk.Must(cmd.Flags().GetString("deadline")))),
			})
		},
	}
)

func init() {
	updateCmd.AddCommand(updateMilestoneCmd)
	rootCmd.AddCommand(updateCmd)

	updateMilestoneCmd.Flags().StringP("title", "t", "", "title of milestone")
	updateMilestoneCmd.Flags().IntP("point", "p", 0, "point of milestone, must be positive")
	updateMilestoneCmd.Flags().StringP("start", "s", "", "milestone start at")
	updateMilestoneCmd.Flags().StringP("deadline", "d", "", "deadline of milestone")
	updateMilestoneCmd.Flags().BoolP("forever", "e", false, "deadline of milestone")
}
