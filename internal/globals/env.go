package globals

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Environments int

const (
	Test Environments = iota
	Dev
	Prod
)

const AppName = "mcpsequencer"

var (
	Environment Environments
	DbName      string
)

func Init() error {
	if isDev() {
		return initEnv(Dev)
	} else {
		return initEnv(Prod)
	}
}

func InitTest() error {
	return initEnv(Test)
}

func initEnv(env Environments) error {
	Environment = env

	if env == Dev {
		DbName = "dev.db"
	} else if env == Test {
		DbName = "test.db"
	} else {
		dbName, err := initProdDbPath()
		if err != nil {
			return err
		}

		DbName = dbName
	}

	return nil
}

func isDev() bool {
	for _, v := range os.Args {
		if v == "--dev" {
			return true
		}
	}
	return false
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func initProdDbPath() (string, error) {
	var dir string

	if isWindows() {
		dir = filepath.Join(os.Getenv("LOCALAPPDATA"), AppName)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		dir = filepath.Join(homeDir, ".local", "share", AppName)
	}

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, fmt.Sprintf("%v.db", AppName))
	return path, err
}
