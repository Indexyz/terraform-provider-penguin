# Tencent Cloud CVM API

The Penguin service exposes a RESTful interface for provisioning and managing
Tencent Cloud CVM instances. All endpoints are protected by the legacy HTTP
Bearer token (see `http.auth` configuration). Requests MAY also include a JWT in
the same `Authorization` header. When both credentials are present (e.g.
`Authorization: Bearer <legacy>, Bearer <jwt>`), the legacy token allows access
while the JWT supplies additional provisioning limits.

## Authentication

- **Legacy token** (required when `http.auth` is configured): continue to send
  the existing bearer token to pass authentication.
- **JWT (optional):** when provided, the service validates signature, expiry,
  and claims. A valid JWT can replace the legacy token and enables the
  fine-grained limits described below. Invalid or expired JWTs yield
  `401 Unauthorized`.
- **Project enforcement:** if the JWT contains `projectId`, every
  `/tencentcloud/...` operation checks the target resource against that project
  and rejects mismatches with
  `403 Forbidden`.

- **Base URL:** `http://<host>:<port>`
- **Namespace:** `/tencentcloud`
- **Content Type:** `application/json`

## Endpoints

### List Zones

- **Method:** `GET`
- **Path:** `/tencentcloud/zones`
- **Description:** Returns the list of Tencent Cloud availability zones accessible to the configured account, grouped by region. The response is sorted lexicographically by region and zone identifier.

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{
  "zones": [
    {
      "region": "ap-guangzhou",
      "regionName": "华南地区(广州)",
      "zone": "ap-guangzhou-1",
      "zoneName": "Guangzhou 1",
      "state": "AVAILABLE"
    },
    {
      "region": "ap-singapore",
      "regionName": "Southeast Asia(Singapore)",
      "zone": "ap-singapore-3",
      "zoneId": "sg-3",
      "state": "UNAVAILABLE"
    }
  ]
}
```

#### Error Codes

- `500 Internal Server Error` – Tencent Cloud API returned an error or responded with an unexpected payload.

### Select Bandwidth Package

- **Method:** `GET`
- **Path:** `/tencentcloud/bandwidth-packages`
- **Description:** Returns the schedulable shared bandwidth package in the region with the highest available capacity. The service queries Tencent Cloud `DescribeBandwidthPackages` with filters `tag:penguin=schedulable` and `network-type=<networkType>` (defaults to `BGP` when omitted). Candidates must be in a usable status and have fewer than 200 bound resources. The response includes the bandwidth package ID and `availableCount = 200 - len(ResourceSet)`; ties are broken by choosing the lexicographically smallest ID.

#### Query Parameters

- `region` (required): Tencent Cloud region, e.g. `ap-guangzhou`.
- `networkType` (optional): Tencent Cloud network type, e.g. `BGP` or `CMCC` (defaults to `BGP`).

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{ "id": "bwp-123", "availableCount": 190 }
```

#### Error Codes

- `400 Bad Request` – missing or invalid query parameters.
- `502 Bad Gateway` – Tencent Cloud API query failed (message includes “failed to query bandwidth packages”).
- `503 Service Unavailable` – no schedulable bandwidth package available (message includes “no schedulable bandwidth package available”).

### Create Virtual Machine

- **Method:** `POST`
- **Path:** `/tencentcloud/vms`
- **Description:** Provision a new CVM instance. Returns the generated virtual
  machine identifier composed of the requested zone and Tencent instance ID.
  When `projectId` is omitted or set to `null`, the Tencent Cloud default
  project (ID `0`) is used. Optionally set `period` (in months) to control the
  prepaid term; the service defaults to `1` month when unset. The same period
  value is forwarded to Tencent Cloud for subsequent renewals when auto-renew
  is enabled. Use the status endpoint after creation to inspect the prepaid
  expiration.
  When `rootLoginPassword` is omitted, Penguin generates a 12-character
  alphanumeric password and persists it for later retrieval via the status
  endpoint. Provide `elasticIpId` to attach an existing elastic public IP
  immediately after provisioning instead of allocating a new address.

When a valid JWT accompanies the request, the following claims are enforced:

