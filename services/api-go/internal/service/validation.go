package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
)

func (s *PlatformService) ValidateCalendarToken(token string) bool {
	return strings.TrimSpace(token) != ""
}

func (s *PlatformService) ValidateKeyword(keyword string) bool {
	return strings.TrimSpace(keyword) != ""
}

func (s *PlatformService) ValidateWebPush(input *dto.WebPushSubscriptionInput) bool {
	if strings.TrimSpace(input.Endpoint) == "" {
		return false
	}
	if strings.TrimSpace(input.Keys["p256dh"]) == "" || strings.TrimSpace(input.Keys["auth"]) == "" {
		return false
	}
	return true
}

func (s *PlatformService) ValidateCommunityPost(input *dto.CommunityPostInput) bool {
	return strings.TrimSpace(input.Title) != "" && strings.TrimSpace(input.Body) != "" && strings.TrimSpace(input.Office) != ""
}

func (s *PlatformService) ValidateReport(input *dto.ReportInput) bool {
	return strings.TrimSpace(input.Reason) != ""
}

func (s *PlatformService) ValidateAdvisorLead(input *dto.AdvisorLeadInput) bool {
	return strings.TrimSpace(input.Name) != "" &&
		strings.TrimSpace(input.Email) != "" &&
		strings.TrimSpace(input.Message) != "" &&
		strings.TrimSpace(input.AdvisorID) != ""
}

func (s *PlatformService) ValidateIngestPayload(input *dto.IngestPayload) bool {
	if strings.TrimSpace(input.Source) == "" || len(input.Rows) == 0 {
		return false
	}
	serialized, err := json.Marshal(input.Rows)
	if err != nil {
		return false
	}
	sum := sha256.Sum256(serialized)
	return hex.EncodeToString(sum[:]) == input.Checksum
}
