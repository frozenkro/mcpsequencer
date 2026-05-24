# Feature Proposal: Phase-Aware Task Ownership & Daily Standup
## Making mcpsequencer a substitute for static markdown planning

**Author:** AI Agent (Hermes)  
**Date:** 2026-05-23  
**Target Branch:** `main` (proposed PR branch: `feature/plan-substitution`)  
**Estimated Effort:** 2–3 working days  

---

## Problem Statement

The current `mcpsequencer` tracks generic project tasks well but lacks the dimensions needed for real-world, multi-week execution planning:

1. **No ownership tag** — cannot distinguish "my hands-on" vs "delegate to agent" vs "collaborative" tasks.
2. **No temporal concept** — `sort` is linear and static; there is no date, phase, or milestone.
3. **No exit criteria** — cannot define "what does done look like for this phase?"
4. **No blockers / daily notes** — no historical timeline per task.
5. **No cross-phase hand-off view** — cannot model "Phase 1 output → Phase 2 input" dependencies.
6. **No daily focus screen** — cannot answer "what am I working on today?"

This proposal adds those dimensions while keeping the existing MCP + TUI architecture intact.

---

## Design Overview

### New Concepts

| Concept | Description |
|---|---|
| **Phase** | High-level milestone with a date range (e.g., "Phase 1: Infrastructure", Sun May 24 → Sat May 30). Contains a checklist of exit criteria. |
| **Owner** | `USER`, `AI_AGENT`, or `COLLAB` per task. Stored in `tasks.owner`. |
| **Blocker** | Free-text explanation of why a task cannot proceed. Stored in `tasks.blocker_text`, nullable. |
| **Task Note** | Chronological daily updates attached to a task. Separate table `task_notes`. |
| **Exit Criterion** | JSON-like checklist item belonging to a phase. Stored in `exit_criteria`. |
| **Cross-phase dependencies** | The existing `dependencies` table is reused as-is. Tasks in any phase can depend on tasks in any other phase (or no phase). The dependency table already supports this because it uses actual `task_id` FKs. A phase simply groups tasks; dependencies span phases naturally. |

### Coexistence of `tasks.notes` and `task_notes`

The existing `tasks.notes` column remains for **static task description / scratch notes**. The new `task_notes` table is for **chronological daily standup entries** — a timestamped log. Both coexist intentionally.

### What Stays the Same

- Projects, tasks, intra-task DAG dependencies (existing `dependencies` table).
- MCP tool definitions and handler patterns.
- TUI Bubble Tea architecture.
- SQLite + sqlc code generation.
- No new dependency table needed — `dependencies` already handles cross-phase within a project.

---

## Schema Changes

### 1. New Table: `phases`

A project can have many phases. This table replaces the implicit "week grouping" in a markdown file.

```sql
CREATE TABLE IF NOT EXISTS phases (
  phase_id INTEGER PRIMARY KEY AUTOINCREMENT,
  project_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  description TEXT NULL,
  start_date TEXT NOT NULL,  -- ISO 8601 date (YYYY-MM-DD)
  end_date TEXT NOT NULL,
  sort INTEGER NOT NULL,     -- display order within project
  FOREIGN KEY(project_id) REFERENCES projects(project_id) ON DELETE CASCADE
);
```

**Validation:** `start_date <= end_date`, enforced at app layer.

### 2. Extend `tasks` Table

SQLite can add nullable columns via `ALTER TABLE ADD COLUMN`, but **cannot add FOREIGN KEY constraints via ALTER TABLE**. Therefore, `phase_id` is added as a plain integer column. Referential integrity is enforced at the **application layer**.

```sql
ALTER TABLE tasks ADD COLUMN owner TEXT NOT NULL DEFAULT 'USER';
ALTER TABLE tasks ADD COLUMN scheduled_date TEXT NULL;    -- ISO 8601 date
ALTER TABLE tasks ADD COLUMN phase_id INTEGER NULL;       -- app-level FK to phases
ALTER TABLE tasks ADD COLUMN blocker_text TEXT NULL;
ALTER TABLE tasks ADD COLUMN blocked_at TEXT NULL;        -- ISO 8601 timestamp
ALTER TABLE tasks ADD COLUMN estimated_hours INTEGER NULL; -- rough sizing
```

- `owner`: App-layer validation `CHECK(owner IN ('USER', 'AI_AGENT', 'COLLAB'))`.
- `phase_id`: Not an enforced DB-level FK. Services must validate that the referenced phase exists before insert/update.

