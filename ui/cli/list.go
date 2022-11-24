package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/service"
)

func PaginateCheckpoint(cmd *cobra.Command, svc service.ICheckpoint, q *service.CheckpointQuery) error {
	t := time.Now()
	q.Timestamp = tk.Pointer(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
	data, err := svc.FindByPage(cmd.Context(), q)
	if err != nil {
		return err
	}
	writer := table.NewWriter()
	writer.AppendHeader(table.Row{"#", "title", "point", "content", "joined at", "checked at", "created at"}, table.RowConfig{AutoMerge: true})
	for _, datum := range data.Data {
		checkedAt := "-"
		if datum.CheckedAt != nil {
			checkedAt = datum.CheckedAt.Format(DefaultTimeLayout)
		}
		writer.AppendRow(table.Row{fmt.Sprintf("%v-%v", datum.MilestoneId, datum.ID), datum.Milestone.Title, datum.Point, datum.Content, datum.JoinedAt.Format(DefaultTimeLayout), checkedAt, datum.CreatedAt.Format(DefaultTimeLayout)})
		writer.AppendSeparator()
	}
	writer.AppendFooter(table.Row{"", "", "", "total page", (data.Total-1)/data.PageSize + 1, "current", data.Page})
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
