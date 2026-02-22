---
name: postgres
host: docker-host-01
criticality: high
storage:
  volume: data
  participates_in_host_backup: true
dependencies:
  - dns
  - network
---

# Overview

Postgres database for internal services.

# Restore

1. Restore volume snapshot
2. Start container
3. Validate replication