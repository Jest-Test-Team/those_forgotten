package dto

type KeywordSubscriptionInput struct {
	Keyword string `json:"keyword"`
}

type WebPushSubscriptionInput struct {
	Endpoint string                 `json:"endpoint"`
	Keys     map[string]string      `json:"keys"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type CheckoutSessionInput struct {
	Kind       string `json:"kind"`
	PlanCode   string `json:"plan_code"`
	CourseSlug string `json:"course_slug"`
}

type CommunityPostInput struct {
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Image   []string `json:"image_urls"`
	Office  string   `json:"office"`
	LotID   string   `json:"lot_id"`
	Author  string   `json:"author"`
	Visible bool     `json:"visible"`
}

type ReportInput struct {
	Reason string `json:"reason"`
}

type AdvisorLeadInput struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
	Category  string `json:"category"`
	AdvisorID string `json:"advisor_id"`
}

type IngestPayload struct {
	Source   string              `json:"source"`
	Checksum string              `json:"checksum"`
	Rows     []NormalizedAuction `json:"rows"`
}

type NormalizedAuction struct {
	AnnouncementNo string   `json:"announcement_no"`
	Office         string   `json:"office"`
	Title          string   `json:"title"`
	Category       string   `json:"category"`
	ClosingAt      string   `json:"closing_at"`
	OriginalLink   string   `json:"original_link"`
	Warnings       []string `json:"warnings"`
}
