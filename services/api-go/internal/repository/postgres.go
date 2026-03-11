package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var demoProfileID = uuid.NewSHA1(uuid.NameSpaceURL, []byte("customs-auction-platform-demo-profile"))

type postgresQuerier interface {
	Close()
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

type PostgresRepository struct {
	pool   postgresQuerier
	memory *MemoryRepository
}

func NewPostgresRepository(ctx context.Context, databaseURL string) (*PostgresRepository, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url: %w", err)
	}

	config.MaxConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &PostgresRepository{
		pool:   pool,
		memory: NewMemoryRepository(),
	}, nil
}

func NewPostgresRepositoryWithPool(pool postgresQuerier) *PostgresRepository {
	return &PostgresRepository{
		pool:   pool,
		memory: NewMemoryRepository(),
	}
}

func (p *PostgresRepository) Close() {
	p.pool.Close()
}

func (p *PostgresRepository) ListAuctions() []model.AuctionLot {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.pool.Query(ctx, `
		SELECT id::text, title, COALESCE(category, ''), COALESCE(closing_at, NOW())::text,
		       COALESCE(viewing_at, NOW())::text, COALESCE(original_link, ''), office,
		       ARRAY['現狀交付','不負瑕疵擔保']
		FROM auction_lots
		ORDER BY closing_at NULLS LAST, created_at DESC
		LIMIT 20
	`)
	if err != nil {
		return p.memory.ListAuctions()
	}
	defer rows.Close()

	auctions := make([]model.AuctionLot, 0, 20)
	for rows.Next() {
		var auction model.AuctionLot
		if err := rows.Scan(
			&auction.ID,
			&auction.Title,
			&auction.Category,
			&auction.ClosingAt,
			&auction.ViewingAt,
			&auction.OfficialURL,
			&auction.CustomsOffice,
			&auction.Disclaimers,
		); err != nil {
			return p.memory.ListAuctions()
		}
		auction.Summary = fmt.Sprintf("%s / %s 標案，請搭配官方文件與現場看貨判斷風險。", auction.CustomsOffice, auction.Category)
		auctions = append(auctions, auction)
	}

	if len(auctions) == 0 {
		return p.memory.ListAuctions()
	}

	return auctions
}

func (p *PostgresRepository) GetAuction(id string) (model.AuctionLot, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var auction model.AuctionLot
	err := p.pool.QueryRow(ctx, `
		SELECT id::text, title, COALESCE(category, ''), COALESCE(closing_at, NOW())::text,
		       COALESCE(viewing_at, NOW())::text, COALESCE(original_link, ''), office,
		       ARRAY['現狀交付','不負瑕疵擔保']
		FROM auction_lots
		WHERE id = $1::uuid
	`, id).Scan(
		&auction.ID,
		&auction.Title,
		&auction.Category,
		&auction.ClosingAt,
		&auction.ViewingAt,
		&auction.OfficialURL,
		&auction.CustomsOffice,
		&auction.Disclaimers,
	)
	if err != nil {
		return p.memory.GetAuction(id)
	}

	auction.Summary = fmt.Sprintf("%s / %s 標案，請搭配官方文件與現場看貨判斷風險。", auction.CustomsOffice, auction.Category)
	return auction, true
}

func (p *PostgresRepository) GetAuctionHistory(id string) []model.AuctionResult {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.pool.Query(ctx, `
		SELECT id::text, final_price, recorded_at::text
		FROM auction_results
		WHERE auction_lot_id = $1::uuid
		ORDER BY recorded_at DESC
		LIMIT 20
	`, id)
	if err != nil {
		return p.memory.GetAuctionHistory(id)
	}
	defer rows.Close()

	results := []model.AuctionResult{}
	for rows.Next() {
		var row model.AuctionResult
		if err := rows.Scan(&row.ID, &row.FinalPrice, &row.RecordedAt); err != nil {
			return p.memory.GetAuctionHistory(id)
		}
		results = append(results, row)
	}

	if len(results) == 0 {
		return p.memory.GetAuctionHistory(id)
	}

	return results
}

