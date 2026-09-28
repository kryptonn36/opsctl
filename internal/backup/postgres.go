package backup

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

)

type Config struct{
	Host string
	Port string
	User string
	DBname string
	Output_dir string
}

func resolve(value string, envName string, defaultvalue string) string{
	if value != ""{
		return value
	}
	env := os.Getenv(envName)
	if env != ""{
		return env
	}
	return defaultvalue
}

func SetConfig(cfg Config) Config{
	cfg.Host = resolve(cfg.Host, "PGHOST", "localhost")
	cfg.Port = resolve(cfg.Port, "PGPORT", "5432")
	cfg.User = resolve(cfg.User, "PGUSER", "")
	cfg.DBname = resolve(cfg.DBname, "PGDATABASE", "")
	return cfg
}

func Run(ctx context.Context, cfg Config) error{
	cfg = SetConfig(cfg)
	if cfg.DBname == ""{
		return fmt.Errorf("Database name is required")
	}
	err := os.MkdirAll(cfg.Output_dir, 0755)
	if err != nil{
		return err
	}

	timestamp := time.Now().Format("20260102_031425")
	fileName := fmt.Sprintf("backup-%s-%s.sql.gz", cfg.DBname, timestamp)
	outputPath := filepath.Join(cfg.Output_dir, fileName)

	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil{
		return fmt.Errorf("create backup file: %w", err)
	}

	success := true

	defer func ()  {
		file.Close()

		if !success{
			os.Remove(outputPath)
		}
	}()

	gzipWriter := gzip.NewWriter(file)

	args := []string{
		"--host", cfg.Host,
		"--port", cfg.Port,
		"--user", cfg.User,
		"--dbname", cfg.DBname,
	}

	pgDump := exec.CommandContext(ctx, "pg_dump", args...)
	
	var stderr bytes.Buffer
	pgDump.Stderr = &stderr
	pgDump.Stdout = gzipWriter

	err = pgDump.Run()

	gziperr := gzipWriter.Close()

	if err != nil{
		if ctx.Err() != nil{
			return fmt.Errorf("backup cancled")
		}
		return fmt.Errorf("pg_dump failed: %w", err)
	}

	if gziperr != nil{
		return fmt.Errorf("gzip failed: %w", gziperr)
	}

	success = true

	fmt.Printf("Backup created: %w", outputPath)
	return nil
}