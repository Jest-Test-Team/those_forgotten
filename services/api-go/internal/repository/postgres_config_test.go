package repository

import "testing"

func TestUsesSupabaseTransactionPooler(t *testing.T) {
	t.Run("detects supabase pooler", func(t *testing.T) {
		databaseURL := "postgresql://postgres.mluxmwdbjunrqgyuqizn:secret@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?sslmode=require"
		if !usesSupabaseTransactionPooler(databaseURL) {
			t.Fatalf("expected pooler detection to be true")
		}
	})

	t.Run("ignores direct postgres host", func(t *testing.T) {
		databaseURL := "postgresql://postgres:secret@db.example.com:5432/postgres?sslmode=require"
		if usesSupabaseTransactionPooler(databaseURL) {
			t.Fatalf("expected pooler detection to be false")
		}
	})
}
