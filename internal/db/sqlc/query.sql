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

-- Tasks CRUD Operations

-- name: CreateTask :one
INSERT INTO tasks (name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes;

-- name: GetTask :one
SELECT task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes
FROM tasks
WHERE task_id = ?;

-- name: GetAllTasks :many
SELECT task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes
FROM tasks
ORDER BY sort;

-- name: GetTasksByProject :many
SELECT task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes
FROM tasks
WHERE project_id = ?
ORDER BY sort;

-- name: GetCompletedTasks :many
SELECT task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes
FROM tasks
WHERE is_completed = 1
ORDER BY sort;

-- name: GetPendingTasks :many
SELECT task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes
FROM tasks
WHERE is_completed = 0 AND is_in_progress = 0
ORDER BY sort;

-- name: UpdateTask :one
UPDATE tasks
SET name = ?, description = ?, sort = ?, dependencies_json = ?, is_completed = ?, is_in_progress = ?, notes = ?
WHERE task_id = ?
RETURNING task_id, description, project_id, sort, is_completed, is_in_progress, notes;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET is_completed = ?, is_in_progress = ?
WHERE task_id = ?
RETURNING task_id, name, description, project_id, sort, dependencies_json, is_completed, is_in_progress, notes;

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

-- name: GetTasksWithProject :many
SELECT 
    t.task_id,
    t.name,
    t.description,
    t.project_id,
    t.sort,
    t.dependencies_json,
    t.is_completed,
    t.is_in_progress,
    t.notes,
    p.name as project_name
FROM tasks t
JOIN projects p ON t.project_id = p.project_id
ORDER BY p.name, t.sort;

-- name: GetProjectWithTaskCount :one
SELECT 
    p.project_id,
    p.name,
    p.description,
    p.absolute_path,
    COUNT(t.task_id) as task_count,
    COUNT(CASE WHEN t.is_completed = 1 THEN 1 END) as completed_count,
    COUNT(CASE WHEN t.is_in_progress = 1 THEN 1 END) as failed_count
FROM projects p
LEFT JOIN tasks t ON p.project_id = t.project_id
WHERE p.project_id = ?
GROUP BY p.project_id, p.name;

-- name: GetAllProjectsWithTaskCounts :many
SELECT 
    p.project_id,
    p.name,
    p.description,
    p.absolute_path,
    COUNT(t.task_id) as task_count,
    COUNT(CASE WHEN t.is_completed = 1 THEN 1 END) as completed_count,
    COUNT(CASE WHEN t.is_in_progress = 1 THEN 1 END) as failed_count
FROM projects p
LEFT JOIN tasks t ON p.project_id = t.project_id
GROUP BY p.project_id, p.name
ORDER BY p.name;
