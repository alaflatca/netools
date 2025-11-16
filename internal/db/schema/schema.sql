CREATE TABLE IF NOT EXISTS ssh_configs (
    name TEXT NOT NULL,
    ip TEXT NOT NULL,
    port TEXT NOT NULL,
    key_path TEXT NOT NULL,
    desc TEXT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
)