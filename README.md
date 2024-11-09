# Kafka Configs Metrics Exporter

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/EladLeev/kafka-config-metrics/release.yml?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/EladLeev/kafka-config-metrics)](https://goreportcard.com/report/github.com/EladLeev/kafka-config-metrics)
[![Renovate](https://img.shields.io/badge/renovate-enabled-%231A1F6C?logo=renovatebot)](https://renovatebot.com)

"Kafka Configs Metrics Exporter" for Prometheus allows you to export some of the Kafka configuration as metrics.

## Motivation

Unlike some other systems, Kafka doesn't expose its configurations as metrics.  
There are few useful configuration parameters that might be beneficial to collect in order to improve the visibility and alerting over Kafka.

A good example might be `log.retention.ms` parameter per topic, which can be integrated into Kafka's dashboards to extend its visibility, or to integrate it into an alerting query to create smarter alerts or automations based on topic retention.

Therefore, I decided to create a Prometheus exporter to collect those metrics.

Read more on [Confluent Blog](https://www.confluent.io/blog/kafka-lag-monitoring-and-metrics-at-appsflyer/)

Table of Contents
-----------------

- [Kafka Configs Metrics Exporter](#kafka-configs-metrics-exporter)
  - [Motivation](#motivation)
  - [Table of Contents](#table-of-contents)
  - [Build from source](#build-from-source)
    - [Prerequisites](#prerequisites)
    - [Building Steps](#building-steps)
  - [Using Docker Image](#using-docker-image)
  - [Helm](#helm)
  - [Configuration](#configuration)
    - [Clusters](#clusters)
    - [Authentication](#authentication)
    - [Prometheus Configuration](#prometheus-configuration)
    - [Endpoints](#endpoints)
  - [Dashboard Example](#dashboard-example)
  - [Contributing](#contributing)
  - [License](#license)

## Build from source

### Prerequisites

- Install Go version 1.20+

### Building Steps

1. Clone this repository

```bash
git clone https://github.com/EladLeev/kafka-config-metrics
cd kafka-config-metrics
```

2. Build the exporter binary

```bash
go build -o kcm-exporter .
```

3. Copy and edit the [config file](https://github.com/EladLeev/kafka-config-metrics/blob/master/kcm.yaml) from this repository and point it into your Kafka clusters.  
   _Use topic filtering as needed._

4. Deploy the binary and run the exporter

```bash
cp ~/my_kcm.yaml /opt/kcm/kcm.yaml
./kcm-exporter
```

The exporter will use `/opt/kcm/kcm.yaml` as default.

## Using Docker Image

1. Clone this repository

```bash
git clone https://github.com/EladLeev/kafka-config-metrics
cd kafka-config-metrics
```

2. Build the Docker image

```bash
docker build . -t kcm-exporter
```

3. Run it with your custom configuration file

```bash
docker run -p 9899:9899 -v ~/my_kcm.yaml:/opt/kcm/kcm.yaml kcm-exporter:latest
```

## Helm

Helm chart is available under the `/charts` dir.

To install the chart:

```bash
helm install kafka-config-metrics ./charts/kafka-config-metrics -f values.yaml
```

## Configuration

This project tried to stand in the Prometheus community [best practices](https://prometheus.io/docs/instrumenting/writing_exporters/) -  
"You should aim for an exporter that requires no custom configuration by the user beyond telling it where the application is".

In fact, you don't really need to change anything beyond the `clusters` configuration.

Example configuration:

```yaml
global:
  port: ":9899"    # Which port to bind
  timeout: 3       # HTTP server timeout in seconds

log:
  format: "text"   # Log format: text or json
  level: "info"    # Log level: info, debug, trace

kafka:
  defaultRefreshRate: 60          # Refresh rate in seconds
  minKafkaVersion: "2.8.0"       # Minimum Kafka version
  adminTimeout: 5                 # Admin client timeout in seconds

clusters:
  prod:
    brokers:
      - "kafka01-prod:9092"
    topicFilter: ""              # Optional regex filter
    tls:
      enabled: false
      # caCert: "/path/to/ca.crt"
      # clientCert: "/path/to/client.crt"
      # clientKey: "/path/to/client.key"
    sasl:
      enabled: false
      usernameEnv: "KAFKA_USERNAME"
      passwordEnv: "KAFKA_PASSWORD"
      mechanism: "PLAIN"         # PLAIN, SCRAM-SHA-256, or SCRAM-SHA-512

  test:
    brokers:
      - "kafka02-test:9092"
      - "kafka03-test:9092"
    topicFilter: "^(qa-|test-).*$"
```

### Authentication

The exporter supports both TLS and SASL authentication:

- TLS: Provide CA certificate and optionally client certificate/key
- SASL: Supports PLAIN, SCRAM-SHA-256, and SCRAM-SHA-512 mechanisms
  - Credentials are loaded from environment variables

### Prometheus Configuration

When setting this exporter in the Prometheus targets, bear in mind that topic configs are not subject to change that often in most use cases.  
Setting a higher `scrape_interval`, let's say to 10 minutes, will lead to lower requests rate to the Kafka cluster while still keeping the exporter functional.

```yaml
scrape_configs:
  - job_name: 'kcm'
    scrape_interval: 600s
    static_configs:
      - targets: ['kcm-prod:9899']
```

### Endpoints

`/metrics` - Metrics endpoint

`/-/healthy` - This endpoint returns 200 and should be used to check the exporter health.

`/-/ready`- This endpoint returns 200 when the exporter is ready to serve traffic.

## Dashboard Example

![Dashboard Sample](doc/dashboard.png)

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details of submitting a pull requests.

## License

This project is licensed under the Apache License - see the [LICENSE.md](LICENSE.md) file for details.