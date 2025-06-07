# Configuration Guide

Speed Checker uses [Viper](https://github.com/spf13/viper) for configuration management, supporting multiple configuration sources with a clear precedence order.

## Configuration Sources (in precedence order)

1. **Environment Variables** (highest precedence)
2. **Configuration Files** 
3. **Default Values** (lowest precedence)

## Environment Variables

All environment variables use the `SPEED_CHECKER_` prefix and follow the structure mapping:

| Environment Variable | Config Path | Default | Description |
|---------------------|-------------|---------|-------------|
| `SPEED_CHECKER_SERVER_HOST` | `server.host` | `localhost` | Server bind address |
| `SPEED_CHECKER_SERVER_PORT` | `server.port` | `8080` | Server port |
| `SPEED_CHECKER_DATABASE_DRIVER` | `database.driver` | `sqlite3` | Database driver |
| `SPEED_CHECKER_DATABASE_DSN` | `database.dsn` | `./speedtest_results.db?_fk=1` | Database connection string |
| `SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL` | `testing.speedtest_interval` | `15m` | Interval between speed tests |
| `SPEED_CHECKER_TESTING_IPERF_INTERVAL` | `testing.iperf_interval` | `10m` | Interval between iperf tests |
| `SPEED_CHECKER_TESTING_IPERF_DURATION` | `testing.iperf_duration` | `10` | Duration of each iperf test in seconds |

### Example Usage

```bash
# Run with custom port and intervals
SPEED_CHECKER_SERVER_PORT=9090 \
SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=30m \
SPEED_CHECKER_TESTING_IPERF_INTERVAL=20m \
./speed-checker

# Run with external database
SPEED_CHECKER_DATABASE_DSN="postgres://user:pass@localhost/speedtest" \
SPEED_CHECKER_DATABASE_DRIVER=postgres \
./speed-checker
```

## Configuration Files

Speed Checker looks for configuration files in the following locations (in order):

1. Current directory: `./config.yaml`
2. System config: `/etc/speed-checker/config.yaml`
3. User config: `$HOME/.speed-checker/config.yaml`

### Supported Formats

- YAML (recommended)
- JSON
- TOML

### Example Configuration File

Create a `config.yaml` file:

```yaml
server:
  host: "0.0.0.0"
  port: "8080"

database:
  driver: "sqlite3"
  dsn: "./speedtest_results.db?_fk=1"

testing:
  speedtest_interval: "15m"
  iperf_interval: "10m"
  iperf_duration: 10
```

### JSON Example

```json
{
  "server": {
    "host": "0.0.0.0",
    "port": "8080"
  },
  "database": {
    "driver": "sqlite3",
    "dsn": "./speedtest_results.db?_fk=1"
  },
  "testing": {
    "speedtest_interval": "15m",
    "iperf_interval": "10m",
    "iperf_duration": 10
  }
}
```

## Configuration Precedence Example

If you have:

1. **config.yaml**: `server.port: "8080"`
2. **Environment**: `SPEED_CHECKER_SERVER_PORT=9090`

The application will use port `9090` (environment variable takes precedence).

## Duration Format

Duration values support Go's duration format:

- `"10s"` - 10 seconds
- `"5m"` - 5 minutes  
- `"2h"` - 2 hours
- `"1h30m"` - 1 hour 30 minutes

## Database Configuration

### SQLite (Default)

```yaml
database:
  driver: "sqlite3"
  dsn: "./speedtest_results.db?_fk=1"
```

### PostgreSQL

```yaml
database:
  driver: "postgres"
  dsn: "postgres://username:password@localhost/speedtest?sslmode=disable"
```

### MySQL

```yaml
database:
  driver: "mysql"
  dsn: "username:password@tcp(localhost:3306)/speedtest?parseTime=true"
```

## Docker Configuration

When running in Docker, use environment variables:

```bash
docker run -d \
  -p 8080:8080 \
  -e SPEED_CHECKER_SERVER_HOST=0.0.0.0 \
  -e SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=30m \
  speed-checker
```

Or with Docker Compose (see `docker-compose.yml`):

```yaml
environment:
  - SPEED_CHECKER_SERVER_HOST=0.0.0.0
  - SPEED_CHECKER_SERVER_PORT=8080
  - SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=15m
```

## Advanced Features

### Automatic Environment Variable Binding

Thanks to Viper's experimental bind struct feature, all configuration fields are automatically bound to environment variables without manual setup.

### Configuration Validation

The application validates configuration on startup and will exit with an error if invalid values are provided.

### Hot Reloading

Configuration files can be watched for changes during development (feature can be enabled in future versions).

## Troubleshooting

### Configuration Not Loading

1. Check if config file exists and has correct permissions
2. Verify YAML/JSON syntax is valid
3. Ensure environment variables use correct prefix (`SPEED_CHECKER_`)
4. Check logs for configuration loading messages

### Environment Variables Not Working

1. Ensure variables use the exact naming convention: `SPEED_CHECKER_SECTION_FIELD`
2. Use underscores for nested fields: `SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL`
3. Check that variables are exported in your shell

### Duration Parsing Errors

Ensure duration strings follow Go's format:
- ✅ `"15m"`, `"30s"`, `"2h30m"`
- ❌ `"15 minutes"`, `"30secs"`, `"2:30"` 