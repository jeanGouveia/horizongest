package observability

import (
	"testing"
	"time"
)

var metricsService *MetricsService

func init() {
	metricsService = NewMetricsService()
}

func TestMetricsService_NewMetricsService(t *testing.T) {
	service := metricsService

	if service == nil {
		t.Error("expected metrics service to be created")
	}
}

func TestMetricsService_RecordHTTPRequest(t *testing.T) {
	service := metricsService

	service.RecordHTTPRequest("GET", "200")
	service.RecordHTTPRequest("POST", "201")
	service.RecordHTTPRequest("GET", "404")

	// Metrics are recorded, we can't easily verify without a real registry
}

func TestMetricsService_RecordHTTPRequestDuration(t *testing.T) {
	service := metricsService

	service.RecordHTTPRequestDuration("GET", "/api/test", 100*time.Millisecond)
	service.RecordHTTPRequestDuration("POST", "/api/users", 200*time.Millisecond)
}

func TestMetricsService_RecordHTTPRequestByEndpoint(t *testing.T) {
	service := metricsService

	service.RecordHTTPRequestByEndpoint("GET", "/api/test")
	service.RecordHTTPRequestByEndpoint("POST", "/api/users")
}

func TestMetricsService_RecordDBQuery(t *testing.T) {
	service := metricsService

	service.RecordDBQuery("SELECT", "users")
	service.RecordDBQuery("INSERT", "products")
}

func TestMetricsService_RecordDBQueryDuration(t *testing.T) {
	service := metricsService

	service.RecordDBQueryDuration("SELECT", "users", 50*time.Millisecond)
	service.RecordDBQueryDuration("INSERT", "products", 100*time.Millisecond)
}

func TestMetricsService_RecordDBError(t *testing.T) {
	service := metricsService

	service.RecordDBError("SELECT", "users")
	service.RecordDBError("INSERT", "products")
}

func TestMetricsService_RecordRedisCommand(t *testing.T) {
	service := metricsService

	service.RecordRedisCommand("GET")
	service.RecordRedisCommand("SET")
}

func TestMetricsService_RecordRedisCommandDuration(t *testing.T) {
	service := metricsService

	service.RecordRedisCommandDuration("GET", 10*time.Millisecond)
	service.RecordRedisCommandDuration("SET", 20*time.Millisecond)
}

func TestMetricsService_RecordRabbitMQMessage(t *testing.T) {
	service := metricsService

	service.RecordRabbitMQMessage("orders", "published")
	service.RecordRabbitMQMessage("notifications", "consumed")
}

func TestMetricsService_RecordRabbitMQMessageDuration(t *testing.T) {
	service := metricsService

	service.RecordRabbitMQMessageDuration("orders", 150*time.Millisecond)
}

func TestMetricsService_RecordJWTTokenIssued(t *testing.T) {
	service := metricsService

	service.RecordJWTTokenIssued("access")
	service.RecordJWTTokenIssued("refresh")
}

func TestMetricsService_RecordJWTTokenValidated(t *testing.T) {
	service := metricsService

	service.RecordJWTTokenValidated(true)
	service.RecordJWTTokenValidated(false)
}

func TestMetricsService_RecordJWTTokenRevoked(t *testing.T) {
	service := metricsService

	service.RecordJWTTokenRevoked()
}

func TestMetricsService_RecordLoginAttempt(t *testing.T) {
	service := metricsService

	service.RecordLoginAttempt("platform")
	service.RecordLoginAttempt("tenant")
}

func TestMetricsService_RecordLoginSuccess(t *testing.T) {
	service := metricsService

	service.RecordLoginSuccess("platform")
	service.RecordLoginSuccess("tenant")
}

func TestMetricsService_RecordLoginFailure(t *testing.T) {
	service := metricsService

	service.RecordLoginFailure("platform", "invalid_credentials")
	service.RecordLoginFailure("tenant", "user_not_found")
}

func TestMetricsService_RecordPasswordResetRequest(t *testing.T) {
	service := metricsService

	service.RecordPasswordResetRequest()
}

func TestMetricsService_RecordPasswordResetSuccess(t *testing.T) {
	service := metricsService

	service.RecordPasswordResetSuccess()
}

func TestMetricsService_RecordInvitationCreated(t *testing.T) {
	service := metricsService

	service.RecordInvitationCreated()
}

func TestMetricsService_RecordInvitationAccepted(t *testing.T) {
	service := metricsService

	service.RecordInvitationAccepted()
}

func TestMetricsService_RecordStockAdjustment(t *testing.T) {
	service := metricsService

	service.RecordStockAdjustment("increase")
	service.RecordStockAdjustment("decrease")
}

func TestMetricsService_RecordOrderCreated(t *testing.T) {
	service := metricsService

	service.RecordOrderCreated()
}

func TestMetricsService_RecordOrderCompleted(t *testing.T) {
	service := metricsService

	service.RecordOrderCompleted()
}

func TestMetricsService_RecordProductCreated(t *testing.T) {
	service := metricsService

	service.RecordProductCreated()
}

func TestMetricsService_RecordProductUpdated(t *testing.T) {
	service := metricsService

	service.RecordProductUpdated()
}

func TestMetricsService_RecordFinanceTransaction(t *testing.T) {
	service := metricsService

	service.RecordFinanceTransaction("sale")
	service.RecordFinanceTransaction("refund")
}
