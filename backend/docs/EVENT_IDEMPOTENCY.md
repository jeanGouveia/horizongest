# Event Idempotency Strategy

**Version:** 1.0
**Date:** 31/07/2026
**Scope:** HorizonGest Event Publishing Infrastructure

---

## 1. Overview

The HorizonGest event publishing infrastructure implements the **Outbox Pattern** with **At Least Once** delivery semantics. This document explains why this choice was made, the implications for consumers, and how to implement idempotent consumers.

---

## 2. Why At Least Once?

### 2.1 Trade-offs in Event Delivery

There are three possible delivery semantics for event systems:

| Semantic | Description | Pros | Cons |
|----------|-------------|------|------|
| **At Most Once** | Events may be lost but never duplicated | Simple, no deduplication needed | Data loss, unacceptable for business events |
| **At Least Once** | Events may be duplicated but never lost | No data loss, reliable | Consumers must be idempotent |
| **Exactly Once** | Events are never lost or duplicated | Ideal | Complex, requires distributed transactions, not always possible |

### 2.2 Why HorizonGest Chose At Least Once

**Business Requirements:**
- Order events MUST NOT be lost
- Stock adjustment events MUST NOT be lost
- Payment events MUST NOT be lost
- Any data loss is unacceptable for business operations

**Technical Constraints:**
- RabbitMQ does not support Exactly Once semantics out of the box
- Exactly Once requires complex coordination between producer and consumer
- Exactly Once often requires two-phase commit or idempotent producers
- For MVP, At Least Once provides the best balance of reliability and complexity

**The Outbox Pattern:**
- Events are persisted in the database within the same transaction as the business operation
- This guarantees that events are never lost if the transaction commits
- The Dispatcher processes events and publishes to RabbitMQ
- If the Dispatcher crashes after publishing but before marking as completed, the event will be republished

**Conclusion:** At Least Once is the correct choice for HorizonGest's business requirements and technical constraints.

---

## 3. Why Consumers MUST Be Idempotent

### 3.1 The Problem: Duplicate Events

**Scenario:**
1. Dispatcher publishes event E to RabbitMQ ✅
2. RabbitMQ confirms publication ✅
3. Dispatcher crashes before marking E as completed ❌
4. Event E remains with status `processing` in the database ❌
5. Next cycle of Dispatcher picks up event E again ❌
6. Event E is published again to RabbitMQ ❌
7. Consumer receives event E twice ❌

**Result:** The same event is delivered twice to the consumer.

### 3.2 The Solution: Idempotent Consumers

An idempotent consumer can process the same event multiple times without side effects.

**Example:**
- Event: `order.created` with `order_id=100`
- First processing: Create order in external system ✅
- Second processing: Check if order already exists, skip creation ✅

**Key Principle:** Consumers MUST check if the operation has already been performed before executing it.

---

## 4. How to Ignore Repeated EventID

### 4.1 Event Identification

Each event in HorizonGest has the following identifiers:

```json
{
  "event_id": 12345,
  "event_type": "order.created",
  "aggregate_type": "order",
  "aggregate_id": 100,
  "tenant_id": 1,
  "event_version": "1.0",
  "payload": {
    "order_id": 100,
    "customer_id": 50,
    "total": 150.00
  }
}
```

**Primary Key for Deduplication:**
- `event_id` - Unique identifier of the event in the Outbox table
- `aggregate_type` + `aggregate_id` + `event_type` - Business key (unique constraint in database)

### 4.2 Deduplication Strategies

#### Strategy 1: Store Processed Event IDs

**Implementation:**
```go
type EmailWorker struct {
    processedEvents map[uint]bool  // event_id -> true
    mu             sync.RWMutex
}

func (w *EmailWorker) ProcessEvent(event domain.OutboxEvent) error {
    w.mu.Lock()
    if w.processedEvents[event.ID] {
        w.mu.Unlock()
        log.Printf("Event %d already processed, skipping", event.ID)
        return nil
    }
    w.processedEvents[event.ID] = true
    w.mu.Unlock()

    // Process event
    return w.sendEmail(event)
}
```

**Pros:**
- Simple to implement
- Fast in-memory check

**Cons:**
- Not persistent across restarts
- Memory usage grows over time

**Use Case:** Suitable for short-lived workers with low event volume.

#### Strategy 2: Database Deduplication Table

