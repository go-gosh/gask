package cmd

import (
	"time"

	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new milestone",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.CreateMilestone(cmd, service.NewMilestone(tk.Must(global.GetDatabase())))
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("title", "t", "", "title (required)")
	cobra.CheckErr(newCmd.MarkFlagRequired("title"))
	newCmd.Flags().IntP("point", "p", 100, "point of milestone, must be positive")
	newCmd.Flags().StringP("start", "s", time.Now().Format(cli.DefaultTimeLayout), "milestone start at")
	newCmd.Flags().StringP("deadline", "d", "", "deadline of milestone (default nil or later than start)")
}
