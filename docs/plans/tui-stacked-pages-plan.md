# TUI Refresh: Stacked Page Navigation + PR#2 Feature Surfacing

## Background

PR#2 (merged) introduced a rich new data model: tasks now have states, owners, scheduling, blockers, dependencies, estimated hours, phases, task notes, and exit criteria. The current TUI is still a cramped 3-pane layout that was built before these features existed and barely surfaces any of them.

This plan replaces the 3-pane horizontal layout with full-screen stacked pages, giving every screen room to show hierarchy, metadata, and rich detail.

## Current State Summary

**Backend (already in `main`):**
- `Task`: Status (NotStarted | InProgress | Completed | Failed), Owner, ScheduledDate, PhaseId, BlockerText/BlockedAt, EstimatedHours, Dependencies
- `Phase`: Name, Description, StartDate, EndDate, Sort
- `TaskNote`: NoteID, Note, CreatedAt
- `ExitCriterion`: CriterionID, Description, IsCompleted, Sort
- Full service layer: `GetPhases`, `AddTaskNote`, `GetTaskNotes`, `SetTaskBlocker`, `GetTodaysTasks`, `GetTasksByOwner`, `GetDownstreamTasks`, `GetUpstreamTasks`, `CompleteExitCriterion`, etc.

**Current TUI Gaps:**
| Feature | Backend | TUI |
|---------|---------|-----|
| Task states (4) | Exists | Status icon only |
| Owner | Exists | Single-char badge only |
| Scheduled date | Exists | Shown as raw text |
| Phase grouping | Exists | Not shown |
| Dependencies | Exists | Not shown |
| Estimated hours | Exists | Not shown |
| Task notes timeline | Exists | Not shown |
| Blocker detail | Exists | Shown as raw text |
| Exit criteria | Exists | Not shown |
| Phase detail | Exists | Not shown |

**Blockers found during review:**
1. `taskdetail.saveChanges()` is a no-op — it mutates local viewmodel but never calls the service layer.
2. No `UpdateTask` service method exists for general edits (only `UpdateTaskState`). The DB query `UpdateTask` exists in sqlc but has no service wrapper.
3. The old 3-pane root model (`internal/tui/model.go`) hardcodes horizontal 1/3 width splits and pane-switching.

---

## Proposed Architecture: Stacked Pages

Replace the 3-pane horizontal layout with a **page stack**:
- `Enter` / `Select` → **pushes** the next page onto the stack
- `Esc` / `h` / `backspace` → **pops** back to the previous page
- Only the top page renders, using full terminal width/height

This gives every page room to show rich information and task hierarchy.

### Page Map

```
ProjectsPage (root)
    └──► TasksPage(projectID)
            ├──► TaskDetailPage(taskID)
            │       ├──► TaskNotesPage(taskID)
            │       └──► TaskEditPage(taskID)
            └──► PhaseDetailPage(phaseID)
                    └──► ExitCriteriaPage(phaseID)
```

---

## Page Designs

### 1. Projects Page (full-screen)
- Project list with name + description + working directory + task count
- Highlighted item gets colored border
- `Enter` → push TasksPage for selected project
- `q` / `ctrl+c` → quit

### 2. Tasks Page (full-screen)
The primary win here is showing **actual task hierarchy** and **phase grouping**.

**Layout:**
```
┌──────────────────────────────────────────────────────┐
│ Project: sqncr                        [All|Today|Mine|Blocked] │
├──────────────────────────────────────────────────────┤
│ Phase: Foundation         2025-06-01 → 2025-06-15  │
│   [U] ○ Setup build pipeline              [2h] [05-28]     │
│   [A] 🔄 Generate sqlc models             [1h] [05-29]     │
│ Phase: Migrations                                    │
│   [U] ○ Add User table migration          [3h] [06-01] (!)│
│ Backlog                                              │
│   [C] ✅ Write docs                         [1h]           │
└──────────────────────────────────────────────────────┘
```

**Features:**
- **Phase grouping**: Phase headers with date range; tasks indented underneath their phase
- **Backlog section**: Tasks with no `PhaseId` (or `nil`) shown in an ungrouped "Backlog" block
- **Dependency-driven indentation**: Tasks indented under the task they depend on (computed from `Dependencies` graph)
- **Columns in delegate**: Owner `[U/A/C]` | Blocker `(!)` | Status icon | Task name | Est. hours | Scheduled date
- **Filters** (toggled via keys):
  - `A` = all tasks
  - `T` = today's tasks (`GetTodaysTasks`)
  - `M` = my tasks (`USER` owner)
  - `I` = AI tasks (`AI_AGENT`)
  - `B` = blocked tasks (has `BlockerText`)
- `Enter` → push TaskDetailPage
- `s` → toggle task status (cycle NotStarted → InProgress → Completed)

### 3. Task Detail Page (full-screen)
Rich information density. This is what the 3-pane layout could never provide.

