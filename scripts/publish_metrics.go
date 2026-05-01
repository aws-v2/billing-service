package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

type Metric struct {
	Timestamp  time.Time              `json:"timestamp"`
	ServiceID  string                 `json:"service_id"`
	MetricName string                 `json:"metric_name"`
	Value      float64                `json:"value"`
	Unit       string                 `json:"unit"`
	TenantID   string                 `json:"tenant_id"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

func main() {
	url := "nats://localhost:4222"
	user := "auth-server"
	password := "auth-secret"
	prefix := "dev.v1"

	nc, err := nats.Connect(url, nats.UserInfo(user, password))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer nc.Close()

	sendMetric := func(service, name string, value float64, unit, tenant string) {
		subject := fmt.Sprintf("%s.billing.metric.%s", prefix, service)
		metric := Metric{
			Timestamp:  time.Now(),
			ServiceID:  service,
			MetricName: name,
			Value:      value,
			Unit:       unit,
			TenantID:   tenant,
		}

		data, _ := json.Marshal(metric)
		if err := nc.Publish(subject, data); err != nil {
			log.Printf("Failed to publish to %s: %v", subject, err)
			return
		}
		fmt.Printf("Published to %s\n", subject)
	}

	sendMetric("ec2", "vmusage", 1.5, "hours", "tenant-123")
	sendMetric("s3", "storage", 100.0, "GB", "tenant-123")
	sendMetric("lambda", "invocations", 5000, "count", "tenant-456")

	nc.Flush()
}
