package db

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	Init()

	file, err := os.Open("projects.db")
	defer file.Close()
	if err != nil {
		t.Errorf("error opening db file: \n%v\n", err)
	}

}
