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
  is_completed INTEGER NOT NULL,
  is_in_progress INTEGER NOT NULL,
  notes TEXT NULL,
  FOREIGN KEY(project_id) REFERENCES projects(project_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS dependencies (
  task_id INTEGER NOT NULL,
  depends_on INTEGER NOT NULL,
  PRIMARY KEY (task_id, depends_on),
  FOREIGN KEY(task_id) REFERENCES tasks(task_id) ON DELETE CASCADE,
  FOREIGN KEY(depends_on) REFERENCES tasks(task_id) ON DELETE CASCADE
);
