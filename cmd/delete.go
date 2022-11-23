package cmd

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
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
				var confirm bool
				err := survey.AskOne(&survey.Confirm{Message: "delete"}, &confirm)
				if err != nil {
					return err
				}
				if !confirm {
					return nil
				}
			}
			return svc.DeleteById(cmd.Context(), id)
		},
	}
)

func init() {
	deleteCmd.AddCommand(deleteMilestoneCmd)
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "delete quietly")
}