- `maxTransferKB`: when non-negative, `totalTransfer` must be `>= 0` and less
  than or equal to the claim value. JWTs that specify `maxTransferKB >= 0`
  disallow `totalTransfer = -1`.
- `allowedInstanceTypes` / `allowedZones`: when the arrays are non-empty, the
  request must use one of the listed instance types or zones.
- `maxBandwidthMbps`: when specified and the request provides a
  `bandwidthLimit`, the requested Mbps must not exceed the claim value.
- `projectId`: overrides the request body `projectId` (and is persisted in the
  Tencent API call) regardless of the user supplied value.

When `sharedBandwidthPackageId` is omitted or空字符串，服务会自动调用
`DescribeBandwidthPackages(Limit=100, NetworkType=BGP)`，仅挑选带标签
`penguin: schedulable`、状态为 `AVAILABLE` 且已绑定资源数小于 200 的带宽包。
若存在多个候选，按照 `BandwidthPackageId` 升序取第一个。查询失败会返回
`502` 并提示 “failed to query bandwidth packages”，缺少可用带宽包时返回
`503` 并提示 “no schedulable bandwidth package available”。

Requests without a JWT behave identically to previous releases.

#### Request Body

```json
{
  "name": "test",
  "zone": "ap-guangzhou-6",
  "instanceType": "SA2.MEDIUM2",
  "securityGroup": "sg-5hilszwp",
  "systemImage": "img-7efla8nv",
  "vpcId": "vpc-oahbq6lh",
  "subnetId": "subnet-95tfs6am",
  "privateIpAddress": "10.0.0.25",
  "systemDiskSize": 20,
  "elasticIpId": "eip-8i1abc23",
  "chargeType": "POSTPAID_BY_HOUR",
  "rootLoginPassword": "ChangeMe123",
  "totalTransfer": 1048576,
  "projectId": 12345,
  "period": 12,
  "cloudInitData": "#cloud-config\npackage_update: true\n",
  "autoRenew": false
}
```

#### Successful Response

- **Code:** `201 Created`
- **Body:**

```json
{ "id": "7f1f1e97-7f45-4c4b-94b4-2f1c248a8b1e" }
```

#### Error Codes

- `400 Bad Request` – payload validation failed (missing fields, invalid FQDN,
  cloud-init payload exceeds 16 KB, etc.).
- `500 Internal Server Error` – upstream Tencent Cloud request failed.

When providing `elasticIpId`, omit `sharedBandwidthPackageId` and
`bandwidthLimit`. The service will bind the supplied elastic IP to the newly
provisioned instance once creation completes.

Set `privateIpAddress` to request a specific private address within the
subnet. When omitted, Tencent Cloud allocates the next available address. Use
the status endpoint to read the current `privateIps` values after creation.

Set `chargeType` to `POSTPAID_BY_HOUR` for postpaid billing; omit the field to
use the default `PREPAID` behaviour. When using a postpaid charge type, the
`period` and `autoRenew` settings are ignored.

When `autoRenew` is omitted or set to `false`, the service disables Tencent
Cloud automatic renewal for the prepaid instance. Setting it to `true` enables
automatic renewal (Tencent Cloud `NOTIFY_AND_AUTO_RENEW`).

Traffic quotas are expressed in kilobytes. A `totalTransfer` value of `-1`
indicates the instance does not enforce a transfer cap. The service tracks
usage in the database and initializes new instances with `usedTransfer`,
`txTransfer`, and `rxTransfer` set to `0` KB.

### Create Elastic IP

- **Method:** `POST`
- **Path:** `/tencentcloud/eips`
- **Description:** Allocates a new elastic public IP address in the specified
  region. The service applies Tencent Cloud defaults for the address type
  (`EIP`), provider (`BGP`), and internet charge mode (`BANDWIDTH_PACKAGE`).
  When `sharedBandwidthPackageId` is omitted, Penguin automatically queries
  schedulable shared bandwidth packages in the region (tagged
  `tag:penguin=schedulable`) and attaches the lowest-ID package with available
  capacity.

#### Request Body

```json
{
  "region": "ap-guangzhou",
  "bandwidthLimit": 20,
  "addressName": "prod-eip-01"
}
```

