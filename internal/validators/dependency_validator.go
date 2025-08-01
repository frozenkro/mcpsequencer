package validators

import (
	"github.com/frozenkro/mcpsequencer/internal/models"
)

type DependencyValidator struct {
	minimalTasks map[int]*minimalTask
}

type minimalTask struct {
	sort     int
	deps     []int
	complete bool
}

func (v DependencyValidator) addDep(dep models.Dependency) {
	if _, ok := v.minimalTasks[dep.Id]; !ok {
		v.minimalTasks[dep.Id] = &minimalTask{
			sort:     dep.Id,
			complete: false,
		}
	}

	task := v.minimalTasks[dep.Id]
	task.deps = append(task.deps, dep.DependsOn)
}

func (v DependencyValidator) Validate(deps []models.Dependency) error {
	for _, d := range deps {
		v.addDep(d)
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
				if _, ok := v.minimalTasks[d]; !ok {
					return InvalidDependencyError{SortID: d}
				}

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
