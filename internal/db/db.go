package db

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"log"

	"github.com/frozenkro/mcpsequencer/internal/globals"
	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	Name  string   `json:"name"`
	Tasks []string `json:"tasks"`
}

//go:embed sqlc/schema.sql
var SchemaSql string

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("sqlite3", globals.DbName)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(SchemaSql)
	if err != nil {
		log.Fatal(err)
	}

}

func InsertProject(p Project) error {
	result, err := DB.Exec(
		"INSERT INTO projects(name) VALUES (?)",
		p.Name,
	)
	if err != nil {
		return err
	}

	taskSql := `
		INSERT INTO tasks(
			description,
			project_id,
			sort,
			is_completed,
			is_in_progress,
		  notes)
		VALUES (?, ?, ?, ?, ?)
	`
	pid, err := result.LastInsertId()
	for i, v := range p.Tasks {
		_, err := DB.Exec(
			taskSql,
			v,
			pid,
			i+1,
			0,
			0,
			"",
		)

		if err != nil {
			return err
		}
	}

	return err
}

func GetAllProjects() ([]Project, error) {
	rows, err := DB.Query("SELECT name, tasks FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		var tasksJSON string
		err = rows.Scan(&p.Name, &tasksJSON)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal([]byte(tasksJSON), &p.Tasks)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}
