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

func TestCreateAdvisorLeadPersistsToPostgres(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	mock.ExpectExec("INSERT INTO advisor_leads").
		WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "會員", "member@example.com", "需要代標協助", "代標需求").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	result := repo.CreateAdvisorLead(&dto.AdvisorLeadInput{
		AdvisorID: "11111111-1111-1111-1111-111111111111",
		Name:      "會員",
		Email:     "member@example.com",
		Message:   "需要代標協助",
		Category:  "代標需求",
	})

	if result.Email != "member@example.com" {
		t.Fatalf("result.Email = %q", result.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() error = %v", err)
	}
}

func TestIngestAuctionsPersistsAnnouncementAndLot(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO auction_announcements").
		WithArgs(
			pgxmock.AnyArg(),
			"臺北關",
			"TP-001",
			"沒入數位相機與鏡頭一批",
			"https://example.com/tp-001",
			"2026-03-16T14:00:00+08:00",
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectExec("INSERT INTO auction_lots").
		WithArgs(
			pgxmock.AnyArg(),
			pgxmock.AnyArg(),
			"沒入數位相機與鏡頭一批",
			"3C",
			"2026-03-16T14:00:00+08:00",
			[]string{"現狀交付", "不負瑕疵擔保"},
		).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	result := repo.IngestAuctions(&dto.IngestPayload{
		Source:   "fixtures",
		Checksum: "demo-checksum",
		Rows: []dto.NormalizedAuction{
			{
				AnnouncementNo: "TP-001",
				Office:         "臺北關",
				Title:          "沒入數位相機與鏡頭一批",
				Category:       "3C",
				ClosingAt:      "2026-03-16T14:00:00+08:00",
				OriginalLink:   "https://example.com/tp-001",
				Warnings:       []string{"現狀交付", "不負瑕疵擔保"},
			},
		},
	})

	if result["received"] != 1 {
		t.Fatalf("result[received] = %v", result["received"])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() error = %v", err)
	}
}
