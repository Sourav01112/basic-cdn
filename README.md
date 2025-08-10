# Basic CDN System

A minimal Content Delivery Network (CDN) implementation demonstrating core caching concepts with Go and Node.js.

## 🚀 What It Does

This CDN shows how major platforms like Netflix and YouTube deliver content faster by caching files closer to users.

**Key Features:**
- **Load Balancer** (Go): Routes requests to edge servers
- **Edge Server** (Go): Caches content with TTL expiration  
- **Origin Server** (Node.js): Serves fresh content
- **TTL-based Caching**: Industry-standard cache expiration
- **Visual Cache Indicators**: JSON responses show cache hits/misses

## 🏗️ Architecture

```
User Request → Load Balancer → Edge Server → Origin Server
                     ↓              ↓
               Routes traffic    Caches content
```

## ⚡ Quick Start

```bash
# Clone and run
git clone <your-repo-url>
cd basic-cdn
docker-compose up

# Test the CDN
curl http://localhost:8080/sample.json
```

**First request**: Cache MISS → fetches from origin  
**Second request**: Cache HIT → serves from cache (notice old timestamp!)

## 🧪 Demo

### Fresh Content (Cache MISS):
```json
{
  "timestamp": "2025-08-10T13:45:23Z",
  "cached": false,
  "requestId": "abc123"
}
```

### Cached Content (Cache HIT):
```json
{
  "timestamp": "2025-08-10T13:45:23Z",  ← Old timestamp proves caching!
  "cached": true,                       ← Modified by edge server
  "cached_at": "2025-08-10T13:47:15Z",  ← When served from cache
  "requestId": "abc123"                 ← Same ID = same cached response
}
```

## 🔧 Configuration

- **Cache TTL**: 60 seconds (configurable)
- **Load Balancer**: Port 8080
- **Edge Server**: Port 8081  
- **Origin Server**: Port 3000

## 🌐 Real-World Equivalent

This demonstrates the same principles used by:
- **Netflix**: Caches movies on edge servers worldwide
- **YouTube**: Stores popular videos closer to users
- **Cloudflare**: Accelerates websites globally


## 🛠️ Tech Stack

- **Go**: High-performance load balancing and caching
- **Node.js**: Dynamic content generation
- **Docker**: Containerized deployment

---

*Built to demonstrate CDN concepts for learning and interviews.*
