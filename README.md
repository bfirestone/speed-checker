# Speed Checker

A comprehensive network performance monitoring tool with automated speed testing, iperf3 testing, and a modern web dashboard.

## Features

- **Automated Speed Testing**: Runs Ookla speedtest every 15 minutes
- **Network Performance Testing**: Automated iperf3 tests against LAN/VPN/remote hosts
- **Host Management**: Add, edit, and delete test hosts with different types
- **Web Dashboard**: Modern SvelteKit frontend with real-time updates
- **Search & Filtering**: Advanced filtering capabilities for test results
- **API**: RESTful API for all operations

## Architecture

- **Backend**: Go with Echo framework and Ent ORM
- **Frontend**: SvelteKit with TypeScript and Tailwind CSS
- **Database**: SQLite with foreign key constraints
- **Services**: Modular service layer for business logic

## Installation

### Prerequisites

- Go 1.21+
- Node.js 18+
- `speedtest` CLI tool (Ookla)
- `iperf3` for network testing

### Setup

1. **Install Go dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up frontend:**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

3. **Build and run:**
   ```bash
   # Build backend
   go build -o speed-checker .
   
   # Run backend (starts on port 8080)
   ./speed-checker
   
   # In another terminal, start frontend
   cd frontend
   npm run dev
   ```

4. **Access the dashboard:**
   - Frontend: http://localhost:5173
   - API: http://localhost:8080/api/v1

## Search & Filtering Features

### Speed Test Filtering

The dashboard provides powerful filtering capabilities for speed test results:

- **Server Name Search**: Filter tests by server name (partial match)
- **Slowest Tests**: Show tests sorted by slowest download speeds
- **Result Limit**: Control number of results (5, 10, 25, 50)

**API Usage:**
```bash
# Search by server name
GET /api/v1/speedtest?server_name=Atlanta

# Get slowest tests
GET /api/v1/speedtest?slowest=true&limit=10

# Combined filters
GET /api/v1/speedtest?server_name=Comcast&limit=25
```

### Iperf Test Filtering

Filter iperf3 test results by various criteria:

- **Host Name Search**: Filter by host name (partial match)
- **Host Type**: Filter by LAN, VPN, or Remote hosts
- **Slowest Tests**: Show tests sorted by slowest received speeds
- **Result Limit**: Control number of results

**API Usage:**
```bash
# Filter by host name
GET /api/v1/iperf?host_name=server

# Filter by host type
GET /api/v1/iperf?host_type=lan

# Get slowest iperf tests
GET /api/v1/iperf?slowest=true&limit=10

# Combined filters
GET /api/v1/iperf?host_type=vpn&limit=25
```

## API Endpoints

### Speed Tests
- `GET /api/v1/speedtest` - Get speed tests (with filtering)
- `POST /api/v1/speedtest/run` - Run manual speed test

### Iperf Tests
- `GET /api/v1/iperf` - Get iperf tests (with filtering)
- `POST /api/v1/iperf/run` - Run manual iperf tests

### Host Management
- `GET /api/v1/hosts` - List all hosts
- `POST /api/v1/hosts` - Add new host
- `PUT /api/v1/hosts/:id` - Update host
- `DELETE /api/v1/hosts/:id` - Delete host

### Dashboard
- `GET /api/v1/dashboard` - Get dashboard summary data

## Database Schema

### SpeedTest
- Timestamp, download/upload speeds, ping, jitter
- Server details, ISP, result URL

### IperfTest  
- Sent/received speeds, RTT, retransmits
- Success status, error messages
- Relationship to Host

### Host
- Name, hostname, port, type (lan/vpn/remote)
- Active status, description

## Configuration

The application uses automatic configuration with sensible defaults:
- Database: `speed_checker.db` (SQLite)
- Backend Port: 8080
- Test Intervals: Speed tests every 15min, iperf every 10min

## Example Hosts Configuration

```sql
-- Add sample hosts for testing
INSERT INTO hosts (name, hostname, port, type, active, description) VALUES
('Local Server', '192.168.1.100', 5201, 'lan', true, 'Main server'),
('VPN Gateway', '10.0.0.1', 5201, 'vpn', true, 'VPN tunnel endpoint'),
('Remote Host', 'remote.example.com', 5201, 'remote', true, 'Internet host');
```

## Performance Results

### Example Speed Test Results
- Download: ~936 Mbps
- Upload: ~350 Mbps  
- Ping: ~14ms

### Example Iperf3 Results
- Bidirectional: 943+ Mbps
- RTT: 2.29ms
- Retransmits: Minimal

## Development

### Adding New Filters

1. **Backend**: Add query parameters to handlers in `internal/handlers/api.go`
2. **Services**: Implement filtering logic in service methods
3. **Frontend**: Add UI controls and API calls in dashboard components

### Architecture Notes

- Services handle business logic and database operations
- Handlers manage HTTP requests/responses and validation  
- Ent provides type-safe database operations
- Echo handles routing and middleware

## Troubleshooting

1. **Database Locked**: Ensure only one instance is running
2. **Permission Denied**: Check iperf3 and speedtest CLI permissions
3. **Connection Refused**: Verify host configurations and network connectivity
4. **Foreign Key Constraints**: Database requires `_fk=1` parameter

## License

MIT License - see LICENSE file for details. 