# Architecture

This repo uses:

- Hexagonal-style separation between core domain, ports, and adapters.
- Go Fiber as the HTTP adapter.
- PostgreSQL as primary storage.
- Redis as cache.
- Kafka adapter wired for event publishing.
- Next.js as the main frontend consuming the API.
