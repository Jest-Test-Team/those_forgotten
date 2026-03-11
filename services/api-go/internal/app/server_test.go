package app

import (
	"bytes"
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
