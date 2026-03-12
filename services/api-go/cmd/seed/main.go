package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if databaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}

	ctx := context.Background()
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse database url: %v\n", err)
		os.Exit(1)
	}

	if usesSupabaseTransactionPooler(databaseURL) {
		config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "create pool: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "begin tx: %v\n", err)
		os.Exit(1)
	}
	defer tx.Rollback(ctx)

	adminProfileID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("seed:admin@example.com"))
	memberProfileID := uuid.NewSHA1(uuid.NameSpaceURL, []byte("seed:member@example.com"))
	announcementID := uuid.MustParse("6c21dc2c-9754-481e-a078-401d124cc2e1")
	lotID := uuid.MustParse("20097b79-9ff5-4cb3-a69c-ef90a0784f11")
	advisorID := uuid.MustParse("f4744a04-fe4a-4ff5-856a-65ff295de2f7")
	courseID := uuid.MustParse("307e0cbb-b524-47fc-b79f-a26a3b46d711")
	postID := uuid.MustParse("1d5eef09-c982-46c0-b48d-a43f1482a1dd")

	mustExec(ctx, tx, `
		INSERT INTO profiles (id, email, full_name, created_at, updated_at)
		VALUES
			($1, 'admin@example.com', 'Admin Reviewer', NOW(), NOW()),
			($2, 'member@example.com', 'Seed Member', NOW(), NOW())
		ON CONFLICT (email) DO UPDATE
		SET full_name = EXCLUDED.full_name,
		    updated_at = NOW()
	`, adminProfileID, memberProfileID)

	mustExec(ctx, tx, `
		INSERT INTO user_roles (id, profile_id, role, created_at)
		VALUES ($1, $2, 'admin', NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), adminProfileID)

	mustExec(ctx, tx, `
		INSERT INTO memberships (id, profile_id, plan_code, status, renews_at, created_at)
		VALUES ($1, $2, 'pro-monthly', 'active', NOW() + INTERVAL '30 days', NOW())
		ON CONFLICT (profile_id, plan_code) DO UPDATE
		SET status = EXCLUDED.status,
		    renews_at = EXCLUDED.renews_at
	`, uuid.New(), memberProfileID)

	mustExec(ctx, tx, `
		INSERT INTO auction_announcements (id, office, announcement_no, title, original_link, status, closing_at, created_at, updated_at)
		VALUES ($1, '臺北關', 'TP-2026-001', '沒入數位相機與鏡頭一批', 'https://web.customs.gov.tw/taipei/htmlList/86208152b05545bdad39749ea730870d', 'upcoming', NOW() + INTERVAL '72 hours', NOW(), NOW())
		ON CONFLICT (announcement_no) DO UPDATE
		SET title = EXCLUDED.title,
		    original_link = EXCLUDED.original_link,
		    closing_at = EXCLUDED.closing_at,
		    updated_at = NOW()
	`, announcementID)

	mustExec(ctx, tx, `
		INSERT INTO auction_lots (id, announcement_id, title, category, viewing_at, closing_at, warning_tags, disclaimers, created_at)
		VALUES ($1, $2, '沒入數位相機與鏡頭一批', '3C', NOW() + INTERVAL '48 hours', NOW() + INTERVAL '72 hours', ARRAY['現狀交付'], ARRAY['現狀交付', '不負瑕疵擔保', '得標後需自行負擔相關稅費'], NOW())
		ON CONFLICT (id) DO UPDATE
		SET category = EXCLUDED.category,
		    viewing_at = EXCLUDED.viewing_at,
		    closing_at = EXCLUDED.closing_at,
		    warning_tags = EXCLUDED.warning_tags,
		    disclaimers = EXCLUDED.disclaimers
	`, lotID, announcementID)

	mustExec(ctx, tx, `
		INSERT INTO auction_results (id, auction_lot_id, final_price, recorded_at, created_at)
		VALUES ($1, $2, 58000, NOW() - INTERVAL '24 hours', NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), lotID)

	mustExec(ctx, tx, `
		INSERT INTO auction_change_log (id, auction_lot_id, checksum, change_summary, created_at)
		VALUES ($1, $2, 'seed-checksum-v1', '{"source":"seed","category":"3C","title":"沒入數位相機與鏡頭一批"}'::jsonb, NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), lotID)

	mustExec(ctx, tx, `
		INSERT INTO knowledge_articles (id, slug, title, summary, body_mdx, published_at, created_at)
		VALUES ($1, 'bid-form-guide', '新手指南：標單怎麼填', '示範如何填寫通信投標標單與押標金流程。', '# 標單填寫\n\n1. 先確認關別與標號。\n2. 押標金依公告繳納。', NOW(), NOW())
		ON CONFLICT (slug) DO UPDATE
		SET title = EXCLUDED.title,
		    summary = EXCLUDED.summary,
		    body_mdx = EXCLUDED.body_mdx,
		    published_at = EXCLUDED.published_at
	`, uuid.New())

	mustExec(ctx, tx, `
		INSERT INTO courses (id, slug, title, description, stripe_price_id, created_at)
		VALUES ($1, 'import-car-practice', '進口車標售實務', '從看車、驗車到領車的完整操作。', 'price_course_change_me', NOW())
		ON CONFLICT (slug) DO UPDATE
		SET title = EXCLUDED.title,
		    description = EXCLUDED.description,
		    stripe_price_id = EXCLUDED.stripe_price_id
	`, courseID)

	mustExec(ctx, tx, `
		INSERT INTO course_lessons (id, course_id, title, body_mdx, sort_order)
		VALUES ($1, $2, '看車前檢核', '列出驗車、文件、拖吊三個核心檢核項。', 1)
		ON CONFLICT DO NOTHING
	`, uuid.New(), courseID)

	mustExec(ctx, tx, `
		INSERT INTO community_posts (id, profile_id, auction_lot_id, title, body, office, status, created_at)
		VALUES ($1, $2, $3, '臺北關相機批次看貨紀錄', '鏡頭外觀有明顯刮痕，建議實地複查。', '臺北關', 'published', NOW())
		ON CONFLICT (id) DO UPDATE
		SET title = EXCLUDED.title,
		    body = EXCLUDED.body,
		    status = EXCLUDED.status
	`, postID, memberProfileID, lotID)

	mustExec(ctx, tx, `
		INSERT INTO community_reports (id, post_id, reporter_profile_id, reason, status, created_at)
		VALUES ($1, $2, $3, '缺少看貨照片佐證', 'pending', NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), postID, adminProfileID)

	mustExec(ctx, tx, `
		INSERT INTO advisor_profiles (id, profile_id, name, specialty, description, contact_email, created_at)
		VALUES ($1, $2, '王顧問', '進口車標售', '協助驗車與相關文件流程。', 'advisor@example.com', NOW())
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    specialty = EXCLUDED.specialty,
		    description = EXCLUDED.description,
		    contact_email = EXCLUDED.contact_email
	`, advisorID, adminProfileID)

	mustExec(ctx, tx, `
		INSERT INTO advisor_leads (id, advisor_id, name, email, message, category, created_at)
		VALUES ($1, $2, '示例會員', 'member@example.com', '需要協助驗車與提領安排。', '進口車驗車', NOW())
		ON CONFLICT DO NOTHING
	`, uuid.New(), advisorID)

	mustExec(ctx, tx, `
		INSERT INTO ad_slots (id, name, location, created_at)
		VALUES ($1, 'homepage-hero', 'home.hero', NOW())
		ON CONFLICT (name) DO NOTHING
	`, uuid.New())

	mustExec(ctx, tx, `
		INSERT INTO crawler_runs (id, source, office, checksum, row_count, status, trigger_source, ran_at, next_run_at)
		VALUES ($1, 'seed', '臺北關', 'seed-checksum-v1', 1, 'healthy', 'seed', NOW() - INTERVAL '10 minutes', NOW() + INTERVAL '20 minutes')
		ON CONFLICT DO NOTHING
	`, uuid.New())

	if err := tx.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "commit seed transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("seed complete at %s\n", time.Now().Format(time.RFC3339))
}

func mustExec(ctx context.Context, tx pgx.Tx, sql string, args ...any) {
	if _, err := tx.Exec(ctx, sql, args...); err != nil {
		fmt.Fprintf(os.Stderr, "seed exec failed: %v\nsql: %s\n", err, strings.TrimSpace(sql))
		os.Exit(1)
	}
}

func usesSupabaseTransactionPooler(databaseURL string) bool {
	parsed, err := url.Parse(databaseURL)
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(parsed.Hostname()), "pooler.supabase.com") && parsed.Port() == "6543"
}
