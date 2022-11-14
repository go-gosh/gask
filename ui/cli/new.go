package cli

import (
	"errors"
	"strconv"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-gosh/gask/app/milestone"
)

const DefaultTimeLayout = "2006-01-02"

var newMilestoneQestions = []*survey.Question{
	{
		Name:      "title",
		Prompt:    &survey.Input{Message: "Title"},
		Validate:  survey.Required,
		Transform: survey.Title,
	},
	{
		Name:   "point",
		Prompt: &survey.Input{Message: "Point", Default: "100"},
		Validate: func(ans interface{}) error {
			s := ans.(string)
			i, err := strconv.Atoi(s)
			if err != nil {
				return err
			}
			if i <= 0 {
				return errors.New("point must greater than zero")
			}
			return nil
		},
		Transform: func(ans interface{}) (newAns interface{}) {
			s := ans.(string)
			i, _ := strconv.Atoi(s)
			return i
		},
	},
	{
		Name:   "content",
		Prompt: &survey.Input{Message: "Content"},
	},
	{
		Name:   "startedAt",
		Prompt: &survey.Input{Message: "Started At", Default: time.Now().Format(DefaultTimeLayout)},
		Validate: func(ans interface{}) error {
			s := ans.(string)
			_, err := time.Parse(DefaultTimeLayout, s)
			return err
		},
		Transform: func(ans interface{}) (newAns interface{}) {
			s := ans.(string)
			t, _ := time.Parse(DefaultTimeLayout, s)
			return t
		},
	},
}

func NewMilestone(svc *milestone.Service) error {
	// input
	input := milestone.Create{}
	err := survey.Ask(newMilestoneQestions, &input)
	if err != nil {
		return err
	}
	// create new
	_, err = svc.CreateMilestone(input)
	if err != nil {
		return err
	}
	return nil
}
