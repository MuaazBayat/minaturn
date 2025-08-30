# Testing Guide for Queue Scheduler

This guide provides comprehensive testing for the queue scheduler service.

## 🧪 Test Types

### 1. Calculation Tests (Standalone)
Tests the core time calculation logic without needing Django API.

```bash
./queue-scheduler --test
```

**What it tests:**
- Average processing time calculation  
- Position-based wait time estimates
- Edge cases (non-existent customers, in-progress customers)

**Expected output:**
```
🧪 Running Calculation Tests
=== Test 1: Queue Statistics ===
Average Process Time: 3m20s
Expected Average: 3m20s
✅ Average calculation correct

=== Test 2: Position Wait Times ===  
Customer 27601234571: Position 1, Wait Time: 3m20s
  ✅ Position correct
  ✅ Wait time estimation reasonable
```

### 2. Integration Tests (With Django)
Tests the complete flow with realistic queue data.

#### Step 1: Start Django
```bash
cd minaturn
python manage.py runserver
```

#### Step 2: Create Test Data
```bash
cd scripts  
python create_test_data.py
```

**This creates:**
- 5 recently completed customers (for avg time calculation)
- 1 customer currently being served
- 8 customers waiting in queue
- Realistic timestamps and service patterns

#### Step 3: Run Scheduler
```bash
cd queue-scheduler
./queue-scheduler
```

**Expected behavior:**
- Calculates average from completed customers
- Sends alerts for customers in positions 1-2
- Shows realistic wait time estimates

## 📊 Expected Test Results

### Calculation Test Results
```bash
./queue-scheduler --test
```
Should show:
- ✅ Average calculation: ~3-4 minutes (from test data)
- ✅ Position calculations: 1, 2, 3 for waiting customers
- ✅ Wait times: Position × Average time

### Integration Test Results
After running `create_test_data.py` and starting scheduler:

```bash
Processing 1 queues
Queue ABC123 (Test Clinic Queue): 8 active, avg process time: 3m30s
📱 WhatsApp Alert to 27601234571: ⏰ You're next! Please be ready.
📱 WhatsApp Alert to 27601234572: 📍 You're #2 in line. Get ready!
```

## 🔧 Customizing Tests

### Modify Test Data
Edit `scripts/create_test_data.py`:
- Change number of completed customers (affects avg calculation)
- Adjust service times (line 47: `service_time = random.randint(2, 8)`)
- Add more waiting customers

### Modify Alert Triggers
Edit `scheduler.go` line 95-108:
- Change position thresholds (currently 1-2)
- Adjust time thresholds (currently 5 minutes)
- Modify alert messages

### Test Different Scenarios

#### Scenario 1: Rush Hour (Long Queues)
```python
# In create_test_data.py, increase waiting customers:
for i in range(15):  # Instead of 8
```

#### Scenario 2: Slow Service (Long Processing Times)
```python
# In create_test_data.py:
service_time = random.randint(8, 15)  # Longer service times
```

#### Scenario 3: Fast Service (Quick Processing)
```python
# In create_test_data.py:
service_time = random.randint(1, 3)  # Quick service times
```

## 🐛 Troubleshooting

### "No recent processed entries found"
- Check that `create_test_data.py` completed successfully
- Verify customers have both `started_at` and `served_at` timestamps
- Ensure completed customers are within 30-minute lookback window

### "API returned status 404"  
- Ensure Django server is running: `python manage.py runserver`
- Check Django logs for errors
- Verify `/queues/all/` endpoint exists

### Alert not triggering
- Check position calculations (customer should be in position 1-2)
- Verify rate limiting isn't blocking (10-minute cooldown)
- Look for "Rate limited" messages in logs

## 📈 Performance Testing

### Load Test with Multiple Queues
```python
# Run create_test_data.py multiple times for different queues
python create_test_data.py  # Queue 1
python create_test_data.py  # Queue 2 
python create_test_data.py  # Queue 3
```

### Memory Usage
```bash
# Monitor memory usage during polling
ps aux | grep queue-scheduler
```

### Response Times
Check scheduler logs for API response times and calculation duration.

---

## 🎯 Test Checklist

- [ ] Standalone calculation tests pass
- [ ] Django integration works
- [ ] Test data creates properly
- [ ] Average processing time calculated correctly
- [ ] Position-based alerts trigger
- [ ] Rate limiting prevents spam
- [ ] Multiple queues handled
- [ ] Graceful shutdown works (Ctrl+C)