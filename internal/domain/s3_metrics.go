package domain

import "time"

type S3StorageMetricDTO struct {
	MetricType string    `json:"metric_type"` // "storage_utilization"
	Timestamp  time.Time `json:"timestamp"`
	BucketID   string    `json:"bucket_id"`
	SizeGB     float64   `json:"size_gb"`
	Region     string    `json:"region"`
	TenantID   string    `json:"tenant_id"`
}

type S3RequestMetricDTO struct {
	MetricType  string    `json:"metric_type"` // "api_request"
	Timestamp   time.Time `json:"timestamp"`
	Operation   string    `json:"operation"`
	RequestTier string    `json:"request_tier"`
	TenantID    string    `json:"tenant_id"`
}

type S3BandwidthMetricDTO struct {
	MetricType string    `json:"metric_type"` // "bandwidth"
	Timestamp  time.Time `json:"timestamp"`
	BytesOut   int64     `json:"bytes_out"`
	Region     string    `json:"region"`
	TenantID   string    `json:"tenant_id"`
}