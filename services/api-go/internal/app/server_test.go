package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAuctions(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auctions", nil)
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestHealthzIncludesRepositoryMode(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}

	payload := map[string]string{}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload["repository"] == "" {
		t.Fatalf("expected repository mode in payload")
	}
}

func TestIngestAuctionsRequiresToken(t *testing.T) {
	t.Setenv("INTERNAL_INGEST_TOKEN", "secret")

	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	body := []byte(`{"source":"taipei","checksum":"abc","rows":[]}`)
	req := httptest.NewRequest(http.MethodPost, "/internal/ingest/auctions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestCalendarFeedRequiresToken(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auctions/calendar.ics", nil)
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestCalendarFeedWithToken(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/auctions/calendar.ics?token=demo", nil)
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestCreateKeywordSubscriptionValidatesKeyword(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/keyword-subscriptions", bytes.NewBufferString(`{"keyword":" "}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestIngestAuctionsValidatesChecksum(t *testing.T) {
	t.Setenv("INTERNAL_INGEST_TOKEN", "secret")

	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/internal/ingest/auctions", bytes.NewBufferString(`{"source":"fixtures","checksum":"abc","rows":[{"announcement_no":"TP-001","office":"臺北關","title":"相機","category":"3C","closing_at":"2026-03-16T14:00:00+08:00","original_link":"https://example.com","warnings":["現狀交付"]}]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Ingest-Token", "secret")
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestCreateWebPushSubscription(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/web-push-subscriptions", bytes.NewBufferString(`{"endpoint":"https://push.example.dev/demo","keys":{"p256dh":"demo","auth":"secret"}}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestCreateCommunityPost(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/community/posts", bytes.NewBufferString(`{"title":"看貨紀錄","body":"鏡頭有刮痕","office":"臺北關"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}

func TestCreateAdvisorLead(t *testing.T) {
	server, err := NewServer()
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/advisor-leads", bytes.NewBufferString(`{"advisor_id":"advisor-001","name":"會員","email":"member@example.com","message":"需要代標協助"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	server.Echo().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("ServeHTTP() code = %d", rec.Code)
	}
}
