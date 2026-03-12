package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}

	migrationsDir := flag.String("dir", "migrations", "path to migrations directory")
	flag.Parse()

	command := "up"
	if flag.NArg() > 0 {
		command = strings.ToLower(strings.TrimSpace(flag.Arg(0)))
	}

	if err := run(command, *migrationsDir, databaseURL); err != nil {
		fmt.Fprintf(os.Stderr, "migration failed: %v\n", err)
		os.Exit(1)
	}
}

func run(command string, migrationsDir string, databaseURL string) error {
	goose.SetBaseFS(nil)

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	absDir, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("resolve migrations dir: %w", err)
	}

	switch command {
	case "up":
		return goose.Up(db, absDir)
	case "status":
		return goose.Status(db, absDir)
	case "down":
		return goose.Down(db, absDir)
	default:
		return fmt.Errorf("unsupported command %q (use up, down, or status)", command)
	}
}
