# ztsfc_proxy

ztsfc_proxy is a Go-based, TLS-first reverse proxy that acts as a Policy Enforcement Point (PEP) for a Zero Trust-style service mesh edge. It terminates incoming TLS (optionally mTLS), chooses a backend service based on the request SNI, and forwards traffic to the configured target while enforcing TLS policy (TLS 1.3, CRL checks, optional client cert validation). The codebase is intentionally small and organized into focused packages so that TLS policy, logging, and proxying behavior are easy to audit.

## High-level request flow
- The frontend HTTP server listens on the configured address with server-side TLS enabled.
- During the TLS handshake, the SNI determines which server certificate is presented.
- If `frontend.tls.client_auth` is enabled, the client certificate is validated against configured CAs and the CRL.
- The PEP handler receives the HTTP request and looks up the target backend using the same SNI value.
- A reverse proxy forwards the request to the target service (HTTP or HTTPS).
- For HTTPS backends, the client-side TLS config is built from the `services.tls` block (certs, CAs, CRL).
- The response is returned to the client with HSTS headers applied.
- Requests and responses are logged with a per-request hash to correlate events in logs.

## Key features
- TLS 1.3 only for both incoming and outgoing connections.
- Optional mTLS for client authentication on the frontend.
- CRL verification for client and server certificates.
- SNI-based certificate selection and backend routing.
- Data plane logging with request/response correlation hashes.
- HSTS response headers are applied to proxied responses.
- HTTP/2 is enforced via the TLS NextProtos configuration.

## Architecture overview
ztsfc_proxy is split into data plane and control plane concepts:
- Data plane: request handling, reverse proxying, and request/response logging.
- Control plane: placeholder for future policy decisions (PDP), not currently active.

Major runtime components:
- Frontend server: configured TLS endpoint that accepts client connections.
- PEP handler: main HTTP handler that chooses a backend based on SNI and proxies traffic.
- Service pool: pre-parsed list of backend URLs keyed by SNI.
- TLS utilities: shared logic for loading certificates, CAs, and CRLs for both sides of the proxy.

## Detailed request lifecycle
1. `cmd/ztsfc_proxy/main.go` loads the YAML config and initializes the system logger.
2. `internal/frontend` builds the server TLS config from `frontend.tls`.
3. The frontend HTTP server starts with a mux that routes all requests to the PEP.
4. The PEP reads `r.TLS.ServerName` (SNI) and finds the matching service in the service pool.
5. The request is hashed (method, URL, headers, body, timestamp) for log correlation.
6. A reverse proxy is created for the target service URL.
7. The outbound transport is selected based on the target scheme:
   - `https`: uses the client TLS config from `services.tls`.
   - `http`: uses a non-TLS transport.
8. The request is forwarded and the response is decorated with HSTS headers.
9. The PEP logs both the forward and response events with the same hash.

## Configuration overview
Configuration is loaded from a YAML file (default: `configs/config.yml`) and mapped to `internal/configs`. Main sections:
- `frontend`: listen address and server-side TLS settings.
  - `addr`: host:port for the frontend server.
  - `tls.certificates`: map of SNI -> cert/key pair used to present the correct server certificate.
  - `tls.client_auth`: enable/disable client certificate validation.
  - `tls.cas`: list of CAs accepted for client certificates.
  - `tls.crl`: CRL used to reject revoked client certificates.
- `data_plane_logger`: output path for request/response logs.
- `control_plane_logger`: output path (wired for future PDP integration).
- `services`: settings for outbound connections to backends.
  - `tls.certificates`: client certs used when connecting to HTTPS backends.
  - `tls.cas`: list of CAs accepted for backend server certificates.
  - `tls.crl`: CRL used to reject revoked backend certificates.
  - `service_pool`: map of SNI -> `service_url` target.

See `configs/config.yml` for a real configuration and `configs/example_config copy.yml` for a template-style example.

### Configuration details and expectations
- `frontend.addr` must be in `host:port` format and should match the TLS certificate names used in `frontend.tls.certificates`.
- `frontend.tls.certificates` is keyed by SNI; the same key is used for routing to a backend.
- `frontend.tls.client_auth` controls whether the proxy requires client certificates. If enabled, clients must present certs signed by `frontend.tls.cas` and not revoked by `frontend.tls.crl`.
- `services.tls.certificates` provides client certificates for mTLS when connecting to HTTPS backends. The map is keyed by SNI, which is derived from the backend URL or presented via the TLS handshake.
- `services.service_pool` is the authoritative routing table. Missing or unknown SNI values result in a 404 from the proxy.

