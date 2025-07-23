package db

import (
	"os"
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/globals"
)

func TestInit(t *testing.T) {
	globals.InitTest()
	Init()

	file, err := os.Open(globals.DbName)
	defer teardown(file)
	if err != nil {
		t.Errorf("error opening db file: \n%v\n", err)
	}

}

func teardown(dbFile *os.File) {
	dbFile.Close()
	os.Remove(globals.DbName)
}
