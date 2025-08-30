package main

import "time"

// QueueEntry represents a customer in a queue, matching Django model
type QueueEntry struct {
	ID        string     `json:"id"`
	MSISDN    string     `json:"msisdn"`
	FullName  *string    `json:"full_name"`
	JoinedAt  time.Time  `json:"joined_at"`
	Left      bool       `json:"left"`
	Status    string     `json:"status"` // "waiting", "in_progress", "served"
	StartedAt *time.Time `json:"started_at"`
	ServedAt  *time.Time `json:"served_at"`
}

// Queue represents a service queue, matching Django model
type Queue struct {
	QueueID     string       `json:"queue_id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time    `json:"created_at"`
	Entries     []QueueEntry `json:"entries"`
}

// APIResponse represents the response from Django /queues/all/ endpoint
type APIResponse struct {
	Queues []Queue `json:"queues"`
}

// QueueStats holds calculated statistics for a queue
type QueueStats struct {
	QueueID            string
	AverageProcessTime time.Duration // Average time from in_progress -> served
	EstimatedWaitTime  time.Duration // Estimated wait based on position
	ActiveEntries      int           // Number of people still waiting
}

// AlertRequest represents a notification to be sent
type AlertRequest struct {
	MSISDN    string
	Message   string
	Channel   string // "whatsapp", "ussd", "websocket"
	QueueID   string
	Timestamp time.Time
}