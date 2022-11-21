package main

import (
	"gorm.io/gen"

	"github.com/go-gosh/gask/app/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./app/query",
		ModelPkgPath: "query",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.ApplyBasic(model.Checkpoint{}, model.Milestone{})
	g.Execute()
}
