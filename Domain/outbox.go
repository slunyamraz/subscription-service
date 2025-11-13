package domain

type OutboxEvent struct {
	AggregateType string      `json:"aggregate_type"`
	AggregateID   string      `json:"aggregate_id"`
	EventType     EventType   `json:"event_type"`
	Payload       interface{} `json:"payload"`
}

type EventType string

const (
	EventTypeSubscriptionCreated EventType = "subscription_created"
	EventTypeSubscriptionUpdated EventType = "subscription_updated"
	EventTypeSubscriptionDeleted EventType = "subscription_deleted"
)

const (
	AggregateTypeSubscription = "subscription"
)
