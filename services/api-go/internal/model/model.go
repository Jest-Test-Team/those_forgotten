package model

type AuctionLot struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	CustomsOffice string   `json:"customsOffice"`
	ClosingAt     string   `json:"closingAt"`
	ViewingAt     string   `json:"viewingAt"`
	Category      string   `json:"category"`
	OfficialURL   string   `json:"officialUrl"`
	Summary       string   `json:"summary"`
	Disclaimers   []string `json:"disclaimers"`
}

type AuctionResult struct {
	ID         string `json:"id"`
	FinalPrice int    `json:"finalPrice"`
	RecordedAt string `json:"recordedAt"`
}

type KeywordSubscription struct {
	ID      string `json:"id"`
	Keyword string `json:"keyword"`
}

type KnowledgeArticle struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type Course struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CommunityPost struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Body    string   `json:"body"`
	Images  []string `json:"images"`
	Office  string   `json:"office"`
	Author  string   `json:"author"`
	Visible bool     `json:"visible"`
}

type CommunityReport struct {
	ID       string `json:"id"`
	PostID   string `json:"postId"`
	Reason   string `json:"reason"`
	Status   string `json:"status"`
	CreateAt string `json:"createdAt"`
}

type AdvisorProfile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Specialty   string `json:"specialty"`
	Description string `json:"description"`
}

type AdvisorLead struct {
	ID        string `json:"id"`
	AdvisorID string `json:"advisorId"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}
