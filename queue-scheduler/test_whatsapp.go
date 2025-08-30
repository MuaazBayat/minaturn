package main

import (
	"log"
	"time"
)

func testWhatsAppAlert() {
	log.Println("🧪 Testing WhatsApp Alert Function")
	
	// Create alert system
	alertSystem := NewAlertSystem()
	
	// Create test alert request
	alert := AlertRequest{
		MSISDN:    "+27797867873", // Your test number
		Message:   "Test message from queue scheduler - your position has been updated!",
		Channel:   "whatsapp",
		QueueID:   "test-queue-1",
		Timestamp: time.Now(),
	}
	
	// Send the alert
	log.Printf("Sending test WhatsApp alert to %s", alert.MSISDN)
	alertSystem.SendAlert(alert)
	log.Println("✅ Test completed")
}