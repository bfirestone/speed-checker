# Docker Quick Start

## 🚀 **Quick Deployment Commands**

### **Single Container (Recommended for most users)**
```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### **Separated Services (For scaling/production)**
```bash
# Start API and daemon separately
docker-compose --profile separated up -d

# Scale API instances (requires override)
docker-compose --profile separated up -d --scale speed-checker-api=3

# View logs
docker-compose logs -f speed-checker-api
docker-compose logs -f speed-checker-daemon
```

### **Manual Docker Commands**
```bash
# Build image
docker build -t speed-checker .

# Run all components
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data speed-checker

# Run API only
docker run -d -p 8080:8080 -v $(pwd)/data:/app/data speed-checker api

# Run daemon only
docker run -d -v $(pwd)/data:/app/data speed-checker daemon
```

## 🎛️ **Management Commands**

```bash
# Execute CLI commands in container
docker exec <container> ./speed-checker hosts list
docker exec <container> ./speed-checker test speed
docker exec <container> ./speed-checker test list

# View real-time logs
docker logs -f <container>

# Access container shell
docker exec -it <container> sh

# Check health
docker inspect <container> --format='{{.State.Health.Status}}'
```

## 🔧 **Environment Variables**

```bash
SPEED_CHECKER_SERVER_HOST=0.0.0.0
SPEED_CHECKER_SERVER_PORT=8080
SPEED_CHECKER_DATABASE_DSN=/app/data/speedtest_results.db?_fk=1
SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=15m
SPEED_CHECKER_TESTING_IPERF_INTERVAL=10m
SPEED_CHECKER_TESTING_IPERF_DURATION=10
```

## 📂 **File Structure**

```
speed-checker/
├── Dockerfile                          # Main image definition
├── docker-compose.yml                  # Multi-service orchestration
├── docker-compose.override.yml.example # Scaling examples
├── data/                               # Database volume (created on first run)
└── DOCKER.md                          # Comprehensive guide
```

## 🔗 **Access Points**

- **Web Dashboard**: http://localhost:8080
- **API Docs**: http://localhost:8080/api/v1/dashboard
- **Container Logs**: `docker-compose logs -f`

For detailed configuration and advanced usage, see [`DOCKER.md`](DOCKER.md). 