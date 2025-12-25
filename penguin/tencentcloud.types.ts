/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

export type VirtualMachineId = string;

export interface CreateVirtualMachineRequest {
  /** Instance name and hostname (must be a valid FQDN). */
  name: string;
  /** Availability zone, e.g. `ap-guangzhou-6`. */
  zone: string;
  /** Machine type identifier, e.g. `SA2.MEDIUM2`. */
  instanceType: string;
  /** Security group identifier. */
  securityGroup: string;
  /** Image identifier, e.g. `img-7efla8nv`. */
  systemImage: string;
  /** Virtual Private Cloud identifier. */
  vpcId: string;
  /** Subnet identifier within the VPC. */
  subnetId: string;
  /** System disk size in GiB (minimum 20). */
  systemDiskSize: number;
  /** Shared bandwidth package identifier; omit when supplying elasticIpId. */
  sharedBandwidthPackageId?: string;
  /** Optional existing elastic IP to associate with the new instance. */
  elasticIpId?: string;
  /** Outbound bandwidth limit in Mbps (required when requesting a new public IP). */
  bandwidthLimit?: number;
  /** Billing mode, defaults to `PREPAID`; set to `POSTPAID_BY_HOUR` for on-demand instances. */
  chargeType?: 'PREPAID' | 'POSTPAID_BY_HOUR';
  /** Root login password (minimum 8 characters). */
  rootLoginPassword: string;
  /** Optional prepaid period in months (defaults to 1). */
  period?: number;
  /** Optional Tencent Cloud project identifier (defaults to 0). */
  projectId?: number;
  /** Total transfer quota in KB (-1 disables the limit). */
  totalTransfer: number;
  /** Optional cloud-init user data (raw content, max 16KB). */
  cloudInitData?: string;
  /** Optional flag to enable Tencent Cloud auto-renewal (default false). */
  autoRenew?: boolean;
}

export interface CreateVirtualMachineResponse {
  /** UUID identifier assigned by the Penguin service. */
  id: VirtualMachineId;
}

export interface CreateElasticIPRequest {
  /** Region where the elastic IP should be provisioned (e.g. `ap-guangzhou`). */
  region: string;
  /** Shared bandwidth package identifier; omit to let the service auto-select a schedulable package in the region. */
  sharedBandwidthPackageId?: string;
  /** Outbound bandwidth cap in Mbps (minimum 1). */
  bandwidthLimit: number;
  /** Human-readable name assigned to the elastic IP. */
  addressName: string;
}

export interface CreateElasticIPResponse {
  /** Tencent Cloud elastic IP identifier (e.g. `eip-12345678`). */
  id: string;
  /** Assigned public IPv4 address, when available. */
  address?: string;
}

export interface RenewVirtualMachineRequest {
  /** Optional prepaid period in months (defaults to 1). */
  period?: number;
  /** Optional flag to enable Tencent Cloud auto-renewal for the renewed term. */
  autoRenew?: boolean;
}

export interface RenewVirtualMachineResponse {
  /** Updated prepaid expiration timestamp, if reported by Tencent Cloud. */
  expiredAt?: string;
}

export interface ReinstallVirtualMachineRequest {
  /** Image identifier to reinstall (e.g. `img-7efla8nv`). */
  imageId: string;
  /** Optional cloud-init user data to apply during reinstall. */
  cloudInitData?: string;
}

export interface MissingVirtualMachine {
  id: string;
  zone: string;
  instanceId: string;
}

export interface VirtualMachineStatus {
  id: VirtualMachineId;
  zone: string;
  instanceId: string;
  instanceType: string;
  /** CVM state; returns `SuspendOverUsage` when traffic suspension is active. */
  instanceState: string;
  restrictState?: string;
  stopChargingMode?: string;
  renewFlag?: string;
  cpu: number;
  memoryGiB: number;
  systemDiskSizeGiB: number;
  privateIps: string[];
  publicIps: string[];
  imageId?: string;
  osName?: string;
  createdAt?: string;
  expiredAt?: string;
  totalTransfer: number;
  usedTransfer: number;
  /** Directional transfer usage in KB; present for forward compatibility. */
  txTransfer?: number;
  rxTransfer?: number;
  remainingTransfer?: number;
}

export interface VirtualMachineVncResponse {
  /** Tencent Cloud websocket URL for launching the CVM console. */
  url: string;
}

export interface SelectBandwidthPackageQuery {
  /** Tencent Cloud region (e.g. `ap-guangzhou`). */
  region: string;
  /** Optional Tencent Cloud network type (defaults to `BGP`). */
  networkType?: string;
}

export interface SelectBandwidthPackageResponse {
  /** Tencent Cloud shared bandwidth package identifier (e.g. `bwp-123`). */
  id: string;
  /** Remaining available bindings before reaching the capacity threshold. */
  availableCount: number;
}

export interface InternalHealthResponse {
  /** Overall application status. */
  status: 'ok' | 'degraded';
  /** Database connectivity status. */
  database: 'ok' | 'error';
}

export interface ApiError {
  /** HTTP status code returned by the API. */
  status: number;
  /** Error message returned by the server. */
  message: string;
}