### 3. New Table: `task_notes`

Chronological log per task. Replaces the manual daily standup section in `plan.md`.

```sql
CREATE TABLE IF NOT EXISTS task_notes (
  note_id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  note TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
);
```

### 4. New Table: `exit_criteria`

Checklist items belonging to a phase. Replaces the `[ ]` checkboxes in markdown phase headers.

```sql
CREATE TABLE IF NOT EXISTS exit_criteria (
  criterion_id INTEGER PRIMARY KEY AUTOINCREMENT,
  phase_id INTEGER NOT NULL,
  description TEXT NOT NULL,
  is_completed INTEGER NOT NULL DEFAULT 0,
  sort INTEGER NOT NULL,
  FOREIGN KEY(phase_id) REFERENCES phases(phase_id) ON DELETE CASCADE
);
```

---

## sqlc Queries (New)

### Phase Queries

```sql
-- name: CreatePhase :one
INSERT INTO phases (project_id, name, description, start_date, end_date, sort)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING phase_id, project_id, name, description, start_date, end_date, sort;

-- name: GetPhasesForProject :many
SELECT phase_id, project_id, name, description, start_date, end_date, sort
FROM phases
WHERE project_id = ?
ORDER BY sort;

-- name: GetPhase :one
SELECT phase_id, project_id, name, description, start_date, end_date, sort
FROM phases
WHERE phase_id = ?;

-- name: UpdatePhase :one
UPDATE phases
SET name = ?, description = ?, start_date = ?, end_date = ?, sort = ?
WHERE phase_id = ?
RETURNING phase_id, project_id, name, description, start_date, end_date, sort;

-- name: DeletePhase :exec
DELETE FROM phases WHERE phase_id = ?;

-- name: UnlinkTasksFromPhase :exec
UPDATE tasks SET phase_id = NULL WHERE phase_id = ?;
```

### Enhanced Task Queries

**IMPORTANT:** `CreateTask` MUST be extended to insert new columns. Otherwise `owner`/`scheduled_date`/`phase_id`/`estimated_hours` passed via MCP will silently drop to defaults.

Extend the existing `CreateTask` query from:

```sql
-- name: CreateTask :one
INSERT INTO tasks (name, description, project_id, sort, is_completed, is_in_progress, notes)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes;
```

To:

```sql
-- name: CreateTask :one
INSERT INTO tasks (name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, estimated_hours)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours;
```

And update `UpdateTask` from:

```sql
-- name: UpdateTask :one
UPDATE tasks
SET name = ?, description = ?, sort = ?, is_completed = ?, is_in_progress = ?, notes = ?
WHERE task_id = ?
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes;
```

To:

```sql
-- name: UpdateTask :one
UPDATE tasks
SET name = ?, description = ?, sort = ?, is_completed = ?, is_in_progress = ?, notes = ?,
    owner = ?, scheduled_date = ?, phase_id = ?, estimated_hours = ?
WHERE task_id = ?
RETURNING task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours;
```

New queries:

```sql
-- name: GetTasksForPhase :many
SELECT task_id, name, description, project_id, sort, is_completed, is_in_progress, notes, owner, scheduled_date, phase_id, blocker_text, blocked_at, estimated_hours
FROM tasks
WHERE phase_id = ?
ORDER BY sort;

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

-- name: SetTaskOwner :exec
UPDATE tasks SET owner = ? WHERE task_id = ?;

-- name: SetTaskBlocker :exec
UPDATE tasks SET blocker_text = ?, blocked_at = datetime('now') WHERE task_id = ?;

-- name: ResolveTaskBlocker :exec
UPDATE tasks SET blocker_text = NULL, blocked_at = NULL WHERE task_id = ?;

-- name: SetTaskSchedule :exec
UPDATE tasks SET scheduled_date = ?, phase_id = ? WHERE task_id = ?;
```

### Task Note Queries

```sql
-- name: AddTaskNote :exec
INSERT INTO task_notes (task_id, note, created_at) VALUES (?, ?, datetime('now'));

-- name: GetTaskNotes :many
SELECT note_id, task_id, note, created_at FROM task_notes WHERE task_id = ? ORDER BY created_at;
```

### Exit Criteria Queries

```sql
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
```

---

## Domain Model Changes

### New / Modified Models

