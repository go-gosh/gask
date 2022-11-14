package cmd

import (
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/milestone"
	"github.com/go-gosh/gask/ui/cli"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := global.GetDatabase()
		if err != nil {
			return err
		}
		page, err := cmd.Flags().GetInt("page")
		if err != nil {
			return err
		}
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			return err
		}
		return cli.PaginateMilestone(milestone.New(db), page, limit)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().IntP("page", "p", 1, "Page of data")
	listCmd.Flags().IntP("limit", "l", 10, "Limit per page")
}
