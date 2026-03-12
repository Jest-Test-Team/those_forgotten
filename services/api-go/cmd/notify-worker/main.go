package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/model"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/repository"
)

type deliveryConfig struct {
	publicKey  string
	privateKey string
	subject    string
}

func (c deliveryConfig) live() bool {
	return c.publicKey != "" && c.privateKey != "" && c.subject != ""
}

func deliver(job model.NotificationJob, cfg deliveryConfig) error {
	subscription := &webpush.Subscription{
		Endpoint: job.Endpoint,
		Keys: webpush.Keys{
			P256dh: job.P256DH,
			Auth:   job.AuthSecret,
		},
	}

	resp, err := webpush.SendNotification([]byte(job.Payload), subscription, &webpush.Options{
		Subscriber:      cfg.subject,
		VAPIDPublicKey:  cfg.publicKey,
		VAPIDPrivateKey: cfg.privateKey,
		TTL:             30,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unexpected push status: %d", resp.StatusCode)
	}

	return nil
}

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

	cfg := deliveryConfig{
		publicKey:  os.Getenv("VAPID_PUBLIC_KEY"),
		privateKey: os.Getenv("VAPID_PRIVATE_KEY"),
		subject:    os.Getenv("VAPID_SUBJECT"),
	}

	jobs := repo.ClaimPendingNotificationJobs(batchSize)
	for _, job := range jobs {
		if cfg.live() {
			if err := deliver(job, cfg); err != nil {
				if markErr := repo.MarkNotificationJobFailed(job.ID, err.Error()); markErr != nil {
					log.Printf("mark failed %s: %v", job.ID, markErr)
				}
				log.Printf("notification failed job=%s endpoint=%s keyword=%s err=%v", job.ID, job.Endpoint, job.Keyword, err)
				continue
			}
		}
		if err := repo.MarkNotificationJobDelivered(job.ID); err != nil {
			log.Printf("mark delivered %s: %v", job.ID, err)
			continue
		}
		mode := "simulate"
		if cfg.live() {
			mode = "vapid"
		}
		log.Printf("notification delivered job=%s endpoint=%s keyword=%s mode=%s", job.ID, job.Endpoint, job.Keyword, mode)
	}

	log.Printf("processed %d notification jobs", len(jobs))
}