**Sections (top to bottom):**
1. **Header bar**: Name (bold, colored), Status pill, Owner badge, Phase name (if assigned), Estimated hours, Scheduled date
2. **Blocker banner**: Red/yellow prominent bar if `BlockerText` exists, showing text + `BlockedAt` timestamp, with `r` key to resolve
3. **Description**: Scrollable wrapped text
4. **Dependencies panel**:
   - Upstream: "Depends on: [Task Name] [Task Name]"
   - Downstream: "Blocking: [Task Name] [Task Name]"
   (Display task names resolved from IDs)
5. **Notes timeline**: Chronological list like a chat log:
   ```
   2025-05-20 09:14  Started refactoring
   2025-05-21 14:30  sqlc generation failing on nullable fields
   ```
   - `n` → opens inline textarea to append a new note
6. **Actions footer**: `e` edit | `s` toggle status | `n` add note | `b` set blocker | `o` change owner | `d` schedule / assign phase | `Esc` back

### 4. Task Edit Page (full form)
Bubble Tea form covering the full width of the terminal.

**Fields:**
| Field | Type | Key |
|-------|------|-----|
| Name | text input | Tab |
| Description | textarea | Tab |
| Owner | cycle USER → AI_AGENT → COLLAB (arrow keys or `o`) | Tab |
| Estimated Hours | number input | Tab |
| Scheduled Date | text input (YYYY-MM-DD) | Tab |
| Phase | selection from project's phases (list) | Tab |
| Dependencies | multi-select from project tasks | Tab |

- `Enter` → save via `svc.UpdateTask` (or chaining individual setters)
- `Esc` → cancel → pop back

### 5. Phase Detail Page
- Phase header: Name, Description, StartDate → EndDate
- List of tasks assigned to this phase (subset of TasksPage delegate, no phase headers)
- `e` → push ExitCriteriaPage
- `Esc` → back

### 6. Exit Criteria Page (checklist)
- Checklist items with `[ ]` / `[x]` prefixes
- `Space` / `Enter` → toggle completion via `svc.CompleteExitCriterion`
- `Esc` → back

---

## Implementation Phases

### Phase 1: Navigation Foundation

**[AI AGENT]**
1. Create `internal/tui/navigation/stack.go`
   - `PageStack` struct with `Push(Page)`, `Pop()`, `Current() Page`
   - `Page` interface: `Init() tea.Cmd`, `Update(tea.Msg) (Page, tea.Cmd)`, `View() string`, `Resize(w, h)`
2. Refactor root model (`internal/tui/model.go`)
   - Replace `ActivePane(ProjectPane|TasksPane|TaskDetailPane)` + 3 sub-models with a single `PageStack`
   - Global keys (quit, resize) handled at root; all others forwarded to `stack.Current().Update(msg)`
   - Remove `handleLeft()`/`handleRight()` pane switching
3. Convert existing 3 panes into Pages
   - `ProjectsPage` (wraps existing projects list model)
   - `TasksPage` (wraps existing tasks list model)
   - `TaskDetailPage` (wraps existing taskdetail model)
   - Wire push/pop: ProjectsPage pushes TasksPage; TasksPage pushes TaskDetailPage
   - Keep feature parity only — no new UI features yet, just full-screen pages

**Acceptance:** Build succeeds, existing behavior works in full-screen mode. `Esc`/`h` navigates back. Test via `go build ./cmd/tui` and manual smoke.

---

### Phase 2: Enrich Tasks Page

**[AI AGENT]**
4. Load phases alongside tasks on TasksPage entry
   - Call `svc.GetPhases(ctx, projectID)` and store in page state
   - Build phase lookup map: `phaseId → PhaseItem`
5. Redesign delegate for phase grouping + hierarchy
   - Detect phase changes between tasks; render phase header row when it changes
   - Compute dependency depth per task (graph traversal on `task.Dependencies`)
   - Add indent to delegate render based on depth
   - Show `EstimatedHours` in delegate line
6. Add filter modes
   - State on TasksPage: `FilterMode` enum (`All`, `Today`, `Mine`, `AI`, `Blocked`)
   - `A` / `T` / `M` / `I` / `B` keys toggle filter
   - Re-fetch appropriate data for each filter (or filter `[]models.Task` in-memory for simplicity)

**Acceptance:** Open TasksPage, see task list grouped under phase headers; press filter keys and the list re-renders with only matching tasks. Build + manual smoke.

---

### Phase 3: Rich Task Detail Page

**[AI AGENT]**
7. Add service wrapper `UpdateTask`
   - Wrap `projectsdb.UpdateTask` into `services.UpdateTask(ctx, args models.UpdateTaskArgs) error`
   - (or chain existing individual setter calls — `SetTaskOwner`, `SetTaskSchedule`, etc.)
8. Load auxiliary data on TaskDetailPage entry
   - `GetTaskNotes(taskID)`, `GetUpstreamTasks(taskID)`, `GetDownstreamTasks(taskID)`, `GetPhases(projectID)` to resolve phase name
   - Store in enhanced viewmodel (`TaskDetailView`)
