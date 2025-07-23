package db

import (
	"log"
	"os"
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/globals"
)

func TestInit(t *testing.T) {
	if err := globals.InitTest(); err != nil {
		log.Fatalf("Application Initialization failed. \nError: %v\n", err.Error())
	}
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
