# S3 Billing Algorithms

This document outlines how S3 metrics are transformed into billing records.

## 1. Storage Calculation
- **Metric**: `storage_utilization` (GB)
- **Algorithm**: Time-Weighted Average (TWA).
- **Frequency**: S3 sends a snapshot every 24 hours.
- **Rounding**: Round up to the nearest whole hour of GB usage.

## 2. Request Tiering
Requests are grouped into tiers based on their computational and business cost.
- **Tier 1 (Expensive)**: PUT, POST, LIST, COPY.
- **Tier 2 (Cheap)**: GET, SELECT, HEAD.
- **Algorithm**: Summation over the billing cycle.

## 3. Data Transfer
- **Metric**: `data_egress` (MB)
- **Algorithm**: Simple summation of all bytes transferred out of the Cloud environment.
