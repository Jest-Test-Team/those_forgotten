package repository

import (
	"fmt"
	"sync"
	"time"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/model"
)

type Repository interface {
	ListAuctions() []model.AuctionLot
	GetAuction(id string) (model.AuctionLot, bool)
	GetAuctionHistory(id string) []model.AuctionResult
	ListKeywordSubscriptions() []model.KeywordSubscription
	CreateKeywordSubscription(keyword string) model.KeywordSubscription
	DeleteKeywordSubscription(id string)
	CreateWebPushSubscription(input *dto.WebPushSubscriptionInput) map[string]any
	GetKnowledgeArticle(slug string) (model.KnowledgeArticle, bool)
	ListCourses() []model.Course
	ListCommunityPosts() []model.CommunityPost
	CreateCommunityPost(input *dto.CommunityPostInput) model.CommunityPost
	ReportCommunityPost(postID string, input *dto.ReportInput) model.CommunityReport
	ListCommunityReports() []model.CommunityReport
	ListAdvisors() []model.AdvisorProfile
	ListAdvisorLeads() []model.AdvisorLead
	ListCrawlerStatuses() []model.CrawlerStatus
	CreateAdvisorLead(input *dto.AdvisorLeadInput) model.AdvisorLead
	IngestAuctions(input *dto.IngestPayload) map[string]any
}

type MemoryRepository struct {
	mu          sync.Mutex
	auctions    []model.AuctionLot
	history     map[string][]model.AuctionResult
	articles    map[string]model.KnowledgeArticle
	courses     []model.Course
	posts       []model.CommunityPost
	advisors    []model.AdvisorProfile
	subs        []model.KeywordSubscription
	webPush     []map[string]any
	reports     []model.CommunityReport
	advisorLead []model.AdvisorLead
	crawlers    []model.CrawlerStatus
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		auctions: []model.AuctionLot{
			{
				ID:            "lot-camera-001",
				Title:         "沒入數位相機與鏡頭一批",
				CustomsOffice: "臺北關",
				ClosingAt:     time.Now().Add(72 * time.Hour).Format(time.RFC3339),
				ViewingAt:     time.Now().Add(48 * time.Hour).Format(time.RFC3339),
				Category:      "3C",
				OfficialURL:   "https://web.customs.gov.tw/taipei/htmlList/86208152b05545bdad39749ea730870d",
				Summary:       "熱門相機批次，適合二手設備商與攝影工作室關注。",
				Disclaimers:   []string{"現狀交付", "不負瑕疵擔保", "得標後需自行負擔相關稅費"},
			},
		},
		history: map[string][]model.AuctionResult{
			"lot-camera-001": {
				{ID: "result-001", FinalPrice: 58000, RecordedAt: time.Now().Add(-24 * time.Hour).Format(time.RFC3339)},
			},
		},
		articles: map[string]model.KnowledgeArticle{
			"bid-form-guide": {
				Slug:    "bid-form-guide",
				Title:   "新手指南：標單怎麼填",
				Summary: "示範如何填寫通信投標標單與押標金流程。",
			},
		},
		courses: []model.Course{
			{ID: "course-import-cars", Title: "進口車標售實務", Description: "從看車、驗車到領車的完整操作。"},
		},
		posts: []model.CommunityPost{
			{ID: "post-001", Title: "臺北關相機批次看貨紀錄", Body: "鏡頭外觀有明顯刮痕，建議實地複查。", Images: []string{}, Office: "臺北關", Author: "系統示例", Visible: true},
		},
		reports: []model.CommunityReport{
			{ID: "report-001", PostID: "post-001", PostTitle: "臺北關相機批次看貨紀錄", Office: "臺北關", Reason: "缺少看貨照片佐證", Status: "pending", CreateAt: time.Now().Add(-2 * time.Hour).Format(time.RFC3339)},
		},
		advisors: []model.AdvisorProfile{
			{ID: "advisor-001", Name: "王顧問", Specialty: "進口車標售", Description: "協助驗車與相關文件流程。"},
		},
		advisorLead: []model.AdvisorLead{
			{ID: "lead-001", AdvisorID: "advisor-001", AdvisorName: "王顧問", Name: "示例會員", Email: "member@example.com", Message: "需要協助驗車與提領安排。", Category: "進口車驗車", CreatedAt: time.Now().Add(-90 * time.Minute).Format(time.RFC3339)},
		},
		crawlers: []model.CrawlerStatus{
			{Office: "基隆關", Status: "healthy", LastRunAt: time.Now().Add(-12 * time.Minute).Format(time.RFC3339), NextRunAt: time.Now().Add(18 * time.Minute).Format(time.RFC3339), LastChecksum: "keelung-demo", LastRowCount: 4, TriggerSource: "schedule"},
			{Office: "臺北關", Status: "healthy", LastRunAt: time.Now().Add(-9 * time.Minute).Format(time.RFC3339), NextRunAt: time.Now().Add(21 * time.Minute).Format(time.RFC3339), LastChecksum: "taipei-demo", LastRowCount: 6, TriggerSource: "schedule"},
			{Office: "臺中關", Status: "warning", LastRunAt: time.Now().Add(-33 * time.Minute).Format(time.RFC3339), NextRunAt: time.Now().Add(-3 * time.Minute).Format(time.RFC3339), LastChecksum: "taichung-demo", LastRowCount: 0, TriggerSource: "retry"},
			{Office: "高雄關", Status: "healthy", LastRunAt: time.Now().Add(-7 * time.Minute).Format(time.RFC3339), NextRunAt: time.Now().Add(23 * time.Minute).Format(time.RFC3339), LastChecksum: "kaohsiung-demo", LastRowCount: 5, TriggerSource: "schedule"},
		},
		subs: []model.KeywordSubscription{
			{ID: "sub-1", Keyword: "相機"},
			{ID: "sub-2", Keyword: "進口車"},
		},
	}
}

