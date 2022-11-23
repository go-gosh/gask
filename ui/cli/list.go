package cli

import (
	"context"
	"fmt"
	"os"
	"time"

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
	return fmt.Sprintf("%s%v->%s:%s", m, c.ID, t, c.Content)
}

func PaginateCheckpoint(svc *service.Milestone, page, limit int, q service.CheckpointQuery) error {
	t := time.Now()
	q.Timestamp = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	data, count, err := svc.PaginateCheckpoints(page, limit, q)
	if err != nil {
		return err
	}
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"#", "title", "point", "content", "joined at", "checked at", "created at"}, table.RowConfig{AutoMerge: true})
	for _, datum := range data {
		checkedAt := "-"
		if datum.CheckedAt != nil {
			checkedAt = datum.CheckedAt.Format(DefaultTimeLayout)
		}
		writer.AppendRow(table.Row{fmt.Sprintf("%v-%v", datum.MilestoneId, datum.ID), datum.Milestone.Title, datum.Point, datum.Content, datum.JoinedAt.Format(DefaultTimeLayout), checkedAt, datum.CreatedAt.Format(DefaultTimeLayout)})
		writer.AppendSeparator()
	}
	writer.AppendFooter(table.Row{"", "", "", "total page", (int(count) + limit - 1) / limit, "current", page})
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	writer.SetStyle(table.StyleLight)
	writer.SetOutputMirror(os.Stdout)
	writer.Render()
	return nil
}

func PaginateMilestone(ctx context.Context, svc service.IMilestone, q *service.MilestoneQuery) error {
	result, err := svc.FindByPage(ctx, q)
	if err != nil {
		return err
	}
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"#", "title", "point", "progress", "content", "started at", "deadline", "created at"}, table.RowConfig{AutoMerge: true})
	for _, datum := range result.Data {
		deadline := "-"
		if datum.Deadline != nil {
			deadline = datum.Deadline.Format(DefaultTimeLayout)
		}
		writer.AppendRow(table.Row{datum.ID, datum.Title, datum.Point, datum.Progress, "-", datum.StartedAt.Format(DefaultTimeLayout), deadline, datum.CreatedAt.Format(DefaultTimeLayout)})
		writer.AppendSeparator()
	}
	writer.AppendFooter(table.Row{"", "", "", "", "total page", (result.Total-1)/result.PageSize + 1, "current", result.Page})
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	writer.SetStyle(table.StyleLight)
	writer.SetOutputMirror(os.Stdout)
	writer.Render()
	return nil
}
