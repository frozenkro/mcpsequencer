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
			name: "ValidateTasksArray-NoDeps",
			tasks: []projectsdb.Task{
				projectsdb.Task{
					Sort:             int64(0),
					DependenciesJson: "[]",
				},
				projectsdb.Task{
					Sort:             int64(1),
					DependenciesJson: "[]",
				},
			},
			succ: true,
		},
		tasksTestCase{
			name: "ValidateTasksArray-SimpleSequence",
			tasks: []projectsdb.Task{
				projectsdb.Task{
					Sort:             int64(0),
					DependenciesJson: "[]",
				},
				projectsdb.Task{
					Sort:             int64(1),
					DependenciesJson: "[0]",
				},
				projectsdb.Task{
					Sort:             int64(2),
					DependenciesJson: "[1]",
				},
			},
			succ: true,
		},
		tasksTestCase{
			name: "ValidateTasksArray-Cyclical",
			tasks: []projectsdb.Task{
				projectsdb.Task{
					Sort:             int64(0),
					DependenciesJson: "[1]",
				},
				projectsdb.Task{
					Sort:             int64(1),
					DependenciesJson: "[0]",
				},
			},
			succ:        false,
			expectedErr: validators.DependencyTreeParseError{},
		},
	}

	for _, test := range tests {
		succ := t.Run(test.name, func(t *testing.T) {
			err := validators.ValidateTasksArray(test.tasks)

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
