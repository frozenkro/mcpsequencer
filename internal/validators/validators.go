package validators

import (
	"encoding/json"

	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type minimalTask struct {
	sort     int
	deps     []int
	complete bool
}

func newMinimalTask(t projectsdb.Task) (*minimalTask, error) {
	depsSl := &[]int{}

	if err := json.Unmarshal([]byte(t.DependenciesJson), depsSl); err != nil {
		return nil, DepsUnmarshalError{DepsJson: t.DependenciesJson, Err: err}
	}

	return &minimalTask{
		sort:     int(t.Sort),
		deps:     *depsSl,
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
			return DupeSortIdError{SortID: mt.sort}
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
					return InvalidDependencyError{SortID: d}
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

		if !anyTaskCompleted && !allTasksComplete {

			completedIds := []int{}
			unreachableIds := []int{}

			for _, t := range minimalTasks {
				if t.complete {
					completedIds = append(completedIds, t.sort)
				} else {
					unreachableIds = append(unreachableIds, t.sort)
				}
			}
			return DependencyTreeParseError{CompletedIds: completedIds, UnreachableIds: unreachableIds}
		}

		if allTasksComplete {
			break
		}
	}

	return nil
}
