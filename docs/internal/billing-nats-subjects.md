# NATS Subject Hierarchy

Our messaging system follows a standardized subject hierarchy to ensure clarity and discoverability across environments.

## Subject Scheme
`env.version.subscriber.subject.action`

### Billing Ingestion Subtree
- **Wildcard**: `dev.v1.billing.metric.>` (Subscribes to all metrics)
- **Service Specifics**:
  - `dev.v1.billing.metric.s3`
  - `dev.v1.billing.metric.ec2`
  - `dev.v1.billing.metric.lambda`

## Prefixes
The `NATS_PREFIX` environment variable (e.g., `dev.v1`, `staging.v1`, `prod.v1`) must be prepended to all subjects to isolate environments reliably.
