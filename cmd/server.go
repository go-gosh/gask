package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/go-gosh/gask/api"
	"github.com/go-gosh/gask/app/conf"
	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/service"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run gask web server",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := global.GetDatabase()
		if err != nil {
			return err
		}
		m := service.NewMilestone(db)
		c := service.NewCheckpoint(db)
		return api.New(m, c).Run(fmt.Sprintf(":%d", conf.GetConfig().Port))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
