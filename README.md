<a href="https://zerodha.tech"><img src="https://zerodha.tech/static/images/github-badge.svg" align="right" /></a>

[![listmonk-logo](https://user-images.githubusercontent.com/547147/231084896-835dba66-2dfe-497c-ba0f-787564c0819e.png)](https://listmonk.app)

listmonk is a standalone, self-hosted, newsletter and mailing list manager. It is fast, feature-rich, and packed into a single binary. It uses a PostgreSQL database as its data store.

[![listmonk-dashboard](https://github.com/user-attachments/assets/689b5fbb-dd25-4956-a36f-e3226a65f9c4)](https://listmonk.app)

Visit [listmonk.app](https://listmonk.app) for more info. Check out the [**live demo**](https://demo.listmonk.app).

## Installation

### Railway

Add a [PostgreSQL service](https://docs.railway.com/databases/postgresql) named `Postgres` and a GitHub service for this repository in the same Railway project. The repository's `railway.json` selects `Dockerfile.railway`, which builds Listmonk from source and runs it behind Caddy. Use the image's default start command.

Paste these references into the Listmonk service's **Variables → Raw Editor**. If your PostgreSQL service has a different name, replace `Postgres` in each reference.

```dotenv
LISTMONK_db__host=${{Postgres.PGHOST}}
LISTMONK_db__port=${{Postgres.PGPORT}}
LISTMONK_db__user=${{Postgres.PGUSER}}
LISTMONK_db__password=${{Postgres.PGPASSWORD}}
LISTMONK_db__database=${{Postgres.PGDATABASE}}
LISTMONK_db__ssl_mode=disable
LISTMONK_app__root_url=https://listmonk-apanage.up.railway.app
```

Listmonk initializes an empty database and applies pending migrations on startup. Caddy listens on Railway's `PORT` and proxies to `127.0.0.1:9000`. If you set a target port manually, use `8080` or match your `PORT` variable.

Generate a public domain for the Listmonk service and open it to set up the administrator. To create the administrator automatically, set `LISTMONK_ADMIN_USER` and `LISTMONK_ADMIN_PASSWORD` before the first deployment. Set `LISTMONK_app__root_url` to your public HTTPS URL without a trailing slash. It overrides the Root URL saved in settings and is used for campaign and authentication links. Redeploy after changing it.

For cross-origin requests, set `LISTMONK_ORIGIN_0` through `LISTMONK_ORIGIN_4` to the allowed origins, such as `https://www.example.com`, without trailing slashes. To persist media stored on the filesystem, mount a [Railway volume](https://docs.railway.com/volumes) at `/listmonk/uploads`.

Configure the monthly SMTP message limit in **Settings → General**, or set `MAX_QUOTA` on the Railway application service to enforce a read-only limit. Each successful SMTP submission uses one slot, including campaign, test, and transactional messages. Counting starts with this upgrade and resets on the first of each month at 00:00 UTC. The **Limit monthly email sending** switch is off by default. Turn it off and save to allow unlimited sending. Removing `MAX_QUOTA` and redeploying restores the saved setting. `MAX_QUOTA` accepts nonnegative integers, including `0` to suspend email delivery.

### Docker

The latest image is available on DockerHub at [`listmonk/listmonk:latest`](https://hub.docker.com/r/listmonk/listmonk/tags?page=1&ordering=last_updated&name=latest).
Download and use the sample [docker-compose.yml](https://github.com/knadh/listmonk/blob/master/docker-compose.yml).


```shell
# Download the compose file to the current directory.
curl -LO https://github.com/knadh/listmonk/raw/master/docker-compose.yml

# Run the services in the background.
docker compose up -d
```
Visit `http://localhost:9000`

See [installation docs](https://listmonk.app/docs/installation)

__________________

### Binary
- Download the [latest release](https://github.com/knadh/listmonk/releases) and extract the listmonk binary.
- `./listmonk --new-config` to generate config.toml. Edit it.
- `./listmonk --install` to setup the Postgres DB (or `--upgrade` to upgrade an existing DB. Upgrades are idempotent and running them multiple times have no side effects).
- Run `./listmonk` and visit `http://localhost:9000`

See [installation docs](https://listmonk.app/docs/installation)
__________________


## Developers
listmonk is free and open source software licensed under AGPLv3. If you are interested in contributing, refer to the [developer setup](https://listmonk.app/docs/developer-setup). The backend is written in Go and the frontend is Vue with Buefy for UI. 


## License
listmonk is licensed under the AGPL v3 license.
