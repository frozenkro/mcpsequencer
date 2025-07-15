package models

type CreateTaskArgs struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	SortId       int    `json:"sortId"`
	Dependencies []int  `json:"dependencies"`
}
