package transformers

import (
	"encoding/json"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type TaskArrayTransformer struct{}

var depTrn DependencyTransformer = DependencyTransformer{}

func (t TaskArrayTransformer) ParseFromJson(tasks string, projectId int) ([]projectsdb.Task, []models.Dependency, error) {
	argsSl := []models.CreateTaskArgs{}
	if err := json.Unmarshal([]byte(tasks), &argsSl); err != nil {
		return nil, nil, TasksUnmarshalError{}
	}

	result := []projectsdb.Task{}
	deps := []models.Dependency{}
	for _, taskArgs := range argsSl {

		task := projectsdb.Task{
			ProjectID:   int64(projectId),
			Name:        taskArgs.Name,
			Description: taskArgs.Description,
			Sort:        int64(taskArgs.SortId),
		}

		for _, d := range depTrn.FromInts(taskArgs.Dependencies, taskArgs.SortId, models.SortId) {
			deps = append(deps, d)
		}

		result = append(result, task)
	}
	return result, deps, nil
}

func (t TaskArrayTransformer) TaskIdMapFromTasks(tasks []projectsdb.Task) models.SortIdTaskIdMap {
	result := models.SortIdTaskIdMap{}
	for _, t := range tasks {
		result[t.Sort] = t.TaskID
	}
	return result
}
