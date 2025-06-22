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

func InsertSSHConfig(ctx context.Context, db DB, sshCfg SSHConfig) error {
	sqlText := `INSERT INTO ssh_configs (name, keypath) VALUES (?, ?)`
	_, err := db.ExecContext(ctx, sqlText, sshCfg.Name, sshCfg.KeyPath)
	if err != nil {
		return fmt.Errorf("[SSH] failed to insert sshconfig: %v", err)
	}
	return nil
}

func SelectSSHConfigs(ctx context.Context, db DB) ([]SSHConfig, error) {
	sqlText := `SELECT name, keypath FROM ssh_configs ORDER BY created_at DESC`
	rows, err := db.QueryContext(ctx, sqlText)
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
