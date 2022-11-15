package cli

import (
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-gosh/gask/app/milestone"
)

func CheckMilestone(svc *milestone.Service) error {
	var id uint
	err := survey.AskOne(&survey.Input{Message: "milestone id"}, &id, survey.WithValidator(survey.Required))
	if err != nil {
		return err
	}
	create := milestone.CheckpointCreate{}
	qs := []*survey.Question{
		{
			Name:   "Point",
			Prompt: &survey.Input{Message: "point", Default: "10"},
		},
		{
			Name:   "Content",
			Prompt: &survey.Input{Message: "content"},
		},
		{
			Name:      "JoinedAt",
			Prompt:    &survey.Input{Message: "joined at", Default: time.Now().Format(DefaultTimeLayout)},
			Validate:  timeValidate,
			Transform: timeTransform,
		},
	}
	err = survey.Ask(qs, &create)
	if err != nil {
		return err
	}
	var checked bool
	err = survey.AskOne(&survey.Confirm{Message: "check it"}, &checked)
	if err != nil {
		return err
	}
	if checked {
		err = survey.Ask([]*survey.Question{
			{
				Name:     "CheckedAt",
				Prompt:   &survey.Input{Message: "checked at", Default: time.Now().Format(DefaultTimeLayout)},
				Validate: timeValidate,
				Transform: func(ans interface{}) (newAns interface{}) {
					s := ans.(string)
					t, _ := time.Parse(DefaultTimeLayout, s)
					return &t
				},
			},
		}, &create.CheckedAt, survey.WithValidator(survey.Required))
		if err != nil {
			return err
		}
	}
	_, err = svc.SplitMilestoneById(id, create)
	if err != nil {
		return err
	}
	return nil
}
