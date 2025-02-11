---
title: "Hadoop"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "hadoop"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Configuration Table

| Parameter               | Type      | Default                                             | Description                                       |
| :---------------------- | :-------- | :-------------------------------------------------- | :------------------------------------------------ |
| enable_metrics          | `bool`    | true                                                | Enable to send metrics.                           |
| collection_interval     | `int`     | 60                                                  | How often (seconds) to scrape for metrics.        |
| address                 | `string`  | localhost                                           | IP address or hostname to scrape for JMX metrics. |
| port                    | `int`     | 8004                                                | Port to scrape for JMX metrics.                   |
| jar_path                | `string`  | "/opt/opentelemetry-java-contrib-jmx-metrics.jar"   | Full path to the JMX metrics jar.                 |
| enable_logs             | `bool`    | true                                                | Enable to send logs.                              |
| enable_datanode_logs    | `bool`    | true                                                | Enable to collect datanode logs.                  |
| datanode_log_path       | `strings` | "/usr/local/hadoop/logs/hadoop-_-datanode-_.log"    | File paths to tail for datanode logs.             |
| enable_resourcemgr_logs | `bool`    | true                                                | Enable to collect resource manager logs.          |
| resourcemgr_log_path    | `strings` | "/usr/local/hadoop/logs/hadoop-_-resourcemgr-_.log" | File paths to tail for resource manager logs.     |
| enable_namenode_logs    | `bool`    | true                                                | Enable to collect namenode logs.                  |
| namenode_log_path       | `strings` | "/usr/local/hadoop/logs/hadoop-_-namenode-_.log"    | File paths to tail for namenode logs.             |
| start_at                | `enum`    | end                                                 | Start reading file from 'beginning' or 'end'.     |
