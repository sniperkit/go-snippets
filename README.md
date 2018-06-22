**Attention**:
The image has been reworked to use Alpine instead of Ubuntu. As a consequence, the UID of the user that the application process runs under has changed from `999` (which was the default in Ubuntu) to `500`. Thereby, if you mount a host directory to persist data and configs, you will need to `chown -R 500:500` your local files before switching to the new version.

# JetBrains Hub
[![](https://images.microbadger.com/badges/image/esycat/jetbrains-hub.svg)](https://microbadger.com/images/esycat/jetbrains-hub "Get your own image badge on microbadger.com")
[![](https://images.microbadger.com/badges/version/esycat/jetbrains-hub.svg)](https://microbadger.com/images/esycat/jetbrains-hub "Get your own version badge on microbadger.com")

[Hub](https://jetbrains.com/hub/) is a system that provides centralized authentication, authorization; users, groups, permissions and project management across multiple installations of [JetBrains](https://jetbrains.com/) team collaboration tools.

Version `2018.1`, build `8691` (released February 21, 2018).

The image is based on [Alpine 3.4 with OpenJDK JRE 8](https://hub.docker.com/r/esycat/java/).

## Persistent Data
The Hub is configured to store all data (including backups, logs and temporary files) under `/var/lib/hub` directory, which is also a Docker volume. In addition, `/opt/hub/conf` directory is used for configuration files.

There are multiple approaches to handling persistent storage with Docker. For detailed information, see [Manage data in containers](https://docs.docker.com/engine/tutorials/dockervolumes/).

If a host directory is used, it should be writable by the application process, which runs as UID `500`.

## Usage

Pull the image, create a new container and start it:

```bash
docker pull esycat/jetbrains-hub
docker create --name jetbrains-hub -p 8080:8080 --restart=always esycat/jetbrains-hub
docker start jetbrains-hub
```

## Virtual Host
Typically, we would like to run the Hub behind a lightweight HTTP server. [`etc/nginx`]() contains an example virtual host configuration for Nginx.
