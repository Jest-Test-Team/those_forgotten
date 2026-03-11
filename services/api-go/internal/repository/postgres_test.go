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

func TestListCommunityReportsReadsFromPostgres(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	rows := pgxmock.NewRows([]string{"id", "post_id", "title", "office", "reason", "status", "created_at"}).
		AddRow(
			"report-001",
			"post-001",
			"臺北關相機批次看貨紀錄",
			"臺北關",
			"缺少看貨照片佐證",
			"pending",
			"2026-03-11T10:00:00Z",
		)

	mock.ExpectQuery("SELECT cr.id::text, cr.post_id::text, cp.title, cp.office, cr.reason, cr.status, cr.created_at::text").
		WillReturnRows(rows)

	result := repo.ListCommunityReports()

	if len(result) != 1 {
		t.Fatalf("len(result) = %d", len(result))
	}
	if result[0].PostTitle != "臺北關相機批次看貨紀錄" {
		t.Fatalf("result[0].PostTitle = %q", result[0].PostTitle)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet() error = %v", err)
	}
}

func TestListAdvisorLeadsReadsFromPostgres(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("pgxmock.NewPool() error = %v", err)
	}
	defer mock.Close()

	repo := NewPostgresRepositoryWithPool(mock)

	rows := pgxmock.NewRows([]string{"id", "advisor_id", "advisor_name", "name", "email", "message", "category", "created_at"}).
		AddRow(
			"lead-001",
			"advisor-001",
			"王顧問",
			"示例會員",
			"member@example.com",
			"需要協助驗車與提領安排。",
			"進口車驗車",
			"2026-03-11T10:30:00Z",
		)

	mock.ExpectQuery("SELECT al.id::text, al.advisor_id::text, ap.name, al.name, al.email, al.message, COALESCE\\(al.category, ''\\), al.created_at::text").
		WillReturnRows(rows)

	result := repo.ListAdvisorLeads()

	if len(result) != 1 {
		t.Fatalf("len(result) = %d", len(result))
	}
	if result[0].AdvisorName != "王顧問" {
		t.Fatalf("result[0].AdvisorName = %q", result[0].AdvisorName)
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
