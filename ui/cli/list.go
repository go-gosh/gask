package cli

import (
	"os"

	"github.com/go-gosh/gask/app/service"
	"github.com/jedib0t/go-pretty/v6/table"
)

func PaginateMilestone(svc *service.Milestone, page, limit int) error {
	data, count, err := svc.Paginate(page, limit)
	if err != nil {
		return err
	}
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"#", "title", "point", "progress", "content", "started at", "deadline", "created at"})
	for _, datum := range data {
		deadline := "-"
		if datum.Deadline != nil {
			deadline = datum.Deadline.Format(DefaultTimeLayout)
		}
		writer.AppendRow(table.Row{datum.ID, datum.Title, datum.Point, datum.Progress, datum.Content, datum.StartedAt.Format(DefaultTimeLayout), deadline, datum.CreatedAt.Format(DefaultTimeLayout)})
	}
	writer.AppendFooter(table.Row{"", "", "", "", "total page", (int(count) + limit - 1) / limit, "current", page})
	writer.SetStyle(table.StyleColoredBright)
	writer.SetOutputMirror(os.Stdout)
	writer.Render()
	return nil
}
