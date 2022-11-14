package main

import (
	"github.com/go-gosh/gask/app/model"
	"gorm.io/gen"
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
