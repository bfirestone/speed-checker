# Speed Checker CLI

The Speed Checker application now provides a comprehensive command-line interface built with Cobra CLI, allowing you to run components separately or together, and manage all aspects of network testing.

## 🏗️ **Architecture Overview**

The application is now modular with separate components:

- **API Server** (`api`) - HTTP REST API and web dashboard
- **Testing Daemon** (`daemon`) - Background speed and iperf testing
- **CLI Tools** (`test`, `hosts`) - One-off testing and host management
- **Combined Mode** (`all`) - Original monolithic behavior

## 📋 **Available Commands**

### **Main Commands**

```bash
# Show all available commands
speed-checker --help

# Run everything together (original behavior)
speed-checker all

# Run only the API server (no background testing)
speed-checker api

# Run only the testing daemon (no web interface)
speed-checker daemon
```

### **Test Management**

```bash
# Run a single speed test
speed-checker test speed

# Run iperf tests against random hosts
speed-checker test iperf

# Run iperf test with custom duration
speed-checker test iperf --duration 30s

# List recent test results
speed-checker test list
speed-checker test list speed --count 5
speed-checker test list iperf --count 10
```

### **Host Management**

```bash
# List all configured hosts
speed-checker hosts list

# Add a new host
speed-checker hosts add --name "Local Server" --hostname 192.168.1.100 --type lan --description "Main server"
speed-checker hosts add --name "VPN Gateway" --hostname vpn.example.com --type vpn --port 5202
speed-checker hosts add --name "Remote Test" --hostname remote.example.com --type remote

# Delete a host by ID
speed-checker hosts delete 4
```

## 🔧 **Configuration**

### **Global Flags**

All commands support these global configuration flags:

```bash
--config string         # Config file path (default: ./config.yaml)
--database-dsn string   # Database connection string
--host string           # Server host (overrides config)
--port string           # Server port (overrides config)
```

### **Environment Variables**

Configuration can be set via environment variables with `SPEED_CHECKER_` prefix:

```bash
export SPEED_CHECKER_SERVER_HOST=0.0.0.0
export SPEED_CHECKER_SERVER_PORT=8080
export SPEED_CHECKER_DATABASE_DSN="file:speed_checker.db?_fk=1"
export SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=15m
export SPEED_CHECKER_TESTING_IPERF_TEST_INTERVAL=10m
```

## 🚀 **Usage Examples**

### **Development Workflow**

```bash
# Start API server only (for frontend development)
speed-checker api --port 8080

# In another terminal, run daemon for background testing
speed-checker daemon

# Run one-off tests
speed-checker test speed
speed-checker test iperf --duration 20s
```

### **Production Deployment**

```bash
# Traditional single-process deployment
speed-checker all

# Or separate services for better scaling
speed-checker api &          # Web interface
speed-checker daemon &       # Background testing
```

### **Host Management**

```bash
# Set up test targets
speed-checker hosts add --name "Main Server" --hostname 192.168.1.10 --type lan --description "Primary LAN server"
speed-checker hosts add --name "VPN Gateway" --hostname 10.0.0.1 --type vpn --description "Office VPN"
speed-checker hosts add --name "CDN Test" --hostname cdn.example.com --type remote --description "CDN endpoint"

# View configuration
speed-checker hosts list

# Test specific setup
speed-checker test iperf
```

### **Monitoring and Maintenance**

```bash
# Check recent results
speed-checker test list --count 20

# Run manual tests
speed-checker test speed
speed-checker test iperf --duration 60s

# Manage hosts
speed-checker hosts list
speed-checker hosts delete 5
```

## 🔄 **Migration from Monolithic**

The original behavior is preserved with the `all` command:

```bash
# Old way (still works)
./speed-checker

# New equivalent
./speed-checker all
```

## 🐳 **Docker Integration**

The updated Docker setup provides multiple deployment options through profiles:

```bash
# Single container (default - preserves original behavior)
docker-compose up -d

# Separated services (API + daemon in separate containers)
docker-compose --profile separated up -d

# Individual services
docker run speed-checker api
docker run speed-checker daemon
docker run speed-checker all

# Custom configuration
docker run -e SPEED_CHECKER_SERVER_PORT=9090 speed-checker all
```

See `DOCKER.md` for comprehensive Docker deployment guide.

## 🎯 **Benefits of New Architecture**

1. **Separation of Concerns** - API and testing logic are independent
2. **Better Scaling** - Run multiple API instances with single daemon
3. **Development Flexibility** - Test components independently
4. **Operational Control** - Start/stop services separately
5. **Professional CLI** - Standard command structure with help and validation
6. **Backward Compatibility** - Existing workflows continue to work

## 🔍 **Command Reference**

### **speed-checker all**
Runs the complete application (API server + testing daemon) in a single process. This preserves the original monolithic behavior.

### **speed-checker api**
Runs only the HTTP API server with web dashboard. Provides REST endpoints and serves the SvelteKit frontend, but does not perform background testing.

### **speed-checker daemon**
Runs only the background testing daemon. Performs scheduled speed tests and iperf tests according to configuration, but provides no web interface.

### **speed-checker test speed**
Runs a single internet speed test using Ookla Speedtest CLI and displays formatted results.

### **speed-checker test iperf**
Runs iperf tests against random hosts from each category (LAN, VPN, remote). Supports custom duration with `--duration` flag.

### **speed-checker test list [type]**
Lists recent test results. Optional type parameter can be `speed` or `iperf`. Supports `--count` flag to limit results.

### **speed-checker hosts list**
Displays all configured iperf test hosts in a formatted table with ID, name, hostname, type, port, active status, and description.

### **speed-checker hosts add**
Adds a new iperf test host using named flags:
- `--name, -n`: Host name (required)
- `--hostname, -H`: Host hostname/IP address (required)  
- `--type, -t`: Host type - must be `lan`, `vpn`, or `remote` (required)
- `--description, -d`: Host description (optional)
- `--port, -p`: Host port (default: 5201)

### **speed-checker hosts delete <host_id>**
Removes an iperf test host by its database ID.

## 🎉 **Next Steps**

This CLI foundation enables future enhancements:

- **Host update commands** - Modify existing host configurations
- **Test scheduling** - Custom test intervals per host type
- **Export functionality** - JSON/CSV export of test results
- **Advanced filtering** - Query tests by date range, host, performance thresholds
- **Configuration management** - CLI-based config file generation
- **Health checks** - System status and connectivity validation 