func (p *PostgresRepository) CreateKeywordSubscription(keyword string) model.KeywordSubscription {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := p.ensureDemoProfile(ctx); err != nil {
		return p.memory.CreateKeywordSubscription(keyword)
	}

	id := uuid.New()
	_, err := p.pool.Exec(ctx, `
		INSERT INTO keyword_subscriptions (id, profile_id, keyword, created_at)
		VALUES ($1, $2, $3, NOW())
	`, id, demoProfileID, keyword)
	if err != nil {
		return p.memory.CreateKeywordSubscription(keyword)
	}

	return model.KeywordSubscription{
		ID:      id.String(),
		Keyword: keyword,
	}
}

func (p *PostgresRepository) DeleteKeywordSubscription(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := p.pool.Exec(ctx, `DELETE FROM keyword_subscriptions WHERE id = $1::uuid`, id)
	if err != nil {
		p.memory.DeleteKeywordSubscription(id)
	}
}

func (p *PostgresRepository) CreateWebPushSubscription(input *dto.WebPushSubscriptionInput) map[string]any {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := p.ensureDemoProfile(ctx); err != nil {
		return p.memory.CreateWebPushSubscription(input)
	}

	id := uuid.New()
	_, err := p.pool.Exec(ctx, `
		INSERT INTO web_push_subscriptions (id, profile_id, endpoint, p256dh, auth_secret, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`, id, demoProfileID, input.Endpoint, input.Keys["p256dh"], input.Keys["auth"])
	if err != nil {
		return p.memory.CreateWebPushSubscription(input)
	}

	return map[string]any{
		"id":       id.String(),
		"endpoint": input.Endpoint,
		"keys":     input.Keys,
	}
}

func (p *PostgresRepository) GetKnowledgeArticle(slug string) (model.KnowledgeArticle, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var article model.KnowledgeArticle
	err := p.pool.QueryRow(ctx, `
		SELECT slug, title, COALESCE(summary, '')
		FROM knowledge_articles
		WHERE slug = $1
	`, slug).Scan(&article.Slug, &article.Title, &article.Summary)
	if err != nil {
		return p.memory.GetKnowledgeArticle(slug)
	}

	return article, true
}

func (p *PostgresRepository) ListCourses() []model.Course {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.pool.Query(ctx, `
		SELECT id::text, title, description
		FROM courses
		ORDER BY created_at DESC
		LIMIT 20
	`)
	if err != nil {
		return p.memory.ListCourses()
	}
	defer rows.Close()

	courses := []model.Course{}
	for rows.Next() {
		var course model.Course
		if err := rows.Scan(&course.ID, &course.Title, &course.Description); err != nil {
			return p.memory.ListCourses()
		}
		courses = append(courses, course)
	}

	if len(courses) == 0 {
		return p.memory.ListCourses()
	}

	return courses
}

func (p *PostgresRepository) ListCommunityPosts() []model.CommunityPost {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.pool.Query(ctx, `
		SELECT id::text, title, body, office, COALESCE(status, 'published')
		FROM community_posts
		ORDER BY created_at DESC
		LIMIT 20
	`)
	if err != nil {
		return p.memory.ListCommunityPosts()
	}
	defer rows.Close()

	posts := []model.CommunityPost{}
	for rows.Next() {
		var post model.CommunityPost
		var status string
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.Office, &status); err != nil {
			return p.memory.ListCommunityPosts()
		}
		post.Author = "站內會員"
		post.Visible = status == "published"
		post.Images = []string{}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		return p.memory.ListCommunityPosts()
	}

	return posts
}

func (p *PostgresRepository) CreateCommunityPost(input *dto.CommunityPostInput) model.CommunityPost {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := uuid.New()
	_, err := p.pool.Exec(ctx, `
		INSERT INTO community_posts (id, title, body, office, status, created_at)
		VALUES ($1, $2, $3, $4, 'published', NOW())
	`, id, input.Title, input.Body, input.Office)
	if err != nil {
		return p.memory.CreateCommunityPost(input)
	}

	return model.CommunityPost{
		ID:      id.String(),
		Title:   input.Title,
		Body:    input.Body,
		Images:  input.Image,
		Office:  input.Office,
		Author:  input.Author,
		Visible: true,
	}
}

