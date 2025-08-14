package validators_test

import (
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
	"github.com/frozenkro/mcpsequencer/internal/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidateTasksArray(t *testing.T) {
	type tasksTestCase struct {
		name        string
		tasks       []projectsdb.Task
		succ        bool
		expectedErr error
	}

	tests := []tasksTestCase{
		tasksTestCase{
			name: "ValidateTasksArray-NoDuplicates",
			tasks: []projectsdb.Task{
				projectsdb.Task{
					Sort: int64(0),
				},
				projectsdb.Task{
					Sort: int64(1),
				},
			},
			succ: true,
		},
		tasksTestCase{
			name: "ValidateTasksArray-WithDuplicates",
			tasks: []projectsdb.Task{
				projectsdb.Task{
					Sort: int64(0),
				},
				projectsdb.Task{
					Sort: int64(0),
				},
			},
			succ:        false,
			expectedErr: validators.DupeSortIdError{},
		},
	}

	sut := validators.TaskArrayValidator{}

	for _, test := range tests {
		succ := t.Run(test.name, func(t *testing.T) {
			err := sut.Validate(test.tasks)

			if test.succ {
				assert.Nil(t, err)
			} else {
				assert.IsType(t, test.expectedErr, err)
			}
		})
		if !succ {
			t.Fatalf("Test '%v' failed", test.name)
		}
	}
}
