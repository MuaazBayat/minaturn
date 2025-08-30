# Queue Scheduler Service

A Go-based scheduling service that monitors Django queue API and sends intelligent notifications to customers about their queue status.

## Features

- **Real-time Queue Monitoring**: Polls Django API every 60 seconds
- **Smart Time Calculations**: Computes average processing time and estimated wait times
- **Multi-channel Alerts**: WhatsApp, USSD, and WebSocket notifications
- **Rate Limiting**: Prevents spam notifications
- **Graceful Shutdown**: Handles SIGINT/SIGTERM properly

## Architecture

```
Django API ← HTTP Client ← Scheduler ← Calculator → Alert System → Notifications
```

## Alert Triggers

1. **Status Change**: When customer status changes to "in_progress"
2. **Position Alert**: When customer is #1 or #2 in line
3. **Time Alert**: When estimated wait time ≤ 5 minutes

## Usage

### Environment Variables

```bash
export DJANGO_BASE_URL=http://127.0.0.1:8000
```

### Run

```bash
cd queue-scheduler
go run .
```

### Build

```bash
go build -o queue-scheduler
./queue-scheduler
```

## Configuration

- **Polling Interval**: 60 seconds (configurable in scheduler.go)
- **Rate Limit Window**: 10 minutes per customer per alert type
- **Lookback Window**: 30 minutes for calculating average processing time
- **Default Process Time**: 5 minutes (fallback when no historical data)

## Integration Points

### Django API Endpoints
- `GET /queues/all/` - Fetch all queues and entries

### Notification Channels (TODO)
- WhatsApp Business API integration
- USSD gateway integration  
- WebSocket server for real-time frontend updates

## Sample Alert Messages

- `"🔔 You're now being served! Please proceed to the counter."`
- `"⏰ You're next! Please be ready."`
- `"📍 You're #2 in line. Get ready!"`
- `"⏱️ Estimated wait time: 3 minutes. Queue: Main Service"`

## Development

The service is modular with clear separation:

- `models.go` - Data structures matching Django API
- `client.go` - HTTP client for Django API
- `calculator.go` - Time calculation logic
- `scheduler.go` - Main polling loop and state management
- `alerts.go` - Notification system with rate limiting
- `main.go` - Service initialization and graceful shutdown