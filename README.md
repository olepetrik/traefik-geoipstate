# Traefik GeoIP State Plugin

A Traefik middleware plugin that adds geolocation headers to incoming requests based on client IP addresses. It detects the user's country and state/subdivision using MaxMind GeoIP2 database.

## Features

- Automatically adds `X-User-Country` and `X-User-State` headers to requests
- Uses MaxMind GeoIP2 City database for accurate geolocation
- Supports both `X-Real-IP` header and `RemoteAddr` for IP detection
- Lightweight and easy to configure

## Headers Added

- `X-User-Country`: ISO country code (e.g., "US", "CA", "GB") or "UNKNOWN"
- `X-User-State`: ISO subdivision/state code (e.g., "CA", "NY", "TX") or "UNKNOWN"

## Installation

### Static Configuration

#### YAML (traefik.yml)

```yaml
experimental:
  plugins:
    geoipstate:
      moduleName: github.com/olepetrik/traefik-geoipstate
      version: v0.3.0
```

#### TOML (traefik.toml)

```toml
[experimental.plugins.geoipstate]
  moduleName = "github.com/olepetrik/traefik-geoipstate"
  version = "v0.3.0"
```

#### Command Line

```bash
--experimental.plugins.geoipstate.modulename=github.com/olepetrik/traefik-geoipstate
--experimental.plugins.geoipstate.version=v0.3.0
```

## Configuration

### File Provider (dynamic.yml)

```yaml
http:
  middlewares:
    geoip-detector:
      plugin:
        geoipstate:
          database: /path/to/GeoLite2-City.mmdb

  routers:
    my-router:
      rule: "Host(`example.com`)"
      service: my-service
      middlewares:
        - geoip-detector

  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://backend:8080"
```

### Docker Labels

```yaml
services:
  my-app:
    image: myapp:latest
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.my-app.rule=Host(`example.com`)"
      - "traefik.http.routers.my-app.middlewares=geoip-detector"
      - "traefik.http.middlewares.geoip-detector.plugin.geoipstate.database=/path/to/GeoLite2-City.mmdb"
```

### Docker Compose Example

```yaml
version: '3'

services:
  traefik:
    image: traefik:v2.11
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--entrypoints.web.address=:80"
      - "--experimental.plugins.geoipstate.modulename=github.com/olepetrik/traefik-geoipstate"
      - "--experimental.plugins.geoipstate.version=v0.3.0"
    ports:
      - "80:80"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./GeoLite2-City.mmdb:/etc/traefik/GeoLite2-City.mmdb:ro

  whoami:
    image: traefik/whoami
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.whoami.rule=Host(`whoami.localhost`)"
      - "traefik.http.routers.whoami.middlewares=geoip"
      - "traefik.http.middlewares.geoip.plugin.geoipstate.database=/etc/traefik/GeoLite2-City.mmdb"
```

## Configuration Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `database` | string | Yes | Path to the MaxMind GeoIP2 City database file (.mmdb format) |

## Getting the GeoIP2 Database

This plugin requires a MaxMind GeoLite2 City database. You can obtain it for free:

1. Create a free account at [MaxMind](https://www.maxmind.com/en/geolite2/signup)
2. Download the GeoLite2 City database in MMDB format
3. Place the `.mmdb` file in a location accessible to Traefik

### Automatic Updates with Docker

```bash
# Using MaxMind's GeoIP Update tool
docker run --rm -v $(pwd):/usr/share/GeoIP \
  -e GEOIPUPDATE_ACCOUNT_ID=your_account_id \
  -e GEOIPUPDATE_LICENSE_KEY=your_license_key \
  -e GEOIPUPDATE_EDITION_IDS=GeoLite2-City \
  maxmindinc/geoipupdate:latest
```

## How It Works

1. The middleware extracts the client IP address from:
   - `X-Real-IP` header (if present)
   - `RemoteAddr` (fallback)

2. Looks up the IP in the GeoIP2 database

3. Adds headers to the request:
   - `X-User-Country`: Country ISO code
   - `X-User-State`: State/subdivision ISO code (primarily for US states)

4. If IP cannot be parsed or found in database, headers are set to "UNKNOWN"

5. The modified request is passed to the next handler

## Use Cases

- Content localization based on user location
- Geographic access control
- Analytics and tracking
- A/B testing by region
- Compliance with regional regulations

## Testing

```bash
# Test with curl
curl -H "X-Real-IP: 8.8.8.8" http://your-domain.com

# Check response headers
curl -I http://your-domain.com
```

## License

This project is open source and available under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/olepetrik/traefik-geoipstate).
