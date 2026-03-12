package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/repository"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	repo, err := repository.NewPostgresRepository(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer repo.Close()

	batchSize := 20
	if value := os.Getenv("NOTIFICATION_BATCH_SIZE"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed > 0 {
			batchSize = parsed
		}
	}

	jobs := repo.ClaimPendingNotificationJobs(batchSize)
	for _, job := range jobs {
		if err := repo.MarkNotificationJobDelivered(job.ID); err != nil {
			log.Printf("mark delivered %s: %v", job.ID, err)
			continue
		}
		log.Printf("notification delivered job=%s endpoint=%s keyword=%s mode=simulate", job.ID, job.Endpoint, job.Keyword)
	}

	log.Printf("processed %d notification jobs", len(jobs))
}
