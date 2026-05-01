# Mock Data Injection Guide

To test the billing pipeline without external service dependencies, you can inject mock metrics directly into NATS.

## Using the Shell Script
Run the automated test script to send a batch of diversified S3 and RDS metrics:
```bash
./scripts/send_metrics.sh
```

## Using the Go Helper
If the NATS CLI is not available, use the Go-based publisher:
```bash
go run scripts/publish_metrics.go
```

## Manual Publication (NATS CLI)
You can publish custom metrics using the `nats pub` command:
```bash
nats pub dev.v1.billing.metric.ec2 '{"timestamp":"'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'","service_id":"ec2","metric_name":"cpu","value":50.0,"unit":"percent","tenant_id":"t1"}'
```
