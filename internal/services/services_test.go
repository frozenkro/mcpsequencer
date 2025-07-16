package services_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/frozenkro/mcpsequencer/internal/db"
	"github.com/frozenkro/mcpsequencer/internal/globals"
	"github.com/frozenkro/mcpsequencer/internal/projectsdb"
	"github.com/frozenkro/mcpsequencer/internal/services"
	"github.com/stretchr/testify/assert"
)

var (
	s    services.Services
	ctx  context.Context
	conn *sql.Conn
	err  error
)

func TestMain(m *testing.M) {
	globals.Init(globals.Test)
	os.Remove(globals.DbName)

	s = services.Services{}
	ctx = context.Background()

	db.Init()

	conn, err = db.DB.Conn(ctx)
	if err != nil {
		log.Fatalf("ERROR: DB Connection Failed during Test Initialization\n%v\n", err.Error())
	}

	code := m.Run()
	os.Exit(code)
}

func TestCreateProject(t *testing.T) {
	projectName := "Test Project Name"
	tasks := []string{"Test task 1", "Test task 2"}

	err := s.CreateProject(ctx, projectName, tasks)
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

func TestGetProjects(t *testing.T) {
	p1 := projectsdb.Project{
		ProjectID: 1,
		Name:      "TestProjectName1",
	}
	p2 := projectsdb.Project{
		ProjectID: 2,
		Name:      "TestProjectName2",
	}

	_, err := conn.ExecContext(ctx, `
		INSERT INTO projects (Name) 
		VALUES (?), (?)`,
		p1.Name, p2.Name)
	assert.Nil(t, err)

	projects, err := s.GetProjects(ctx)
	assert.Nil(t, err)

	var project1Found = false
	var project2Found = false
	for _, v := range projects {
		if v.Name == p1.Name {
			project1Found = true
		} else if v.Name == p2.Name {
			project2Found = true
		}
	}

	assert.True(t, project1Found)
	assert.True(t, project2Found)
}

func TestRenameProject(t *testing.T) {
	oldName := "BeforeRenameTest"
	newName := "AfterRenameTest"

	res, err := conn.ExecContext(ctx, `
		INSERT INTO projects (Name)
		VALUES (?)
		`,
		oldName,
	)
	assert.Nil(t, err)

	projectId, err := res.LastInsertId()
	assert.Nil(t, err)

	err = s.RenameProject(ctx, projectId, newName)
	assert.Nil(t, err)

	row := conn.QueryRowContext(ctx,
		"SELECT project_id, name FROM projects WHERE project_id = ?",
		projectId,
	)
	p := projectsdb.Project{}
	err = row.Scan(&p.ProjectID, &p.Name)

	assert.Nil(t, err)

	assert.Equal(t, newName, p.Name)
}

func TestDeleteProject(t *testing.T) {
	res, err := conn.ExecContext(ctx, `
		INSERT INTO projects (Name)
		VALUES (?)
		`,
		"To Be Deleted",
	)
	assert.Nil(t, err)
	projectId, err := res.LastInsertId()
	assert.Nil(t, err)

	err = s.DeleteProject(ctx, projectId)

	row := conn.QueryRowContext(ctx,
		"SELECT project_id, name FROM projects WHERE project_id = ?",
		projectId,
	)
	p := projectsdb.Project{}
	err = row.Scan(&p.ProjectID, &p.Name)
	assert.ErrorIs(t, sql.ErrNoRows, err)
}

func TestGetTasksByProject(t *testing.T) {
	projectName := "Test Project With Tasks"
	taskNames := []string{"Test task 1", "Test task 2"}

	res, err := conn.ExecContext(ctx,
		`INSERT INTO projects (name)
		VALUES (?)`,
		projectName,
	)
	assert.Nil(t, err)

	projectId, err := res.LastInsertId()
	assert.Nil(t, err)

	_, err = conn.ExecContext(ctx,
		`INSERT INTO tasks 
		(description, project_id, sort, is_completed, is_in_progress)
		VALUES (?, ?, ?, ?, ?), (?, ?, ?, ?, ?)`,
		taskNames[0], projectId, 0, 0, 0,
		taskNames[1], projectId, 1, 0, 0,
	)
	assert.Nil(t, err)

	tasks, err := s.GetTasksByProject(ctx, projectId)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(tasks))
	task1Found := false
	task2Found := false
	for _, v := range tasks {
		if v.Description == taskNames[0] {
			task1Found = true
		} else if v.Description == taskNames[1] {
			task2Found = true
		}
	}
	assert.True(t, task1Found)
	assert.True(t, task2Found)
}
