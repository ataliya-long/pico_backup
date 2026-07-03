package backup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Postgres struct{}

func (p *Postgres) Dump(opts Options) (string, error) {
	name := backupName(opts.Database, opts.Host, "sql")
	outPath := filepath.Join(opts.Dir, name)

	args := []string{
		"-h", opts.Host,
		"-p", fmt.Sprintf("%d", opts.Port),
		"-U", opts.User,
		"-d", opts.Database,
		"--no-password",
		"-f", outPath,
	}

	cmd := exec.Command("pg_dump", args...)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+opts.Password)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("[postgres] pg_dump -h %s -p %d -U %s -d %s ...\n",
		opts.Host, opts.Port, opts.User, opts.Database)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("pg_dump: %w", err)
	}
	fmt.Printf("[postgres] dump -> %s\n", outPath)
	return outPath, nil
}