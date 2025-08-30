package main

import (
	"fmt"
	"log"
	"time"
)

// testCalculations tests our time calculation logic with sample data
func testCalculations() {
	log.Println("🧪 Testing Queue Calculations")

	// Create test data
	now := time.Now()
	
	// Sample queue with realistic entries
	queue := Queue{
		QueueID: "TEST123",
		Name:    "Test Queue",
		Entries: []QueueEntry{
			// Recently completed customers (for avg calculation)
			{
				ID:        "COMP1",
				MSISDN:    "27601234567",
				Status:    "served",
				JoinedAt:  now.Add(-25 * time.Minute),
				StartedAt: timePtr(now.Add(-20 * time.Minute)),
				ServedAt:  timePtr(now.Add(-17 * time.Minute)), // 3 min service
			},
			{
				ID:        "COMP2", 
				MSISDN:    "27601234568",
				Status:    "served",
				JoinedAt:  now.Add(-22 * time.Minute),
				StartedAt: timePtr(now.Add(-18 * time.Minute)),
				ServedAt:  timePtr(now.Add(-13 * time.Minute)), // 5 min service
			},
			{
				ID:        "COMP3",
				MSISDN:    "27601234569", 
				Status:    "served",
				JoinedAt:  now.Add(-15 * time.Minute),
				StartedAt: timePtr(now.Add(-10 * time.Minute)),
				ServedAt:  timePtr(now.Add(-8 * time.Minute)), // 2 min service
			},
			// Currently being served
			{
				ID:        "CURR1",
				MSISDN:    "27601234570",
				Status:    "in_progress", 
				JoinedAt:  now.Add(-8 * time.Minute),
				StartedAt: timePtr(now.Add(-3 * time.Minute)),
				ServedAt:  nil,
			},
			// Waiting customers
			{
				ID:       "WAIT1",
				MSISDN:   "27601234571", 
				Status:   "waiting",
				JoinedAt: now.Add(-5 * time.Minute),
				Left:     false,
			},
			{
				ID:       "WAIT2",
				MSISDN:   "27601234572",
				Status:   "waiting", 
				JoinedAt: now.Add(-3 * time.Minute),
				Left:     false,
			},
			{
				ID:       "WAIT3",
				MSISDN:   "27601234573",
				Status:   "waiting",
				JoinedAt: now.Add(-1 * time.Minute),
				Left:     false,
			},
		},
	}

	calculator := NewQueueCalculator()
	
	// Test 1: Calculate queue statistics
	fmt.Println("\n=== Test 1: Queue Statistics ===")
	stats := calculator.CalculateQueueStats(queue)
	
	fmt.Printf("Queue ID: %s\n", stats.QueueID)
	fmt.Printf("Average Process Time: %v\n", stats.AverageProcessTime)
	fmt.Printf("Active Entries: %d\n", stats.ActiveEntries)
	fmt.Printf("Estimated Wait Time: %v\n", stats.EstimatedWaitTime)
	
	// Expected: Avg = (3+5+2)/3 = 3.33 minutes
	expectedAvg := (3*time.Minute + 5*time.Minute + 2*time.Minute) / 3
	fmt.Printf("Expected Average: %v\n", expectedAvg)
	
	if abs(stats.AverageProcessTime - expectedAvg) < time.Second {
		fmt.Println("✅ Average calculation correct")
	} else {
		fmt.Println("❌ Average calculation incorrect")
	}
	
	// Test 2: Position-specific wait times
	fmt.Println("\n=== Test 2: Position Wait Times ===")
	
	testCustomers := []string{"27601234571", "27601234572", "27601234573"}
	expectedPositions := []int{1, 2, 3}
	
	for i, msisdn := range testCustomers {
		waitTime := calculator.EstimateWaitTimeForPosition(queue, msisdn)
		position := calculator.findPositionInQueue(queue, msisdn)
		
		fmt.Printf("Customer %s: Position %d, Wait Time: %v\n", 
			msisdn, position, waitTime)
			
		if position == expectedPositions[i] {
			fmt.Printf("  ✅ Position correct\n")
		} else {
			fmt.Printf("  ❌ Position incorrect (expected %d)\n", expectedPositions[i])
		}
		
		// Wait time should be roughly position * avg process time
		expectedWait := time.Duration(position) * stats.AverageProcessTime
		if abs(waitTime - expectedWait) < time.Minute {
			fmt.Printf("  ✅ Wait time estimation reasonable\n")
		} else {
			fmt.Printf("  ❌ Wait time seems off (expected ~%v)\n", expectedWait)
		}
	}
	
	// Test 3: Edge cases
	fmt.Println("\n=== Test 3: Edge Cases ===")
	
	// Customer not in queue
	waitTime := calculator.EstimateWaitTimeForPosition(queue, "27999999999")
	fmt.Printf("Non-existent customer wait time: %v\n", waitTime)
	if waitTime == 0 {
		fmt.Println("✅ Non-existent customer handled correctly")
	}
	
	// Customer already being served
	waitTime = calculator.EstimateWaitTimeForPosition(queue, "27601234570") 
	fmt.Printf("In-progress customer wait time: %v\n", waitTime)
	if waitTime == 0 {
		fmt.Println("✅ In-progress customer handled correctly")
	}
}

// Helper function to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}

// Helper function to calculate absolute difference
func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

// Call this function from main with --test flag