CREATE TABLE projects (
  project_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);
	
CREATE TABLE tasks (
  task_id INTEGER PRIMARY KEY AUTOINCREMENT,
  description TEXT NOT NULL,
  project_id INTEGER NOT NULL,
  sort INTEGER NOT NULL,
  is_completed INTEGER NOT NULL,
  is_failed INTEGER NOT NULL,
  notes TEXT NULL,
  FOREIGN KEY(project_id) REFERENCES projects(project_id)
);
