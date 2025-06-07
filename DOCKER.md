# Docker Deployment Guide

This guide covers all Docker deployment options for the Speed Checker application with modular frontend and backend builds.

## 🏗 **Architecture Overview**

### **Modular Build System**
- **Backend**: `Dockerfile` - API server + daemon (Go 1.24, SQLite, OpenAPI)
- **Frontend**: `frontend/Dockerfile` - SvelteKit application (Node.js 22, standalone)

### **Docker Compose Options**
- **Main**: `docker-compose.yml` - All deployment profiles
- **API-focused**: `docker-compose.api.yml` - Streamlined API-centric setup

---

## 🚀 **Quick Start**

### **API-Centric (Recommended)**
```bash
# Backend + daemon with API communication
docker-compose -f docker-compose.api.yml up

# Include standalone frontend
docker-compose -f docker-compose.api.yml --profile frontend up
```

### **Traditional Monolithic**
```bash
# Single container with everything
docker-compose up
```

---

## 📦 **Individual Builds**

### **Backend Only**
```bash
# Build
docker build -t speed-checker .

# Run API server
docker run -p 8080:8080 -v ./data:/app/data speed-checker api

# Run daemon (API mode)
docker run speed-checker daemon --api-endpoint http://api-host:8080

# Run daemon (legacy direct DB)
docker run -v ./data:/app/data speed-checker daemon --legacy
```

### **Frontend Only**
```bash
# Build
docker build -t speed-checker-frontend ./frontend

# Run
docker run -p 3000:3000 speed-checker-frontend
```

---

## 🔧 **Docker Compose Profiles**

### **Profile: `default` / `single`**
**Traditional monolithic deployment**

```bash
docker-compose up
# OR
docker-compose --profile single up
```

**Services:**
- `speed-checker-all`: API + daemon + database in one container

**Use Cases:**
- Development/testing
- Simple single-host deployments
- Backward compatibility

---

### **Profile: `api-centric`** ⭐ **Recommended**
**Modern microservices architecture**

```bash
docker-compose --profile api-centric up
```

**Services:**
- `speed-checker-api`: HTTP API server (port 8080)
- `speed-checker-daemon-api`: Daemon using API communication

