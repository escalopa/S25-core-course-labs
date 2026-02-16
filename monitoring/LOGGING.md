# Logging Stack Documentation

## Overview

This logging stack uses **Promtail**, **Loki**, and **Grafana** to collect, aggregate, and visualize logs from Docker containers. This is a lightweight, cost-effective alternative to the ELK stack (Elasticsearch, Logstash, Kibana).

## Architecture

```text
Docker Containers (Python + Go Apps)
            ↓
    Docker JSON Logs
            ↓
  Promtail (Log Collector)
            ↓
    Loki (Log Aggregation)
            ↓
 Grafana (Visualization)
```

## Components

### 1. Loki

**Purpose**: Log aggregation and storage system

**Key Features**:
- Indexes only metadata (labels), not full-text
- Cost-effective storage (stores logs compressed)
- Horizontally scalable
- Query logs using LogQL (similar to PromQL)

**Configuration** (`loki-config.yaml`):
- **Port**: 3100
- **Storage**: Filesystem-based (local volume)
- **Retention**: 7 days (168 hours)
- **Schema**: BoltDB shipper with filesystem object store

**Why Loki**:
- More efficient than Elasticsearch (indexes labels, not content)
- Perfect for container logs
- Designed for cloud-native environments

### 2. Promtail

**Purpose**: Log collector agent that ships logs to Loki

**Key Features**:
- Discovers targets automatically
- Attaches labels to log streams
- Supports multiple input sources
- Tails log files in real-time

**Configuration** (`promtail-config.yaml`):
- **Port**: 9080
- **Docker Socket**: `/var/run/docker.sock` (reads container logs)
- **Jobs**: Docker containers + system logs

**How It Works**:
1. Connects to Docker socket
2. Discovers all running containers with label `service`
3. Reads logs from Docker's JSON log driver
4. Attaches labels: container_name, app, service, image
5. Ships logs to Loki

**Label Strategy**:
- `container_name`: Container identifier
- `app`: Application type (moscow-time, wordle-game)
- `service`: Service type (python-app, go-app)
- `stream`: stdout or stderr

### 3. Grafana

**Purpose**: Web-based visualization and dashboarding

**Key Features**:
- Query and explore logs
- Create custom dashboards
- Set up alerts
- Multiple data source support

**Configuration**:
- **Port**: 3000
- **Default Credentials**: admin/admin (change on first login)
- **Datasource**: Loki (auto-provisioned)
- **Storage**: Persistent volume for dashboards

**Access**: <http://localhost:3000>

## Applications

### Python App (Moscow Time Display)

- **Port**: 5000
- **Image**: `escalopax/moscow-time-app:latest`
- **Logs**: Uvicorn server logs, FastAPI request logs
- **Labels**: `app=moscow-time`, `service=python-app`

### Go App (Wordle Game)

- **Port**: 8080
- **Image**: `escalopax/wordle-game:latest`
- **Logs**: Game creation logs, HTTP request logs
- **Labels**: `app=wordle-game`, `service=go-app`

## Quick Start

### Start the Stack

```bash
cd monitoring
docker compose up -d
```

### Access Services

- **Grafana**: <http://localhost:3000> (admin/admin)
- **Python App**: <http://localhost:5000>
- **Go App**: <http://localhost:8080>
- **Loki**: <http://localhost:3100> (API)

### Stop the Stack

```bash
docker compose down
```

### View Logs

```bash
docker compose logs -f
```

## Using Grafana

### 1. Login

1. Open <http://localhost:3000>
2. Login with: username `admin`, password `admin`
3. Change password (or skip)

### 2. Explore Logs

1. Click **Explore** (compass icon in sidebar)
2. Select **Loki** datasource
3. Use LogQL to query logs

### 3. Example Queries

**View all logs from Python app**:

```logql
{app="moscow-time"}
```

**View all logs from Go app**:

```logql
{app="wordle-game"}
```

**View logs from specific container**:

```logql
{container_name="moscow-time-app"}
```

**Filter by log level (if structured)**:

```logql
{app="moscow-time"} |= "ERROR"
```

**View stderr logs only**:

```logql
{stream="stderr"}
```

**Count requests per minute**:

```logql
rate({app="moscow-time"}[1m])
```

### 4. Create Dashboard

1. Click **Dashboards** → **New** → **New Dashboard**
2. Add panel
3. Select Loki datasource
4. Enter LogQL query
5. Configure visualization
6. Save dashboard

## Log Flow Diagram

```text
1. App writes logs to stdout/stderr
   ↓
2. Docker captures logs (JSON log driver)
   ↓
3. Promtail reads from Docker socket
   ↓
4. Promtail adds labels (container, app, service)
   ↓
5. Promtail pushes to Loki
   ↓
6. Loki indexes labels and stores logs
   ↓
7. Grafana queries Loki using LogQL
   ↓
8. User views logs in Grafana UI
```

## Configuration Details

### Loki Storage

- **Type**: Filesystem-based
- **Location**: Docker volume `loki-data`
- **Retention**: 7 days
- **Compression**: Enabled

### Promtail Discovery

