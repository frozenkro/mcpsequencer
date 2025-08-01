package transformers

import (
	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type DependencyTransformer struct{}

func (t DependencyTransformer) FromInts(depTaskIds []int, taskId int, disc models.DependencyDiscriminator) []models.Dependency {
	result := []models.Dependency{}
	for _, d := range depTaskIds {
		result = append(result, models.Dependency{
			Id:            taskId,
			DependsOn:     d,
			Discriminator: disc,
		})
	}
	return result
}

func (t DependencyTransformer) FromDbRows(rows []projectsdb.Dependency) []models.Dependency {
	result := []models.Dependency{}
	for _, r := range rows {
		result = append(result, models.NewDependencyWithTaskIds(int(r.TaskID), int(r.DependsOn)))
	}
	return result
}
