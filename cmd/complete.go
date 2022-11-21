package cmd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/go-gosh/gask/app/query"
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
		svc := service.NewMilestone(query.Q)
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		timeStr, err := cmd.Flags().GetString("time")
		if err != nil {
			return err
		}
		timestamp, err := time.Parse(cli.DefaultTimeLayout, timeStr)
		if err != nil {
			return err
		}
		err = svc.CompleteCheckpointById(uint(id), timestamp)
		if err != nil {
			return err
		}
		log.Printf("Completed checkpoint<%v> at %s", id, timestamp.Format(cli.DefaultTimeLayout))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)

	completeCmd.Flags().StringP("time", "t", time.Now().Format(cli.DefaultTimeLayout), "completed time")
}
