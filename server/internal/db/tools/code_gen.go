package main

import (
	"github.com/mats0319/unnamed_plan/server/internal/db/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../dal/",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
		ModelPkgPath: "dal",
	})

	g.ApplyBasic(
		model.User{},
		model.CloudFile{},
	)

	g.Execute()
}
