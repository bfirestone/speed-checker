# Docker Setup for Speed Checker

This document explains how to run the Speed Checker application using Docker and Docker Compose.

## Prerequisites

- Docker
- Docker Compose
- Network access for speed tests and iperf3

## Quick Start

1. **Build and run with Docker Compose:**
   ```bash
   docker-compose up --build
   ```

2. **Access the application:**
   - Open your browser to `http://localhost:8080`
   - The dashboard will show speed tests and iperf tests

3. **Stop the application:**
   ```bash
   docker-compose down
   ```

## Configuration

The application can be configured using environment variables in the `docker-compose.yml` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `SPEED_CHECKER_SERVER_HOST` | `0.0.0.0` | Server bind address |
| `SPEED_CHECKER_SERVER_PORT` | `8080` | Server port |
| `SPEED_CHECKER_DATABASE_DSN` | `/app/data/speedtest_results.db?_fk=1` | SQLite database path |
| `SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL` | `15m` | Interval between speed tests |
| `SPEED_CHECKER_TESTING_IPERF_INTERVAL` | `10m` | Interval between iperf tests |
| `SPEED_CHECKER_TESTING_IPERF_DURATION` | `10` | Duration of each iperf test in seconds |

## Data Persistence

The SQLite database is stored in a Docker volume mounted at `./data:/app/data`. This ensures your test data persists between container restarts.

## Building Manually

If you want to build the Docker image manually:

```bash
# Build the image
docker build -t speed-checker .

# Run the container
docker run -d \
  --name speed-checker \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  speed-checker
```

## Health Checks

The container includes health checks that verify the application is responding on the dashboard endpoint. You can check the health status with:

```bash
docker-compose ps
```

## Logs

View application logs:

```bash
# All logs
docker-compose logs

# Follow logs
docker-compose logs -f

# Specific service logs
docker-compose logs speed-checker
```

## Troubleshooting

### iperf3 Tests Failing
- Ensure the target hosts are accessible from the container
- Check that iperf3 servers are running on the target hosts
- Verify firewall rules allow iperf3 traffic

### Speed Tests Failing
- Ensure internet connectivity from the container
- The speedtest-cli tool is included in the container

### Database Issues
- Check that the `./data` directory has proper permissions
- Ensure the volume mount is working correctly

## Development

For development, you can run the container with additional volume mounts:

```bash
docker run -d \
  --name speed-checker-dev \
  -p 8080:8080 \
  -v $(pwd)/data:/app/data \
  -v $(pwd):/workspace \
  speed-checker
```

## Container Details

- **Base Image:** Alpine Linux (minimal footprint)
- **Included Tools:** speedtest-cli, iperf3, sqlite
- **User:** Runs as non-root user (appuser:appgroup)
- **Ports:** 8080 (HTTP)
- **Health Check:** HTTP GET to `/api/v1/dashboard` 