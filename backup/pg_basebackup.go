package backup

import (
	"fmt"
	"os"
	"path/filepath"
)

type PostgresBaseBackup struct{}

func (p *PostgresBaseBackup) BaseBackup(opts basebackup_options) (string, error) {
	name := backupName(opts.Database, opts.Host, "tar")
	outPath := filepath.Join(opts.BackupDir, name)

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}

    

	return outPath, nil
}