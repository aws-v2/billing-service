package domain

import "time"

// Metric represents a standardized usage record from an AWS-like service.
type Metric struct {
	Timestamp  time.Time              `json:"timestamp"`
	ServiceID  string                 `json:"service_id"`
	MetricName string                 `json:"metric_name"`
	Value      float64                `json:"value"`
	Unit       string                 `json:"unit"`
	TenantID   string                 `json:"tenant_id"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}
