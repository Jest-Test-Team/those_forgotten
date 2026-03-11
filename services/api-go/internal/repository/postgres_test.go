package repository

import (
	"testing"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	pgxmock "github.com/pashagolub/pgxmock/v4"
)

func TestCreateKeywordSubscriptionPersistsToPostgres(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	mock.ExpectExec("INSERT INTO profiles").
		WithArgs(demoProfileID, "demo@customs.local", "Demo Member").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO keyword_subscriptions").
		WithArgs(pgxmock.AnyArg(), demoProfileID, "相機").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	result := repo.CreateKeywordSubscription("相機")

	if result.Keyword != "相機" {
		t.Fatalf("result.Keyword = %q", result.Keyword)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() error = %v", err)
	}
}

func TestCreateWebPushSubscriptionPersistsToPostgres(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	mock.ExpectExec("INSERT INTO profiles").
		WithArgs(demoProfileID, "demo@customs.local", "Demo Member").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO web_push_subscriptions").
		WithArgs(pgxmock.AnyArg(), demoProfileID, "https://push.example.dev/demo", "demo-key", "secret-key").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	result := repo.CreateWebPushSubscription(&dto.WebPushSubscriptionInput{
		Endpoint: "https://push.example.dev/demo",
		Keys: map[string]string{
			"p256dh": "demo-key",
			"auth":   "secret-key",
		},
	})

	if result["endpoint"] != "https://push.example.dev/demo" {
		t.Fatalf("result[endpoint] = %v", result["endpoint"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() error = %v", err)
	}
}

func TestReportCommunityPostPersistsToPostgres(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	mock.ExpectExec("INSERT INTO community_reports").
		WithArgs(pgxmock.AnyArg(), "11111111-1111-1111-1111-111111111111", "疑似內容不實").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	result := repo.ReportCommunityPost("11111111-1111-1111-1111-111111111111", &dto.ReportInput{
		Reason: "疑似內容不實",
	})

	if result.Status != "pending" {
		t.Fatalf("result.Status = %q", result.Status)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() error = %v", err)
	}
}
