package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

func IsDev() bool {
	for _, v := range os.Args {
		if v == "--dev" {
			return true
		}
	}
	return false
}

func ParseTasksArray(tasks []string, projectId int) ([]projectsdb.Task, error) {
	result := []projectsdb.Task{}
	for _, jsonT := range tasks {
		args := models.CreateTaskArgs{}

		if err := json.Unmarshal([]byte(jsonT), args); err != nil {
			return nil, fmt.Errorf("Error unmarshaling task '%v'\n%v", jsonT, err.Error())
		}

		jsonDeps, err := json.Marshal(args.Dependencies)
		if err != nil {
			return nil, fmt.Errorf("Error marshaling dependency array '%v'\n%v", args.Dependencies, err.Error())
		}

		task := projectsdb.Task{
			ProjectID:        int64(projectId),
			Name:             args.Name,
			Description:      args.Description,
			Sort:             int64(args.SortId),
			DependenciesJson: string(jsonDeps),
		}

		result = append(result, task)
	}
	return result, nil
}

type minimalTask struct {
	sort     int
	deps     []int
	complete bool
}

func newMinimalTask(t projectsdb.Task) (*minimalTask, error) {
	depsSl := []int{}

	if err := json.Unmarshal([]byte(t.DependenciesJson), depsSl); err != nil {
		return nil, fmt.Errorf("Error unmarshaling dependency list: %v\n%v\n", t.DependenciesJson, err.Error())
	}

	return &minimalTask{
		sort:     int(t.Sort),
		deps:     depsSl,
		complete: false,
	}, nil
}

func ValidateTasksArray(tasks []projectsdb.Task) error {
	minimalTasks := make(map[int]*minimalTask)

	for _, t := range tasks {
		mt, err := newMinimalTask(t)
		if err != nil {
			return err
		}

		if _, ok := minimalTasks[mt.sort]; ok {
			return fmt.Errorf("Found duplicate sort ID %v\n", mt.sort)
		}

		minimalTasks[mt.sort] = mt
	}

	for {
		// Was any task able to be completed in this iteration?
		anyTaskCompleted := false

		// Are there any tasks remaining to complete
		allTasksComplete := true

		for _, t := range minimalTasks {

			if t.complete {
				continue
			}
			allTasksComplete = false

			// Is this task still waiting on a dependent task?
			taskLocked := false

			for _, d := range t.deps {
				if _, ok := minimalTasks[d]; !ok {
					return fmt.Errorf("Dependency '%v' not a valid Sort ID in this list", d)
				}

				if !minimalTasks[d].complete {
					taskLocked = true
				}
			}

			if !taskLocked {
				t.complete = true
				anyTaskCompleted = true
			}
		}

		if !anyTaskCompleted && allTasksComplete {

			completedIds := []int{}
			unreachableIds := []int{}

			for _, t := range minimalTasks {
				if t.complete {
					completedIds = append(completedIds, t.sort)
				} else {
					unreachableIds = append(unreachableIds, t.sort)
				}
			}

			completedStr := strings.Join(strings.Fields(fmt.Sprint(completedIds)), ",")
			unreachableStr := strings.Join(strings.Fields(fmt.Sprint(unreachableIds)), ",")

			return fmt.Errorf(`
			Dependency tree cannot be walked, there is a cyclical dependency or deadlock.
			Completed Sort IDs: %v
			Unreachable Sort IDs: %v
			`, completedStr, unreachableStr)
		}

		if allTasksComplete {
			break
		}
	}

	return nil
}
