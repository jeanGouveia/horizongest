package observability

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MetricsService handles Prometheus metrics collection
// FASE A.3: B20 - Prometheus metrics
type MetricsService struct {
	// HTTP metrics
	httpRequestsTotal      *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpRequestsByEndpoint *prometheus.CounterVec

	// Database metrics
	dbQueriesTotal      *prometheus.CounterVec
	dbQueryDuration     *prometheus.HistogramVec
	dbErrorsTotal       *prometheus.CounterVec

	// Redis metrics
	redisCommandsTotal  *prometheus.CounterVec
	redisCommandDuration *prometheus.HistogramVec

	// RabbitMQ metrics
	rabbitmqMessagesTotal *prometheus.CounterVec
	rabbitmqMessageDuration *prometheus.HistogramVec

	// JWT metrics
	jwtTokensIssued    *prometheus.CounterVec
	jwtTokensValidated *prometheus.CounterVec
	jwtTokensRevoked   *prometheus.CounterVec

	// Login metrics
	loginAttemptsTotal *prometheus.CounterVec
	loginSuccessTotal  *prometheus.CounterVec
	loginFailureTotal  *prometheus.CounterVec

	// Password reset metrics
	passwordResetRequests *prometheus.CounterVec
	passwordResetSuccess  *prometheus.CounterVec

	// Invitation metrics
	invitationCreated *prometheus.CounterVec
	invitationAccepted *prometheus.CounterVec

	// Stock metrics
	stockAdjustments *prometheus.CounterVec

	// Orders metrics
	ordersCreated *prometheus.CounterVec
	ordersCompleted *prometheus.CounterVec

	// Products metrics
	productsCreated *prometheus.CounterVec
	productsUpdated *prometheus.CounterVec

	// Finance metrics
	financeTransactions *prometheus.CounterVec
}

// NewMetricsService creates a new metrics service
func NewMetricsService() *MetricsService {
	return &MetricsService{
		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		httpRequestsByEndpoint: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_by_endpoint_total",
				Help: "Total number of HTTP requests by endpoint",
			},
			[]string{"method", "endpoint"},
		),

		// Database metrics
		dbQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "table"},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation", "table"},
		),
		dbErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_errors_total",
				Help: "Total number of database errors",
			},
			[]string{"operation", "table"},
		),

		// Redis metrics
		redisCommandsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_commands_total",
				Help: "Total number of Redis commands",
			},
			[]string{"command"},
		),
		redisCommandDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "redis_command_duration_seconds",
				Help:    "Redis command duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"command"},
		),

		// RabbitMQ metrics
		rabbitmqMessagesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "rabbitmq_messages_total",
				Help: "Total number of RabbitMQ messages",
			},
			[]string{"queue", "direction"},
		),
		rabbitmqMessageDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "rabbitmq_message_duration_seconds",
				Help:    "RabbitMQ message processing duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"queue"},
		),

		// JWT metrics
		jwtTokensIssued: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_tokens_issued_total",
				Help: "Total number of JWT tokens issued",
			},
			[]string{"type"},
		),
		jwtTokensValidated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_tokens_validated_total",
				Help: "Total number of JWT tokens validated",
			},
			[]string{"valid"},
		),
		jwtTokensRevoked: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_tokens_revoked_total",
				Help: "Total number of JWT tokens revoked",
			},
			[]string{},
		),

		// Login metrics
		loginAttemptsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "login_attempts_total",
				Help: "Total number of login attempts",
			},
			[]string{"type"},
		),
		loginSuccessTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "login_success_total",
				Help: "Total number of successful logins",
			},
			[]string{"type"},
		),
		loginFailureTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "login_failure_total",
				Help: "Total number of failed logins",
			},
			[]string{"type", "reason"},
		),

		// Password reset metrics
		passwordResetRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "password_reset_requests_total",
				Help: "Total number of password reset requests",
			},
			[]string{},
		),
		passwordResetSuccess: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "password_reset_success_total",
				Help: "Total number of successful password resets",
			},
			[]string{},
		),

		// Invitation metrics
		invitationCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "invitation_created_total",
				Help: "Total number of invitations created",
			},
			[]string{},
		),
		invitationAccepted: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "invitation_accepted_total",
				Help: "Total number of invitations accepted",
			},
			[]string{},
		),

		// Stock metrics
		stockAdjustments: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "stock_adjustments_total",
				Help: "Total number of stock adjustments",
			},
			[]string{"type"},
		),

		// Orders metrics
		ordersCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "orders_created_total",
				Help: "Total number of orders created",
			},
			[]string{},
		),
		ordersCompleted: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "orders_completed_total",
				Help: "Total number of orders completed",
			},
			[]string{},
		),

		// Products metrics
		productsCreated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "products_created_total",
				Help: "Total number of products created",
			},
			[]string{},
		),
		productsUpdated: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "products_updated_total",
				Help: "Total number of products updated",
			},
			[]string{},
		),

		// Finance metrics
		financeTransactions: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "finance_transactions_total",
				Help: "Total number of finance transactions",
			},
			[]string{"type"},
		),
	}
}

