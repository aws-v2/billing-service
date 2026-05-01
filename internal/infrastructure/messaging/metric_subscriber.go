package messaging

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/Qarani-m/billing-service/internal/domain"
	"github.com/Qarani-m/billing-service/pkg/messaging"
)

type MetricSubscriber struct {
	subscriber *messaging.NATSSubscriber
}

func NewMetricSubscriber(subscriber *messaging.NATSSubscriber) *MetricSubscriber {
	return &MetricSubscriber{subscriber: subscriber}
}

func (s *MetricSubscriber) StartIngestion(billingService domain.BillingService) error {
	log.Println("[subscriber] Starting metric ingestion...")

	_, err := s.subscriber.Subscribe(func(subject string, data []byte) error {
		log.Printf("[subscriber] Received message on subject: %s", subject)

		parts := strings.Split(subject, ".")
		serviceID := "unknown"
		if len(parts) > 0 {
			serviceID = parts[len(parts)-1]
		}

		log.Printf("[subscriber] Resolved serviceID: %s", serviceID)

		metric, err := parseMetric(serviceID, data)
		if err != nil {
			log.Printf("[subscriber] Failed to parse metric for service=%s: %v", serviceID, err)
			return err
		}

		if metric.ServiceID == "" {
			log.Printf("[subscriber] Empty metric after parsing for service=%s, skipping", serviceID)
			return nil
		}

		log.Printf("[subscriber] Ingested metric: service=%s metric=%s value=%.2f %s tenant=%s",
			metric.ServiceID, metric.MetricName, metric.Value, metric.Unit, metric.TenantID)

		if err := billingService.RecordMetric(metric); err != nil {
			log.Printf("[subscriber] Failed to record metric: service=%s metric=%s tenant=%s error=%v",
				metric.ServiceID, metric.MetricName, metric.TenantID, err)
			return err
		}

		log.Printf("[subscriber] Metric recorded OK: service=%s tenant=%s", metric.ServiceID, metric.TenantID)
		return nil
	})

	if err != nil {
		log.Printf("[subscriber] Failed to start ingestion: %v", err)
	}

	return err
}

func parseMetric(serviceID string, data []byte) (domain.Metric, error) {
	switch serviceID {
	case "s3":
		return parseS3Metric(data)
	case "ec2":
		return parseEC2Metric(data)
	case "rds":
		return parseRDSMetric(data)
	default:
		log.Printf("[subscriber] Unknown serviceID=%s, attempting generic parse", serviceID)
		var metric domain.Metric
		if err := json.Unmarshal(data, &metric); err != nil {
			return domain.Metric{}, err
		}
		return metric, nil
	}
}

func parseS3Metric(data []byte) (domain.Metric, error) {
	var s3Storage domain.S3StorageMetricDTO
	if err := json.Unmarshal(data, &s3Storage); err == nil && s3Storage.BucketID != "" {
		log.Printf("[subscriber] Parsed S3 storage metric: bucket=%s size=%.2fGB tenant=%s",
			s3Storage.BucketID, s3Storage.SizeGB, s3Storage.TenantID)
		return domain.Metric{
			Timestamp:  s3Storage.Timestamp,
			ServiceID:  "s3",
			MetricName: "storage_utilization",
			Value:      s3Storage.SizeGB,
			Unit:       "GB",
			TenantID:   s3Storage.TenantID,
		}, nil
	}

	var s3Req domain.S3RequestMetricDTO
	if err := json.Unmarshal(data, &s3Req); err == nil {
		log.Printf("[subscriber] Parsed S3 request metric: tenant=%s", s3Req.TenantID)
		return domain.Metric{
			Timestamp:  s3Req.Timestamp,
			ServiceID:  "s3",
			MetricName: "api_request",
			Value:      1,
			Unit:       "count",
			TenantID:   s3Req.TenantID,
		}, nil
	}

	log.Println("[subscriber] Failed to parse S3 metric")
	return domain.Metric{}, nil
}

func parseEC2Metric(data []byte) (domain.Metric, error) {
	var ec2Compute domain.EC2ComputeMetricDTO
	if err := json.Unmarshal(data, &ec2Compute); err == nil && ec2Compute.InstanceID != "" && ec2Compute.InstanceType != "" {
		log.Printf("[subscriber] Parsed EC2 compute metric: instance=%s type=%s tenant=%s",
			ec2Compute.InstanceID, ec2Compute.InstanceType, ec2Compute.TenantID)
		return domain.Metric{
			Timestamp:  ec2Compute.Timestamp,
			ServiceID:  "ec2",
			MetricName: "compute_usage",
			Value:      1,
			Unit:       "state",
			TenantID:   ec2Compute.TenantID,
		}, nil
	}

	var ec2Storage domain.EC2StorageMetricDTO
	if err := json.Unmarshal(data, &ec2Storage); err == nil && ec2Storage.VolumeID != "" {
		log.Printf("[subscriber] Parsed EC2 storage metric: volume=%s size=%.2fGB tenant=%s",
			ec2Storage.VolumeID, ec2Storage.SizeGB, ec2Storage.TenantID)
		return domain.Metric{
			Timestamp:  ec2Storage.Timestamp,
			ServiceID:  "ec2",
			MetricName: "ebs_storage",
			Value:      ec2Storage.SizeGB,
			Unit:       "GB",
			TenantID:   ec2Storage.TenantID,
		}, nil
	}

	log.Println("[subscriber] Failed to parse EC2 metric")
	return domain.Metric{}, nil
}

func parseRDSMetric(data []byte) (domain.Metric, error) {
	var rdsInstance domain.RDSInstanceMetricDTO
	if err := json.Unmarshal(data, &rdsInstance); err == nil && rdsInstance.InstanceClass != "" {
		log.Printf("[subscriber] Parsed RDS instance metric: class=%s tenant=%s",
			rdsInstance.InstanceClass, rdsInstance.TenantID)
		return domain.Metric{
			Timestamp:  rdsInstance.Timestamp,
			ServiceID:  "rds",
			MetricName: "instance_usage",
			Value:      1,
			Unit:       "state",
			TenantID:   rdsInstance.TenantID,
		}, nil
	}

	var rdsStorage domain.RDSStorageMetricDTO
	if err := json.Unmarshal(data, &rdsStorage); err == nil && rdsStorage.AllocatedStorageGB != 0 {
		log.Printf("[subscriber] Parsed RDS storage metric: size=%.2fGB tenant=%s",
			rdsStorage.AllocatedStorageGB, rdsStorage.TenantID)
		return domain.Metric{
			Timestamp:  rdsStorage.Timestamp,
			ServiceID:  "rds",
			MetricName: "storage_usage",
			Value:      rdsStorage.AllocatedStorageGB,
			Unit:       "GB",
			TenantID:   rdsStorage.TenantID,
		}, nil
	}

	log.Println("[subscriber] Failed to parse RDS metric")
	return domain.Metric{}, nil
}