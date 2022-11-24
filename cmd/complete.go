package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/service"
	"github.com/go-gosh/gask/ui/cli"
)

func checkArgsId(name string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires %s id in first args, only received 0", name)
		}
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if id <= 0 {
			return fmt.Errorf("%s id must greater than 0", name)
		}
		return nil
	}
}

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete <checkpoint id>",
	Short: "Complete an unchecked checkpoint of milestone",
	Args:  checkArgsId("checkpoint"),
	RunE: func(cmd *cobra.Command, args []string) error {
		svc := service.NewCheckpoint(tk.Must(global.GetDatabase()))
		return svc.UpdateById(
			cmd.Context(),
			uint(tk.Must(strconv.Atoi(args[0]))),
			&service.CheckpointUpdate{
				IsChecked: tk.Pointer(true),
				Point:     nil,
				Content:   nil,
				CheckedAt: tk.Must(tk.ParseTimePointer(tk.Must(cmd.Flags().GetString("time")))),
				JoinedAt:  nil,
			},
		)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	completeCmd.Flags().StringP("time", "t", time.Now().Format(cli.DefaultTimeLayout), "completed time")
}