Optionally include `sharedBandwidthPackageId` to force a specific shared
bandwidth package instead of the automatically selected one.

#### Successful Response

- **Code:** `201 Created`
- **Body:**

```json
{ "id": "eip-12345678", "address": "203.0.113.10" }
```

#### Error Codes

- `400 Bad Request` – payload validation failed or the JSON body is invalid.
- `500 Internal Server Error` – Tencent Cloud API returned an error or
  responded with an unexpected payload.

### Delete Elastic IP

- **Method:** `DELETE`
- **Path:** `/tencentcloud/eips/:id`
- **Query Parameters:**
  - `region` – Tencent Cloud region where the elastic IP resides.
- **Description:** Releases an elastic public IP address and returns the
  capacity to the shared bandwidth package in the specified region.

#### Successful Response

- **Code:** `204 No Content`

#### Error Codes

- `400 Bad Request` – missing or invalid `id` or `region` values.
- `404 Not Found` – the elastic IP was not found in Tencent Cloud.
- `500 Internal Server Error` – Tencent Cloud API returned an error or
  responded with an unexpected payload.

### Issue JWT

- **Method:** `POST`
- **Path:** `/auth/jwt`
- **Authentication:** must include the legacy bearer token. JWTs cannot be used
  to access this endpoint.
- **Description:** Generates a signed JWT containing optional provisioning
  limits. The token can later be supplied in the `Authorization` header to
  enforce limits during `/tencentcloud/...` requests.

#### Request Body

```json
  {
    "maxTransferKB": 1048576,
    "allowedInstanceTypes": ["SA2.MEDIUM2"],
    "allowedZones": ["ap-guangzhou-6"],
    "maxBandwidthMbps": 100,
    "projectId": 12345,
    "ttlMinutes": 60
  }
```

- `maxTransferKB`: `-1` for unlimited, or a non-negative limit (KB).
- `allowedInstanceTypes`, `allowedZones`: optional allow lists. Empty or
  omitted arrays do not impose a restriction.
- `projectId`: when set, the JWT enforces this project on resource creation and
  subsequent operations.
- `ttlMinutes`: lifetime for the token, must be positive and ≤ configured
  `jwt.maxTTLMinutes` (default 5,256,000 minutes / 10 years).

#### Successful Response

- **Code:** `201 Created`
- **Body:**

```json
{ "token": "<jwt>", "expiresAt": "2025-01-02T03:04:05Z" }
```

#### Error Codes

- `400 Bad Request` – request validation failed or `ttlMinutes` exceeded the
  configured maximum.
- `401 Unauthorized` – legacy bearer token missing or invalid.

### Delete Virtual Machine

- **Method:** `DELETE`
- **Path:** `/tencentcloud/vms/:id`
- **Description:** Enqueues termination of a previously created CVM instance.
  The job releases any non-cascading elastic IPs before deleting the instance.
  If the Tencent Cloud instance no longer exists, Penguin simply removes its
  record and still completes successfully.

#### Path Parameters

- `id` – Identifier returned by the create endpoint (UUID v4 string).

#### Responses

- `202 Accepted` – delete request scheduled for asynchronous processing.
- `400 Bad Request` – invalid identifier format.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Adjust Virtual Machine Bandwidth

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/bandwidth`
- **Description:** Updates the instance internet egress bandwidth limit。服务会先调用
  `DescribeAddresses` 区分公网 IP 类型：
  - `WanIP` 继续通过 `ResetInstancesInternetMaxBandwidth` 更新实例带宽；
  - `EIP` / `AnycastEIP` 通过 `ModifyAddressesBandwidth` 调整 EIP 带宽；
  - 若未查询到公网 IP，则返回 `404 Not Found` 并提示 “no public address bound to instance”。
  若后续腾讯云 API 调用失败，将根据错误类型返回 `502/500` 等响应。

#### Path Parameters

- `id` – Identifier returned by the create endpoint (UUID v4 string).

#### Request Body

```json
{ "bandwidthLimit": 50 }
```

#### Responses

- `202 Accepted` – bandwidth update requested。
- `400 Bad Request` – invalid identifier or payload。
- `404 Not Found` – 虚拟机不存在或未绑定公网 IP。
- `502 Bad Gateway` – 调用腾讯云网络 API（如 `DescribeAddresses`、`ModifyAddressesBandwidth`）失败。
- `500 Internal Server Error` – 腾讯云 API 返回错误且无法分类。

### Get Image By Name

- **Method:** `GET`
- **Path:** `/tencentcloud/images/:name`
- **Query Parameters:**
  - `region` – Tencent Cloud region to query.
- **Description:** Looks up an image by its `image-name` attribute using the
  Tencent Cloud CVM `DescribeImages` API.

#### Responses

- `200 OK` – returns the matching image metadata.
- `400 Bad Request` – missing `name` or `region` parameters.
- `404 Not Found` – no image matches the specified name in the region.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Start Virtual Machine

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/start`
- **Description:** Powers on the specified CVM instance.

