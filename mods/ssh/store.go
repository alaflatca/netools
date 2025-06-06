package ssh

import (
	"context"
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
	// sqlText := `
	// INSERT INTO
	// `
	// store.db.ExecContext()
	// store.db.ExecContext(ctx, sqlText, )

}