// HTTP metrics methods
func (m *MetricsService) RecordHTTPRequest(method, status string) {
	m.httpRequestsTotal.WithLabelValues(method, status).Inc()
}

func (m *MetricsService) RecordHTTPRequestDuration(method, endpoint string, duration time.Duration) {
	m.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func (m *MetricsService) RecordHTTPRequestByEndpoint(method, endpoint string) {
	m.httpRequestsByEndpoint.WithLabelValues(method, endpoint).Inc()
}

// Database metrics methods
func (m *MetricsService) RecordDBQuery(operation, table string) {
	m.dbQueriesTotal.WithLabelValues(operation, table).Inc()
}

func (m *MetricsService) RecordDBQueryDuration(operation, table string, duration time.Duration) {
	m.dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

func (m *MetricsService) RecordDBError(operation, table string) {
	m.dbErrorsTotal.WithLabelValues(operation, table).Inc()
}

// Redis metrics methods
func (m *MetricsService) RecordRedisCommand(command string) {
	m.redisCommandsTotal.WithLabelValues(command).Inc()
}

func (m *MetricsService) RecordRedisCommandDuration(command string, duration time.Duration) {
	m.redisCommandDuration.WithLabelValues(command).Observe(duration.Seconds())
}

// RabbitMQ metrics methods
func (m *MetricsService) RecordRabbitMQMessage(queue, direction string) {
	m.rabbitmqMessagesTotal.WithLabelValues(queue, direction).Inc()
}

func (m *MetricsService) RecordRabbitMQMessageDuration(queue string, duration time.Duration) {
	m.rabbitmqMessageDuration.WithLabelValues(queue).Observe(duration.Seconds())
}

// JWT metrics methods
func (m *MetricsService) RecordJWTTokenIssued(tokenType string) {
	m.jwtTokensIssued.WithLabelValues(tokenType).Inc()
}

func (m *MetricsService) RecordJWTTokenValidated(valid bool) {
	m.jwtTokensValidated.WithLabelValues(strconv.FormatBool(valid)).Inc()
}

func (m *MetricsService) RecordJWTTokenRevoked() {
	m.jwtTokensRevoked.WithLabelValues().Inc()
}

// Login metrics methods
func (m *MetricsService) RecordLoginAttempt(loginType string) {
	m.loginAttemptsTotal.WithLabelValues(loginType).Inc()
}

func (m *MetricsService) RecordLoginSuccess(loginType string) {
	m.loginSuccessTotal.WithLabelValues(loginType).Inc()
}

func (m *MetricsService) RecordLoginFailure(loginType, reason string) {
	m.loginFailureTotal.WithLabelValues(loginType, reason).Inc()
}

// Password reset metrics methods
func (m *MetricsService) RecordPasswordResetRequest() {
	m.passwordResetRequests.WithLabelValues().Inc()
}

func (m *MetricsService) RecordPasswordResetSuccess() {
	m.passwordResetSuccess.WithLabelValues().Inc()
}

// Invitation metrics methods
func (m *MetricsService) RecordInvitationCreated() {
	m.invitationCreated.WithLabelValues().Inc()
}

func (m *MetricsService) RecordInvitationAccepted() {
	m.invitationAccepted.WithLabelValues().Inc()
}

// Stock metrics methods
func (m *MetricsService) RecordStockAdjustment(adjustmentType string) {
	m.stockAdjustments.WithLabelValues(adjustmentType).Inc()
}

// Orders metrics methods
func (m *MetricsService) RecordOrderCreated() {
	m.ordersCreated.WithLabelValues().Inc()
}

func (m *MetricsService) RecordOrderCompleted() {
	m.ordersCompleted.WithLabelValues().Inc()
}

// Products metrics methods
func (m *MetricsService) RecordProductCreated() {
	m.productsCreated.WithLabelValues().Inc()
}

func (m *MetricsService) RecordProductUpdated() {
	m.productsUpdated.WithLabelValues().Inc()
}

// Finance metrics methods
func (m *MetricsService) RecordFinanceTransaction(transactionType string) {
	m.financeTransactions.WithLabelValues(transactionType).Inc()
}
