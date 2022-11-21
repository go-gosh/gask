package cli

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/service"
)

func checkpointToString(c *model.Checkpoint) string {
	if c == nil {
		return ""
	}
	m := "x"
	t := c.JoinedAt.Format(DefaultTimeLayout)
	if c.CheckedAt != nil {
		m = "v"
		t = c.CheckedAt.Format(DefaultTimeLayout)
	}
	return fmt.Sprintf("%s%v:%s", m, c.ID, t)
}
func PaginateMilestone(svc *service.Milestone, page, limit int, check bool) error {
	data, count, err := svc.Paginate(page, limit)
	if err != nil {
		return err
	}
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"#", "title", "point", "progress", "content", "started at", "deadline", "created at"}, table.RowConfig{AutoMerge: true})
	for _, datum := range data {
		deadline := "-"
		if datum.Deadline != nil {
			deadline = datum.Deadline.Format(DefaultTimeLayout)
		}
		content := ""
		for _, checkpoint := range datum.Checkpoints {
			if content != "" {
				content += "\n"
			}
			content += checkpointToString(checkpoint)
			if !check {
				break
			}
		}
		if content == "" {
			content = "-"
		}
		writer.AppendRow(table.Row{datum.ID, datum.Title, datum.Point, datum.Progress, content, datum.StartedAt.Format(DefaultTimeLayout), deadline, datum.CreatedAt.Format(DefaultTimeLayout)})
		writer.AppendSeparator()
	}
	writer.AppendFooter(table.Row{"", "", "", "", "total page", (int(count) + limit - 1) / limit, "current", page})
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	writer.SetStyle(table.StyleLight)
	writer.SetOutputMirror(os.Stdout)
	writer.Render()
	return nil
}
