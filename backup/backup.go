package backup

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Options struct {
	DBType   string
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Dir      string
}

type Backuper interface {
	Dump(opts Options) (dumpPath string, err error)
}

func New(dbType string) (Backuper, error) {
	if dbType != "postgres" {
		return nil, fmt.Errorf("unsupported db type: %s", dbType)
	}
	return &Postgres{}, nil
}

func Run(b Backuper, opts Options) (string, error) {
	dumpPath, err := b.Dump(opts)
	if err != nil {
		return "", err
	}
	return gzipFile(dumpPath)
}

func backupName(database, host, ext string) string {
	ts := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s_%s_%s.%s", database, host, ts, ext)
}

func gzipFile(src string) (string, error) {
	dst := src + ".gz"
	in, err := os.Open(src)
	if err != nil {
		return "", fmt.Errorf("open for gzip: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", fmt.Errorf("create gz: %w", err)
	}
	defer out.Close()

	gw := gzip.NewWriter(out)
	defer gw.Close()
	gw.Name = filepath.Base(src)

	if _, err := io.Copy(gw, in); err != nil {
		os.Remove(dst)
		return "", fmt.Errorf("gzip: %w", err)
	}
	gw.Close()
	out.Close()
	os.Remove(src)
	return dst, nil
}
