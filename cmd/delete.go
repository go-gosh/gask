package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/service"
)

var quiet bool

var (
	// deleteCmd represents the delete command
	deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resources.",
	}
	// deleteMilestoneCmd represents the delete milestone command
	deleteMilestoneCmd = &cobra.Command{
		Use:   "milestone <id>",
		Short: "Delete milestone by id",
		Args:  checkArgsId("milestone"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := uint(tk.Must(strconv.Atoi(args[0])))
			svc := service.NewMilestoneV2(tk.Must(global.GetDatabase()))
			entity := tk.Must(svc.OneById(cmd.Context(), id))
			fmt.Printf("Will delete milestone:\nid:%v\ttitle:%v\n\n", entity.ID, entity.Title)
			if !quiet {
				var confirm string
				fmt.Printf("Confirm delete milestone? (y/N)")
				_, _ = fmt.Scanln(&confirm)
				if strings.TrimSpace(confirm) != "y" {
					return nil
				}
			}
			return svc.DeleteById(cmd.Context(), id)
		},
	}

	// deleteCheckpointCmd represents the delete checkpoint command
	deleteCheckpointCmd = &cobra.Command{
		Use:   "checkpoint <id>",
		Short: "Delete checkpoint by id",
		Args:  checkArgsId("checkpoint"),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := uint(tk.Must(strconv.Atoi(args[0])))
			svc := service.NewCheckpoint(tk.Must(global.GetDatabase()))
			entity := tk.Must(svc.OneById(cmd.Context(), id))
			fmt.Printf("Will delete checkpoint:\nid:%v\tcontent:%v\n\n", entity.ID, entity.Content)
			if !quiet {
				var confirm string
				fmt.Printf("Confirm delete milestone? (y/N)")
				_, _ = fmt.Scanln(&confirm)
				if strings.TrimSpace(confirm) != "y" {
					return nil
				}
			}
			return svc.DeleteById(cmd.Context(), id)
		},
	}
)

func init() {
	deleteCmd.AddCommand(deleteMilestoneCmd)
	deleteCmd.AddCommand(deleteCheckpointCmd)
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "delete quietly")
}
