package domain

import "time"

// RDSInstanceMetricDTO represents the runtime usage of an RDS instance.
type RDSInstanceMetricDTO struct {
	Timestamp     time.Time `json:"timestamp"`
	DBInstanceID  string    `json:"db_instance_id"`
	InstanceClass string    `json:"instance_class"` // e.g. db.t3.medium
	Engine        string    `json:"engine"`         // e.g. postgres, mysql
	Status        string    `json:"status"`         // e.g. available, stopped
	TenantID      string    `json:"tenant_id"`
}

// RDSStorageMetricDTO represents the storage capacity of an RDS instance.
type RDSStorageMetricDTO struct {
	Timestamp          time.Time `json:"timestamp"`
	DBInstanceID       string    `json:"db_instance_id"`
	AllocatedStorageGB float64   `json:"allocated_storage_gb"`
	StorageType        string    `json:"storage_type"` // e.g. gp2, io1
	TenantID           string    `json:"tenant_id"`
}
