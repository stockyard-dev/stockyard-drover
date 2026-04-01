# Stockyard Drover

**Server inventory and asset tracker — servers, domains, SaaS subscriptions, renewal dates, costs**

Part of the [Stockyard](https://stockyard.dev) family of self-hosted developer tools.

## Quick Start

```bash
docker run -p 9190:9190 -v drover_data:/data ghcr.io/stockyard-dev/stockyard-drover
```

Or with docker-compose:

```bash
docker-compose up -d
```

Open `http://localhost:9190` in your browser.

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `9190` | HTTP port |
| `DATA_DIR` | `./data` | SQLite database directory |
| `DROVER_LICENSE_KEY` | *(empty)* | Pro license key |

## Free vs Pro

| | Free | Pro |
|-|------|-----|
| Limits | 25 assets | Unlimited assets and reminders |
| Price | Free | $2.99/mo |

Get a Pro license at [stockyard.dev/tools/](https://stockyard.dev/tools/).

## Category

Operations & Teams

## License

Apache 2.0
