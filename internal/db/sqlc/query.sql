-- name: CreateProject :one
INSERT INTO projects (name, description, absolute_path)
VALUES (?, ?, ?)
RETURNING project_id, name, description, absolute_path;

-- name: GetProject :one
SELECT project_id, name, description, absolute_path
FROM projects
WHERE project_id = ?;

-- name: GetAllProjects :many
SELECT project_id, name, description, absolute_path
FROM projects
ORDER BY name;

-- name: UpdateProject :one
UPDATE projects
SET name = ?,
description = ?,
absolute_path = ?
WHERE project_id = ?
RETURNING project_id, name, description, absolute_path;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = ?;

-- name: CreateTask :one
INSERT INTO tasks (name, description, project_id, sort, is_completed, is_in_progress, notes)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes;

-- name: GetTask :one
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes
FROM tasks
WHERE task_id = ?;

-- name: GetAllTasks :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes
FROM tasks
ORDER BY sort;

-- name: GetTasksByProject :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes
FROM tasks
WHERE project_id = ?
ORDER BY sort;

-- name: UpdateTask :one
UPDATE tasks
SET name = ?, description = ?, sort = ?, is_completed = ?, is_in_progress = ?, notes = ?
WHERE task_id = ?
RETURNING task_id, description, project_id, sort, is_completed, is_in_progress, notes;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET is_completed = ?, is_in_progress = ?
WHERE task_id = ?
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes;

-- name: UpdateTaskSort :exec
UPDATE tasks
SET sort = ?
WHERE task_id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE task_id = ?;

-- name: DeleteTasksByProject :exec
DELETE FROM tasks
WHERE project_id = ?;

-- name: GetDependenciesForTask :many
SELECT task_id, depends_on
FROM dependencies
WHERE task_id = ?;

-- name: AddDependencyForTask :exec
INSERT INTO dependencies (task_id, depends_on)
VALUES (?, ?);

-- name: RemoveDependency :exec
DELETE FROM dependencies
WHERE task_id = ? AND depends_on = ?;
