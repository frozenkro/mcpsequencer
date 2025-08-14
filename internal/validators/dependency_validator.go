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

func (v DependencyValidator) initMinimalTasks(allIds []int) {
	v.minimalTasks = map[int]*minimalTask{}

	for _, i := range allIds {
		v.minimalTasks[i] = &minimalTask{
			entId:    i,
			complete: false,
		}
	}
}

func (v DependencyValidator) addDep(dep models.Dependency) error {
	if _, ok := v.minimalTasks[dep.Id]; !ok {
		return InvalidDependencyError{SortID: dep.Id}
	}

	task := v.minimalTasks[dep.Id]
	task.deps = append(task.deps, dep.DependsOn)
	return nil
}

// Will return err if dependency deadlock is possible, or if ids are missing.
//
// deps: list of dependencies; associations between tasks that need to be performed sequentially
//
// allIds: represents all sortIds or taskIds (depending on Task.Discriminator) for the task list being validated against
func (v DependencyValidator) Validate(deps []models.Dependency, allIds []int) error {
	v.initMinimalTasks(allIds)

	for _, d := range deps {
		if err := v.addDep(d); err != nil {
			return err
		}
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
