// db/sqlite_ssh.go 같은 파일

package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SSHConfig struct {
	ID        int64
	Name      string
	IP        string
	User      string
	Password  string
	Port      string
	KeyPath   string
	Desc      string
	CreatedAt time.Time
}

func InsertSSHConfig(ctx context.Context, db *DB, cfg SSHConfig) (int64, error) {
	baseSQL := `
        INSERT INTO ssh_configs (name, ip, user, password, port, key_path, desc)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, baseSQL,
		cfg.Name,
		cfg.IP,
		cfg.User,
		cfg.Password,
		cfg.Port,
		cfg.KeyPath,
		cfg.Desc,
	)
	if err != nil {
		return 0, fmt.Errorf("[SSH] failed to insert config: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("[SSH] failed to get last insert id: %w", err)
	}
	return id, nil
}

func SelectSSHConfigs(ctx context.Context, db *DB) ([]SSHConfig, error) {
	sqlText := `
        SELECT id, name, ip, user, password, port, key_path, desc, created_at
        FROM ssh_configs
        ORDER BY created_at DESC
    `

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := db.conn.QueryContext(ctx, sqlText)
	if err != nil {
		return nil, fmt.Errorf("[SSH] failed to query configs: %w", err)
	}
	defer rows.Close()

	var sshConfigs []SSHConfig
	for rows.Next() {
		var sshConfig SSHConfig
		if err := rows.Scan(
			&sshConfig.ID,
			&sshConfig.Name,
			&sshConfig.IP,
			&sshConfig.User,
			&sshConfig.Password,
			&sshConfig.Port,
			&sshConfig.KeyPath,
			&sshConfig.Desc,
			&sshConfig.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("[SSH] failed to scan row: %w", err)
		}
		sshConfigs = append(sshConfigs, sshConfig)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[SSH] rows error: %w", err)
	}

	return sshConfigs, nil
}

func UpdateSSHConfig(ctx context.Context, db *DB, cfg SSHConfig) error {
	if cfg.ID == 0 {
		return fmt.Errorf("[SSH] UpdateSSHConfig: missing ID")
	}

	sqlText := `
        UPDATE ssh_configs
        SET name = ?, ip = ?, user = ?, password = ?, port = ?, key_path = ?, desc = ?
        WHERE id = ?
    `

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, sqlText,
		cfg.Name,
		cfg.IP,
		cfg.User,
		cfg.Password,
		cfg.Port,
		cfg.KeyPath,
		cfg.Desc,
		cfg.ID,
	)
	if err != nil {
		return fmt.Errorf("[SSH] failed to update config: %w", err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteSSHConfig(ctx context.Context, db *DB, id int64) error {
	sqlText := `DELETE FROM ssh_configs WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := db.conn.ExecContext(ctx, sqlText, id)
	if err != nil {
		return fmt.Errorf("[SSH] failed to delete config: %w", err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return sql.ErrNoRows
	}
	return nil
}
