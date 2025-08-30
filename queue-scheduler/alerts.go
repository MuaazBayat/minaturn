package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// WhatsAppMessage represents the message structure for WhatsApp Business API
type WhatsAppMessage struct {
	MessagingProduct string           `json:"messaging_product"`
	RecipientType    string           `json:"recipient_type"`
	To               string           `json:"to"`
	Type             string           `json:"type"`
	Text             WhatsAppTextBody `json:"text"`
}

type WhatsAppTextBody struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}

// AlertSystem handles sending notifications through various channels
type AlertSystem struct {
	// Rate limiting to prevent spam
	sentAlerts map[string]time.Time
	
	// Configuration
	rateLimitWindow time.Duration
}

// NewAlertSystem creates a new alert system
func NewAlertSystem() *AlertSystem {
	return &AlertSystem{
		sentAlerts:      make(map[string]time.Time),
		rateLimitWindow: 5 * time.Minute, // Allow position updates every 5 minutes
	}
}

// SendAlert processes and sends an alert through appropriate channel
func (a *AlertSystem) SendAlert(alert AlertRequest) {
	// Create unique key for rate limiting
	alertKey := fmt.Sprintf("%s:%s:%s", alert.MSISDN, alert.QueueID, alert.Channel)
	
	// Check rate limiting
	if lastSent, exists := a.sentAlerts[alertKey]; exists {
		if time.Since(lastSent) < a.rateLimitWindow {
			log.Printf("Rate limited: skipping alert for %s", alert.MSISDN)
			return
		}
	}

	// Route to appropriate channel
	switch alert.Channel {
	case "whatsapp":
		a.sendWhatsApp(alert)
	case "ussd":
		a.sendUSSD(alert)
	case "websocket":
		a.sendWebSocket(alert)
	default:
		log.Printf("Unknown alert channel: %s", alert.Channel)
		return
	}

	// Record sent alert for rate limiting
	a.sentAlerts[alertKey] = alert.Timestamp
	
	// Cleanup old rate limit entries (prevent memory leak)
	a.cleanupOldAlerts()
}

// sendWhatsApp sends alert via WhatsApp Business API
func (a *AlertSystem) sendWhatsApp(alert AlertRequest) {
	log.Printf("📱 WhatsApp Alert to %s: %s", alert.MSISDN, alert.Message)
	
	// WhatsApp Business API configuration
	const (
		whatsappURL = "https://graph.facebook.com/v23.0/228390687031431/messages"
		accessToken = "EAALG7bsFLaEBO6jdfDo26teUHL7t2UzZChoglQOXlATZAEbJIFRmr9Xh7nRdKVkDZAwrKVZAvw0Pi2lnAeTJDwmZBxjALGU030PewqDwMkiwWL9rb9KzTPv7s3jTsQtUwLZBG4QGKmdMKvTgZBfTkAZBZCNfuvvjp00D8ZAezgYEBRB5VZCam2ZCuZC0A7pTqcnGIudhD9gZDZD"
	)
	
	// Create WhatsApp message payload
	message := WhatsAppMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               alert.MSISDN,
		Type:             "text",
		Text: WhatsAppTextBody{
			PreviewURL: true,
			Body:       alert.Message,
		},
	}
	
	// Marshal JSON payload
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WhatsApp message: %v", err)
		return
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", whatsappURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating WhatsApp request: %v", err)
		return
	}
	
	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	// Send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending WhatsApp message: %v", err)
		return
	}
	defer resp.Body.Close()
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading WhatsApp response: %v", err)
		return
	}
	
	// Check response status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ WhatsApp message sent successfully to %s", alert.MSISDN)
		fmt.Printf("WHATSAPP SUCCESS: %s -> %s\n", alert.MSISDN, alert.Message)
	} else {
		log.Printf("❌ WhatsApp API error (status %d): %s", resp.StatusCode, string(body))
		fmt.Printf("WHATSAPP ERROR: %s -> %s (Status: %d)\n", alert.MSISDN, alert.Message, resp.StatusCode)
	}
}

// sendUSSD sends alert via USSD gateway
func (a *AlertSystem) sendUSSD(alert AlertRequest) {
	log.Printf("📟 USSD Alert to %s: %s", alert.MSISDN, alert.Message)
	
	// TODO: Integrate with USSD gateway
	// For now, just log the alert
	fmt.Printf("USSD: %s -> %s\n", alert.MSISDN, alert.Message)
}

// sendWebSocket sends real-time update via WebSocket
func (a *AlertSystem) sendWebSocket(alert AlertRequest) {
	log.Printf("🔌 WebSocket Alert to %s: %s", alert.MSISDN, alert.Message)
	
	// TODO: Integrate with WebSocket server
	// For now, just log the alert
	fmt.Printf("WEBSOCKET: %s -> %s\n", alert.MSISDN, alert.Message)
}

// cleanupOldAlerts removes rate limit entries older than the window
func (a *AlertSystem) cleanupOldAlerts() {
	cutoff := time.Now().Add(-a.rateLimitWindow)
	
	for key, timestamp := range a.sentAlerts {
		if timestamp.Before(cutoff) {
			delete(a.sentAlerts, key)
		}
	}
}

// GetAlertStats returns statistics about sent alerts
func (a *AlertSystem) GetAlertStats() map[string]int {
	stats := make(map[string]int)
	
	cutoff := time.Now().Add(-24 * time.Hour) // Last 24 hours
	
	for _, timestamp := range a.sentAlerts {
		if timestamp.After(cutoff) {
			stats["alerts_24h"]++
		}
	}
	
	return stats
}