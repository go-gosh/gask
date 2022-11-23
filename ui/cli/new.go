package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/service"
)

const DefaultTimeLayout = tk.TimeLayoutFormatMinute

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

func CreateMilestone(cmd *cobra.Command, svc service.IMilestone) error {
	var deadline *time.Time
	d := tk.Must(cmd.Flags().GetString("deadline"))
	if d != "" {
		deadline = tk.Pointer(tk.Must(tk.ParseTime(d)))
	}
	input := service.MilestoneCreate{
		Point:     tk.Must(cmd.Flags().GetInt("point")),
		Title:     tk.Must(cmd.Flags().GetString("title")),
		StartedAt: tk.Must(tk.ParseTime(tk.Must(cmd.Flags().GetString("start")))),
		Deadline:  deadline,
	}

	_, err := svc.Create(cmd.Context(), input)
	return err
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
