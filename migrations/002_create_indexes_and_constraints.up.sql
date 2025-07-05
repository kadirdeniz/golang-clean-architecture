CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);
CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
CREATE INDEX IF NOT EXISTS idx_todos_updated_at ON todos(updated_at);
CREATE INDEX IF NOT EXISTS idx_todos_priority ON todos(priority);
CREATE INDEX IF NOT EXISTS idx_todos_due_date ON todos(due_date);

CREATE INDEX IF NOT EXISTS idx_todos_completed_created_at ON todos(completed, created_at);
CREATE INDEX IF NOT EXISTS idx_todos_priority_due_date ON todos(priority, due_date);
CREATE INDEX IF NOT EXISTS idx_todos_completed_priority ON todos(completed, priority);

CREATE INDEX IF NOT EXISTS idx_todos_active ON todos(created_at) WHERE completed = FALSE; 