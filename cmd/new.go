package cmd

import (
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/milestone"
	"github.com/go-gosh/gask/ui/cli"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
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
		svc := milestone.New(db)
		return cli.NewMilestone(svc)
	},
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(newCmd)
}
