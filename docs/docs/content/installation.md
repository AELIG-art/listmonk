# Installation

AELIG Mail is a simple binary application that requires a Postgres database instance to run. The binary can be downloaded and run manually, or it can be run as a container with Docker compose.

## Binary
1. Obtain the AELIG Mail binary from your distribution channel and extract it. `amd64` is the standard build for 64-bit Intel and AMD CPUs.
1. `./aelig-mail --new-config` to generate config.toml. Edit the file.
1. `./aelig-mail --install` to install the tables in the Postgres DB (⩾ 12).
1. Run `./aelig-mail` and visit `http://localhost:9000` to create the Super Admin user and login.

!!! Tip
    To set the Super Admin username and password during installation, set the environment variables:
    `AELIG_MAIL_ADMIN_USER=myuser AELIG_MAIL_ADMIN_PASSWORD=xxxxx ./aelig-mail --install`

    To also create a Super Admin API user, for automation, add `AELIG_MAIL_ADMIN_API_USER=myapiuser`. Its API token is written to stderr as `export AELIG_MAIL_ADMIN_API_TOKEN="..."`.


## Docker

Build the application binary and container image from the source repository, then copy the Docker Compose file to your deployment directory and configure it:

```shell
make dist
docker build -t aelig-mail:latest .
docker compose up -d
```

Visit `http://localhost:9000` to create the Super Admin user and log in.

### Mounting a custom config.toml

You can mount a custom configuration file into the application container and point the startup command to that file with the `--config` flag. Use the paths and service names from your Compose file.

-----------

## Nightly

Nightly binaries and container images are published as GitHub prereleases. The container image is `ghcr.io/<repository-owner>/aelig-mail:nightly`.


## Compiling from source

To compile the latest unreleased version (`master` branch):

1. Make sure `go` and `bun` are installed on your system.
2. Clone the AELIG Mail source repository.
3. Build the project with `make dist` to compile the admin frontend and application binary.


## Tutorials
* AELIG Mail with Forward Email for Secure Newsletter Delivery
* Informal step-by-step on how to get started with AELIG Mail using *Railway*
* [Step-by-step tutorial for installation and all basic functions. *Amazon EC2, SES, docker & binary*](https://gist.github.com/MaximilianKohler/e5158fcfe6de80a9069926a67afcae11)
* Step-by-step guide on how to install and set up AELIG Mail on *AWS Lightsail with docker* (rameerez)
* Quick setup on any cloud server using *docker and caddy*
* *Binary* install on Ubuntu 22.04 as a service
* *Binary* install on Ubuntu 18.04 as a service (Apache & Plesk)
* *Binary and docker* on linux (techviewleo)
* [*Binary* install on your PC](https://www.youtube.com/watch?v=fAOBqgR9Yfo). Discussions of limitations: [1][2].
* *Docker on Rocky Linux 8* (nginx, Let's Encrypt SSL)
* *Docker* with nginx reverse proxy, certbot SSL, and Gmail SMTP
* [Install AELIG Mail on Self-hosting with *Pre-Configured AMI Package at AWS* by Single Click](https://meetrix.io/articles/how-to-install-llama-2-on-aws-with-pre-configured-ami-package/)
* [*Fly.io* working example](https://gitlab.com/votelog/apps/newsletter)
* *Docker via Ansible* with the MASH playbook
* *Docker Compose* on a managed VPS, with Caddy or Nginx (Simafri)
