package db

import (
	"context"
	"fmt"
	"log"
	"time"
)

type SSHConfig struct {
	Name      string
	IP        string
	Port      string
	KeyPath   string
	Desc      string
	CreatedAt time.Time
}

func InsertSSHConfig(cfg SSHConfig) error {
	log.Printf("insert cfg: %+v\n", cfg)
	baseSQL := `INSERT INTO ssh_configs (name, ip, port, key_path, desc) VALUES (?, ?, ?, ?, ?)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := Get().ExecContext(ctx, baseSQL, cfg.Name, cfg.IP, cfg.Port, cfg.KeyPath, cfg.Desc)
	if err != nil {
		return fmt.Errorf("[SSH] failed to insert config: %v", err)
	}
	return nil
}

func SelectSSHConfigs() ([]SSHConfig, error) {
	sqlText := `SELECT name, ip, port, key_path, desc FROM ssh_configs ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := Get().QueryContext(ctx, sqlText)
	if err != nil {
		return nil, err
	}

	sshConfigs := []SSHConfig{}
	for rows.Next() {
		var sshConfig SSHConfig
		err := rows.Scan(&sshConfig.Name, &sshConfig.IP, &sshConfig.Port, &sshConfig.KeyPath, &sshConfig.Desc)
		if err != nil {
			return nil, err
		}
		sshConfigs = append(sshConfigs, sshConfig)
	}

	return sshConfigs, nil
}
