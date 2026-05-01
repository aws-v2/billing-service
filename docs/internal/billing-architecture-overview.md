# ----------------------------------
# Billing System Architecture Overview

The Billing Service is an asynchronous, event-driven system designed to calculate costs based on service usage metrics.

## Core Components
1. **Ingestion Layer**: Receives raw metrics via NATS from various services (S3, EC2, Lambda).
2. **Persistence**: Stores metrics in a time-series optimized database.
3. **Pricing Engine**: Applies rate cards and billing algorithms to aggregated usage data.
4. **Invoicing**: Generates and manages user invoices.

## Data Flow
`Services` -> `NATS (Subjects: dev.v1.billing.metric.*)` -> `Billing Ingestor` -> `Postgres/TimescaleDB` -> `Pricing Engine` -> `Invoices`
