package cmd

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/go-gosh/gask/app/query"
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
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			svc := service.NewMilestone(query.Q)
			entity, err := svc.RetrieveById(uint(id))
			if err != nil {
				return err
			}
			fmt.Printf("Will delete milestone:\nid:%v\ttitle:%v\n\n", entity.ID, entity.Title)
			if !quiet {
				var confirm bool
				err := survey.AskOne(&survey.Confirm{Message: "delete"}, &confirm)
				if err != nil {
					return err
				}
				if !confirm {
					return nil
				}
			}
			return svc.DeleteById(uint(id))
		},
	}
)

func init() {
	deleteCmd.AddCommand(deleteMilestoneCmd)
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "delete quietly")
}
