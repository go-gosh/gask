package cmd

import (
	"strconv"
	"time"

	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check <milestone id>",
	Short: "Create a checkpoint for milestone",
	Args:  checkArgsId("milestone"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cli.CheckMilestone(cmd, service.NewCheckpoint(tk.Must(global.GetDatabase())), uint(tk.Must(strconv.Atoi(args[0]))))
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().IntP("point", "p", 10, "point of checkpoint")
	checkCmd.Flags().StringP("content", "c", "", "content of checkpoint")
	checkCmd.Flags().StringP("joined", "j", time.Now().Format(cli.DefaultTimeLayout), "joined at")
	checkCmd.Flags().StringP("checked", "o", "", "checked at")
}
