package db

import (
	"context"
	"fmt"
	"time"
)

type SSHConfig struct {
	Name      string
	KeyPath   string
	CreatedAt time.Time
}

func InsertSSHConfig(cfg SSHConfig) error {
	sqlText := `INSERT INTO ssh_configs (name, key_path) VALUES (?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := Get().ExecContext(ctx, sqlText, cfg.Name, cfg.KeyPath)
	if err != nil {
		return fmt.Errorf("[SSH] failed to insert config: %v", err)
	}
	return nil
}

func SelectSSHConfigs() ([]SSHConfig, error) {
	sqlText := `SELECT name, key_path FROM ssh_configs ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := Get().QueryContext(ctx, sqlText)
	if err != nil {
		return nil, err
	}

	sshConfigs := []SSHConfig{}
	for rows.Next() {
		var sshConfig SSHConfig
		err := rows.Scan(&sshConfig.Name, &sshConfig.KeyPath)
		if err != nil {
			return nil, err
		}

		sshConfigs = append(sshConfigs, sshConfig)
	}

	return sshConfigs, nil
}
