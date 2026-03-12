package controller

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/dennislee928/those_forgotten/services/api-go/internal/dto"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/model"
	"github.com/dennislee928/those_forgotten/services/api-go/internal/service"
	"github.com/golang-jwt/jwt/v5"
	echo "github.com/labstack/echo/v4"
)

type Controller struct {
	service     *service.PlatformService
	ingestToken string
	jwtSecret   []byte
}

func New(service *service.PlatformService, ingestToken string, jwtSecret string) *Controller {
	return &Controller{
		service:     service,
		ingestToken: ingestToken,
		jwtSecret:   []byte(strings.TrimSpace(jwtSecret)),
	}
}

func RegisterRoutes(e *echo.Echo, ctl *Controller) {
	api := e.Group("/v1")
	api.GET("/auctions", ctl.ListAuctions)
	api.GET("/auctions/:id", ctl.GetAuction)
	api.GET("/auctions/:id/history", ctl.GetAuctionHistory)
	api.GET("/auctions/calendar.ics", ctl.GetCalendarFeed)
	api.GET("/auth/context", ctl.GetAuthContext)
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
	api.POST("/admin/community-reports/:id/resolve", ctl.ResolveCommunityReport)
	api.GET("/admin/crawler-status", ctl.ListCrawlerStatuses)
	api.GET("/advisors", ctl.ListAdvisors)
	api.GET("/admin/advisor-leads", ctl.ListAdvisorLeads)
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

func (ctl *Controller) GetAuthContext(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.GetAuthContext(ctl.actorEmail(c))})
}

func (ctl *Controller) actorEmailFromBearer(c echo.Context) (string, bool) {
	if len(ctl.jwtSecret) == 0 {
		return "", false
	}

	authHeader := strings.TrimSpace(c.Request().Header.Get("Authorization"))
	if authHeader == "" {
		return "", false
	}

	scheme, tokenString, ok := strings.Cut(authHeader, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(tokenString) == "" {
		return "", false
	}

	token, err := jwt.Parse(strings.TrimSpace(tokenString), func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}

		return ctl.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return "", false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", false
	}

	normalized := strings.TrimSpace(strings.ToLower(email))
	return normalized, normalized != ""
}

func (ctl *Controller) actorEmail(c echo.Context) string {
	if bearerEmail, ok := ctl.actorEmailFromBearer(c); ok {
		return bearerEmail
	}

	headerEmail := strings.TrimSpace(c.Request().Header.Get("X-Actor-Email"))
	if headerEmail != "" {
		return strings.ToLower(headerEmail)
	}

	return strings.ToLower(strings.TrimSpace(c.QueryParam("email")))
}

func (ctl *Controller) protectedActorEmail(c echo.Context) string {
	if len(ctl.jwtSecret) > 0 {
		if bearerEmail, ok := ctl.actorEmailFromBearer(c); ok {
			return bearerEmail
		}
		return ""
	}

	return ctl.actorEmail(c)
}

func (ctl *Controller) requireAdmin(c echo.Context) error {
	context := ctl.service.GetAuthContext(ctl.protectedActorEmail(c))
	if context.Role != "admin" {
		return c.JSON(http.StatusForbidden, map[string]any{"error": "admin role required"})
	}

	return nil
}

func (ctl *Controller) requireMember(c echo.Context) (model.AuthContext, error) {
	context := ctl.service.GetAuthContext(ctl.protectedActorEmail(c))
	if context.Role == "guest" {
		return context, c.JSON(http.StatusUnauthorized, map[string]any{"error": "authenticated member required"})
	}

	return context, nil
}

func (ctl *Controller) ListKeywordSubscriptions(c echo.Context) error {
	if _, err := ctl.requireMember(c); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListKeywordSubscriptions()})
}

func (ctl *Controller) CreateKeywordSubscription(c echo.Context) error {
	if _, err := ctl.requireMember(c); err != nil {
		return err
	}

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
	if _, err := ctl.requireMember(c); err != nil {
		return err
	}

	ctl.service.DeleteKeywordSubscription(c.Param("id"))
	return c.NoContent(http.StatusNoContent)
}

func (ctl *Controller) CreateWebPushSubscription(c echo.Context) error {
	if _, err := ctl.requireMember(c); err != nil {
		return err
	}

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
	auth, err := ctl.requireMember(c)
	if err != nil {
		return err
	}

	input := new(dto.CheckoutSessionInput)
	if c.Request().ContentLength > 0 {
		if err := c.Bind(input); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
		}
	}
	if !ctl.service.ValidateCheckoutSession(input) {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "checkout kind requires plan_code or course_slug"})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.CheckoutSession(auth.Email, input)})
}

func (ctl *Controller) StripeWebhook(c echo.Context) error {
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
	}

	result, err := ctl.service.HandleStripeWebhook(c.Request().Header.Get("Stripe-Signature"), payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, result)
}

func (ctl *Controller) ListCommunityPosts(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListCommunityPosts()})
}

func (ctl *Controller) CreateCommunityPost(c echo.Context) error {
	if _, err := ctl.requireMember(c); err != nil {
		return err
	}

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
	if _, err := ctl.requireMember(c); err != nil {
		return err
	}

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
	if err := ctl.requireAdmin(c); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListCommunityReports()})
}

func (ctl *Controller) ResolveCommunityReport(c echo.Context) error {
	if err := ctl.requireAdmin(c); err != nil {
		return err
	}

	report, ok := ctl.service.ResolveCommunityReport(c.Param("id"))
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "report not found"})
	}

	return c.JSON(http.StatusOK, map[string]any{"data": report})
}

func (ctl *Controller) ListCrawlerStatuses(c echo.Context) error {
	if err := ctl.requireAdmin(c); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListCrawlerStatuses()})
}

func (ctl *Controller) ListAdvisors(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListAdvisors()})
}

func (ctl *Controller) ListAdvisorLeads(c echo.Context) error {
	if err := ctl.requireAdmin(c); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]any{"data": ctl.service.ListAdvisorLeads()})
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
