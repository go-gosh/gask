package cli

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/go-gosh/gask/app/service"
)

func PaginateMilestone(svc *service.Milestone, page, limit int) error {
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
		content := "x"
		if len(datum.Checkpoints) > 0 {
			c := datum.Checkpoints[0]
			if c.CheckedAt != nil {
				content = c.CheckedAt.Format(DefaultTimeLayout)
			}
			content += ":" + c.Content
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
