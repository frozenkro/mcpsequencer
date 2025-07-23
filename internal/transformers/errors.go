package transformers

import (
	"fmt"
	"strings"
)

type TasksUnmarshalError struct {
	TasksJson string
	Err       error
}

func (e TasksUnmarshalError) Error() string {
	return fmt.Sprintf("Error unmarshaling tasks '%v'\n%v", e.TasksJson, e.Err.Error())
}

type DepsMarshalError struct {
	Deps []int
	Err  error
}

func (e DepsMarshalError) Error() string {
	depsStr := strings.Join(strings.Fields(fmt.Sprint(e.Deps)), ",")
	return fmt.Sprintf("Error marshaling dependency array '%v'\n%v", depsStr, e.Err.Error())
}
