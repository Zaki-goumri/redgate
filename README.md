# Redgate

Redis regulator/proxy in Go with modular architecture.

## Current status

Project scaffold is initialized with module directories and local dev tooling.

## Quick start

### Run locally

```bash
make run
```

### Build

```bash
make build
```

### Test

```bash
make test
```

### Run with Docker Compose (2 Redis backends + redgate)

```bash
make docker-up
```

## Layout

```text
cmd/redgate          main entrypoint
internal/resp        RESP protocol parsing
internal/proxy       client TCP handling
internal/router      routing decisions
internal/dispatcher  command dispatching
internal/pool        backend connection pools
internal/merger      fan-out response merge
internal/namespace   tenant key isolation
internal/admin       admin interface
pkg/config           shared typed config package
pkg/metrics          shared metrics package
```

## Next implementation order

1. `internal/resp`
2. `internal/pool`
3. `internal/proxy`
4. `internal/router`
5. `internal/dispatcher`
6. `internal/merger`
7. `internal/namespace`
8. `pkg/config`
9. `pkg/metrics`
10. `internal/admin`
