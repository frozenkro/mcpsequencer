package db

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Project struct {
	Name  string   `json:"name"`
	Tasks []string `json:"tasks"`
}

const dbFileName = "projects.db"

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	createProjectsTableQuery := `
	CREATE TABLE IF NOT EXISTS projects (
		project_id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);
	`
	createTasksTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		task_id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		project_id INTEGER NOT NULL,
		order INTEGER NOT NULL,
		is_completed INTEGER NOT NULL,
		is_failed INTEGER NOT NULL,
		FOREIGN KEY(project_id) REFERENCES projects(project_id)
	);
	`
	_, err = db.Exec(createProjectsTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTasksTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func InsertProject(p Project) error {
	result, err := db.Exec(
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
			order,
			is_completed,
			is_failed)
		VALUES (?, ?, ?, ?, ?)
	`
	pid, err := result.LastInsertId()
	for i, v := range p.Tasks {
		_, err := db.Exec(
			taskSql,
			v,
			pid,
			i+1,
			0,
			0,
		)

		if err != nil {
			return err
		}
	}

	return err
}

func GetProjects() ([]Project, error) {
	rows, err := db.Query("SELECT name, tasks FROM projects")
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
