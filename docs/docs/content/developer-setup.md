# Developer setup
The app is a Go backend that server-renders the admin UI (HTML templates with JS, built with [bun](https://bun.sh)).


### Pre-requisites
- `go`
- `bun` (if you are working on the admin frontend)
- Postgres database, locally installed or through the supplied Docker development environment.


### First time setup
Clone the AELIG Mail source repository into a directory outside the Go `src` path. The project uses Go modules.

1. Copy `config.toml.sample` as `config.toml` and add your config.
2. Build the frontend and the binary: `make dist` (builds the SSR admin and the Go binary, embedding the assets). Once the binary is built, run `./aelig-mail --install` to run the DB setup. For subsequent dev runs, use `make run`.

> [mailhog](https://github.com/mailhog/MailHog) is an excellent standalone mock SMTP server (with a UI) for testing and dev.


### Running the Docker dev environment

- Run `make build-dev-docker` once to build the dev container which bundles Go, Postgres, Bun etc.
- Run `make run-dev-docker` to start it and enter its shell.
- Inside the shell, run `make run` which builds static assets, initializes a fresh DB, and starts AELIG Mail on `:9000`.
- For subsequent changes, simply Ctrl-C and re-run `make run` for changes to Go and static files to be picked up.


Run `make stop-dev-docker` to stop the container and `make rm-dev-docker` to completely remove it.

### Frontend tests

AELIG Mail maintains a Playwright suite in `tests/` for end-to-end testing the admin UI. It runs in the Docker development environment, which bundles Postgres, MailHog, and Bun.

#### Running

Follow the steps described above to run and enter the Docker dev environment.

Inside the dev container shell at `/app`:

```sh
make dist

cd tests
bun install
bun playwright install --with-deps chromium
bun run test
```

# Production build
Run `make dist` to build the SSR admin frontend and the Go binary, embedding the static assets into the self-contained `aelig-mail` executable.
