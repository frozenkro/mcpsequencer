package utils_test

import (
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
	"github.com/frozenkro/mcpsequencer/internal/utils"
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
		tasksTestCase{},
		tasksTestCase{},
		tasksTestCase{},
	}

	for _, test := range tests {
		succ := t.Run(test.name, func(t *testing.T) {
			err := utils.ValidateTasksArray(test.tasks)

			if test.succ {
				assert.Nil(t, err)
			} else {
				assert.ErrorIs(t, test.expectedErr, err)
			}
		})
		if !succ {
			t.Fatalf("Test '%v' failed", test.name)
		}
	}
}
