# Speed Checker Dashboard

A comprehensive network performance monitoring tool that combines internet speed testing with iperf3 network testing, featuring a modern web dashboard built with SvelteKit.

## Features

### 🚀 Speed Testing
- **Automated Internet Speed Tests**: Uses Ookla's speedtest CLI for accurate internet speed measurements
- **Scheduled Testing**: Configurable intervals for automatic testing (default: every 15 minutes)
- **Comprehensive Metrics**: Download/upload speeds, ping, jitter, server information, and ISP details

### 🌐 Network Performance Testing
- **iperf3 Integration**: Test network performance against local LAN, VPN, and remote hosts
- **Random Host Selection**: Automatically selects random hosts from each category for testing
- **Configurable Test Duration**: Customizable test duration (default: 10 seconds)
- **Multiple Host Types**: Support for LAN, VPN, and remote host categories

### 📊 Modern Web Dashboard
- **Real-time Dashboard**: Beautiful, responsive interface built with SvelteKit and Tailwind CSS
- **Live Data Updates**: Dashboard refreshes automatically every 30 seconds
- **Host Management**: Add, view, and manage iperf3 test hosts
- **Manual Test Triggers**: Run speed tests and iperf tests on-demand
- **Historical Data**: View recent test results and performance trends

### 🏗️ Architecture
- **Go Backend**: Built with Echo framework for high-performance REST API
- **Ent ORM**: Type-safe database operations with automatic migrations
- **SQLite Database**: Lightweight, embedded database for data persistence
- **SvelteKit Frontend**: Modern, fast, and responsive user interface
- **RESTful API**: Clean API design for easy integration and extension

## Prerequisites

- **Go 1.21+**: For the backend application
- **Node.js 18+**: For the SvelteKit frontend
- **speedtest CLI**: Ookla's official speedtest command-line tool
- **iperf3**: Network performance testing tool

### Installing Prerequisites

#### Speedtest CLI
```bash
# Ubuntu/Debian
curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | sudo bash
sudo apt-get install speedtest

# macOS
brew install speedtest-cli

# Or download from: https://www.speedtest.net/apps/cli
```

#### iperf3
```bash
# Ubuntu/Debian
sudo apt-get install iperf3

# macOS
brew install iperf3

# CentOS/RHEL
sudo yum install iperf3
```

## Installation

1. **Clone the repository**:
```bash
git clone <repository-url>
cd speed-checker
```

2. **Install Go dependencies**:
```bash
go mod tidy
```

3. **Install and build frontend**:
```bash
cd frontend
npm install
npm run build
cd ..
```

4. **Run the application**:
```bash
go run main.go
```

The application will start on `http://localhost:8080`

## Configuration

The application uses sensible defaults but can be configured by modifying the `config.Default()` function in `internal/config/config.go`:

```go
Testing: TestingConfig{
    SpeedTestInterval: 15 * time.Minute,  // Speed test frequency
    IperfTestInterval: 10 * time.Minute,  // iperf test frequency
    IperfTestDuration: 10,                // iperf test duration in seconds
},
Server: ServerConfig{
    Port: "8080",                         // Server port
    Host: "localhost",                    // Server host
},
```

## API Endpoints

### Speed Tests
- `GET /api/v1/speedtest` - Get recent speed tests
- `GET /api/v1/speedtest/range?start=<RFC3339>&end=<RFC3339>` - Get speed tests in date range
- `POST /api/v1/speedtest/run` - Run a speed test manually

### iperf Tests
- `GET /api/v1/iperf` - Get recent iperf tests
- `POST /api/v1/iperf/run` - Run iperf tests manually

### Host Management
- `GET /api/v1/hosts` - Get all configured hosts
- `POST /api/v1/hosts` - Add a new host

### Dashboard
- `GET /api/v1/dashboard` - Get dashboard summary data

## Setting Up iperf3 Hosts

To test network performance, you need to set up iperf3 servers on your target hosts:

### On the target host (server):
```bash
# Run iperf3 in server mode
iperf3 -s -p 5201

# Or run as a daemon
iperf3 -s -D -p 5201
```

### Add hosts via the web interface:
1. Navigate to the "Hosts" page in the dashboard
2. Click "Add Host"
3. Fill in the host details:
   - **Name**: Friendly name (e.g., "Home Router")
   - **Hostname/IP**: IP address or hostname
   - **Port**: iperf3 server port (default: 5201)
   - **Type**: LAN, VPN, or Remote
   - **Description**: Optional description

## Database Schema

The application uses three main entities:

- **SpeedTest**: Internet speed test results
- **IperfTest**: Network performance test results
- **Host**: iperf3 server configurations

Database migrations are handled automatically by Ent ORM.

## Development

### Project Structure
```
speed-checker/
├── main.go                 # Application entry point
├── ent/                    # Ent ORM generated code
│   └── schema/            # Database schema definitions
├── internal/
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP handlers
│   └── services/          # Business logic
├── frontend/              # SvelteKit frontend
│   └── src/
│       └── routes/        # SvelteKit pages
└── README.md
```

### Running in Development

1. **Backend** (with hot reload using air):
```bash
# Install air for hot reload
go install github.com/air-verse/air@latest

# Run with hot reload
air
```

2. **Frontend** (development server):
```bash
cd frontend
npm run dev
```

### Building for Production

```bash
# Build frontend
cd frontend
npm run build
cd ..

# Build Go binary
go build -o speed-checker main.go

# Run
./speed-checker
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **speedtest command not found**:
   - Install Ookla's speedtest CLI (see Prerequisites)

2. **iperf3 tests failing**:
   - Ensure iperf3 servers are running on target hosts
   - Check firewall settings and port accessibility
   - Verify host configurations in the dashboard

3. **Frontend not loading**:
   - Ensure the frontend was built (`npm run build` in frontend directory)
   - Check that the build output exists in `frontend/build/`

4. **Database errors**:
   - The SQLite database is created automatically
   - Check file permissions in the application directory

### Logs

The application logs all activities to stdout. Key events include:
- Speed test executions and results
- iperf test executions and results
- API requests and errors
- Database operations

For more detailed logging, you can modify the log level in the application code. 