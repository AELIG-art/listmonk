#!/bin/bash

set -euo pipefail

# Support legacy deployment variables while preferring the AELIG Mail names.
: "${AELIG_MAIL_db__host:=${LISTMONK_db__host:-}}"
: "${AELIG_MAIL_db__user:=${LISTMONK_db__user:-}}"
: "${AELIG_MAIL_db__password:=${LISTMONK_db__password:-}}"
: "${AELIG_MAIL_db__database:=${LISTMONK_db__database:-}}"
: "${AELIG_MAIL_db__host:?Set AELIG_MAIL_db__host to the PostgreSQL host}"
: "${AELIG_MAIL_db__user:?Set AELIG_MAIL_db__user to the PostgreSQL user}"
: "${AELIG_MAIL_db__password:?Set AELIG_MAIL_db__password to the PostgreSQL password}"
: "${AELIG_MAIL_db__database:?Set AELIG_MAIL_db__database to the PostgreSQL database}"
export AELIG_MAIL_db__port="${AELIG_MAIL_db__port:-${LISTMONK_db__port:-5432}}"
export PORT="${PORT:-8080}"

# Copy legacy values to the new names for the application process.
export AELIG_MAIL_db__ssl_mode="${AELIG_MAIL_db__ssl_mode:-${LISTMONK_db__ssl_mode:-disable}}"
if [[ -z "${AELIG_MAIL_app__address:-}" && -n "${LISTMONK_app__address:-}" ]]; then
    export AELIG_MAIL_app__address="$LISTMONK_app__address"
fi
export AELIG_MAIL_app__address="${AELIG_MAIL_app__address:-127.0.0.1:9000}"
for index in 0 1 2 3 4; do
    new_var="AELIG_MAIL_ORIGIN_${index}"
    legacy_var="LISTMONK_ORIGIN_${index}"
    if [[ -z "${!new_var:-}" && -n "${!legacy_var:-}" ]]; then
        export "$new_var=${!legacy_var}"
    fi
done

echo 'Waiting for PostgreSQL...'
deadline=$((SECONDS + 120))
until PGSSLMODE="${AELIG_MAIL_db__ssl_mode}" pg_isready -q -t 2 \
    -h "$AELIG_MAIL_db__host" -p "$AELIG_MAIL_db__port" \
    -U "$AELIG_MAIL_db__user" -d "$AELIG_MAIL_db__database"; do
    if (( SECONDS >= deadline )); then
        echo 'PostgreSQL did not become ready within 120 seconds.' >&2
        exit 1
    fi
    sleep 1
done

./aelig-mail --install --idempotent --yes
./aelig-mail --upgrade --yes

cleanup() {
    kill "$aelig_mail_pid" "$caddy_pid" 2>/dev/null || true
    wait "$aelig_mail_pid" "$caddy_pid" 2>/dev/null || true
}

./aelig-mail &
aelig_mail_pid=$!
caddy run --config /etc/caddy/Caddyfile --adapter caddyfile &
caddy_pid=$!

trap cleanup EXIT
trap 'exit 143' TERM
trap 'exit 130' INT
wait -n "$aelig_mail_pid" "$caddy_pid"

# A clean exit from either server still leaves the service unavailable.
exit 1