#### Responses

- `202 Accepted` – start operation dispatched.
- `400 Bad Request` – invalid identifier format.
- `409 Conflict` – transfer quota exceeded; instance remains stopped.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Shutdown Virtual Machine

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/shutdown`
- **Description:** Gracefully stops the specified CVM instance.

#### Responses

- `202 Accepted` – shutdown request dispatched.
- `400 Bad Request` – invalid identifier format.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Renew Virtual Machine

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/renew`
- **Description:** Extends the prepaid term of the specified CVM instance. When `period` is omitted the term defaults to one month. Set `autoRenew` to enable or disable Tencent Cloud automatic renewal for the renewed term.

#### Request Body (optional)

```json
{
  "period": 12,
  "autoRenew": true
}
```

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{ "expiredAt": "2025-12-01T00:00:00Z" }
```

#### Error Codes

- `400 Bad Request` – invalid identifier or malformed body.
- `404 Not Found` – the virtual machine record does not exist.
- `409 Conflict` – renewal rejected (e.g. instance not in prepaid mode or expiration unchanged).
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Get Virtual Machine Status

- **Method:** `GET`
- **Path:** `/tencentcloud/vms/:id/status`
- **Description:** Returns the tracked transfer usage together with the latest CVM instance metadata and power state. Directional fields `txTransfer` (upload) and `rxTransfer` (download) are included; `usedTransfer` remains the sum for backward compatibility. When Penguin has suspended the instance after exceeding its transfer quota, the `instanceState` field is reported as `SuspendOverUsage`.

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{
  "id": "7f1f1e97-7f45-4c4b-94b4-2f1c248a8b1e",
  "zone": "ap-guangzhou-6",
  "instanceId": "ins-a1b2c3",
  "instanceType": "S2.SMALL1",
  "instanceState": "RUNNING",
  "cpu": 2,
  "memoryGiB": 4,
  "systemDiskSizeGiB": 60,
  "privateIps": ["10.0.0.12"],
  "publicIps": ["43.1.2.3"],
  "imageId": "img-7efla8nv",
  "osName": "Ubuntu Server 22.04 LTS",
  "createdAt": "2024-01-01T00:00:00Z",
  "expiredAt": "2025-01-01T00:00:00Z",
  "totalTransfer": 1048576,
  "usedTransfer": 512000,
  "txTransfer": 400000,
  "rxTransfer": 112000,
  "remainingTransfer": 536576,
  "password": "StatusPass1",
  "defaultLoginUser": "root"
}
```

#### Error Codes

- `400 Bad Request` – identifier missing or malformed.
- `404 Not Found` – the virtual machine record does not exist.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Get Virtual Machine Metrics

- **Method:** `GET`
- **Path:** `/tencentcloud/vms/:id/metrics`
- **Query Parameters:**
  - `range` – Optional. One of `day`, `week`, or `month`. Defaults to `day`.
