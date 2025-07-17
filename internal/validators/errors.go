package validators

import (
	"fmt"
	"strings"
)

type DepsUnmarshalError struct {
	DepsJson string
	Err      error
}

func (e DepsUnmarshalError) Error() string {
	return fmt.Sprintf("Error unmarshaling dependency array '%v'\n%v", e.DepsJson, e.Err.Error())
}

type DupeSortIdError struct {
	SortID int
}

func (e DupeSortIdError) Error() string {
	return fmt.Sprintf("Found duplicate sort ID %v\n", e.SortID)
}

type InvalidDependencyError struct {
	SortID int
}

func (e InvalidDependencyError) Error() string {
	return fmt.Sprintf("Dependency '%v' not a valid Sort ID in this list", e.SortID)
}

type DependencyTreeParseError struct {
	CompletedIds   []int
	UnreachableIds []int
}

func (e DependencyTreeParseError) Error() string {
	completedStr := strings.Join(strings.Fields(fmt.Sprint(e.CompletedIds)), ",")
	unreachableStr := strings.Join(strings.Fields(fmt.Sprint(e.UnreachableIds)), ",")

	return fmt.Sprintf(`
	Dependency tree cannot be walked, there is a cyclical dependency or deadlock.
	Completed Sort IDs: %v
	Unreachable Sort IDs: %v
	`, completedStr, unreachableStr)
}
