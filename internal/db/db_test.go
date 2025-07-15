package db

import (
	"os"
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/globals"
)

func TestInit(t *testing.T) {
	globals.Init(globals.Test)
	Init()

	file, err := os.Open(globals.DbName)
	defer file.Close()
	if err != nil {
		t.Errorf("error opening db file: \n%v\n", err)
	}

}
