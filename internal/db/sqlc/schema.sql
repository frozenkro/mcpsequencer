CREATE TABLE IF NOT EXISTS projects (
  project_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  description TEXT NULL,
  absolute_path TEXT NULL
);
	
CREATE TABLE IF NOT EXISTS tasks (
  task_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  project_id INTEGER NOT NULL,
  sort INTEGER NOT NULL,
  is_completed INTEGER NOT NULL DEFAULT 0,
  is_in_progress INTEGER NOT NULL DEFAULT 0,
  notes TEXT NULL,
  owner TEXT NOT NULL DEFAULT 'USER',
  scheduled_date TEXT NULL,
  phase_id INTEGER NULL,
  blocker_text TEXT NULL,
  blocked_at TEXT NULL,
  estimated_hours INTEGER NULL,
  FOREIGN KEY(project_id) REFERENCES projects(project_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS dependencies (
  task_id INTEGER NOT NULL,
  depends_on INTEGER NOT NULL,
  PRIMARY KEY (task_id, depends_on),
  FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
  FOREIGN KEY(depends_on) REFERENCES tasks(task_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS phases (
  phase_id INTEGER PRIMARY KEY AUTOINCREMENT,
  project_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  description TEXT NULL,
  start_date TEXT NOT NULL,
  end_date TEXT NOT NULL,
  sort INTEGER NOT NULL,
  FOREIGN KEY(project_id) REFERENCES projects(project_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS task_notes (
  note_id INTEGER PRIMARY KEY AUTOINCREMENT,
  task_id INTEGER NOT NULL,
  note TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS exit_criteria (
  criterion_id INTEGER PRIMARY KEY AUTOINCREMENT,
  phase_id INTEGER NOT NULL,
  description TEXT NOT NULL,
  is_completed INTEGER NOT NULL DEFAULT 0,
  sort INTEGER NOT NULL,
  FOREIGN KEY(phase_id) REFERENCES phases(phase_id) ON DELETE CASCADE
);
