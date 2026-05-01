# EC2 Billing Metric Design

This document defines the metrics required to calculate costs for Elastic Compute Cloud (EC2) services.

## 1. Compute Utilization
Tracks the runtime and state of virtual machines.
- **Metric**: `ec2_compute_usage`
- **Fields**:
  - `instance_id`: Unique ID.
  - `instance_type`: e.g., t3.medium.
  - `state`: running | stopped | terminated.
- **Algorithm**: Time-Weighted Average (TWA) or Duration Summation.

## 2. EBS Block Storage
Tracks the provisioned storage capacity attached to instances.
- **Metric**: `ec2_ebs_volume_usage`
- **Fields**:
  - `volume_id`: Unique ID.
  - `volume_type`: gp3 | io2.
  - `provisioned_size_gb`: Capacity in GB.
- **Algorithm**: Max value per month (GB-Month).

## 3. Network Egress
Tracks data transferred out of the instance.
- **Metric**: `ec2_network_egress`
- **Fields**:
  - `instance_id`: Source ID.
  - `bytes_sent`: Bytes.
  - `destination_type`: internet | regional.
- **Algorithm**: Simple Summation.
