package main

import (
	"flag"
	"fmt"
	"os"
	"db-backup/backup"
	"db-backup/config"
)

func main() {
	// ---- cli flags ----
	cfgPath := flag.String("config", "config.yaml", "YAML config file")
	dbType := flag.String("db", "postgres", "database type (postgres)")
	host := flag.String("host", "", "database host")
	port := flag.Int("port", 0, "database port")
	user := flag.String("user", "", "database user")
	passwd := flag.String("passwd", "", "database password")
	database := flag.String("database", "", "database name")
	outDir := flag.String("dir", "", "output directory")
	flag.Parse()

	if *dbType == "" {
		fmt.Fprintln(os.Stderr, "Usage: backup --db postgres --host <h> --user <u> --database <db> --passwd <p> [--port <p>] [--config <f>] [--dir <d>]")
		os.Exit(1)
	}

	// ---- load config ----
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config warning: %v (using built-in defaults)\n", err)
		cfg = config.Default()
	}

	// ---- merge defaults ----
	defaults := cfg.Postgres

	h := pick(*host, defaults.Host)
	p := pickInt(*port, defaults.Port, 5432)
	u := pick(*user, defaults.User)
	pw := pick(*passwd, defaults.Password)
	d := *database
	dir := *outDir; if dir == "" { dir = cfg.BackupDir }; if dir == "" { dir = "." }

	if h == "" || u == "" || d == "" || pw == "" {
		fmt.Fprintln(os.Stderr, "Error: --host, --user, --database, --passwd are required")
		os.Exit(1)
	}

	// ---- run backup ----
	b, err := backup.New(*dbType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	opts := backup.Options{
		DBType:   *dbType,
		Host:     h,
		Port:     p,
		User:     u,
		Password: pw,
		Database: d,
		Dir:      dir,
	}

	fmt.Printf("=== starting full backup: %s/%s@%s:%d ===\n", u, d, h, p)
	gzPath, err := backup.Run(b, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "backup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("=== done: %s ===\n", gzPath)
}

func pick(val, fallback string) string {
	if val != "" {
		return val
	}
	return fallback
}

func pickInt(val, fallback1, fallback2 int) int {
	if val != 0 {
		return val
	}
	if fallback1 != 0 {
		return fallback1
	}
	return fallback2
}
