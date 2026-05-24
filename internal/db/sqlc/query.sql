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
INSERT INTO tasks (name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, estimated_hours)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours;

-- name: GetTask :one
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
WHERE task_id = ?;

-- name: GetAllTasks :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
ORDER BY sort;

-- name: GetTasksByProject :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
WHERE project_id = ?
ORDER BY sort;

-- name: UpdateTask :one
UPDATE tasks
SET name = ?, description = ?, sort = ?, is_completed = ?, is_in_progress = ?, notes = ?,
    owner = ?, scheduled_date = ?, phase_id = ?, estimated_hours = ?
WHERE task_id = ?
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours;

-- name: UpdateTaskStatus :one
UPDATE tasks
SET is_completed = ?, is_in_progress = ?
WHERE task_id = ?
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours;

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

-- name: CreatePhase :one
INSERT INTO phases (project_id, name, description, start_date, end_date, sort)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING phase_id, project_id, name, description, start_date, end_date, sort;

-- name: GetPhase :one
SELECT phase_id, project_id, name, description, start_date, end_date, sort
FROM phases
WHERE phase_id = ?;

-- name: GetPhasesForProject :many
SELECT phase_id, project_id, name, description, start_date, end_date, sort
FROM phases
WHERE project_id = ?
ORDER BY sort;

-- name: UpdatePhase :one
UPDATE phases
SET name = ?, description = ?, start_date = ?, end_date = ?, sort = ?
WHERE phase_id = ?
RETURNING phase_id, project_id, name, description, start_date, end_date, sort;

-- name: DeletePhase :exec
DELETE FROM phases WHERE phase_id = ?;

-- name: UnlinkTasksFromPhase :exec
UPDATE tasks SET phase_id = NULL WHERE phase_id = ?;

-- name: GetTasksForPhase :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
WHERE phase_id = ?
ORDER BY sort;

-- name: SetTaskOwner :exec
UPDATE tasks SET owner = ? WHERE task_id = ?;

-- name: SetTaskBlocker :exec
UPDATE tasks SET blocker_text = ?, blocked_at = datetime('now') WHERE task_id = ?;

-- name: ResolveTaskBlocker :exec
UPDATE tasks SET blocker_text = NULL, blocked_at = NULL WHERE task_id = ?;

-- name: SetTaskSchedule :exec
UPDATE tasks SET scheduled_date = ?, phase_id = ? WHERE task_id = ?;

-- name: GetTasksForOwner :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
WHERE project_id = ? AND owner = ?
ORDER BY sort;

-- name: GetTodaysTasks :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
WHERE scheduled_date <= date('now') AND is_completed = 0 AND project_id = ?
ORDER BY COALESCE(phase_id, 999999), sort;

-- name: AddTaskNote :exec
INSERT INTO task_notes (task_id, note, created_at) VALUES (?, ?, datetime('now'));

-- name: GetTaskNotes :many
SELECT note_id, task_id, note, created_at FROM task_notes WHERE task_id = ? ORDER BY created_at;

-- name: CreateExitCriterion :one
INSERT INTO exit_criteria (phase_id, description, is_completed, sort)
VALUES (?, ?, ?, ?)
RETURNING criterion_id, phase_id, description, is_completed, sort;

-- name: GetExitCriteriaForPhase :many
SELECT criterion_id, phase_id, description, is_completed, sort
FROM exit_criteria
WHERE phase_id = ?
ORDER BY sort;

-- name: CompleteExitCriterion :exec
UPDATE exit_criteria SET is_completed = 1 WHERE criterion_id = ?;

-- name: UncompleteExitCriterion :exec
UPDATE exit_criteria SET is_completed = 0 WHERE criterion_id = ?;

-- name: GetDownstreamTasks :many
SELECT task_id FROM dependencies WHERE depends_on = ?;

-- name: GetUpstreamTasks :many
SELECT depends_on FROM dependencies WHERE task_id = ?;
