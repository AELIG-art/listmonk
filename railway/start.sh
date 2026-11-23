#!/bin/bash

set -euo pipefail

: "${LISTMONK_db__host:?Set LISTMONK_db__host to the PostgreSQL host}"
: "${LISTMONK_db__user:?Set LISTMONK_db__user to the PostgreSQL user}"
: "${LISTMONK_db__password:?Set LISTMONK_db__password to the PostgreSQL password}"
: "${LISTMONK_db__database:?Set LISTMONK_db__database to the PostgreSQL database}"
export LISTMONK_db__port="${LISTMONK_db__port:-5432}"
export PORT="${PORT:-8080}"

echo 'Waiting for PostgreSQL...'
deadline=$((SECONDS + 120))
until PGSSLMODE="${LISTMONK_db__ssl_mode:-disable}" pg_isready -q -t 2 \
    -h "$LISTMONK_db__host" -p "$LISTMONK_db__port" \
    -U "$LISTMONK_db__user" -d "$LISTMONK_db__database"; do
    if (( SECONDS >= deadline )); then
        echo 'PostgreSQL did not become ready within 120 seconds.' >&2
        exit 1
    fi
    sleep 1
done

./listmonk --install --idempotent --yes
./listmonk --upgrade --yes

cleanup() {
    kill "$listmonk_pid" "$caddy_pid" 2>/dev/null || true
    wait "$listmonk_pid" "$caddy_pid" 2>/dev/null || true
}

./listmonk &
listmonk_pid=$!
caddy run --config /etc/caddy/Caddyfile --adapter caddyfile &
caddy_pid=$!

trap cleanup EXIT
trap 'exit 143' TERM
trap 'exit 130' INT
wait -n "$listmonk_pid" "$caddy_pid"

# A clean exit from either server still leaves the service unavailable.
exit 1
