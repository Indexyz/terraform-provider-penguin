// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package penguin

import "fmt"

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e == nil {
		return "penguin API error"
	}
	if e.Message == "" {
		return fmt.Sprintf("penguin API error %d", e.Status)
	}
	return fmt.Sprintf("penguin API error %d: %s", e.Status, e.Message)
}

type InternalHealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

type Zone struct {
	Region     string `json:"region"`
	RegionName string `json:"regionName"`
	Zone       string `json:"zone"`
	ZoneName   string `json:"zoneName"`
	ZoneID     string `json:"zoneId,omitempty"`
	State      string `json:"state"`
}

type ZonesResponse struct {
	Zones []Zone `json:"zones"`
}

type BandwidthPackageSelectionResponse struct {
	ID             string `json:"id"`
	AvailableCount int64  `json:"availableCount"`
}

type CreateVirtualMachineRequest struct {
	Name                     string  `json:"name"`
	Zone                     string  `json:"zone"`
	InstanceType             string  `json:"instanceType"`
	SecurityGroup            string  `json:"securityGroup"`
	SystemImage              string  `json:"systemImage"`
	VPCID                    string  `json:"vpcId"`
	SubnetID                 string  `json:"subnetId"`
	PrivateIPAddress         *string `json:"privateIpAddress,omitempty"`
	SystemDiskSizeGiB        int64   `json:"systemDiskSize"`
	SharedBandwidthPackageID *string `json:"sharedBandwidthPackageId,omitempty"`
	ElasticIPID              *string `json:"elasticIpId,omitempty"`
	BandwidthLimitMbps       *int64  `json:"bandwidthLimit,omitempty"`
	ChargeType               *string `json:"chargeType,omitempty"`
	RootLoginPassword        *string `json:"rootLoginPassword,omitempty"`
	TotalTransferKB          int64   `json:"totalTransfer"`
	ProjectID                *int64  `json:"projectId,omitempty"`
	PeriodMonths             *int64  `json:"period,omitempty"`
	CloudInitData            *string `json:"cloudInitData,omitempty"`
	AutoRenew                *bool   `json:"autoRenew,omitempty"`
}

type CreateVirtualMachineResponse struct {
	ID string `json:"id"`
}

type CreateElasticIPRequest struct {
	Region                   string  `json:"region"`
	SharedBandwidthPackageID *string `json:"sharedBandwidthPackageId,omitempty"`
	BandwidthLimitMbps       int64   `json:"bandwidthLimit"`
	AddressName              string  `json:"addressName"`
}

type CreateElasticIPResponse struct {
	ID      string `json:"id"`
	Address string `json:"address,omitempty"`
}

type RenewVirtualMachineRequest struct {
	PeriodMonths *int64 `json:"period,omitempty"`
	AutoRenew    *bool  `json:"autoRenew,omitempty"`
}

type RenewVirtualMachineResponse struct {
	ExpiredAt *string `json:"expiredAt,omitempty"`
}

type ReinstallVirtualMachineRequest struct {
	ImageID       string  `json:"imageId"`
	CloudInitData *string `json:"cloudInitData,omitempty"`
}

type ResetVirtualMachinePasswordRequest struct {
	ForceStop *bool `json:"forceStop,omitempty"`
}

type ResetVirtualMachinePasswordResponse struct {
	Password string `json:"password"`
}

type MissingVirtualMachine struct {
	ID         string `json:"id"`
	Zone       string `json:"zone"`
	InstanceID string `json:"instanceId"`
}

type VirtualMachineStatus struct {
	ID                string   `json:"id"`
	Zone              string   `json:"zone"`
	InstanceID        string   `json:"instanceId"`
	InstanceType      string   `json:"instanceType"`
	InstanceState     string   `json:"instanceState"`
	RestrictState     *string  `json:"restrictState,omitempty"`
	StopChargingMode  *string  `json:"stopChargingMode,omitempty"`
	RenewFlag         *string  `json:"renewFlag,omitempty"`
	CPU               int64    `json:"cpu"`
	MemoryGiB         int64    `json:"memoryGiB"`
	SystemDiskSizeGiB int64    `json:"systemDiskSizeGiB"`
	PrivateIPs        []string `json:"privateIps"`
	PublicIPs         []string `json:"publicIps"`
	ImageID           *string  `json:"imageId,omitempty"`
	OSName            *string  `json:"osName,omitempty"`
	CreatedAt         *string  `json:"createdAt,omitempty"`
	ExpiredAt         *string  `json:"expiredAt,omitempty"`
	TotalTransfer     int64    `json:"totalTransfer"`
	UsedTransfer      int64    `json:"usedTransfer"`
	TxTransfer        *int64   `json:"txTransfer,omitempty"`
	RxTransfer        *int64   `json:"rxTransfer,omitempty"`
	RemainingTransfer *int64   `json:"remainingTransfer,omitempty"`
	Password          *string  `json:"password,omitempty"`
	DefaultLoginUser  *string  `json:"defaultLoginUser,omitempty"`
}

type VirtualMachineVNCResponse struct {
	URL string `json:"url"`
}

type VirtualMachineMetricsResponse struct {
	Range                string  `json:"range"`
	CPUAveragePercent    float64 `json:"cpuAveragePercent"`
	MemoryAveragePercent float64 `json:"memoryAveragePercent"`
	NetworkOutKB         int64   `json:"networkOutKB"`
	NetworkInKB          int64   `json:"networkInKB"`
	Start                string  `json:"start"`
	End                  string  `json:"end"`
}

type AdjustBandwidthRequest struct {
	BandwidthLimitMbps int64 `json:"bandwidthLimit"`
}

type IssueJWTRequest struct {
	MaxTransferKB        *int64   `json:"maxTransferKB,omitempty"`
	AllowedInstanceTypes []string `json:"allowedInstanceTypes,omitempty"`
	AllowedZones         []string `json:"allowedZones,omitempty"`
	MaxBandwidthMbps     *int64   `json:"maxBandwidthMbps,omitempty"`
	ProjectID            *int64   `json:"projectId,omitempty"`
	TTLMinutes           int64    `json:"ttlMinutes"`
}

type IssueJWTResponse struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expiresAt"`
}
