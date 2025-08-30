package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Scheduler handles periodic polling and processing of queue data
type Scheduler struct {
	client     *APIClient
	calculator *QueueCalculator
	alerter    *AlertSystem
	interval   time.Duration
	
	// Store previous state to detect changes
	previousQueues map[string]Queue
}

// NewScheduler creates a new scheduler with default 60-second interval
func NewScheduler(apiClient *APIClient, calculator *QueueCalculator, alerter *AlertSystem) *Scheduler {
	return &Scheduler{
		client:         apiClient,
		calculator:     calculator,
		alerter:        alerter,
		interval:       60 * time.Second,
		previousQueues: make(map[string]Queue),
	}
}

// Start begins the periodic polling loop
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("Starting queue scheduler with %v interval", s.interval)
	
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Initial run
	s.processQueues()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopping...")
			return
		case <-ticker.C:
			s.processQueues()
		}
	}
}

// processQueues fetches queue data and processes alerts
func (s *Scheduler) processQueues() {
	log.Println("Fetching queue data...")
	
	response, err := s.client.GetAllQueues()
	if err != nil {
		log.Printf("Error fetching queues: %v", err)
		return
	}

	log.Printf("Processing %d queues", len(response.Queues))

	for _, queue := range response.Queues {
		s.processQueue(queue)
		
		// Update previous state
		s.previousQueues[queue.QueueID] = queue
	}
}

// processQueue handles a single queue's processing and alerts
func (s *Scheduler) processQueue(queue Queue) {
	stats := s.calculator.CalculateQueueStats(queue)
	
	log.Printf("Queue %s (%s): %d active, avg process time: %v", 
		queue.QueueID, queue.Name, stats.ActiveEntries, stats.AverageProcessTime)

	// Process each active entry for potential alerts
	for _, entry := range queue.Entries {
		if entry.Left || entry.Status == "served" {
			continue // Skip inactive entries
		}

		s.checkForAlerts(queue, entry, stats)
	}
}

// checkForAlerts determines if alerts should be sent for a queue entry
func (s *Scheduler) checkForAlerts(queue Queue, entry QueueEntry, stats QueueStats) {
	waitTime := s.calculator.EstimateWaitTimeForPosition(queue, entry.MSISDN)
	position := s.calculator.findPositionInQueue(queue, entry.MSISDN)

	shouldAlert := false
	var message string

	if entry.Status == "in_progress" && !s.wasInProgress(queue.QueueID, entry.MSISDN) {
		shouldAlert = true
		message = "🔔 You're now being served! Please proceed to the counter."
	} else if entry.Status == "waiting" {
		// Send update to ALL waiting customers with their position and wait time
		shouldAlert = true
		
		if position == 1 {
			message = "⏰ You're NEXT! Please be ready. Estimated wait: 0-2 minutes"
		} else if position == 2 {
			message = fmt.Sprintf("📍 Position #%d - You're almost up! Estimated wait: %d minutes", 
				position, int(waitTime.Minutes()))
		} else {
			message = fmt.Sprintf("📋 Position #%d in %s. Estimated wait: %d minutes", 
				position, queue.Name, int(waitTime.Minutes()))
		}
	}

	if shouldAlert {
		alert := AlertRequest{
			MSISDN:    entry.MSISDN,
			Message:   message,
			Channel:   "whatsapp", // Default channel
			QueueID:   queue.QueueID,
			Timestamp: time.Now(),
		}

		s.alerter.SendAlert(alert)
	}
}

// wasInProgress checks if a customer was already in progress in previous state
func (s *Scheduler) wasInProgress(queueID, msisdn string) bool {
	prevQueue, exists := s.previousQueues[queueID]
	if !exists {
		return false
	}

	for _, entry := range prevQueue.Entries {
		if entry.MSISDN == msisdn && entry.Status == "in_progress" {
			return true
		}
	}
	
	return false
}