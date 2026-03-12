package service

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/model"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/repository"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v82"
	stripecheckoutsession "github.com/stripe/stripe-go/v82/checkout/session"
	"github.com/stripe/stripe-go/v82/webhook"
)

type PlatformService struct {
	repo                  repository.Repository
	adminEmails           []string
	stripeCheckoutBaseURL string
	stripeSecretKey       string
	stripeWebhookSecret   string
	stripeSuccessURL      string
	stripeCancelURL       string
	stripeMembershipPrice string
	stripeCoursePrice     string
}

func NewPlatformService(
	repo repository.Repository,
	adminEmails []string,
	stripeCheckoutBaseURL string,
	stripeSecretKey string,
	stripeWebhookSecret string,
	stripeSuccessURL string,
	stripeCancelURL string,
	stripeMembershipPrice string,
	stripeCoursePrice string,
) *PlatformService {
	return &PlatformService{
		repo:                  repo,
		adminEmails:           adminEmails,
		stripeCheckoutBaseURL: stripeCheckoutBaseURL,
		stripeSecretKey:       stripeSecretKey,
		stripeWebhookSecret:   stripeWebhookSecret,
		stripeSuccessURL:      stripeSuccessURL,
		stripeCancelURL:       stripeCancelURL,
		stripeMembershipPrice: stripeMembershipPrice,
		stripeCoursePrice:     stripeCoursePrice,
	}
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

func (s *PlatformService) CheckoutSession(email string, input *dto.CheckoutSessionInput) (map[string]string, error) {
	baseURL := strings.TrimSpace(s.stripeCheckoutBaseURL)
	if baseURL == "" {
		baseURL = "https://checkout.stripe.com/pay/cs_test_example"
	}

	kind := strings.TrimSpace(strings.ToLower(input.Kind))
	priceID := ""
	mode := stripe.CheckoutSessionModePayment
	switch kind {
	case "membership":
		priceID = strings.TrimSpace(s.stripeMembershipPrice)
		mode = stripe.CheckoutSessionModeSubscription
	case "course":
		priceID = strings.TrimSpace(s.stripeCoursePrice)
	default:
		return nil, errors.New("unsupported checkout kind")
	}

	if strings.TrimSpace(s.stripeSecretKey) != "" &&
		priceID != "" &&
		strings.TrimSpace(s.stripeSuccessURL) != "" &&
		strings.TrimSpace(s.stripeCancelURL) != "" {
		stripe.Key = s.stripeSecretKey
		params := &stripe.CheckoutSessionParams{
			SuccessURL:    stripe.String(s.stripeSuccessURL),
			CancelURL:     stripe.String(s.stripeCancelURL),
			CustomerEmail: stripe.String(strings.TrimSpace(strings.ToLower(email))),
			Mode:          stripe.String(string(mode)),
			Metadata: map[string]string{
				"kind":        kind,
				"email":       strings.TrimSpace(strings.ToLower(email)),
				"plan_code":   strings.TrimSpace(input.PlanCode),
				"course_slug": strings.TrimSpace(input.CourseSlug),
			},
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				{
					Price:    stripe.String(priceID),
					Quantity: stripe.Int64(1),
				},
			},
		}

		session, err := stripecheckoutsession.New(params)
		if err != nil {
			return nil, err
		}

		return map[string]string{
			"provider":  "stripe",
			"url":       session.URL,
			"reference": session.ID,
			"kind":      kind,
			"mode":      "live",
		}, nil
	}

	values := url.Values{}
	values.Set("reference", uuid.NewString())
	values.Set("prefilled_email", strings.TrimSpace(strings.ToLower(email)))
	values.Set("kind", kind)
	if strings.TrimSpace(input.PlanCode) != "" {
		values.Set("plan_code", strings.TrimSpace(input.PlanCode))
	}
	if strings.TrimSpace(input.CourseSlug) != "" {
		values.Set("course_slug", strings.TrimSpace(input.CourseSlug))
	}

	separator := "?"
	if strings.Contains(baseURL, "?") {
		separator = "&"
	}

	return map[string]string{
		"provider":  "stripe",
		"url":       baseURL + separator + values.Encode(),
		"reference": values.Get("reference"),
		"kind":      values.Get("kind"),
		"mode":      "fallback",
	}, nil
}

func (s *PlatformService) HandleStripeWebhook(signature string, payload []byte) (map[string]any, error) {
	eventType := ""
	rawObject := payload

	if strings.TrimSpace(s.stripeWebhookSecret) != "" {
		event, err := webhook.ConstructEvent(payload, signature, s.stripeWebhookSecret)
		if err != nil {
			return nil, err
		}
		eventType = string(event.Type)
		rawObject = event.Data.Raw
	} else {
		var fallback struct {
			Type string `json:"type"`
			Data struct {
				Object json.RawMessage `json:"object"`
			} `json:"data"`
		}
		if err := json.Unmarshal(payload, &fallback); err != nil {
			return nil, err
		}
		eventType = fallback.Type
		rawObject = fallback.Data.Object
	}

	if eventType == "" {
		return nil, errors.New("stripe event type required")
	}
	if eventType != "checkout.session.completed" {
		return map[string]any{"received": true, "type": eventType, "ignored": true}, nil
	}

	var session struct {
		Metadata        map[string]string `json:"metadata"`
		CustomerDetails struct {
			Email string `json:"email"`
		} `json:"customer_details"`
	}
	if err := json.Unmarshal(rawObject, &session); err != nil {
		return nil, err
	}

	email := strings.TrimSpace(strings.ToLower(session.Metadata["email"]))
	if email == "" {
		email = strings.TrimSpace(strings.ToLower(session.CustomerDetails.Email))
	}
	if email == "" {
		return nil, errors.New("stripe session email required")
	}

	switch strings.TrimSpace(strings.ToLower(session.Metadata["kind"])) {
	case "membership":
		planCode := strings.TrimSpace(session.Metadata["plan_code"])
		if planCode == "" {
			return nil, errors.New("membership plan_code required")
		}
		if err := s.repo.UpsertMembershipByEmail(email, planCode, "active", time.Now().Add(30*24*time.Hour)); err != nil {
			return nil, err
		}
	case "course":
		courseSlug := strings.TrimSpace(session.Metadata["course_slug"])
		if courseSlug == "" {
			return nil, errors.New("course_slug required")
		}
		if err := s.repo.GrantCourseAccessByEmail(email, courseSlug, "stripe-webhook"); err != nil {
			return nil, err
		}
	default:
		return map[string]any{"received": true, "type": eventType, "ignored": true}, nil
	}

	return map[string]any{"received": true, "type": eventType, "processed": true, "customer": email}, nil
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
