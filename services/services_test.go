package services_test

import (
	"context"
	"testing"

	"github.com/frozenkro/mcpsequencer/db"
	"github.com/frozenkro/mcpsequencer/globals"
	"github.com/frozenkro/mcpsequencer/projectsdb"
	"github.com/frozenkro/mcpsequencer/services"
	"github.com/stretchr/testify/assert"
)

func TestCreateProject(t *testing.T) {
	globals.Init(globals.Test)
	db.Init()
	s := services.Services{}
	ctx := context.Background()

	projectName := "Test Project Name"
	tasks := []string{"Test task 1", "Test task 2"}

	err := s.CreateProject(ctx, projectName, tasks)
	assert.Nil(t, err)

	conn, err := db.DB.Conn(ctx)
	assert.Nil(t, err)

	rows, err := conn.QueryContext(ctx, "SELECT * FROM projects WHERE Name = ?", projectName)
	assert.Nil(t, err)

	p := projectsdb.Project{}
	err = rows.Scan(p)
	assert.Nil(t, err)

	assert.False(t, rows.Next())

	assert.Equal(t, projectName, p.Name)
}