| Model | Fields |
|---|---|
| `models.Phase` (new) | `PhaseId int`, `ProjectId int`, `Name string`, `Description string`, `StartDate string`, `EndDate string`, `Sort int` |
| `models.Task` (modified) | Adds: `Owner Owner`, `ScheduledDate *string`, `PhaseId *int`, `BlockerText *string`, `BlockedAt *string`, `EstimatedHours *int` |
| `models.Owner` (new) | `USER`, `AI_AGENT`, `COLLAB` |
| `models.TaskNote` (new) | `NoteId int`, `TaskId int`, `Note string`, `CreatedAt string` |
| `models.ExitCriterion` (new) | `CriterionId int`, `PhaseId int`, `Description string`, `IsCompleted bool`, `Sort int` |

### Updated `models.NewTask` factory

`internal/models/task.go` line 27 must be extended to populate the new fields from the DB-generated struct:

```go
func NewTask(dbTask projectsdb.Task, deps []Dependency) Task {
    status := NotStarted
    if dbTask.IsCompleted != int64(0) { status = Completed }
    else if dbTask.IsInProgress != int64(0) { status = InProgress }

    notes, _ := dbTask.Notes.(string)
    owner := Owner(dbTask.Owner)  // new
    
    var scheduledDate *string  // new
    if dbTask.ScheduledDate.Valid { scheduledDate = &dbTask.ScheduledDate.String }
    
    var phaseId *int  // new
    if dbTask.PhaseID.Valid { v := int(dbTask.PhaseID.Int64); phaseId = &v }
    
    var blockerText *string  // new
    if dbTask.BlockerText.Valid { blockerText = &dbTask.BlockerText.String }
    
    var blockedAt *string  // new
    if dbTask.BlockedAt.Valid { blockedAt = &dbTask.BlockedAt.String }
    
    var estimatedHours *int  // new
    if dbTask.EstimatedHours.Valid { v := int(dbTask.EstimatedHours.Int64); estimatedHours = &v }

    return Task{
        TaskId: int(dbTask.TaskID), Name: dbTask.Name, Description: dbTask.Description,
        ProjectId: int(dbTask.ProjectID), Sort: int(dbTask.Sort), Status: status,
        Notes: notes, Dependencies: deps,
        Owner: owner, ScheduledDate: scheduledDate, PhaseId: phaseId,
        BlockerText: blockerText, BlockedAt: blockedAt, EstimatedHours: estimatedHours,
    }
}
```

### New `CreateTaskArgs` Fields

Extend the JSON task schema used by `createProject`. All new fields are optional.

```go
type CreateTaskArgs struct {
    Name         string `json:"name"`
    Description  string `json:"description"`
    SortId       int    `json:"sortId"`
    Dependencies []int  `json:"dependencies"`
    Owner        *string `json:"owner,omitempty"`          // defaults to "USER"
    ScheduledDate *string `json:"scheduledDate,omitempty"`  // ISO date
    PhaseId      *int   `json:"phaseId,omitempty"`        // belongs to phase
    EstimatedHours *int `json:"estimatedHours,omitempty"` // rough sizing
}
```

The `TaskArrayTransformer.ParseFromJson` must be updated to map these optional fields. If omitted, DB DEFAULTs handle it.

---

## MCP Tool Additions

| Tool | Params | Description |
|---|---|---|
| `createPhase` | `projectId`, `name`, `description`, `startDate`, `endDate`, `sort` | Create a new phase in a project. |
| `getPhasesForProject` | `projectId` | List all phases. |
| `updatePhase` | `phaseId`, `name`, `description`, `startDate`, `endDate`, `sort` | Modify a phase. |
| `deletePhase` | `phaseId` | Unlink all tasks in phase, then delete phase. |
| `setTaskOwner` | `taskId`, `owner` (`USER`\| `AI_AGENT` \| `COLLAB`) | Assign ownership. |
| `setTaskSchedule` | `taskId`, `scheduledDate`, `phaseId` | Move task to a date/phase. |
| `setTaskBlocker` | `taskId`, `blockerText` | Block a task. |
| `resolveTaskBlocker` | `taskId` | Clear blocker. |
| `addTaskNote` | `taskId`, `note` | Append a daily note. |
| `getTaskNotes` | `taskId` | Fetch all notes for a task. |
| `getTodaysTasks` | `projectId` | Tasks where `scheduledDate <= today` and not completed. |
| `getTasksByOwner` | `projectId`, `owner` | Filter tasks by ownership tag. |
| `addExitCriterion` | `phaseId`, `description`, `sort` | Add an exit criterion to a phase. |
| `completeExitCriterion` | `criterionId` | Mark as done. |
| `getExitCriteriaForPhase` | `phaseId` | List criteria + completion status. |
| `getDownstreamTasks` | `taskId` | What tasks depend on this one? (uses existing `dependencies` table; already exists as `GetDependenciesForTask` query — expose via MCP) |
| `getPhaseSummary` | `phaseId` | Completion %, owner breakdown, exit criteria status. |