**Implementation:**
```sql
CREATE TABLE processed_events (
    event_id INTEGER PRIMARY KEY,
    processed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

```go
func (w *EmailWorker) ProcessEvent(event domain.OutboxEvent) error {
    // Check if already processed
    var count int
    err := w.db.QueryRow("SELECT COUNT(*) FROM processed_events WHERE event_id = ?", event.ID).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        log.Printf("Event %d already processed, skipping", event.ID)
        return nil
    }

    // Mark as processed
    _, err = w.db.Exec("INSERT INTO processed_events (event_id) VALUES (?)", event.ID)
    if err != nil {
        return err
    }

    // Process event
    return w.sendEmail(event)
}
```

**Pros:**
- Persistent across restarts
- Bounded memory usage

**Cons:**
- Additional database query per event
- Requires additional table

**Use Case:** Suitable for production with high reliability requirements.

#### Strategy 3: Business-Level Deduplication

**Implementation:**
```go
func (w *EmailWorker) ProcessEvent(event domain.OutboxEvent) error {
    // Check if email was already sent for this order
    var count int
    err := w.db.QueryRow(
        "SELECT COUNT(*) FROM sent_emails WHERE order_id = ? AND event_type = ?",
        event.AggregateID,
        event.EventType,
    ).Scan(&count)
    if err != nil {
        return err
    }
    if count > 0 {
        log.Printf("Email already sent for order %d, skipping", event.AggregateID)
        return nil
    }

    // Send email
    err = w.sendEmail(event)
    if err != nil {
        return err
    }

    // Record sent email
    _, err = w.db.Exec(
        "INSERT INTO sent_emails (order_id, event_type, sent_at) VALUES (?, ?, ?)",
        event.AggregateID,
        event.EventType,
        time.Now(),
    )
    return err
}
```

**Pros:**
- Leverages existing business data
- No additional deduplication infrastructure

**Cons:**
- Requires business-specific logic
- May not work for all event types

**Use Case:** Suitable when business state already tracks the operation.

---

## 5. Best Practices by Consumer Type

### 5.1 Email Worker

**Scenario:** Send confirmation email when order is created.

**Idempotency Strategy:** Business-Level Deduplication

**Implementation:**
```go
type EmailWorker struct {
    db *gorm.DB
}

func (w *EmailWorker) ProcessOrderCreated(event domain.OutboxEvent) error {
    // Check if email was already sent
    var sentEmail EmailLog
    err := w.db.Where("order_id = ? AND event_type = ?", event.AggregateID, event.EventType).First(&sentEmail).Error
    if err == nil {
        log.Printf("Email already sent for order %d", event.AggregateID)
        return nil
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    // Parse payload
    var payload OrderCreatedPayload
    if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
        return err
    }

    // Send email
    if err := w.sendEmail(payload.CustomerEmail, "Order Confirmation", generateEmailBody(payload)); err != nil {
        return err
    }

    // Record sent email
    sentEmail = EmailLog{
        OrderID:   event.AggregateID,
        EventType: event.EventType,
        EventID:   event.ID,
        SentAt:    time.Now(),
    }
    return w.db.Create(&sentEmail).Error
}
```

**Key Points:**
- Check `sent_emails` table before sending
- Record successful sends in `sent_emails` table
- Use `event_id` for correlation

---

### 5.2 Webhook Worker

**Scenario:** Send webhook to external system when order status changes.

**Idempotency Strategy:** Database Deduplication Table

**Implementation:**
```go
type WebhookWorker struct {
    db *gorm.DB
}

func (w *WebhookWorker) ProcessOrderStatusChanged(event domain.OutboxEvent) error {
    // Check if webhook was already sent
    var webhookLog WebhookLog
    err := w.db.Where("event_id = ?", event.ID).First(&webhookLog).Error
    if err == nil {
        log.Printf("Webhook already sent for event %d", event.ID)
        return nil
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    // Parse payload
    var payload OrderStatusChangedPayload
    if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
        return err
    }

    // Send webhook
    response, err := w.sendWebhook(payload.WebhookURL, payload)
    if err != nil {
        return err
    }

    // Record webhook call
    webhookLog = WebhookLog{
        EventID:     event.ID,
        URL:         payload.WebhookURL,
        ResponseCode: response.StatusCode,
        SentAt:      time.Now(),
    }
    return w.db.Create(&webhookLog).Error
}
```

**Key Points:**
- Use `event_id` as primary deduplication key
- Record webhook response for debugging
- Retry on failure (with backoff)

---

### 5.3 iFood Worker

**Scenario:** Send order to iFood when order is confirmed.

**Idempotency Strategy:** Business-Level Deduplication (iFood API)

**Implementation:**
```go
type iFoodWorker struct {
    iFoodClient *iFood.Client
    db          *gorm.DB
}

