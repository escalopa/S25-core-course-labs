# Monitoring with Prometheus

This document describes the metrics collection and monitoring setup using Prometheus and Grafana for the application stack.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Components](#components)
- [Quick Start](#quick-start)
- [Prometheus Configuration](#prometheus-configuration)
- [Application Metrics](#application-metrics)
- [Grafana Dashboards](#grafana-dashboards)
- [Resource Management](#resource-management)
- [Screenshots](#screenshots)

## Overview

The monitoring stack consists of:

- **Prometheus**: Time-series metrics database and monitoring system
- **Grafana**: Visualization and dashboarding platform
- **Application Metrics**: Custom metrics from Python and Go applications
- **Infrastructure Metrics**: Metrics from Loki, Promtail, and Grafana

## Architecture

```text
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│  Python App     │────▶│   Prometheus    │────▶│    Grafana      │
│  (port 5000)    │     │   (port 9090)   │     │   (port 3000)   │
└─────────────────┘     └─────────────────┘     └─────────────────┘
                               ▲
┌─────────────────┐            │
│  Go App         │────────────┤
│  (port 8080)    │            │
└─────────────────┘            │
                               │
┌─────────────────┐            │
│  Loki           │────────────┤
│  (port 3100)    │            │
└─────────────────┘            │
                               │
┌─────────────────┐            │
│  Promtail       │────────────┘
│  (port 9080)    │
└─────────────────┘
```

## Components

### Prometheus

**Purpose**: Collects and stores time-series metrics from all services

**Configuration**:

- Scrape interval: 15 seconds
- Data retention: 15 days
- Storage: Docker volume (`prometheus-data`)

**Scraped Targets**:

- `prometheus` (self-monitoring): `localhost:9090`
- `loki`: `loki:3100/metrics`
- `promtail`: `promtail:9080/metrics`
- `grafana`: `grafana:3000/metrics`
- `python-app`: `python-app:5000/metrics`
- `go-app`: `go-app:8080/metrics`

### Grafana

**Purpose**: Visualizes metrics from Prometheus and logs from Loki

**Datasources** (auto-provisioned):

- Loki (default)
- Prometheus

**Access**:

- URL: `http://localhost:3000`
- Username: `admin`
- Password: `admin`

## Quick Start

Start the monitoring stack:

```bash
make monitoring-up
```

Stop the monitoring stack:

```bash
make monitoring-down
```

View logs:

```bash
make monitoring-logs
```

Restart services:

```bash
make monitoring-restart
```

Clean up:

```bash
make monitoring-clean
```

## Prometheus Configuration

The Prometheus configuration is stored in `prometheus.yml`:

**Global Settings**:

- `scrape_interval`: 15s (how often to scrape targets)
- `evaluation_interval`: 15s (how often to evaluate rules)

**Scrape Configs**:

Each service has a dedicated job with appropriate labels and metrics path.

## Application Metrics

### Python Application Metrics

**Library**: `prometheus-client`

**Endpoint**: `http://localhost:5000/metrics`

**Metrics Collected**:

Metric Name | Type | Description | Labels
------------|------|-------------|-------
`http_requests_total` | Counter | Total HTTP requests | `method`, `endpoint`, `status`
`http_request_duration_seconds` | Histogram | Request duration | `method`, `endpoint`
`http_errors_total` | Counter | Total HTTP errors | `method`, `endpoint`, `status`

**Example Query**:

```promql
# Request rate per second
rate(http_requests_total{service="python-app"}[5m])

# 95th percentile response time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
```

### Go Application Metrics

**Library**: `prometheus/client_golang`

**Endpoint**: `http://localhost:8080/metrics` (via port 9090)

**Metrics Collected**:

Metric Name | Type | Description | Labels
------------|------|-------------|-------
`http_requests_total` | Counter | Total HTTP requests | `method`, `endpoint`, `status`
`http_request_duration_seconds` | Histogram | Request duration | `method`, `endpoint`
`games_created_total` | Counter | Total games created | -
`guesses_total` | Counter | Total guesses made | -
`games_won_total` | Counter | Total games won | -
`games_lost_total` | Counter | Total games lost | -

**Example Queries**:

```promql
# Games won vs lost ratio
games_won_total / (games_won_total + games_lost_total)

# Average guesses per game
rate(guesses_total[5m]) / rate(games_created_total[5m])

# Request rate for game endpoints
rate(http_requests_total{service="go-app",endpoint="/game/"}[5m])
```

## Grafana Dashboards

### Importing Dashboards

You can import pre-built dashboards for Loki and Prometheus monitoring:

**Recommended Dashboards**:

1. **Loki Dashboard** (ID: 13407)
   - Logs rate by service
   - Log volume over time
   - Error rate trends

2. **Prometheus 2.0 Stats** (ID: 3662)
   - Prometheus metrics overview
   - Query performance
   - Storage usage

**Manual Import Steps**:

1. Login to Grafana at `http://localhost:3000`
2. Navigate to **Dashboards** → **Import**
3. Enter dashboard ID (e.g., `13407`)
4. Select the appropriate datasource (Loki or Prometheus)
5. Click **Import**

### Custom Dashboards

Create custom dashboards to monitor:

- Application-specific metrics (games, requests, errors)
- Infrastructure health (memory, CPU via cAdvisor if added)
- Business metrics (user engagement, game completion rates)

**Example Panels**:

- **Request Rate**: `rate(http_requests_total[5m])`
- **Error Rate**: `rate(http_errors_total[5m])`
- **Response Time (p95)**: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))`
- **Game Success Rate**: `games_won_total / (games_won_total + games_lost_total)`

## Resource Management

### Memory Limits

Memory limits ensure services don't consume excessive resources:

Service | Memory Limit | Memory Reservation
--------|--------------|-------------------
Loki | 1 GB | 512 MB
Prometheus | 1 GB | 512 MB
Grafana | 512 MB | 256 MB
Promtail | 512 MB | 256 MB
Python App | 512 MB | 256 MB
Go App | 512 MB | 256 MB

**Memory Reservation**: Guaranteed minimum memory allocation

**Memory Limit**: Hard cap on memory usage (container will be OOM-killed if exceeded)

### Log Rotation

All services use Docker's JSON file logging driver with rotation:

```yaml
logging:
  driver: json-file
  options:
    max-size: "10m"    # Maximum size of each log file
    max-file: "3"      # Maximum number of log files to keep
```

**Total log storage per service**: ~30 MB (3 files × 10 MB)

### Health Checks

All services include health checks for monitoring and dependency management:

Parameter | Value | Description
----------|-------|------------
`interval` | 30s | Time between health checks
`timeout` | 10s | Max time to wait for response
`retries` | 3 | Failed checks before unhealthy
`start_period` | 5-10s | Grace period on startup

**Health Check Endpoints**:

- Loki: `http://localhost:3100/ready`
- Prometheus: `http://localhost:9090/-/healthy`
- Grafana: `http://localhost:3000/api/health`
- Promtail: `http://localhost:9080/ready`
- Python App: `http://localhost:5000/health`
- Go App: `http://localhost:8080/health`

## Screenshots

Screenshots demonstrating the monitoring setup and dashboards will be provided in the pull request comments.

**Required Screenshots**:

1. Prometheus Targets page showing all scraped services
2. Prometheus Graph showing application metrics
3. Grafana Loki Dashboard (ID: 13407)
4. Grafana Prometheus Dashboard (ID: 3662)
5. Custom application metrics visualization

## Troubleshooting

### Prometheus Not Scraping Targets

Check if services are healthy:

```bash
docker ps
```

Verify Prometheus targets:

```text
http://localhost:9090/targets
```

### Application Metrics Not Appearing

Verify metrics endpoint is accessible:

```bash
curl http://localhost:5000/metrics  # Python app
curl http://localhost:9090/metrics  # Go app (via external port)
```

### High Memory Usage

Monitor resource usage:

```bash
docker stats
```

Adjust memory limits in `docker-compose.yml` if needed.

### Grafana Datasource Connection Failed

Ensure Prometheus and Loki are running:

```bash
docker logs prometheus
docker logs loki
```

Check datasource configuration in Grafana:

- Navigate to **Configuration** → **Data Sources**
- Test connection for both Loki and Prometheus

## Best Practices

1. **Metric Cardinality**: Avoid high-cardinality labels (e.g., user IDs, timestamps)
2. **Query Optimization**: Use recording rules for frequently accessed queries
3. **Retention**: Adjust Prometheus retention based on storage capacity
4. **Alerting**: Set up Alertmanager for critical metric thresholds
5. **Backup**: Regularly backup Prometheus data volume
6. **Security**: Use authentication and TLS in production environments

## Summary

This monitoring setup provides:

- Comprehensive metrics collection from all services
- Centralized visualization with Grafana
- Resource management with memory limits and log rotation
- Health monitoring with automated dependency checks
- Application-specific metrics for business insights

For questions or issues, refer to the official documentation:

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
