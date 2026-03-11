package controller

import (
	"net/http"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/service"
	echo "github.com/labstack/echo/v4"
)

type Controller struct {
	service     *service.PlatformService
	ingestToken string
}

func New(service *service.PlatformService, ingestToken string) *Controller {
	return &Controller{service: service, ingestToken: ingestToken}
}

func RegisterRoutes(e *echo.Echo, ctl *Controller) {
	api := e.Group("/v1")
	api.GET("/auctions", ctl.ListAuctions)
	api.GET("/auctions/:id", ctl.GetAuction)
	api.GET("/auctions/:id/history", ctl.GetAuctionHistory)
	api.GET("/auctions/calendar.ics", ctl.GetCalendarFeed)
	api.GET("/keyword-subscriptions", ctl.ListKeywordSubscriptions)
	api.POST("/keyword-subscriptions", ctl.CreateKeywordSubscription)
	api.DELETE("/keyword-subscriptions/:id", ctl.DeleteKeywordSubscription)
	api.POST("/web-push-subscriptions", ctl.CreateWebPushSubscription)
	api.GET("/knowledge/articles/:slug", ctl.GetKnowledgeArticle)
	api.GET("/courses", ctl.ListCourses)
	api.POST("/stripe/checkout", ctl.CreateCheckoutSession)
	api.POST("/stripe/webhook", ctl.StripeWebhook)
	api.GET("/community/posts", ctl.ListCommunityPosts)
	api.POST("/community/posts", ctl.CreateCommunityPost)
	api.POST("/community/posts/:id/report", ctl.ReportCommunityPost)
	api.GET("/admin/community-reports", ctl.ListCommunityReports)
	api.GET("/advisors", ctl.ListAdvisors)
	api.POST("/advisor-leads", ctl.CreateAdvisorLead)

	internal := e.Group("/internal")
	internal.POST("/ingest/auctions", ctl.IngestAuctions)
}

func (ctl *Controller) ListAuctions(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListAuctions()})
}

func (ctl *Controller) GetAuction(c echo.Context) error {
	auction, ok := ctl.service.GetAuction(c.Param("id"))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "auction not found"})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": auction})
}

func (ctl *Controller) GetAuctionHistory(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.GetAuctionHistory(c.Param("id"))})
}

func (ctl *Controller) GetCalendarFeed(c echo.Context) error {
	if !ctl.service.ValidateCalendarToken(c.QueryParam("token")) {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "calendar token required"})
	}

	return c.Blob(http.StatusOK, "text/calendar", []byte(ctl.service.CalendarFeed()))
}

func (ctl *Controller) ListKeywordSubscriptions(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListKeywordSubscriptions()})
}

func (ctl *Controller) CreateKeywordSubscription(c echo.Context) error {
	input := new(dto.KeywordSubscriptionInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}
	if !ctl.service.ValidateKeyword(input.Keyword) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "keyword required"})
	}

	return c.JSON(http.StatusCreated, map[string]any{"data": ctl.service.CreateKeywordSubscription(input.Keyword)})
}

func (ctl *Controller) DeleteKeywordSubscription(c echo.Context) error {
	ctl.service.DeleteKeywordSubscription(c.Param("id"))
	return c.NoContent(http.StatusNoContent)
}

func (ctl *Controller) CreateWebPushSubscription(c echo.Context) error {
	input := new(dto.WebPushSubscriptionInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}
	if !ctl.service.ValidateWebPush(input) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid web push subscription"})
	}

	return c.JSON(http.StatusCreated, map[string]any{"data": ctl.service.CreateWebPushSubscription(input)})
}

func (ctl *Controller) GetKnowledgeArticle(c echo.Context) error {
	article, ok := ctl.service.GetKnowledgeArticle(c.Param("slug"))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "article not found"})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": article})
}

func (ctl *Controller) ListCourses(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListCourses()})
}

func (ctl *Controller) CreateCheckoutSession(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.CheckoutSession()})
}

func (ctl *Controller) StripeWebhook(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"received": true})
}

func (ctl *Controller) ListCommunityPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListCommunityPosts()})
}

func (ctl *Controller) CreateCommunityPost(c echo.Context) error {
	input := new(dto.CommunityPostInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}
	if !ctl.service.ValidateCommunityPost(input) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "title, body, and office are required"})
	}

	return c.JSON(http.StatusCreated, map[string]any{"data": ctl.service.CreateCommunityPost(input)})
}

func (ctl *Controller) ReportCommunityPost(c echo.Context) error {
	input := new(dto.ReportInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}
	if !ctl.service.ValidateReport(input) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "reason required"})
	}

	return c.JSON(http.StatusCreated, map[string]any{"data": ctl.service.ReportCommunityPost(c.Param("id"), input)})
}

func (ctl *Controller) ListCommunityReports(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListCommunityReports()})
}

func (ctl *Controller) ListAdvisors(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListAdvisors()})
}

func (ctl *Controller) CreateAdvisorLead(c echo.Context) error {
	input := new(dto.AdvisorLeadInput)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}
	if !ctl.service.ValidateAdvisorLead(input) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "advisor lead is incomplete"})
	}

	return c.JSON(http.StatusCreated, map[string]any{"data": ctl.service.CreateAdvisorLead(input)})
}

func (ctl *Controller) IngestAuctions(c echo.Context) error {
	if ctl.ingestToken != "" && c.Request().Header.Get("X-Ingest-Token") != ctl.ingestToken {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "invalid ingest token"})
	}

	input := new(dto.IngestPayload)
	if err := c.Bind(input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}
	if !ctl.service.ValidateIngestPayload(input) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "ingest payload failed checksum validation"})
	}

	return c.JSON(http.StatusAccepted, map[string]any{"data": ctl.service.IngestAuctions(input)})
}