func (p *PostgresRepository) ReportCommunityPost(postID string, input *dto.ReportInput) model.CommunityReport {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := uuid.New()
	_, err := p.pool.Exec(ctx, `
		INSERT INTO community_reports (id, post_id, reporter_profile_id, reason, status, created_at)
		VALUES ($1, $2::uuid, NULL, $3, 'pending', NOW())
	`, id, postID, input.Reason)
	if err != nil {
		return p.memory.ReportCommunityPost(postID, input)
	}

	return model.CommunityReport{
		ID:       id.String(),
		PostID:   postID,
		Reason:   input.Reason,
		Status:   "pending",
		CreateAt: time.Now().Format(time.RFC3339),
	}
}

func (p *PostgresRepository) ListAdvisors() []model.AdvisorProfile {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.pool.Query(ctx, `
		SELECT id::text, name, specialty, description
		FROM advisor_profiles
		ORDER BY created_at DESC
		LIMIT 20
	`)
	if err != nil {
		return p.memory.ListAdvisors()
	}
	defer rows.Close()

	advisors := []model.AdvisorProfile{}
	for rows.Next() {
		var advisor model.AdvisorProfile
		if err := rows.Scan(&advisor.ID, &advisor.Name, &advisor.Specialty, &advisor.Description); err != nil {
			return p.memory.ListAdvisors()
		}
		advisors = append(advisors, advisor)
	}

	if len(advisors) == 0 {
		return p.memory.ListAdvisors()
	}

	return advisors
}

func (p *PostgresRepository) CreateAdvisorLead(input *dto.AdvisorLeadInput) model.AdvisorLead {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id := uuid.New()
	advisorID, err := uuid.Parse(input.AdvisorID)
	if err != nil {
		return p.memory.CreateAdvisorLead(input)
	}

	_, err = p.pool.Exec(ctx, `
		INSERT INTO advisor_leads (id, advisor_id, name, email, message, category, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, id, advisorID, input.Name, input.Email, input.Message, input.Category)
	if err != nil {
		return p.memory.CreateAdvisorLead(input)
	}

	return model.AdvisorLead{
		ID:        id.String(),
		AdvisorID: input.AdvisorID,
		Name:      input.Name,
		Email:     input.Email,
		Message:   input.Message,
	}
}

func (p *PostgresRepository) IngestAuctions(input *dto.IngestPayload) map[string]any {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return p.memory.IngestAuctions(input)
	}
	defer tx.Rollback(ctx)

	for _, row := range input.Rows {
		announcementID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("%s:%s", input.Source, row.AnnouncementNo)))
		lotID := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("%s:%s:lot", input.Source, row.AnnouncementNo)))

		_, err := tx.Exec(ctx, `
			INSERT INTO auction_announcements (id, office, announcement_no, title, original_link, closing_at, updated_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6::timestamptz, NOW(), NOW())
			ON CONFLICT (announcement_no) DO UPDATE
			SET office = EXCLUDED.office,
			    title = EXCLUDED.title,
			    original_link = EXCLUDED.original_link,
			    closing_at = EXCLUDED.closing_at,
			    updated_at = NOW()
		`, announcementID, row.Office, row.AnnouncementNo, row.Title, row.OriginalLink, row.ClosingAt)
		if err != nil {
			return p.memory.IngestAuctions(input)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO auction_lots (id, announcement_id, title, category, closing_at, warning_tags, disclaimers, created_at)
			VALUES ($1, $2, $3, $4, $5::timestamptz, $6, ARRAY['現狀交付','不負瑕疵擔保'], NOW())
			ON CONFLICT (id) DO UPDATE
			SET title = EXCLUDED.title,
			    category = EXCLUDED.category,
			    closing_at = EXCLUDED.closing_at,
			    warning_tags = EXCLUDED.warning_tags
		`, lotID, announcementID, row.Title, row.Category, row.ClosingAt, row.Warnings)
		if err != nil {
			return p.memory.IngestAuctions(input)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return p.memory.IngestAuctions(input)
	}

	return map[string]any{
		"source":   input.Source,
		"checksum": input.Checksum,
		"received": len(input.Rows),
	}
}

func (p *PostgresRepository) ensureDemoProfile(ctx context.Context) error {
	_, err := p.pool.Exec(ctx, `
		INSERT INTO profiles (id, email, full_name, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, demoProfileID, "demo@customs.local", "Demo Member")
	return err
}
