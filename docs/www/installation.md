---
title: Installation
category: 636c08d51eb043000f8ce20e
slug: installation
hidden: false
---

# Installation

### Install BindPlane OP

BindPlane OP does not have any dependencies and can run on Windows, Linux, or macOS although we recommend installing BindPlane OP server on a supported Linux distribution.

The BindPlane CLI (included in the BindPlane OP binary) will run on Linux, Windows, and macOS. For a detailed list of commands and installation instructions for remote clients, see our [CLI](doc:cli) page.

The OpenTelemetry Agent can be downloaded and installed using the UI. More instruction on agent installation can be found in our getting started guide [here](doc:getting-started#step-3-install-an-agent).

To download packages directly, see our [Downloads](doc:downloads) page or visit the [BindPlane OP Github Repository](https://github.com/observIQ/bindplane-op).

> 📘 Note: initialize your server
>
> After installing BindPlane OP, make sure to run the `bindplane init` command. The specific command for your OS can be found when running the installer.

### Linux

The following distributions are officially supported:

* Red Hat, Centos, Oracle Linux 7 and 8
* Debian 10 and 11
* Ubuntu LTS 18.04, 20.04
* Suse Linux 12 and 15
* Alma and Rocky Linux

While BindPlane OP will generally run on any modern distribution of Linux, systemd is the only supported init system. BindPlane OP will install on a non-systemd system, however, service management will be up to the user and is not a supported solution.

#### Server

Debian and RHEL style packages are available for BindPlane Server.

An installation script is available to simplify installation.

```bash
curl -fsSlL https://github.com/observiq/bindplane-op/releases/latest/download/install-linux.sh | bash -s --
```

Once installed, you can check the service.

```bash
sudo systemctl status bindplane
```

#### Client

Debian, RHEL, and Alpine packages are available for BindPlane Client. The packages will install the same binary included with the BindPlane Server package, but will not create a user, config, log, storage directory, or service. To see a full list of supported packages, see the [Downloads](doc:downloads) page.

Once installed, the `bindplane` command will be available and can be used to connect to a BindPlane Server. See the [Configuration](doc:configuration) page for configuration instructions.

**Installing Client on Debian / Ubuntu**

Example (amd64):

```bash
curl -L -o bindplane.deb https://github.com/observIQ/bindplane-op/releases/download/v1.0.8/bindplanectl_1.0.1_linux_amd64.deb
sudo apt-get install -f ./bindplane.deb
```

**Installing Client on Centos / RHEL**

Example (amd64):

```bash
sudo dnf install https://github.com/observIQ/bindplane-op/releases/download/v1.0.1/bindplanectl_1.0.1_linux_amd64.rpm
```

### macOS

Any macOS version 10.13 or newer is supported by BindPlane OP.

#### Client

A script is available for installation on macOS.

```bash
curl -fsSlL https://github.com/observiq/bindplane-op/releases/latest/download/install-macos.sh | bash -s --
```

Note: The macOS client includes some server configuration, however BindPlane Server is not officially supported.

### Docker

BindPlane OP can run as a container using Docker. The following commands will:

* Name container `bindplane`
* Keep persistent data in a volume named `bindplane`
* Expose port 3001 (REST and Websocket)

Be sure to replace `latest` with the disired release image tag. E.g: `observiq/bindplane:1.3.0`.

```bash
docker volume create bindplane

docker run -d \
    --name bindplane \
    --restart always \
    --mount source=bindplane,target=/data \
    -e BINDPLANE_CONFIG_USERNAME=admin \
    -e BINDPLANE_CONFIG_PASSWORD=admin \
    -e BINDPLANE_CONFIG_SERVER_URL=http://localhost:3001 \
    -e BINDPLANE_CONFIG_REMOTE_URL=ws://localhost:3001 \
    -e BINDPLANE_CONFIG_SESSIONS_SECRET=2c23c9d3-850f-4062-a5c8-3f9b814ae144 \
    -e BINDPLANE_CONFIG_SECRET_KEY=8a5353f7-bbf4-4eea-846d-a6d54296b781 \
    -e BINDPLANE_CONFIG_LOG_OUTPUT=stdout \
    -p 3001:3001 \
    observiq/bindplane:latest
```

### Upgrade BindPlane OP

To upgrade to the latest version of BindPlane OP, rerun the installer command or package for your operating system.

### Uninstall BindPlane OP

#### Linux

**Server**

1. Stop the process:

```bash
sudo systemctl disable bindplane && sudo systemctl stop bindplane
```

2. Remove the package

* Debian and Ubuntu:

```bash
sudo apt-get remove bindplane -y
```

* CentOS and RHEL 8 and newer (use yum for anything older)

```bash
sudo dnf remove bindplane -y
```

3. Optionally remove leftover data

```bash
sudo rm -rf /etc/bindplane /var/lib/bindplane /var/log/bindplane
```

**Client**

1. Remove the package

* Debian and Ubuntu:

```bash
sudo apt-get remove bindplane -y
```

* CentOS and RHEL 8 and newer (use yum for anything older)

```bash
sudo dnf remove bindplane -y
```

2. Remove local data

```bash
rm -rf ~/.bindplane
```

**Agent**

Run the following command:

```bash
sudo sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_unix.sh)" install_unix.sh -r
```

#### macOS

**Client**

1. Remove the binary

```bash
rm -f ~/bin/bindplane
```

2. Optionally remove profiles, data, and logs

```bash
rm -rf ~/.bindplane
```

**Agent**

Run the following command:

```bash
sh -c "$(curl -fsSlL https://github.com/observiq/observiq-otel-collector/releases/latest/download/install_macos.sh)" install_macos.sh -r
```

#### Windows

**Agent**

1. Navigate to the control panel, then to the "Uninstall a program" dialog.
2. Locate the `observIQ OpenTelemetry Collector` entry, and select uninstall.
3. Follow the wizard to complete removal of the collector.

Alternatively, run the Powershell command below:

```powershell
(Get-WmiObject -Class Win32_Product -Filter "Name = 'observIQ OpenTelemetry Collector'").Uninstall()
```
