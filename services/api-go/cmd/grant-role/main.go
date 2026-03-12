package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	email := flag.String("email", "", "user email")
	role := flag.String("role", "admin", "role to grant")
	name := flag.String("name", "", "optional full name")
	flag.Parse()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(*email))
	normalizedRole := strings.TrimSpace(strings.ToLower(*role))
	if normalizedEmail == "" {
		log.Fatal("--email is required")
	}
	if normalizedRole == "" {
		log.Fatal("--role is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pool.Close()

	profileID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("profile:"+normalizedEmail))
	roleID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("role:"+normalizedEmail+":"+normalizedRole))

	tx, err := pool.Begin(ctx)
	if err != nil {
		log.Fatalf("begin tx: %v", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO profiles (id, email, full_name, created_at, updated_at)
		VALUES ($1, $2, NULLIF($3, ''), NOW(), NOW())
		ON CONFLICT (email) DO UPDATE
		SET updated_at = NOW(),
		    full_name = COALESCE(NULLIF(EXCLUDED.full_name, ''), profiles.full_name)
	`, profileID, normalizedEmail, strings.TrimSpace(*name))
	if err != nil {
		log.Fatalf("upsert profile: %v", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO user_roles (id, profile_id, role, created_at)
		SELECT $1, $2, $3, NOW()
		WHERE NOT EXISTS (
			SELECT 1
			FROM user_roles
			WHERE profile_id = $2
			  AND role = $3
		)
	`, roleID, profileID, normalizedRole)
	if err != nil {
		log.Fatalf("grant role: %v", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("commit tx: %v", err)
	}

	fmt.Printf("granted role %q to %s\n", normalizedRole, normalizedEmail)
}
