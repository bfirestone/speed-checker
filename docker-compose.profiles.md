# Docker Compose Deployment Profiles

This project supports multiple deployment architectures using Docker Compose profiles.

## Available Profiles

### 1. **Default/Single Service** (Profile: `default` or `single`)
**Traditional monolithic deployment - everything in one container**

```bash
# Start single service (default)
docker-compose up

# Or explicitly
docker-compose --profile single up
```

**Services:**
- `speed-checker-all`: API server + daemon + database in one container

**Use Cases:**
- Simple deployments
- Development/testing
- Small-scale monitoring

---

### 2. **API-Centric Architecture** (Profile: `api-centric`) ⭐ **Recommended**
**Modern microservices architecture with API communication**

```bash
# Start API-centric deployment
docker-compose --profile api-centric up
```

**Services:**
- `speed-checker-api`: HTTP API server only
- `speed-checker-daemon-api`: Daemon that submits via HTTP API

**Benefits:**
- ✅ **Process isolation**: API and daemon failures are independent
- ✅ **Horizontal scaling**: Run multiple daemons
- ✅ **Type-safe communication**: Generated SDK with OpenAPI contracts
- ✅ **Production ready**: Clear service boundaries

---

### 3. **Legacy Separated** (Profile: `legacy`)
**Separated services with direct database access**

```bash
# Start legacy separated deployment
docker-compose --profile legacy up speed-checker-api speed-checker-daemon-legacy
```

**Services:**
- `speed-checker-api`: HTTP API server
- `speed-checker-daemon-legacy`: Daemon with direct DB access (legacy mode)

**Use Cases:**
- Migration from monolithic to API-centric
- Backward compatibility testing

---

### 4. **Horizontal Scaling** (Profile: `scaling`)
**Multiple daemons with different test intervals**

```bash
# Start API server + multiple daemons
docker-compose --profile api-centric --profile scaling up
```

**Services:**
- `speed-checker-api`: API server
- `speed-checker-daemon-api`: Primary daemon (15m speed, 10m iperf)
- `speed-checker-daemon-api-2`: Secondary daemon (20m speed, 15m iperf)

**Use Cases:**
- High-frequency monitoring
- Geographic distribution
- Load balancing across multiple network interfaces

## Configuration Examples

### Environment Variables

**API Server:**
```bash
SPEED_CHECKER_DATABASE_DSN=/app/data/speedtest_results.db?_fk=1
SPEED_CHECKER_SERVER_HOST=0.0.0.0
SPEED_CHECKER_SERVER_PORT=8080
```

**API Daemon:**
```bash
SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=15m
SPEED_CHECKER_TESTING_IPERF_INTERVAL=10m
SPEED_CHECKER_TESTING_IPERF_DURATION=10
```

### Custom Networks

**Connect daemon to external API:**
```bash
docker run --network host \
  speed-checker daemon --api-endpoint http://api.example.com:8080
```

### Health Checks

All services include health checks:
- **API Server**: HTTP endpoint check (`/api/v1/dashboard`)
- **Daemons**: Process check (`pgrep speed-checker`)

## Deployment Recommendations

### 🏠 **Home Lab / Single Host**
```bash
docker-compose --profile api-centric up
```

### 🏢 **Production / Multi-Host**
```bash
# Host 1: API Server
docker-compose up speed-checker-api

# Host 2+: Daemons
docker run speed-checker daemon --api-endpoint http://api-host:8080
```

### ☁️ **Cloud / Kubernetes**
Use the API-centric profile as a base for Kubernetes manifests:
- `speed-checker-api` → API Deployment + Service
- `speed-checker-daemon-api` → Daemon Deployment (multiple replicas)

## Migration Path

**From Monolithic to API-Centric:**

1. **Current**: `docker-compose up` (single service)
2. **Transition**: `docker-compose --profile legacy up` (separated with direct DB)
3. **Target**: `docker-compose --profile api-centric up` (API communication)

Each step maintains data compatibility and zero downtime. 