func (w *iFoodWorker) ProcessOrderConfirmed(event domain.OutboxEvent) error {
    // Check if order was already sent to iFood
    var iFoodLog iFoodLog
    err := w.db.Where("order_id = ?", event.AggregateID).First(&iFoodLog).Error
    if err == nil {
        log.Printf("Order already sent to iFood: %s", iFoodLog.IFoodOrderID)
        return nil
    }
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    // Parse payload
    var payload OrderConfirmedPayload
    if err := json.Unmarshal([]byte(event.Payload), &payload); err != nil {
        return err
    }

    // Send to iFood
    iFoodOrderID, err := w.iFoodClient.CreateOrder(payload)
    if err != nil {
        return err
    }

    // Record iFood order ID
    iFoodLog = iFoodLog{
        OrderID:      event.AggregateID,
        IFoodOrderID: iFoodOrderID,
        EventID:      event.ID,
        SentAt:       time.Now(),
    }
    return w.db.Create(&iFoodLog).Error
}
```

**Key Points:**
- Check local database before calling iFood API
- Store iFood's order ID for correlation
- iFood API may have its own idempotency (check their documentation)

---

### 5.4 Future Consumers

**General Guidelines:**

1. **Always check before processing:**
   - Use `event_id` or business key to check if already processed
   - Never assume an event is new

2. **Record processing:**
   - Store successful processing in a persistent store
   - Include `event_id` for correlation

3. **Handle failures gracefully:**
   - Retry with backoff
   - Don't mark as processed until operation succeeds
   - Move to dead letter after max retries

4. **Log everything:**
   - Log when an event is skipped (duplicate)
   - Log when an event is processed
   - Log errors with full context

5. **Test idempotency:**
   - Write tests that process the same event twice
   - Verify no side effects on second processing
   - Verify first processing still works correctly

---

## 6. Common Pitfalls

### 6.1 Assuming Events Are Unique

**Wrong:**
```go
func (w *Worker) ProcessEvent(event domain.OutboxEvent) error {
    // Assumes event is new, no check
    return w.doSomething(event)
}
```

**Correct:**
```go
func (w *Worker) ProcessEvent(event domain.OutboxEvent) error {
    // Always check if already processed
    if w.isAlreadyProcessed(event.ID) {
        return nil
    }
    return w.doSomething(event)
}
```

### 6.2 Using In-Memory State Only

**Wrong:**
```go
var processedEvents = make(map[uint]bool)  // Lost on restart
```

**Correct:**
```go
// Use persistent storage (database, Redis, etc.)
```

### 6.3 Not Recording Processing

**Wrong:**
```go
func (w *Worker) ProcessEvent(event domain.OutboxEvent) error {
    if w.isAlreadyProcessed(event.ID) {
        return nil
    }
    // Process but don't record
    return w.doSomething(event)
}
```

**Correct:**
```go
func (w *Worker) ProcessEvent(event domain.OutboxEvent) error {
    if w.isAlreadyProcessed(event.ID) {
        return nil
    }
    err := w.doSomething(event)
    if err != nil {
        return err
    }
    // Record successful processing
    return w.recordProcessed(event.ID)
}
```

### 6.4 Ignoring Event Type in Deduplication

**Wrong:**
```go
// Only checks order_id, ignores event type
if w.emailSentForOrder(event.AggregateID) {
    return nil
}
```

**Correct:**
```go
// Checks both order_id and event type
if w.emailSentForOrder(event.AggregateID, event.EventType) {
    return nil
}
```

---

## 7. Testing Idempotency

### 7.1 Unit Test Example

```go
func TestEmailWorker_Idempotency(t *testing.T) {
    worker := NewEmailWorker(db)

    event := domain.OutboxEvent{
        ID:            1,
        AggregateType: "order",
        AggregateID:   100,
        EventType:     "order.created",
        Payload:       `{"order_id":100}`,
    }

    // First processing
    err := worker.ProcessEvent(event)
    if err != nil {
        t.Fatalf("First processing failed: %v", err)
    }

    // Verify email was sent
    var count int
    db.QueryRow("SELECT COUNT(*) FROM sent_emails WHERE event_id = ?", event.ID).Scan(&count)
    if count != 1 {
        t.Errorf("Expected 1 sent email, got %d", count)
    }

    // Second processing (should be idempotent)
    err = worker.ProcessEvent(event)
    if err != nil {
        t.Fatalf("Second processing failed: %v", err)
    }

    // Verify no additional email was sent
    db.QueryRow("SELECT COUNT(*) FROM sent_emails WHERE event_id = ?", event.ID).Scan(&count)
    if count != 1 {
        t.Errorf("Expected 1 sent email (idempotent), got %d", count)
    }
}
```

---

## 8. Summary

**Key Takeaways:**

1. **HorizonGest uses At Least Once delivery** - Events may be duplicated but never lost
2. **Consumers MUST be idempotent** - Handle duplicate events gracefully
3. **Check before processing** - Use `event_id` or business key to deduplicate
4. **Record processing** - Store successful processing in persistent storage
5. **Test thoroughly** - Verify idempotency with unit tests

**Failure to implement idempotency will result in:**
- Duplicate emails sent to customers
- Duplicate webhooks called
- Duplicate orders sent to iFood
- Data inconsistency in external systems

**Resources:**
- Outbox Pattern: https://microservices.io/patterns/data/transactional-outbox
- Idempotency Patterns: https://www.enterpriseintegrationpatterns.com/patterns/messaging/idempotent-messenger.html

---

**END OF DOCUMENT**
