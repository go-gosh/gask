package cli

import (
	"fmt"
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
		writer.AppendRow(table.Row{datum.ID, datum.Title, datum.Point, datum.Progress, "-", datum.StartedAt.Format(DefaultTimeLayout), deadline, datum.CreatedAt.Format(DefaultTimeLayout)})
		for _, c := range datum.Checkpoints {
			checked := "-"
			if c.CheckedAt != nil {
				checked = c.CheckedAt.Format(DefaultTimeLayout)
			}
			writer.AppendRow(table.Row{datum.ID, fmt.Sprintf("  ->%v", c.ID), c.Point, c.CheckedAt != nil, c.Content, c.JoinedAt.Format(DefaultTimeLayout), checked, c.CreatedAt.Format(DefaultTimeLayout)})
		}
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