func (m *MemoryRepository) ListAuctions() []model.AuctionLot {
	return m.auctions
}

func (m *MemoryRepository) GetAuction(id string) (model.AuctionLot, bool) {
	for _, auction := range m.auctions {
		if auction.ID == id {
			return auction, true
		}
	}

	return model.AuctionLot{}, false
}

func (m *MemoryRepository) GetAuctionHistory(id string) []model.AuctionResult {
	return m.history[id]
}

func (m *MemoryRepository) ListKeywordSubscriptions() []model.KeywordSubscription {
	return append([]model.KeywordSubscription{}, m.subs...)
}

func (m *MemoryRepository) CreateKeywordSubscription(keyword string) model.KeywordSubscription {
	m.mu.Lock()
	defer m.mu.Unlock()

	sub := model.KeywordSubscription{
		ID:      fmt.Sprintf("sub-%d", len(m.subs)+1),
		Keyword: keyword,
	}
	m.subs = append(m.subs, sub)
	return sub
}

func (m *MemoryRepository) DeleteKeywordSubscription(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	filtered := make([]model.KeywordSubscription, 0, len(m.subs))
	for _, sub := range m.subs {
		if sub.ID != id {
			filtered = append(filtered, sub)
		}
	}
	m.subs = filtered
}

func (m *MemoryRepository) CreateWebPushSubscription(input *dto.WebPushSubscriptionInput) map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()

	record := map[string]any{
		"id":       fmt.Sprintf("push-%d", len(m.webPush)+1),
		"endpoint": input.Endpoint,
		"keys":     input.Keys,
	}
	m.webPush = append(m.webPush, record)
	return record
}

func (m *MemoryRepository) GetKnowledgeArticle(slug string) (model.KnowledgeArticle, bool) {
	article, ok := m.articles[slug]
	return article, ok
}

func (m *MemoryRepository) ListCourses() []model.Course {
	return m.courses
}

func (m *MemoryRepository) ListCommunityPosts() []model.CommunityPost {
	return m.posts
}

func (m *MemoryRepository) CreateCommunityPost(input *dto.CommunityPostInput) model.CommunityPost {
	m.mu.Lock()
	defer m.mu.Unlock()

	post := model.CommunityPost{
		ID:      fmt.Sprintf("post-%d", len(m.posts)+1),
		Title:   input.Title,
		Body:    input.Body,
		Images:  input.Image,
		Office:  input.Office,
		Author:  input.Author,
		Visible: input.Visible,
	}
	m.posts = append(m.posts, post)
	return post
}

func (m *MemoryRepository) ReportCommunityPost(postID string, input *dto.ReportInput) model.CommunityReport {
	m.mu.Lock()
	defer m.mu.Unlock()

	postTitle := ""
	office := ""
	for _, post := range m.posts {
		if post.ID == postID {
			postTitle = post.Title
			office = post.Office
			break
		}
	}

	report := model.CommunityReport{
		ID:        fmt.Sprintf("report-%d", len(m.reports)+1),
		PostID:    postID,
		PostTitle: postTitle,
		Office:    office,
		Reason:    input.Reason,
		Status:    "pending",
		CreateAt:  time.Now().Format(time.RFC3339),
	}
	m.reports = append(m.reports, report)
	return report
}

func (m *MemoryRepository) ListCommunityReports() []model.CommunityReport {
	return append([]model.CommunityReport{}, m.reports...)
}

func (m *MemoryRepository) ListAdvisors() []model.AdvisorProfile {
	return m.advisors
}

func (m *MemoryRepository) ListAdvisorLeads() []model.AdvisorLead {
	return append([]model.AdvisorLead{}, m.advisorLead...)
}

func (m *MemoryRepository) ListCrawlerStatuses() []model.CrawlerStatus {
	return append([]model.CrawlerStatus{}, m.crawlers...)
}

func (m *MemoryRepository) CreateAdvisorLead(input *dto.AdvisorLeadInput) model.AdvisorLead {
	m.mu.Lock()
	defer m.mu.Unlock()

	advisorName := ""
	for _, advisor := range m.advisors {
		if advisor.ID == input.AdvisorID {
			advisorName = advisor.Name
			break
		}
	}

	lead := model.AdvisorLead{
		ID:          fmt.Sprintf("lead-%d", len(m.advisorLead)+1),
		AdvisorID:   input.AdvisorID,
		AdvisorName: advisorName,
		Name:        input.Name,
		Email:       input.Email,
		Message:     input.Message,
		Category:    input.Category,
		CreatedAt:   time.Now().Format(time.RFC3339),
	}
	m.advisorLead = append(m.advisorLead, lead)
	return lead
}

func (m *MemoryRepository) IngestAuctions(input *dto.IngestPayload) map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()

	return map[string]any{
		"source":   input.Source,
		"checksum": input.Checksum,
		"received": len(input.Rows),
	}
}
