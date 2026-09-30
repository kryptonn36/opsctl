package backup

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"os/exec"
)

func Restore(ctx context.Context, cfg Config, backupPath string) error{
	file, err := os.Open(backupPath)
	if err != nil{
		return fmt.Errorf("Unable to open file: %v",err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil{
		return fmt.Errorf("Error in reading: %v", err)
	}
	defer gzipReader.Close()

	args := []string{
		"--host", cfg.Host,
		"--port", cfg.Port,
		"--user", cfg.User,
		"--dbname", cfg.DBname,
	}

	psql := exec.CommandContext(ctx, "psql", args...)

	psql.Stdin = gzipReader

	var stderr bytes.Buffer
	psql.Stderr = &stderr

	err = psql.Run()
	if err != nil{
		if ctx.Err() != nil{
			return fmt.Errorf("Restore Cancelled")
		}
		return fmt.Errorf("Error in restore: %v", err)
	}
	fmt.Printf("Database restored successfully from path: %s", backupPath)
	return nil
}