package cli

import (
	"time"

	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/service"
)

const DefaultTimeLayout = tk.TimeLayoutFormatMinute

func CreateMilestone(cmd *cobra.Command, svc service.IMilestone) error {
	var deadline *time.Time
	d := tk.Must(cmd.Flags().GetString("deadline"))
	if d != "" {
		deadline = tk.Pointer(tk.Must(tk.ParseTime(d)))
	}
	input := service.MilestoneCreate{
		Point:     tk.Must(cmd.Flags().GetInt("point")),
		Title:     tk.Must(cmd.Flags().GetString("title")),
		StartedAt: tk.Must(tk.ParseTime(tk.Must(cmd.Flags().GetString("start")))),
		Deadline:  deadline,
	}

	_, err := svc.Create(cmd.Context(), input)
	return err
}