**Modified Tools:**

| Tool | Change |
|---|---|
| `createProject` | Accept new optional task fields in the `Tasks` JSON array. |
| `getTasksForProject` | Return richer task model including `owner`, `scheduledDate`, `phaseId`, `blockerText`, `estimatedHours`. |
| `addTask` | Accept new optional params. |
| `getTaskListInstructions` | Update instructions text to document new optional fields. |

### Existing Bug Fixes (do as part of this PR)

| File | Line | Issue | Fix |
|---|---|---|---|
| `handlers.go` | 155 | `BeginTaskHandler` returns `"Task completed successfully"` | Return `"Task in progress"` |
| `handlers.go` | 117 | `AddTaskHandler` uses `globals.ProjectId` instead of `globals.Description` | Use correct param |

---

## Service Layer Additions

### `internal/services/services.go`

Add the following service methods to the existing `Services` struct. Keep existing methods untouched.

```go
func (s *Services) CreatePhase(ctx context.Context, args models.Phase) error
func (s *Services) GetPhases(ctx context.Context, projectId int64) ([]models.Phase, error)
func (s *Services) UpdatePhase(ctx context.Context, args models.UpdatePhaseArgs) error
func (s *Services) DeletePhase(ctx context.Context, phaseId int64) error

func (s *Services) SetTaskOwner(ctx context.Context, taskId int64, owner models.Owner) error
func (s *Services) SetTaskSchedule(ctx context.Context, taskId int64, date string, phaseId *int) error
func (s *Services) SetTaskBlocker(ctx context.Context, taskId int64, text string) error
func (s *Services) ResolveTaskBlocker(ctx context.Context, taskId int64) error
func (s *Services) GetTodaysTasks(ctx context.Context, projectId int64) ([]models.Task, error)
func (s *Services) GetTasksByOwner(ctx context.Context, projectId int64, owner models.Owner) ([]models.Task, error)

func (s *Services) AddTaskNote(ctx context.Context, taskId int64, note string) error
func (s *Services) GetTaskNotes(ctx context.Context, taskId int64) ([]models.TaskNote, error)

func (s *Services) AddExitCriterion(ctx context.Context, phaseId int64, desc string, sort int) error
func (s *Services) GetExitCriteria(ctx context.Context, phaseId int64) ([]models.ExitCriterion, error)
func (s *Services) CompleteExitCriterion(ctx context.Context, criterionId int64) error
```

### `DeletePhase` Service Logic

1. Call `s.Queries.UnlinkTasksFromPhase(ctx, phaseId)` to set `tasks.phase_id = NULL`.
2. Call `s.Queries.DeletePhase(ctx, phaseId)`.

This order avoids any referential ambiguity.

---

## TUI Changes

### TUI Key Scoping Note

The TUI currently dispatches `tea.KeyMsg` globally. New single-character keybindings (`b`, `o`, `n`) must be gated on `m.ActivePane` to avoid fire inside edit/input modes.

### 1. Phase List Screen

New Bubble Tea screen accessible from the project detail view.

- Shows phases as a vertical list ordered by `sort`.
- Each phase item shows: `name | start_date → end_date | [exit criteria: 3/7]`.
- Enter on a phase → Task list filtered to that phase.

### 2. Task List Enhancements

Modify `internal/tui/components/tasks/model.go`:

- **Owner badge:** render `[U]`, `[A]`, or `[C]` next to task name with different colors.
- **Blocker indicator:** if `blocker_text` is non-null, prepend `(!)` to the title.
- **Date badge:** if `scheduled_date` is set, append `[May 24]` to description.
- **Key binding:** `B` = set blocker (prompt), `O` = cycle owner, `N` = add note. Only active when `ActivePane == TaskListPane`.

### 3. Task Detail Screen Enhancements

Modify `internal/tui/components/taskdetail/model.go`:

