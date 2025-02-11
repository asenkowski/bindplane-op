---
title: "PostgreSQL"
category: 636c08e1212e49001e7a3032
parentDoc: 636c0a0c46142d00a50b384d
slug: "postgresql"
hidden: false
---
## Supported Platforms

| Platform | Metrics | Logs | Traces |
| :------- | :------ | :--- | :----- |
| Linux    | ✓       | ✓    |        |
| Windows  | ✓       | ✓    |        |
| macOS    | ✓       | ✓    |        |

## Prerequisites

This source supports PostgreSQL versions 10.18 and higher.

The monitoring user must be granted `SELECT` on `pg_stat_database`.

## Configuration Table

| Parameter           | Type  | Default | Description                                |
| :------------------ | :---- | :------ | :----------------------------------------- |
| enable_metrics | `bool` | true | Enable to collect metrics. |
| username* | `string` | | Username used to authenticate. |
| password* | `string` | | Password used to authenticate. |
| endpoint | `string` | localhost:5432 | The endpoint of the postgres server. If transport is set to unix, the endpoint will internally be translated from host:port to /host.s.PGSQL.port. |
| transport | `enum` | tcp | The transport protocol being used to connect to Postgres. Valid values are `tcp`, or `unix`. |
| databases | `strings` | | The list of databases for which the receiver will attempt to collect statistics. If an empty list is provided, the receiver will attempt to collect statistics for all databases. |
| collection_interval | `int` | 60 | How often (seconds) to scrape for metrics. |
| enable_tls | `bool` | false | Whether or not to use TLS. |
| enable_tlsinsecure_skip_verify | `bool` | false | Enable to skip TLS certificate verification. |
| ca_file | `string` | | Certificate authority used to validate the database server's TLS certificate. |
| cert_file | `string` | | A TLS certificate used for client authentication, if mutual TLS is enabled. |
| key_file | `string` | | A TLS private key used for client authentication, if mutual TLS is enabled. |
| enable_logs | `bool` | true | Enable to collect logs. |
| postgresql_log_path | `strings` | For CentOS / RHEL: `- \"/var/log/postgresql/postgresql_.log\"`  \nFor SLES: `- \"/var/lib/pgsql/data/log/postgresql_.log\"`  \nFor Debian / Ubuntu: `- \"/var/lib/pgsql/_/data/log/postgresql_.log\"` | Path to Postgres log file(s). |
| start_at | `enum` | end | Start reading file from 'beginning' or 'end'. |

<span style="color:red">\*_required field_</span>
