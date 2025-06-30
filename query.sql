-- name: CreateProject :one
INSERT INTO projects (name)
VALUES (?)
RETURNING project_id, name;

-- name: GetProject :one
SELECT project_id, name
FROM projects
WHERE project_id = ?;

-- name: GetAllProjects :many
SELECT project_id, name
FROM projects
ORDER BY name;

-- name: UpdateProject :one
UPDATE projects
SET name = ?
WHERE project_id = ?
RETURNING project_id, name;

-- name: DeleteProject :exec
DELETE FROM projects
WHERE project_id = ?;

-- Tasks CRUD Operations

-- name: CreateTask :one
INSERT INTO tasks (description, project_id, sort, is_completed, is_failed, notes)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING task_id, description, project_id, sort, is_completed, is_failed, notes;

-- name: GetTask :one
SELECT task_id, description, project_id, sort, is_completed, is_failed, notes
FROM tasks
WHERE task_id = ?;

-- name: GetAllTasks :many
SELECT task_id, description, project_id, sort, is_completed, is_failed, notes
FROM tasks
ORDER BY sort;

-- name: GetTasksByProject :many
SELECT task_id, description, project_id, sort, is_completed, is_failed, notes
FROM tasks
WHERE project_id = ?
ORDER BY sort;

-- name: GetCompletedTasks :many
SELECT task_id, description, project_id, sort, is_completed, is_failed, notes
FROM tasks
WHERE is_completed = 1
ORDER BY sort;

-- name: GetFailedTasks :many
SELECT task_id, description, project_id, sort, is_completed, is_failed, notes
FROM tasks
WHERE is_failed = 1
ORDER BY sort;

-- name: GetPendingTasks :many
SELECT task_id, description, project_id, sort, is_completed, is_failed, notes
FROM tasks
WHERE is_completed = 0 AND is_failed = 0
ORDER BY sort;

-- name: UpdateTask :one
UPDATE tasks
SET description = ?, sort = ?, is_completed = ?, is_failed = ?, notes = ?
WHERE task_id = ?
RETURNING task_id, description, project_id, sort, is_completed, is_failed, notes;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET is_completed = ?, is_failed = ?
WHERE task_id = ?
RETURNING task_id, description, project_id, sort, is_completed, is_failed, notes;

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

-- Joined queries for richer data

-- name: GetTasksWithProject :many
SELECT 
    t.task_id,
    t.description,
    t.project_id,
    t.sort,
    t.is_completed,
    t.is_failed,
    t.notes,
    p.name as project_name
FROM tasks t
JOIN projects p ON t.project_id = p.project_id
ORDER BY p.name, t.sort;

-- name: GetProjectWithTaskCount :one
SELECT 
    p.project_id,
    p.name,
    COUNT(t.task_id) as task_count,
    COUNT(CASE WHEN t.is_completed = 1 THEN 1 END) as completed_count,
    COUNT(CASE WHEN t.is_failed = 1 THEN 1 END) as failed_count
FROM projects p
LEFT JOIN tasks t ON p.project_id = t.project_id
WHERE p.project_id = ?
GROUP BY p.project_id, p.name;

-- name: GetAllProjectsWithTaskCounts :many
SELECT 
    p.project_id,
    p.name,
    COUNT(t.task_id) as task_count,
    COUNT(CASE WHEN t.is_completed = 1 THEN 1 END) as completed_count,
    COUNT(CASE WHEN t.is_failed = 1 THEN 1 END) as failed_count
FROM projects p
LEFT JOIN tasks t ON p.project_id = t.project_id
GROUP BY p.project_id, p.name
ORDER BY p.name;