9. Redesign TaskDetailPage layout
   - Top: metadata bar (status color, owner badge, phase name, scheduled, est. hours)
   - Blocker banner if active (red, prominent, with resolve action)
   - Middle: Description (scrollable text)
   - Bottom-left: Notes timeline (compact, last 5-6 notes)
   - Bottom-right: Dependencies (upstream + downstream task names)
10. Wire edit actions to real service calls
    - `taskdetail.saveChanges()` → call `svc.UpdateTask(...)` or chain individual setters
    - `n` key → inline textarea for new note → `svc.AddTaskNote(...)`
    - `b` key → prompt for blocker text → `svc.SetTaskBlocker(...)`
    - `o` key → cycle owner → `svc.SetTaskOwner(...)`
    - `d` key → prompt for date/phase → `svc.SetTaskSchedule(...)`
    - `r` key (when blocked) → `svc.ResolveTaskBlocker(...)`

**Acceptance:** Select a task from TasksPage, see full detail with notes and dependencies. Press `s`, `o`, `b`, `n`, `d` and changes persist to DB (verify via DB inspection or reload). Build + manual smoke.

---

### Phase 4: Phase & Exit Criteria Pages

**[AI AGENT]**
11. PhaseDetailPage
    - Entry from TasksPage: when cursor is on a phase header row, `Enter` pushes PhaseDetailPage
    - Shows phase metadata + tasks in that phase + exit criteria count summary
    - `e` → push ExitCriteriaPage
12. ExitCriteriaPage
    - Custom list delegate showing `[ ]` / `[x]` checkboxes
    - `Space` toggles completion via `svc.CompleteExitCriterion(...)`
    - Re-fetch and re-render on toggle

**Acceptance:** Navigate to a phase, press `e`, see checklist, toggle items with space, they persist in DB. Build + manual smoke.

---

### Phase 5: Task Edit Form Expansion

**[USER] or [COLLAB]**
13. Expand the Task Edit form from name+description to all new fields
    - Owner cycling (arrow keys)
    - Estimated Hours (numeric text input)
    - Scheduled Date (text input, YYYY-MM-DD)
    - Phase selection (list of project phases, arrow keys to pick)
    - Dependencies (multi-select list of project tasks)
    - This is UI-heavy and involves Bubble Tea form patterns; likely needs user review of interaction style

**Acceptance:** Press `e` on a task, edit all fields, save, reload, fields persisted.

---

## File Changes

| File | Action | Description |
|------|--------|-------------|
| `internal/tui/navigation/stack.go` | **NEW** | Page stack abstraction |
| `internal/tui/pages/projects/page.go` | **NEW** | Full-screen Projects page |
| `internal/tui/pages/tasks/page.go` | **NEW** | Full-screen Tasks page (with phase grouping + filters) |
| `internal/tui/pages/tasks/delegate.go` | **NEW** | Task list delegate (indent, columns, icons) |
| `internal/tui/pages/taskdetail/page.go` | **NEW** | Full-screen Task Detail (notes, deps, actions) |
| `internal/tui/pages/taskedit/page.go` | **NEW** | Task edit form (expanded fields) |
| `internal/tui/pages/phase/page.go` | **NEW** | Phase detail page |
| `internal/tui/pages/exitcriteria/page.go` | **NEW** | Exit criteria checklist page |
| `internal/tui/model.go` | **REFACTOR** | Root model uses stack, removes 3-pane layout |
| `internal/tui/app.go` | **REFACTOR** | Init returns stack-based model |
| `internal/tui/styles.go` | **MODIFY** | Add page-level styles, remove 3-pane styles |
| `internal/tui/constants/keybindings.go` | **MODIFY** | Add filter keys, edit keys, back keys |
| `internal/tui/constants/messages.go` | **MODIFY** | Add page-transition messages if needed |
| `internal/tui/viewmodels/viewmodels.go` | **MODIFY** | Add `TaskDetailView` with notes/deps/phase name |
| `internal/services/services.go` | **ADD** | `UpdateTask` wrapper around `projectsdb.UpdateTask` |
| `internal/tui/components/*` | **DELETE** | Old 3-pane component models (superseded by pages) |

---

## Open Questions

1. **Entry point for Task Edit**: Inline within TaskDetailPage (expandable section) or push a dedicated TaskEditPage? Bubble Tea forms are cleaner as dedicated pages.
2. **Dependency multi-select**: Bubble Tea's `list` component supports single-select out of the box; multi-select needs custom behavior or a `github.com/charmbracelet/bubbles` multi-select pattern (maybe a checklist-style list). Phase 5 should spike this before committing.
3. **Notes scroll vs paginate**: Should the Notes timeline in TaskDetailPage scroll infinitely, or show last N with a "more..." prompt? Recommend scrollable for simplicity.
4. **Phase header as selectable row**: In TasksPage, should phase headers be cursor-selectable rows that push PhaseDetailPage, or just decorative headers? Making them selectable is natural — pressing Enter on a phase header opens its detail page.
