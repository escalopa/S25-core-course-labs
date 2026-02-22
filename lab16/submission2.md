# Lab 16: IPFS and Fleek — Submission

## Task 1: IPFS Gateway via Docker

### Setup

Ran the IPFS node using the `ipfs/kubo` Docker image (the current successor to the deprecated `ipfs/go-ipfs`):

```sh
docker run -d --name ipfs_host \
  -v /path/to/lab16:/export \
  -v ipfs_data:/data/ipfs \
  -p 8080:8080 \
  -p 4001:4001 \
  -p 5001:5001 \
  ipfs/kubo:release
```

Verified the container is running:

```sh
$ docker ps
CONTAINER ID   IMAGE             COMMAND                  STATUS          PORTS
e9eabd22de45   ipfs/kubo:release "/sbin/tini -- /usr/…"  Up 2 minutes    0.0.0.0:4001->4001/tcp, 0.0.0.0:5001->5001/tcp, 0.0.0.0:8080->8080/tcp   ipfs_host
```

Daemon startup confirmation from logs:

```
RPC API server listening on /ip4/0.0.0.0/tcp/5001
WebUI: http://127.0.0.1:5001/webui
Gateway server listening on /ip4/0.0.0.0/tcp/8080
Daemon is ready
```

### Node Identity

```
$ docker exec ipfs_host ipfs id
{
  "ID": "12D3KooWBcNCib1yEm8K8DzLyvp2eqdgB73DfeTHT7vKdPLU2p7k",
  "AgentVersion": "kubo/0.39.0/2896aed/docker",
  "Protocols": [
    "/ipfs/bitswap",
    "/ipfs/bitswap/1.0.0",
    "/ipfs/bitswap/1.1.0",
    "/ipfs/bitswap/1.2.0",
    "/ipfs/id/1.0.0",
    "/ipfs/id/push/1.0.0",
    "/ipfs/lan/kad/1.0.0",
    "/ipfs/ping/1.0.0",
    ...
  ]
}
```

### Connected Peers

```
$ docker exec ipfs_host ipfs swarm peers
/ip4/135.181.3.221/udp/4001/quic-v1/p2p/12D3KooWMH4hRLwnNMu6JDZCFRFqYBXEyo8bfYYoT4sqi2Nx48NS
/ip4/141.95.145.190/tcp/4001/p2p/12D3KooWKSMTgHEZWv82tVE51XSw5PAJLRE3rLvfVWy1nc319oiD
/ip4/152.70.127.248/udp/4001/webrtc-direct/p2p/12D3KooWHVHW3AFv6qXxecJTFjEDcWZ2Bnf6e3ptZKu4QYKijqhF
/ip4/176.215.114.161/tcp/48888/p2p/12D3KooWG2EbeXnYeKS5kBAt1g9eVPCJEbb18XqvzhwTQh6VYKnb
/ip4/217.210.25.123/tcp/48888/p2p/12D3KooWGxnNJyTY5Nv2DDz3uGyFGrhoPSWP3VWPRfRo7bTpWr96
/ip4/65.109.31.54/udp/4001/quic-v1/p2p/12D3KooWBjY5aKaFzBF4hUyv1TF5sfDfeSqaNSVcqnSDfAJUWAMu
/ip4/65.109.61.242/udp/4001/quic-v1/p2p/12D3KooWRcwwFBPJiQsmcqHvFpCw3bzPNBoQwhXCpsiKNeWbXGW2
/ip4/65.21.233.138/udp/4001/quic-v1/p2p/12D3KooWQtFYh18NLaykBs8A7t8uYL6cJtPXSLE3NT53DZQBPnDj
/ip4/66.151.34.140/tcp/4001/p2p/12D3KooWNWZfG3NueANPMcR7SAB2uYHYyB5oHtQtAbVZ6gqLbCXo
/ip4/95.84.138.239/tcp/48888/p2p/12D3KooWJKqDH5K6p5PpK1PGNdxPmSRSV2yRHzqMafejuV3cBbwW
```

**Total peers connected: 10**

### Bandwidth Statistics

```
$ docker exec ipfs_host ipfs stats bw
Bandwidth
TotalIn:  9.2 MB
TotalOut: 1.8 MB
RateIn:   2.2 kB/s
RateOut:  0 B/s
```

The node received ~9.2 MB during DHT bootstrap and peer discovery. Outbound rate is near zero since no content is being actively served yet.

### File Upload

Uploaded `lab16/index.html` (the DevOps course landing page):

```sh
$ docker exec ipfs_host ipfs add /export/index.html
added QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj index.html
 8.88 KiB / 8.88 KiB  100.00%
```

**CID: `QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj`**

### Gateway Verification URLs

The file can be accessed via public IPFS gateways using the CID above:

| Gateway | URL |
|---|---|
| IPFS.io | https://ipfs.io/ipfs/QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj |
| Cloudflare | https://cloudflare-ipfs.com/ipfs/QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj |
| Local gateway | http://127.0.0.1:8080/ipfs/QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj |

> Note: Public gateway availability may vary depending on network propagation time and gateway load.

---

## Task 2: Fleek Deployment

### What is IPFS?

IPFS (InterPlanetary File System) is a peer-to-peer, content-addressed distributed file system. Unlike traditional HTTP where files are identified by their *location* (a server URL), IPFS identifies files by their *content* (a cryptographic hash — the CID). This means:
- Any node holding the file can serve it — no single point of failure
- Content cannot be silently modified (the hash would change)
- Files are automatically deduplicated across the network

### What is Fleek?

Fleek is a Web3 hosting platform built on IPFS and other decentralized protocols. It provides a developer-friendly interface to deploy static websites and apps to IPFS with features like:
- GitHub integration for automatic deployments on push
- Custom domains with automatic SSL
- IPNS (InterPlanetary Name System) for a stable, updateable link to the latest deployment
- CDN-like global edge network for fast IPFS content delivery

### Deployment

The `lab16/index.html` page (DevOps Engineering course landing page) was deployed to Fleek pointing at the repository.

**Fleek project configuration:**
- Source: GitHub repository `escalopa/S25-core-course-labs`
- Build directory: `lab16/`
- Publish directory: `lab16/`
- Framework: Other (static HTML)

**Deployed IPFS CID:** `QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj`

**IPFS Link:** https://ipfs.io/ipfs/QmZawNCTzWst7BViJFt6VZu92NsfoHh6jxTRcynBh97cGj

**Fleek Domain:** https://escalopa-s25-core-course-labs.on-fleek.app
