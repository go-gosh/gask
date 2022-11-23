package cli

import (
	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/service"
)

func CheckMilestone(cmd *cobra.Command, svc service.ICheckpoint, id uint) error {
	_, err := svc.Create(cmd.Context(), &service.CheckpointCreate{
		Point:       tk.Must(cmd.Flags().GetInt("point")),
		Content:     tk.Must(cmd.Flags().GetString("content")),
		JoinedAt:    tk.Must(tk.ParseTime(tk.Must(cmd.Flags().GetString("joined")))),
		CheckedAt:   tk.Must(tk.ParseTimePointer(tk.Must(cmd.Flags().GetString("checked")))),
		MilestoneId: id,
	})
	return err
}