### Minimal config example
```yaml
frontend:
  addr: "proxy.example.com:443"
  tls:
    certificates:
      proxy.example.com:
        cert_file: "/etc/letsencrypt/live/proxy.example.com/fullchain.pem"
        key_file: "/etc/letsencrypt/live/proxy.example.com/privkey.pem"
    client_auth: true
    cas:
      - "/etc/ssl/certs/client_ca.crt"
    crl: "/etc/ssl/crl/client_ca.crl"

data_plane_logger:
  output: "./logs/ztsfc_proxy_dp.log"

control_plane_logger:
  output: "./logs/ztsfc_proxy_cp.log"

services:
  tls:
    certificates:
      backend.internal:
        cert_file: "/etc/ssl/certs/pep_client.crt"
        key_file: "/etc/ssl/private/pep_client.key"
    cas:
      - "/etc/ssl/certs/backend_ca.crt"
    crl: "/etc/ssl/crl/backend_ca.crl"
  service_pool:
    proxy.example.com:
      service_url: "https://backend.internal:8443"
```

## Repository structure
- `cmd/ztsfc_proxy/main.go`: application entrypoint; loads config, sets up loggers, starts frontend.
- `internal/configs`: YAML-backed config structs for frontend, services, loggers, and TLS.
- `internal/frontend`: HTTP server wiring, TLS setup, and handler registration.
- `internal/pep`: Policy Enforcement Point implementation; SNI routing and reverse proxy behavior.
- `internal/pdp`: placeholder Policy Decision Point (currently not wired).
- `internal/service`: service pool model; parses target URLs and initializes client TLS.
- `internal/security/tlsutil`: TLS helpers (cert maps, CA pools, CRL validation, mTLS).
- `internal/security/hashutil`: request hashing for log correlation.
- `internal/logger`: system/data/control plane loggers.
- `internal/web`: simple HTTP error response helpers.
- `configs`: example and active YAML configurations.
- `build`: output directory for the compiled binary.

## Logging and observability
- System logger (`internal/logger/system_logger.go`) writes startup and configuration diagnostics.
- Data plane logger logs request forwarding and response delivery with a short hash.
- The hash is derived from method, headers, body, and a timestamp to help correlate logs.
- Error logging from the reverse proxy is routed to the data plane logger when available.

## Error handling behavior
- Unknown SNI: returns a 404 response via `internal/web`.
- Unsupported backend scheme (non-HTTP/HTTPS): returns a 501 response.
- Internal errors creating the reverse proxy: returns a 500 response.
- TLS/CRL errors during startup cause initialization to fail.

## Build and run
This project uses a simple Makefile for building the binary. The default build output is `./build/ztsfc_proxy`.

Build:
```bash
make build
```

Start:
```bash
./build/ztsfc_proxy -c ./configs/config.yml
```

Notes:
- By default, the proxy expects valid certificate and CRL paths. Use the example config as a template.
- Use `-d` or `-e` to adjust the system log level, and `-o` to redirect system logs.

## CLI flags
The entrypoint is `cmd/ztsfc_proxy/main.go`, which supports:
- `-c`: path to YAML config (default `./configs/config.yml`)
- `-o`: system logger output (default `stdout`)
- `-d`: debug log level (mutually exclusive with `-e`)
- `-e`: error-only log level (mutually exclusive with `-d`)
- `-f`: system logger format (`text` or `json`)

## Security notes
- TLS 1.3 is enforced by configuration (MinVersion/MaxVersion).
- CRL validation is performed on both inbound client certs and outbound server certs.
- If `client_auth` is disabled, client certificates are not required and CRL checking on the client side is not used.
- The proxy does not currently implement dynamic policy decisions; all routing is static via SNI.

## Notes on current behavior
- Only the PEP is active; PDP wiring is currently commented out in `internal/frontend/frontend.go`.
- Backend selection is strictly based on SNI from the client TLS handshake.
- HTTP/2 is enforced via `NextProtos` in TLS configuration.
