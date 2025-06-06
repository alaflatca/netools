package ssh

import (
	"context"
	"fmt"
	"sshtn/internal/db"
)

type SSHStore struct {
	db db.DB
}

func NewSSHStore(db db.DB) *SSHStore {
	return &SSHStore{db: db}
}

type SSHConfig struct {
	Name    string
	KeyPath string
}

func (store *SSHStore) SaveConfig(ctx context.Context, cfg SSHConfig) error {
	sqlText := `
	INSERT INTO ssh_configs (name, key_path) VALUES (?, ?)
	`
	if _, err := store.db.ExecContext(ctx, sqlText, cfg.Name, cfg.KeyPath); err != nil {
		return fmt.Errorf("[SSH-SaveConfig] failed to exec db: %v", err)
	}

	return nil
}