- Show full note history under a "Notes" section.
- Show exit criteria if task belongs to a phase.
- Show upstream/downstream task names (via existing `GetDependenciesForTask` query).

### 4. Today's Focus Screen

New top-level view model `internal/tui/components/today/`:

- Queries `GetTodaysTasks` grouped by owner.
- Shows count: "USER (2), AI_AGENT (3), COLLAB (1)".
- Lists tasks with phase name + project name for context.
- Hotkey: `T` from project browser. Only active when `ActivePane == ProjectListPane`.

### 5. Exit Criteria Popup

From phase list, press `C` (checklist, only when in PhaseListPane) → modal popup showing exit criteria with `[x]`/`[ ]` checkboxes. Use `Space` to toggle completion.

---

## sqlc Generation Workflow

After editing `schema.sql` and `query.sql`:

```bash
sqlc generate
```

This regenerates:
- `internal/projectsdb/db.go`
- `internal/projectsdb/query.sql.go`
- `internal/projectsdb/models.go`

All existing Go code that imports `internal/projectsdb` will need to be updated for the new struct fields. This is the largest mechanical change.

---

## Test Plan

1. **DB migration:** Create a fresh `projects.db`, verify all tables are created with `IF NOT EXISTS`.
2. **Phase CRUD:** Create project → create phase → list phases → update phase → delete phase (verify tasks are unlinked and phase removed).
3. **Task with new fields:** Create task with `owner='AI_AGENT'`, `scheduledDate='2026-05-24'`, `phase_id=1`. Verify TUI renders owner badge.
4. **Today's focus:** Set `scheduledDate` to today. Verify `getTodaysTasks` returns it. Verify TUI screen shows it.
5. **Blocker:** `setTaskBlocker(1, "waiting for ARM64 image")` → TUI shows `(!)`. `resolveTaskBlocker(1)` → `(!)` cleared.
6. **Task notes:** `addTaskNote(1, "Refactored Dockerfile today")` → task detail shows it.
7. **Exit criteria:** `addExitCriterion(phaseId=1, "ARM64 image builds", sort=0)` → `getExitCriteriaForPhase` shows it. Complete it. Verify popup updates.
8. **Phase deletion safety:** Create phase, add tasks to phase, delete phase → tasks survive with `phase_id = NULL`.
9. **Existing tests:** Run `go test ./...` — no existing tests should fail.

---

## Backwards Compatibility

- All new SQL table columns have `DEFAULT` values. Existing databases will show `owner='USER'`, `phase_id=NULL`, etc.
- The JSON task schema in `createProject` gets **optional** new fields. Existing JSON payloads without them still parse.
- The MCP `getTasksForProject` tool now returns richer JSON but with the same top-level shape. Any external client that ignores unknown keys will continue working.
- **No `cross_project_dependencies` table added** — existing `dependencies` table is reused, preserving all existing cross-task relationships without structural changes.

---

## Files to Touch

| File | Change |
|---|---|
| `internal/db/sqlc/schema.sql` | Add `phases`, `task_notes`, `exit_criteria` tables; add columns to `tasks`. |
| `internal/db/sqlc/query.sql` | Add all new queries; extend `CreateTask`, `UpdateTask`, `GetTasksByProject`. |
| `internal/models/*.go` | Add new model structs; extend `Task`, `NewTask`, `CreateTaskArgs`. |
| `internal/projectsdb/*` | Regenerated by sqlc. |
| `internal/services/services.go` | Add service methods; fix existing minor bugs. |
| `internal/mcp/tools/tools.go` | Add tool definitions. |
| `internal/mcp/handlers/handlers.go` | Add handler implementations; fix existing bugs. |
| `internal/tui/components/tasks/*` | Owner badge, blocker, date rendering. |
| `internal/tui/components/taskdetail/*` | Notes, exit criteria, upstream/downstream display. |
| `internal/tui/components/today/*` | New screen: today's focus. |
| `internal/tui/components/phases/*` | New screen: phase list + exit criteria popup. |
| `internal/tui/app.go` | Wire new screens into navigation flow. |
| `internal/globals/mcp_args.go` | Add new MCP arg constants. |

---

## Rollback Plan

All changes are additive schema + additive code. To rollback:
1. Revert the PR.
2. The SQLite DB will have extra tables/columns — existing code will still work because it only reads/writes the columns it knows about. For strict rollback, run a manual DB migration to drop extra tables.
