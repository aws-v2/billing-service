#!/bin/bash

# Configuration
NATS_URL=${NATS_URL:-"nats://localhost:4222"}
NATS_USER=${NATS_USER:-"auth-server"}
NATS_PASSWORD=${NATS_PASSWORD:-"auth-secret"}
NATS_PREFIX=${NATS_PREFIX:-"dev.v1"}

# Check if nats CLI is installed
if ! command -v nats &> /dev/null
then
    echo "NATS CLI not found. You can install it with: brew install nats-server or go install github.com/nats-io/natscli/nats@latest"
    echo "Alternatively, I will create a Go-based sender for you."
    exit 1
fi

TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# 1. S3 Specialized Metrics
echo "Sending S3 metrics..."
nats pub --user "$NATS_USER" --password "$NATS_PASSWORD" -s "$NATS_URL" $NATS_PREFIX.billing.metric.s3 \
'{"timestamp":"'$TIMESTAMP'","bucket_id":"app-uploads","size_gb":120.5,"region":"us-east-1","tenant_id":"tenant-123"}'

nats pub --user "$NATS_USER" --password "$NATS_PASSWORD" -s "$NATS_URL" $NATS_PREFIX.billing.metric.s3 \
'{"timestamp":"'$TIMESTAMP'","operation":"PutObject","request_tier":"tier_1","tenant_id":"tenant-123"}'

# 2. EC2 Specialized Metrics
echo "Sending EC2 metrics..."
nats pub --user "$NATS_USER" --password "$NATS_PASSWORD" -s "$NATS_URL" $NATS_PREFIX.billing.metric.ec2 \
'{"timestamp":"'$TIMESTAMP'","instance_id":"i-0abcd1234efgh5678","instance_type":"t3.medium","state":"running","tenant_id":"tenant-123"}'

nats pub --user "$NATS_USER" --password "$NATS_PASSWORD" -s "$NATS_URL" $NATS_PREFIX.billing.metric.ec2 \
'{"timestamp":"'$TIMESTAMP'","volume_id":"vol-0987654321","volume_type":"gp3","size_gb":100,"tenant_id":"tenant-123"}'

# 3. RDS Specialized Metrics
echo "Sending RDS metrics..."
nats pub --user "$NATS_USER" --password "$NATS_PASSWORD" -s "$NATS_URL" $NATS_PREFIX.billing.metric.rds \
'{"timestamp":"'$TIMESTAMP'","db_instance_id":"db-prod-01","instance_class":"db.t3.large","engine":"postgres","status":"available","tenant_id":"tenant-123"}'

nats pub --user "$NATS_USER" --password "$NATS_PASSWORD" -s "$NATS_URL" $NATS_PREFIX.billing.metric.rds \
'{"timestamp":"'$TIMESTAMP'","db_instance_id":"db-prod-01","allocated_storage_gb":50.0,"storage_type":"gp2","tenant_id":"tenant-123"}'
