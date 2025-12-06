CREATE TABLE IF NOT EXISTS ssh_configs (
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    name      TEXT NOT NULL,
    ip        TEXT NOT NULL,
    user      TEXT,
    password  TEXT,
    port      TEXT NOT NULL,
    key_path  TEXT,
    desc      TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);