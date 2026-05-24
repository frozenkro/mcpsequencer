package validators_test

import (
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/models"
	"github.com/frozenkro/mcpsequencer/internal/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidateDependencyArray(t *testing.T) {
	type testCase struct {
		name        string
		deps        []models.Dependency
		allIds      []int
		succ        bool
		expectedErr error
	}

	tests := []testCase{
		testCase{
			name:   "ValidateDependencyArray-NoDeps",
			deps:   []models.Dependency{},
			allIds: []int{0, 1, 2},
			succ:   true,
		},
		testCase{
			name: "ValidateDependencyArray-SimpleSequence",
			deps: []models.Dependency{
				models.Dependency{
					Id:        1,
					DependsOn: 0,
				},
				models.Dependency{
					Id:        2,
					DependsOn: 1,
				},
			},
			allIds: []int{0, 1, 2},
			succ:   true,
		},
		testCase{
			name: "ValidateDependencyArray-Cyclical",
			deps: []models.Dependency{
				models.Dependency{
					Id:        0,
					DependsOn: 1,
				},
				models.Dependency{
					Id:        1,
					DependsOn: 0,
				},
			},
			allIds:      []int{0, 1, 2},
			succ:        false,
			expectedErr: validators.DependencyTreeParseError{},
		},
		testCase{
			name: "ValidateDependencyArray-BadDepId",
			deps: []models.Dependency{
				models.Dependency{
					Id:        0,
					DependsOn: 3,
				},
				models.Dependency{
					Id:        1,
					DependsOn: 0,
				},
			},
			allIds:      []int{0, 1, 2},
			succ:        false,
			expectedErr: validators.InvalidDependencyError{},
		},
	}

	sut := validators.DependencyValidator{}

	for _, test := range tests {
		succ := t.Run(test.name, func(t *testing.T) {
			err := sut.Validate(test.deps, test.allIds)

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
