package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type TaskUnmarshalError struct {
	TaskJson string
	Err      error
}

func (e *TaskUnmarshalError) Error() string {
	return fmt.Sprintf("Error unmarshaling task '%v'\n%v", e.TaskJson, e.Err.Error())
}

type DepsMarshalError struct {
	Deps []int
	Err  error
}

func (e *DepsMarshalError) Error() string {
	depsStr := strings.Join(strings.Fields(fmt.Sprint(e.Deps)), ",")
	return fmt.Sprintf("Error marshaling dependency array '%v'\n%v", depsStr, e.Err.Error())
}

type DepsUnmarshalError struct {
	DepsJson string
	Err      error
}

func (e *DepsUnmarshalError) Error() string {
	return fmt.Sprintf("Error unmarshaling dependency array '%v'\n%v", e.DepsJson, e.Err.Error())
}

type DupeSortIdError struct {
	SortID int
}

func (e *DupeSortIdError) Error() string {
	return fmt.Sprintf("Found duplicate sort ID %v\n", e.SortID)
}

type InvalidDependencyError struct {
	SortID int
}

func (e *InvalidDependencyError) Error() string {
	return fmt.Sprintf("Dependency '%v' not a valid Sort ID in this list", e.SortID)
}

type DependencyTreeParseError struct {
	CompletedIds   []int
	UnreachableIds []int
}

func (e *DependencyTreeParseError) Error() string {
	completedStr := strings.Join(strings.Fields(fmt.Sprint(e.CompletedIds)), ",")
	unreachableStr := strings.Join(strings.Fields(fmt.Sprint(e.UnreachableIds)), ",")

	return fmt.Sprintf(`
	Dependency tree cannot be walked, there is a cyclical dependency or deadlock.
	Completed Sort IDs: %v
	Unreachable Sort IDs: %v
	`, completedStr, unreachableStr)
}

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
			return nil, &TaskUnmarshalError{TaskJson: jsonT, Err: err}
		}

		jsonDeps, err := json.Marshal(args.Dependencies)
		if err != nil {
			return nil, &DepsMarshalError{Deps: args.Dependencies, Err: err}
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
		return nil, &DepsUnmarshalError{DepsJson: t.DependenciesJson, Err: err}
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
			return &DupeSortIdError{SortID: mt.sort}
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
					return &InvalidDependencyError{SortID: d}
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
			return &DependencyTreeParseError{CompletedIds: completedIds, UnreachableIds: unreachableIds}
		}

		if allTasksComplete {
			break
		}
	}

	return nil
}
