package globals

type Environments int

const (
	Test Environments = iota
	Dev
	Prod
)

var (
	Environment Environments
	DbName      string
)

func Init(env Environments) {
	Environment = env

	if env == Dev {
		DbName = "dev.db"
	} else if env == Test {
		DbName = "test.db"
	} else {
		DbName = "projects.db"
	}
}
