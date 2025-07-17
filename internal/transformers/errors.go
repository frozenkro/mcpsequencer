package transformers

import (
	"fmt"
	"strings"
)

type TaskUnmarshalError struct {
	TaskJson string
	Err      error
}

func (e TaskUnmarshalError) Error() string {
	return fmt.Sprintf("Error unmarshaling task '%v'\n%v", e.TaskJson, e.Err.Error())
}

type DepsMarshalError struct {
	Deps []int
	Err  error
}

func (e DepsMarshalError) Error() string {
	depsStr := strings.Join(strings.Fields(fmt.Sprint(e.Deps)), ",")
	return fmt.Sprintf("Error marshaling dependency array '%v'\n%v", depsStr, e.Err.Error())
}
