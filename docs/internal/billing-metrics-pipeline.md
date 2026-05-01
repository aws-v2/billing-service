# Metric Ingestion Pipeline

The metric ingestion pipeline ensures that usage data from all services reaches the billing database reliably.

## Ingestion Flow
1. **Producer**: Services like EC2 or S3 publish JSON metrics to NATS.
2. **Subject**: `dev.v1.billing.metric.<service_id>`
3. **Queue Group**: `billing_ingestion` (Ensures load balanced processing).
4. **Subscriber**: The Billing Service pulls metrics and validates them against domain DTOs.
5. **Persistence**: Validated metrics are saved to the `metrics` table.

## Scalability
The use of NATS Queue Groups allows the ingestion layer to scale horizontally by adding more instances of the Billing Service.
