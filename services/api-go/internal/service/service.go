package service

import (
	"strings"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/model"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/repository"
)

type PlatformService struct {
	repo        repository.Repository
	adminEmails []string
}

func NewPlatformService(repo repository.Repository, adminEmails []string) *PlatformService {
	return &PlatformService{repo: repo, adminEmails: adminEmails}
}

func (s *PlatformService) ListAuctions() []model.AuctionLot {
	return s.repo.ListAuctions()
}

func (s *PlatformService) GetAuction(id string) (model.AuctionLot, bool) {
	return s.repo.GetAuction(id)
}

func (s *PlatformService) GetAuctionHistory(id string) []model.AuctionResult {
	return s.repo.GetAuctionHistory(id)
}

func (s *PlatformService) CalendarFeed() string {
	return strings.Join([]string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Customs Auction Platform//EN",
		"BEGIN:VEVENT",
		"UID:lot-camera-001@customs-auction-platform",
		"SUMMARY:臺北關 沒入數位相機與鏡頭一批",
		"DTSTART:20260314T020000Z",
		"DTEND:20260314T030000Z",
		"DESCRIPTION:現狀交付，不負瑕疵擔保。",
		"END:VEVENT",
		"END:VCALENDAR",
	}, "\r\n")
}

func (s *PlatformService) ListKeywordSubscriptions() []model.KeywordSubscription {
	return s.repo.ListKeywordSubscriptions()
}

func (s *PlatformService) CreateKeywordSubscription(keyword string) model.KeywordSubscription {
	return s.repo.CreateKeywordSubscription(keyword)
}

func (s *PlatformService) DeleteKeywordSubscription(id string) {
	s.repo.DeleteKeywordSubscription(id)
}

func (s *PlatformService) CreateWebPushSubscription(input *dto.WebPushSubscriptionInput) map[string]any {
	return s.repo.CreateWebPushSubscription(input)
}

func (s *PlatformService) GetKnowledgeArticle(slug string) (model.KnowledgeArticle, bool) {
	return s.repo.GetKnowledgeArticle(slug)
}

func (s *PlatformService) ListCourses() []model.Course {
	return s.repo.ListCourses()
}

func (s *PlatformService) CheckoutSession() map[string]string {
	return map[string]string{
		"provider": "stripe",
		"url":      "https://checkout.stripe.com/pay/cs_test_example",
	}
}

func (s *PlatformService) ListCommunityPosts() []model.CommunityPost {
	return s.repo.ListCommunityPosts()
}

func (s *PlatformService) CreateCommunityPost(input *dto.CommunityPostInput) model.CommunityPost {
	if input.Author == "" {
		input.Author = "匿名會員"
	}
	input.Visible = true
	return s.repo.CreateCommunityPost(input)
}

func (s *PlatformService) ReportCommunityPost(postID string, input *dto.ReportInput) model.CommunityReport {
	return s.repo.ReportCommunityPost(postID, input)
}

func (s *PlatformService) ListCommunityReports() []model.CommunityReport {
	return s.repo.ListCommunityReports()
}

func (s *PlatformService) ResolveCommunityReport(id string) (model.CommunityReport, bool) {
	return s.repo.ResolveCommunityReport(id)
}

func (s *PlatformService) ListAdvisors() []model.AdvisorProfile {
	return s.repo.ListAdvisors()
}

func (s *PlatformService) ListAdvisorLeads() []model.AdvisorLead {
	return s.repo.ListAdvisorLeads()
}

func (s *PlatformService) ListCrawlerStatuses() []model.CrawlerStatus {
	return s.repo.ListCrawlerStatuses()
}

func (s *PlatformService) CreateAdvisorLead(input *dto.AdvisorLeadInput) model.AdvisorLead {
	return s.repo.CreateAdvisorLead(input)
}

func (s *PlatformService) IngestAuctions(input *dto.IngestPayload) map[string]any {
	return s.repo.IngestAuctions(input)
}

func (s *PlatformService) GetAuthContext(email string) model.AuthContext {
	normalized := strings.TrimSpace(strings.ToLower(email))
	context := model.AuthContext{
		Email:        normalized,
		Role:         "guest",
		Capabilities: []string{"browse"},
		Source:       "allowlist",
	}

	if normalized == "" {
		return context
	}

	if resolved, ok := s.repo.ResolveAuthContext(normalized); ok {
		return resolved
	}

	context.Role = "member"
	context.Capabilities = []string{"browse", "member"}

	for _, adminEmail := range s.adminEmails {
		if normalized == adminEmail {
			context.Role = "admin"
			context.Capabilities = []string{"browse", "member", "admin"}
			break
		}
	}

	return context
}
