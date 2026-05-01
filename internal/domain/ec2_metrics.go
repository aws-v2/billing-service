package domain

import "time"

// EC2ComputeMetricDTO represents the compute runtime usage of an instance.
type EC2ComputeMetricDTO struct {
	Timestamp    time.Time `json:"timestamp"`
	InstanceID   string    `json:"instance_id"`
	InstanceType string    `json:"instance_type"` // e.g. t3.medium
	State        string    `json:"state"`         // e.g. running, stopped
	TenantID     string    `json:"tenant_id"`
}

// EC2NetworkMetricDTO represents the network data transfer for an EC2 instance.
type EC2NetworkMetricDTO struct {
	Timestamp       time.Time `json:"timestamp"`
	InstanceID      string    `json:"instance_id"`
	BytesSent       int64     `json:"bytes_sent"`
	DestinationType string    `json:"destination_type"` // e.g. internet, regional
	TenantID        string    `json:"tenant_id"`
}

// EC2StorageMetricDTO represents the EBS volume provisioned capacity.
type EC2StorageMetricDTO struct {
	Timestamp    time.Time `json:"timestamp"`
	VolumeID     string    `json:"volume_id"`
	VolumeType   string    `json:"volume_type"` // e.g. gp3, io2
	SizeGB       float64   `json:"size_gb"`
	TenantID     string    `json:"tenant_id"`
}
