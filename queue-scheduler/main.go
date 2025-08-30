package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Parse command line flags
	testMode := flag.Bool("test", false, "Run calculation tests instead of scheduler")
	whatsappTest := flag.Bool("whatsapp", false, "Test WhatsApp alert function")
	flag.Parse()

	if *testMode {
		log.Println("🧪 Running Calculation Tests")
		testCalculations()
		return
	}

	if *whatsappTest {
		testWhatsAppAlert()
		return
	}

	log.Println("🚀 Starting Queue Scheduler Service")

	// Configuration from environment or defaults
	djangoBaseURL := getEnv("DJANGO_BASE_URL", "http://127.0.0.1:8000")
	
	log.Printf("Django API URL: %s", djangoBaseURL)

	// Initialize components
	apiClient := NewAPIClient(djangoBaseURL)
	calculator := NewQueueCalculator()
	alertSystem := NewAlertSystem()
	scheduler := NewScheduler(apiClient, calculator, alertSystem)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v", sig)
		log.Println("Initiating graceful shutdown...")
		cancel()
	}()

	// Test API connectivity
	log.Println("Testing API connectivity...")
	if _, err := apiClient.GetAllQueues(); err != nil {
		log.Printf("⚠️  Warning: Could not connect to Django API: %v", err)
		log.Println("Service will continue and retry on next poll...")
	} else {
		log.Println("✅ API connectivity test successful")
	}

	// Start the scheduler
	log.Println("Starting scheduler...")
	scheduler.Start(ctx)
	
	log.Println("Queue Scheduler Service stopped")
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}