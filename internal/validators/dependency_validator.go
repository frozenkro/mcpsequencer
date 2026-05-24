package validators

import (
	"github.com/frozenkro/mcpsequencer/internal/models"
)

type DependencyValidator struct {
	minimalTasks map[int]*minimalTask
}

type minimalTask struct {
	entId    int
	deps     []int
	complete bool
}

func (v DependencyValidator) Validate(deps []models.Dependency, allIds []int) error {
	v.minimalTasks = make(map[int]*minimalTask)
	for _, i := range allIds {
		v.minimalTasks[i] = &minimalTask{
			entId:    i,
			complete: false,
		}
	}

	for _, d := range deps {
		if _, ok := v.minimalTasks[d.Id]; !ok {
			return InvalidDependencyError{SortID: d.Id}
		}
		if _, ok := v.minimalTasks[d.DependsOn]; !ok {
			return InvalidDependencyError{SortID: d.DependsOn}
		}
		task := v.minimalTasks[d.Id]
		task.deps = append(task.deps, d.DependsOn)
	}

	for {
		// Was any task able to be completed in this iteration?
		anyTaskCompleted := false

		// Are there any tasks remaining to complete
		allTasksComplete := true

		for _, t := range v.minimalTasks {

			if t.complete {
				continue
			}
			allTasksComplete = false

			// Is this task still waiting on a dependent task?
			taskLocked := false

			for _, d := range t.deps {
				if !v.minimalTasks[d].complete {
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

			for _, t := range v.minimalTasks {
				if t.complete {
					completedIds = append(completedIds, t.entId)
				} else {
					unreachableIds = append(unreachableIds, t.entId)
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
