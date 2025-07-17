CREATE TABLE IF NOT EXISTS projects (
  project_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE
);
	
CREATE TABLE IF NOT EXISTS tasks (
  task_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  project_id INTEGER NOT NULL,
  sort INTEGER NOT NULL,
  dependencies_json TEXT NOT NULL,
  is_completed INTEGER NOT NULL,
  is_in_progress INTEGER NOT NULL,
  notes TEXT NULL,
  FOREIGN KEY(project_id) REFERENCES projects(project_id)
);
