# Health Endpoints

The Penguin service exposes both authenticated and unauthenticated health
probes. Authenticated probes require the global bearer token when `http.auth`
is configured, while unauthenticated probes live under the `/_internal`
namespace and remain accessible without credentials.

## Authenticated Service Availability

- **Method:** `GET`
- **Path:** `/health`
- **Description:** Lightweight availability probe that validates request
  routing and (if configured) bearer authentication. Returns `200 OK` without
  touching downstream dependencies.
- **Successful Response:**

```json
{
  "status": "ok"
}
```

## Internal Health Check

- **Method:** `GET`
- **Path:** `/_internal/health`
- **Description:** Returns the application status together with database
  connectivity.
- **Successful Response:**

```json
{
  "status": "ok",
  "database": "ok"
}
```

- **Failure Response:**

```json
{
  "status": "degraded",
  "database": "error"
}
```

The endpoint returns `503 Service Unavailable` when the database probe fails.

## Prometheus Metrics

- **Method:** `GET`
- **Path:** `/_internal/metrics`
- **Description:** Prometheus exporter containing runtime metrics, including HTTP
  request counters/latency histograms and the Tencent transfer update job
  instrumentation. The handler emits standard Prometheus text exposition format.
