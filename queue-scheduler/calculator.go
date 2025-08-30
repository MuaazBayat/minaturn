package main

import (
	"log"
	"time"
)

// QueueCalculator handles time calculations for queues
type QueueCalculator struct {
	// lookbackWindow defines how far back to look for historical data
	lookbackWindow time.Duration
}

// NewQueueCalculator creates a new calculator with default 30-minute lookback
func NewQueueCalculator() *QueueCalculator {
	return &QueueCalculator{
		lookbackWindow: 30 * time.Minute,
	}
}

// CalculateQueueStats computes statistics for a given queue
func (calc *QueueCalculator) CalculateQueueStats(queue Queue) QueueStats {
	stats := QueueStats{
		QueueID: queue.QueueID,
	}

	// Calculate average processing time from recent served customers
	avgProcessTime := calc.calculateAverageProcessTime(queue.Entries)
	stats.AverageProcessTime = avgProcessTime

	// Count active entries (not left, not served)
	activeCount := 0
	for _, entry := range queue.Entries {
		if !entry.Left && entry.Status != "served" {
			activeCount++
		}
	}
	stats.ActiveEntries = activeCount

	// Estimate wait time based on position and average process time
	if avgProcessTime > 0 && activeCount > 0 {
		// Rough estimate: average process time * people ahead in line
		stats.EstimatedWaitTime = time.Duration(activeCount) * avgProcessTime
	}

	return stats
}

// calculateAverageProcessTime computes average time from in_progress -> served
// Only considers entries from the last 30 minutes that have been fully processed
func (calc *QueueCalculator) calculateAverageProcessTime(entries []QueueEntry) time.Duration {
	now := time.Now()
	cutoff := now.Add(-calc.lookbackWindow)
	
	var totalProcessTime time.Duration
	var processedCount int

	for _, entry := range entries {
		// Only consider entries that:
		// 1. Were served recently (within lookback window)
		// 2. Have both started_at and served_at timestamps
		if entry.ServedAt == nil || entry.StartedAt == nil {
			continue
		}
		
		if entry.ServedAt.Before(cutoff) {
			continue // Too old, skip
		}

		// Calculate processing time (in_progress -> served)
		processTime := entry.ServedAt.Sub(*entry.StartedAt)
		if processTime > 0 {
			totalProcessTime += processTime
			processedCount++
		}
	}

	if processedCount == 0 {
		log.Printf("No recent processed entries found for queue calculation")
		return 5 * time.Minute // Default fallback
	}

	avgTime := totalProcessTime / time.Duration(processedCount)
	log.Printf("Calculated average process time: %v from %d entries", avgTime, processedCount)
	
	return avgTime
}

// EstimateWaitTimeForPosition calculates wait time for a specific position in queue
func (calc *QueueCalculator) EstimateWaitTimeForPosition(queue Queue, targetMSISDN string) time.Duration {
	stats := calc.CalculateQueueStats(queue)
	
	// Find the target customer's position
	position := calc.findPositionInQueue(queue, targetMSISDN)
	if position <= 0 {
		return 0 // Not found or already being served
	}

	// Estimate: position * average process time
	if stats.AverageProcessTime > 0 {
		return time.Duration(position) * stats.AverageProcessTime
	}

	// Fallback estimate
	return time.Duration(position) * 5 * time.Minute
}

// findPositionInQueue determines customer's position among waiting customers
func (calc *QueueCalculator) findPositionInQueue(queue Queue, targetMSISDN string) int {
	position := 0
	
	for _, entry := range queue.Entries {
		if entry.Left || entry.Status == "served" || entry.Status == "in_progress" {
			continue // Skip customers who left, were served, or are currently being served
		}
		
		position++
		
		if entry.MSISDN == targetMSISDN {
			return position
		}
	}
	
	return 0 // Not found in waiting queue
}