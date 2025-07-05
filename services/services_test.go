package services_test

import (
	"context"
	"os"
	"testing"

	"github.com/frozenkro/mcpsequencer/db"
	"github.com/frozenkro/mcpsequencer/globals"
	"github.com/frozenkro/mcpsequencer/projectsdb"
	"github.com/frozenkro/mcpsequencer/services"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	globals.Init(globals.Test)
	os.Remove(globals.DbName)
	db.Init()
	code := m.Run()
	os.Exit(code)
}

func TestCreateProject(t *testing.T) {
	s := services.Services{}
	ctx := context.Background()

	projectName := "Test Project Name"
	tasks := []string{"Test task 1", "Test task 2"}

	err := s.CreateProject(ctx, projectName, tasks)
	assert.Nil(t, err)

	conn, err := db.DB.Conn(ctx)
	assert.Nil(t, err)

	rows, err := conn.QueryContext(ctx, "SELECT project_id, name FROM projects WHERE Name = ?", projectName)
	assert.Nil(t, err)

	p := projectsdb.Project{}
	rows.Next()
	err = rows.Scan(&p.ProjectID, &p.Name)
	assert.Nil(t, err)

	assert.False(t, rows.Next())

	assert.Equal(t, projectName, p.Name)

	rows, err = conn.QueryContext(ctx, `
		SELECT task_id, 
		description,
		sort
		FROM tasks
		WHERE project_id = ?
	`, p.ProjectID)

	for rows.Next() {
		task := projectsdb.Task{}
		err = rows.Scan(&task.TaskID, &task.Description, &task.Sort)
		assert.Nil(t, err)

		assert.Equal(t, tasks[task.Sort], task.Description)
	}
}
