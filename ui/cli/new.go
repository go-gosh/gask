package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/go-gosh/gask/app/service"
)

const DefaultTimeLayout = "2006-01-02 15:04:05"

func timeValidate(ans interface{}) error {
	s := ans.(string)
	_, err := time.Parse(DefaultTimeLayout, s)
	return err
}

func timeTransform(ans interface{}) (newAns interface{}) {
	s := ans.(string)
	t, _ := time.Parse(DefaultTimeLayout, s)
	return t
}

func NewMilestone(svc *service.Milestone, cmd *cobra.Command) error {
	// input
	input := service.Create{
		Point:     getIntFromFlags(cmd.Flags(), "point"),
		Title:     getStringFromFlags(cmd.Flags(), "title"),
		StartedAt: getTimeFromFlags(cmd.Flags(), "start"),
		Deadline:  getTimePFromFlags(cmd.Flags(), "deadline"),
	}

	// create new
	_, err := svc.CreateMilestone(input)
	if err != nil {
		return err
	}
	return nil
}

func getTimePFromFlags(flags *pflag.FlagSet, name string) *time.Time {
	s := getStringFromFlags(flags, name)
	if s == "" {
		return nil
	}
	t, err := time.Parse(DefaultTimeLayout, s)
	cobra.CheckErr(err)
	return &t
}

func getTimeFromFlags(flags *pflag.FlagSet, name string) time.Time {
	s := getTimePFromFlags(flags, name)
	if s == nil {
		cobra.CheckErr(fmt.Sprintf("flag %s is requried", name))
	}
	return *s
}

func getIntFromFlags(flags *pflag.FlagSet, name string) int {
	s, err := flags.GetInt(name)
	cobra.CheckErr(err)
	return s
}

func getStringFromFlags(flags *pflag.FlagSet, name string) string {
	s, err := flags.GetString(name)
	cobra.CheckErr(err)
	return s
}