**Benefits:**
- Process isolation (failures don't cascade)
- Horizontal scaling (multiple daemons)
- Type-safe API communication
- Production-ready architecture

---

### **Profile: `legacy`**
**Separated services with direct DB access**

```bash
docker-compose --profile legacy up
```

**Services:**
- `speed-checker-api`: HTTP API server
- `speed-checker-daemon-legacy`: Daemon with direct database access

**Use Cases:**
- Migration testing
- Backward compatibility verification

---

### **Profile: `scaling`**
**Multiple daemons with different intervals**

```bash
docker-compose --profile api-centric --profile scaling up
```

**Services:**
- `speed-checker-api`: API server
- `speed-checker-daemon-api`: Primary daemon (15m speed, 10m iperf)
- `speed-checker-daemon-api-2`: Secondary daemon (20m speed, 15m iperf)

**Use Cases:**
- High-frequency monitoring
- Load distribution
- Multiple network interfaces

---

### **Profile: `frontend`**
**Standalone SvelteKit frontend**

```bash
docker-compose --profile frontend up speed-checker-frontend
# OR with API backend
docker-compose -f docker-compose.api.yml --profile frontend up
```

**Services:**
- `speed-checker-frontend`: SvelteKit server (port 3000)

**Use Cases:**
- Separate frontend deployment
- CDN/reverse proxy setups
- Development with external API

---

## 🌐 **Production Deployment Examples**

### **Single Host Setup**
```bash
# Complete stack with frontend
docker-compose -f docker-compose.api.yml --profile frontend up -d

# Access points:
# - API: http://localhost:8080
# - Frontend: http://localhost:3000
```

### **Multi-Host Setup**

**Host 1 (API Server):**
```bash
docker run -d --name speed-checker-api \
  -p 8080:8080 \
  -v /data/speed-checker:/app/data \
  speed-checker api
```

**Host 2+ (Daemons):**
```bash
docker run -d --name speed-checker-daemon-1 \
  -e SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=10m \
  speed-checker daemon --api-endpoint http://api-host:8080

docker run -d --name speed-checker-daemon-2 \
  -e SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=30m \
  speed-checker daemon --api-endpoint http://api-host:8080
```

**Host 3 (Frontend):**
```bash
docker run -d --name speed-checker-frontend \
  -p 3000:3000 \
  speed-checker-frontend
```

### **Cloud/Kubernetes**
```yaml
# API Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: speed-checker-api
spec:
  replicas: 1
  selector:
    matchLabels:
      app: speed-checker-api
  template:
    spec:
      containers:
      - name: api
        image: speed-checker:latest
        args: ["api"]
        ports:
        - containerPort: 8080

---
# Daemon Deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: speed-checker-daemon
spec:
  replicas: 3  # Multiple daemons
  selector:
    matchLabels:
      app: speed-checker-daemon
  template:
    spec:
      containers:
      - name: daemon
        image: speed-checker:latest
        args: ["daemon", "--api-endpoint", "http://speed-checker-api:8080"]
```

---

## 🔧 **Configuration**

### **Environment Variables**

**API Server:**
```bash
SPEED_CHECKER_DATABASE_DSN=/app/data/speedtest_results.db?_fk=1
SPEED_CHECKER_SERVER_HOST=0.0.0.0
SPEED_CHECKER_SERVER_PORT=8080
```

**Daemon:**
```bash
SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=15m
SPEED_CHECKER_TESTING_IPERF_INTERVAL=10m
SPEED_CHECKER_TESTING_IPERF_DURATION=10
```

**Frontend:**
```bash
HOST=0.0.0.0
PORT=3000
```

### **Volume Mounts**

**API Server (Database persistence):**
```bash
-v ./data:/app/data
```

**Time synchronization (all containers):**
```bash
-v /etc/localtime:/etc/localtime:ro
```

---

## 🩺 **Health Checks**

### **API Server**
```bash
# HTTP endpoint check
wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/dashboard
```

### **Daemon**
```bash
# Process check
pgrep -f "speed-checker daemon"
```

### **Frontend**
```bash
# HTTP endpoint check
wget --no-verbose --tries=1 --spider http://localhost:3000
```

---

## 🔄 **Migration Path**

### **From Monolithic to API-Centric**

1. **Current**: `docker-compose up` (all in one)
2. **Transition**: `docker-compose --profile legacy up` (separated with direct DB)
3. **Target**: `docker-compose --profile api-centric up` (API communication)

Each step maintains data compatibility with zero downtime.

---

## 🐛 **Troubleshooting**

### **Build Issues**

**Code generation fails:**
```bash
# Ensure generated files are present
make generate
```

**Frontend build errors:**
```bash
cd frontend
npm install --legacy-peer-deps
npm run build
```

### **Runtime Issues**

**Daemon can't connect to API:**
```bash
# Check network connectivity
docker exec daemon-container wget -O- http://api-container:8080/api/v1/dashboard
```

**Database permissions:**
```bash
# Fix ownership
docker exec api-container chown -R appuser:appgroup /app/data
```

**Frontend not accessible:**
```bash
# Check SvelteKit adapter
# Ensure @sveltejs/adapter-node is installed
cd frontend && npm list @sveltejs/adapter-node
```

---

## 📊 **Performance Considerations**

### **Resource Requirements**

**API Server:**
- CPU: 1 core
- Memory: 512MB
- Disk: 10GB (database growth)

**Daemon:**
- CPU: 0.5 cores
- Memory: 256MB
- Network: Requires external connectivity

**Frontend:**
- CPU: 0.5 cores
- Memory: 256MB
- Disk: 100MB

### **Scaling Guidelines**

**Horizontal (Multiple Daemons):**
```bash
# Different intervals to spread load
-e SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=10m  # Fast
-e SPEED_CHECKER_TESTING_SPEEDTEST_INTERVAL=30m  # Slow
```

**Vertical (Resource Limits):**
```yaml
resources:
  limits:
    memory: "512Mi"
    cpu: "1000m"
  requests:
    memory: "256Mi"
    cpu: "500m"
```

This modular Docker architecture supports everything from simple development setups to complex production deployments with full separation of concerns! 🎉 