**Docker Service Discovery**:
- Automatically discovers containers with label `service`
- Refreshes every 5 seconds
- Attaches metadata as labels

**Log Parsing**:
- Reads Docker JSON logs
- Preserves timestamps
- Separates stdout/stderr

### Grafana Provisioning

**Datasources**: Loki configured automatically on startup

**Dashboards**: Can be added to `grafana/provisioning/dashboards/`

## Troubleshooting

### Loki not receiving logs

Check Promtail logs:

```bash
docker compose logs promtail
```

Verify Docker socket is accessible:

```bash
docker compose exec promtail ls -la /var/run/docker.sock
```

### Grafana can't connect to Loki

Check network connectivity:

```bash
docker compose exec grafana ping loki
docker compose exec grafana curl http://loki:3100/ready
```

### No logs appearing in Grafana

1. Verify containers are running: `docker compose ps`
2. Check Promtail is scraping: `docker compose logs promtail`
3. Query Loki directly: `curl http://localhost:3100/loki/api/v1/labels`
4. Ensure apps have proper labels in docker-compose.yml

### Can't access Grafana

- Verify port 3000 is not in use
- Check Grafana logs: `docker compose logs grafana`
- Try accessing: <http://localhost:3000>

## Best Practices

### 1. Label Strategy

Always add meaningful labels to containers:

```yaml
labels:
  app: "service-name"
  service: "service-type"
  environment: "dev"
```

### 2. Log Retention

- Default: 7 days (good for development)
- Production: Adjust based on compliance requirements
- Consider using object storage (S3, GCS) for production

### 3. Query Performance

- Use specific labels in queries
- Avoid regex when possible
- Use time ranges to limit data scanned

### 4. Security

- Change default Grafana password
- Use authentication for Loki in production
- Limit network access to monitoring services
- Don't expose Loki port in production

## Testing the Stack

### 1. Generate Logs

**Python App**:

```bash
# Access the app multiple times
for i in {1..10}; do curl http://localhost:5000; done
```

**Go App**:

```bash
# Play some Wordle games
curl http://localhost:8080
```

### 2. View Logs in Grafana

1. Go to Grafana Explore
2. Query: `{container_name=~".*"}`
3. See logs from all containers
4. Filter by specific app or service

### 3. Verify Promtail is Collecting

```bash
# Check Promtail targets
curl http://localhost:9080/targets
```

## Stack Components Summary

Component | Port | Purpose | Storage
--- | --- | --- | ---
Loki | 3100 | Log aggregation | Volume: loki-data
Promtail | 9080 | Log collection | None (agent)
Grafana | 3000 | Visualization | Volume: grafana-data
Python App | 5000 | Moscow time | None
Go App | 8080 | Wordle game | None

## LogQL Examples

**Basic Queries**:

```logql
# All logs from Python app
{app="moscow-time"}

# All logs from Go app
{app="wordle-game"}

# Combined
{app=~"moscow-time|wordle-game"}

# Last 5 minutes
{app="moscow-time"}[5m]
```

**Filtering**:

```logql
# Contains "error"
{app="moscow-time"} |= "error"

# Case-insensitive
{app="moscow-time"} |~ "(?i)error"

# Exclude health checks
{app="moscow-time"} != "/health"
```

**Aggregations**:

```logql
# Log rate per second
rate({app="moscow-time"}[1m])

# Count logs per app
count_over_time({app=~".*"}[5m])
```

## Advantages of This Stack

**vs ELK Stack**:
- **Simpler**: Fewer components to manage
- **Lighter**: Lower resource usage
- **Faster**: Better query performance for recent logs
- **Cheaper**: Minimal indexing reduces storage costs

**vs Cloud Solutions**:
- **Free**: No per-GB costs
- **Local**: No data leaves your infrastructure
- **Flexible**: Full control over retention and configuration

## Production Considerations

For production deployments:

1. **High Availability**: Run multiple Loki instances
2. **Object Storage**: Use S3/GCS instead of filesystem
3. **Security**: Enable authentication, use TLS
4. **Retention**: Adjust based on compliance needs
5. **Alerting**: Configure Loki ruler for log-based alerts
6. **Scaling**: Add more Promtail agents as needed

## Useful Resources

- [Grafana Loki Documentation](https://grafana.com/docs/loki/latest/)
- [LogQL Query Language](https://grafana.com/docs/loki/latest/logql/)
- [Promtail Configuration](https://grafana.com/docs/loki/latest/clients/promtail/configuration/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)

## Screenshots

Screenshots demonstrating the successful operation of the logging stack are included in the pull request comments, showing:

- Grafana dashboard with logs from both applications
- Loki datasource configuration
- LogQL queries in action
- Container log streams
- Real-time log visualization

## Conclusion

This logging stack provides:

- **Centralized Logging**: All container logs in one place
- **Easy Queries**: LogQL for powerful log searching
- **Visualization**: Grafana for dashboards and exploration
- **Low Overhead**: Efficient resource usage
- **Production Ready**: Scalable architecture

The stack is designed to be simple yet powerful, perfect for monitoring containerized applications in both development and production environments.
