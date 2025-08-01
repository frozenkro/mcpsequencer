package models

type DependencyDiscriminator int

const (
	SortId DependencyDiscriminator = iota
	TaskId
)

// Id: Either Sort ID or Task ID of Dependent Task
// DependsOn: Either Sort ID or Task ID of Task depended on
// Discriminator: Indicates whether Ids are Sort IDs or Task IDs
type Dependency struct {
	Id            int                     `json:"id"`
	DependsOn     int                     `json:"depends_on"`
	Discriminator DependencyDiscriminator `json:"discriminator"`
}

func NewDependencyWithSortIds(id int, dependsOn int) Dependency {
	return Dependency{
		Id:            id,
		DependsOn:     dependsOn,
		Discriminator: SortId,
	}
}

func NewDependencyWithTaskIds(id int, dependsOn int) Dependency {
	return Dependency{
		Id:            id,
		DependsOn:     dependsOn,
		Discriminator: TaskId,
	}
}

type SortIdTaskIdMap map[int64]int64
