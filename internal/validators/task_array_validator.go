package validators

import (
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
)

type TaskArrayValidator struct{}

func (v TaskArrayValidator) Validate(tasks []projectsdb.Task) error {
	sortIds := make(map[int64]bool)

	for _, t := range tasks {
		if _, ok := sortIds[t.Sort]; ok {
			return DupeSortIdError{SortID: t.Sort}
		}

		sortIds[t.Sort] = true
	}

	return nil
}
