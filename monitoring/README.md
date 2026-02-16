# Monitoring Stack

## Overview

Complete logging and monitoring solution using Promtail, Loki, and Grafana (PLG Stack) for Docker containers.

## Quick Start

```bash
cd monitoring
docker compose up -d
```

**Access**:

- **Grafana**: <http://localhost:3000> (admin/admin)
- **Python App**: <http://localhost:5000>
- **Go App**: <http://localhost:8080>

## Using Makefile

```bash
make monitoring-up       # Start stack
make monitoring-down     # Stop stack
make monitoring-logs     # View logs
make monitoring-restart  # Restart services
make monitoring-clean    # Clean volumes
```

## Components

Service | Port | Purpose
--- | --- | ---
Loki | 3100 | Log aggregation
Promtail | 9080 | Log collection
Grafana | 3000 | Visualization
Python App | 5000 | Moscow time display
Go App | 8080 | Wordle game

## Viewing Logs

1. Open Grafana: <http://localhost:3000>
2. Login: admin/admin
3. Go to **Explore** (compass icon)
4. Select **Loki** datasource
5. Query logs:

```logql
{app="moscow-time"}     # Python app logs
{app="wordle-game"}     # Go app logs
{container_name=~".*"}  # All logs
```

## Configuration Files

- `docker-compose.yml` - Service definitions
- `loki-config.yaml` - Loki configuration
- `promtail-config.yaml` - Log collection rules
- `grafana/provisioning/` - Auto-configured datasources

## Features

- ✅ Automatic log collection from all containers
- ✅ Label-based log organization
- ✅ 7-day log retention
- ✅ Pre-configured Loki datasource
- ✅ Real-time log streaming
- ✅ Powerful LogQL queries

## Documentation

See [LOGGING.md](./LOGGING.md) for detailed documentation including:

- Architecture overview
- Component details
- LogQL query examples
- Troubleshooting guide
- Best practices

## Testing

Generate logs and view them in Grafana:

```bash
# Generate Python app logs
for i in {1..10}; do curl http://localhost:5000; done

# Generate Go app logs
curl http://localhost:8080

# View in Grafana Explore
# Query: {container_name=~".*"}
```

## Troubleshooting

**No logs appearing**:

```bash
# Check Promtail status
docker compose logs promtail

# Verify Loki is ready
curl http://localhost:3100/ready

# Check available labels
curl http://localhost:3100/loki/api/v1/labels
```

**Grafana issues**:

```bash
# Check Grafana logs
docker compose logs grafana

# Verify Loki datasource
curl http://localhost:3000/api/datasources
```

## Log Retention

Default retention: 7 days

To change retention, edit `loki-config.yaml`:

```yaml
limits_config:
  retention_period: 168h  # Change this value
```