- **Description:** Returns average CPU/memory utilization and aggregate inbound/outbound transfer for the requested window.

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{
  "range": "day",
  "cpuAveragePercent": 32.4,
  "memoryAveragePercent": 58.1,
  "networkOutKB": 145678,
  "networkInKB": 98321,
  "start": "2025-10-01T08:00:00Z",
  "end": "2025-10-02T08:00:00Z"
}
```

#### Error Codes

- `400 Bad Request` – identifier missing/malformed or unsupported range.
- `404 Not Found` – the virtual machine record does not exist.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Get Virtual Machine VNC URL

- **Method:** `GET`
- **Path:** `/tencentcloud/vms/:id/vnc`
- **Description:** Returns the Tencent Cloud VNC websocket URL for the specified CVM instance. Append the returned value to the Tencent Cloud web console (`https://img.qcloud.com/qcloud/app/active_vnc/index.html?InstanceVncUrl=<value>`) to launch the remote console.

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{ "url": "wss://bjvnc.qcloud.com:26789/vnc?s=..." }
```

#### Error Codes

- `400 Bad Request` – identifier missing or malformed.
- `404 Not Found` – the virtual machine record does not exist.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Reinstall Virtual Machine

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/reinstall`
- **Description:** Reinstalls the operating system of the specified CVM instance using the provided image ID.

#### Request Body

```json
{
  "imageId": "img-7efla8nv",
  "cloudInitData": "#cloud-config\nwrite_files:\n- path: /etc/motd\n  content: Hello\n"
}
```

#### Responses

- `202 Accepted` – reinstall task dispatched.
- `400 Bad Request` – invalid identifier or image ID.
- `404 Not Found` – the virtual machine record does not exist.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Reset Virtual Machine Password

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/reset-password`
- **Description:** Generates a new 12-character alphanumeric password, resets it on the Tencent Cloud instance, and stores it for future status queries. Set `forceStop` to `true` to allow forced shutdown when the instance is running.

#### Request Body (optional)

```json
{
  "forceStop": true
}
```

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
{
  "password": "Ab3kLmN9PqRs"
}
```

#### Error Codes

- `400 Bad Request` – invalid identifier or malformed body.
- `404 Not Found` – the virtual machine record does not exist.
- `500 Internal Server Error` – Tencent Cloud API returned an error.

### List Missing Virtual Machines

- **Method:** `GET`
- **Path:** `/tencentcloud/vms/missing`
- **Description:** Returns virtual machines tracked in Penguin whose Tencent
  Cloud instances are no longer present. The service cross-checks the stored
  instance identifiers with the `DescribeInstances` API and reports any that
  are missing from Tencent Cloud.

#### Successful Response

- **Code:** `200 OK`
- **Body:**

```json
[
  { "id": "7f1f1e97-7f45-4c4b-94b4-2f1c248a8b1e", "zone": "ap-guangzhou-6", "instanceId": "ins-abc123" }
]
```

#### Error Codes

- `500 Internal Server Error` – Tencent Cloud API returned an error.

### Reset Transfer Usage

- **Method:** `POST`
- **Path:** `/tencentcloud/vms/:id/reset-transfer`
- **Description:** Resets the recorded transfer usage for the current billing period to zero. Use this after granting additional quota or clearing an exceeded state. Each reset writes an audit row to `virtual_machine_reset_log` capturing the VM’s project ID, cloud instance ID, pre-reset `txTransfer`/`rxTransfer`, and the reset timestamp (UTC).

#### Responses

- `204 No Content` – usage reset successfully.
- `400 Bad Request` – invalid identifier format.
- `404 Not Found` – the virtual machine record does not exist.
- `500 Internal Server Error` – database update failed.

## Notes

- The request payload mirrors the internal `CreateVirtualMachineRequest`
  structure in Go. Fields such as `name` must be valid FQDNs, `systemDiskSize`
  must be at least 20 GiB, and `bandwidthLimit` must be at least 1 Mbps.
- The service automatically enables security, monitor, and automation services
  on new CVM instances. Control prepaid duration and auto-renew behaviour via
  the create and renew endpoints.
- Identifiers produced by the create endpoint are UUID v4 strings mapped to the
  underlying Tencent Cloud zone and instance ID. Keep them safe and reuse them
  for lifecycle operations.
- `rootLoginPassword` is optional. When omitted, Penguin generates a
  12-character password using the `[a-zA-Z0-9]` character set and stores it for
  later retrieval via the status endpoint.
- Optional `cloudInitData` accepts raw cloud-init content; the service encodes
  it as base64 for Tencent Cloud. Payloads larger than 16 KB are rejected.
