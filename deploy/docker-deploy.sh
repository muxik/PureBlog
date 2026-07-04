#!/usr/bin/env bash
set -Eeuo pipefail

APP_NAME="pureblog"
DEFAULT_REPO_URL="https://github.com/muxik/PureBlog.git"
DEFAULT_BRANCH="main"

log() {
  printf '[%s] %s\n' "$APP_NAME" "$*"
}

warn() {
  printf '[%s] warning: %s\n' "$APP_NAME" "$*" >&2
}

die() {
  printf '[%s] error: %s\n' "$APP_NAME" "$*" >&2
  exit 1
}

have() {
  command -v "$1" >/dev/null 2>&1
}

usage() {
  cat <<'EOF'
Usage:
  docker-deploy.sh [command]

Commands:
  deploy    Prepare source, start services, and show recent logs (default)
  prepare   Create deployment directory, fetch source, and create deploy/.env
  up        Build and start Docker Compose services
  logs      Follow Docker Compose logs
  tail      Show recent Docker Compose logs without following
  status    Show Docker Compose service status
  restart   Restart Docker Compose services
  down      Stop Docker Compose services

Environment:
  PUREBLOG_DEPLOY_DIR           Deployment directory, default: local repo or /opt/pureblog
  PUREBLOG_REPO_URL             Git repository URL, default: https://github.com/muxik/PureBlog.git
  PUREBLOG_BRANCH               Git branch, default: main
  PUREBLOG_PREPARE_SCRIPT_URL   Optional extra preparation script URL to download and run
  LOG_LINES                     Lines shown by deploy/tail, default: 120
EOF
}

resolve_local_repo_root() {
  local source_path="${BASH_SOURCE[0]:-}"
  local script_dir

  [[ -n "$source_path" && "$source_path" != "bash" && -f "$source_path" ]] || return 1
  script_dir="$(cd "$(dirname "$source_path")" && pwd)"

  if [[ -f "$script_dir/docker-compose.yml" && -d "$script_dir/../backend" && -d "$script_dir/../frontend" ]]; then
    cd "$script_dir/.." && pwd
    return 0
  fi

  return 1
}

LOCAL_REPO_ROOT="$(resolve_local_repo_root || true)"
DEPLOY_DIR="${PUREBLOG_DEPLOY_DIR:-${LOCAL_REPO_ROOT:-/opt/pureblog}}"
REPO_URL="${PUREBLOG_REPO_URL:-$DEFAULT_REPO_URL}"
BRANCH="${PUREBLOG_BRANCH:-$DEFAULT_BRANCH}"
PREPARE_SCRIPT_URL="${PUREBLOG_PREPARE_SCRIPT_URL:-}"
LOG_LINES="${LOG_LINES:-120}"

COMPOSE_FILE="$DEPLOY_DIR/deploy/docker-compose.yml"
ENV_FILE="$DEPLOY_DIR/deploy/.env"
ENV_EXAMPLE="$DEPLOY_DIR/deploy/.env.example"

require_command() {
  have "$1" || die "missing required command: $1"
}

compose() {
  if docker compose version >/dev/null 2>&1; then
    docker compose "$@"
    return
  fi

  if have docker-compose; then
    docker-compose "$@"
    return
  fi

  die "Docker Compose is required. Install Docker Compose v2 or docker-compose."
}

compose_args() {
  printf '%s\0' --env-file "$ENV_FILE" -f "$COMPOSE_FILE"
}

run_compose() {
  local args=()
  local item

  while IFS= read -r -d '' item; do
    args+=("$item")
  done < <(compose_args)

  compose "${args[@]}" "$@"
}

create_deployment_dir() {
  log "Creating deployment directory: $DEPLOY_DIR"
  mkdir -p "$DEPLOY_DIR"
}

is_empty_dir() {
  [[ -d "$1" ]] || return 1
  [[ -z "$(find "$1" -mindepth 1 -maxdepth 1 -print -quit)" ]]
}

fetch_source() {
  if [[ -n "$LOCAL_REPO_ROOT" && "$DEPLOY_DIR" == "$LOCAL_REPO_ROOT" ]]; then
    log "Using local repository: $DEPLOY_DIR"
    return
  fi

  require_command git

  if [[ -d "$DEPLOY_DIR/.git" ]]; then
    log "Updating repository in $DEPLOY_DIR"
    git -C "$DEPLOY_DIR" fetch origin "$BRANCH"
    git -C "$DEPLOY_DIR" checkout "$BRANCH"
    git -C "$DEPLOY_DIR" pull --ff-only origin "$BRANCH"
    return
  fi

  if is_empty_dir "$DEPLOY_DIR"; then
    log "Cloning $REPO_URL#$BRANCH into $DEPLOY_DIR"
    git clone --branch "$BRANCH" "$REPO_URL" "$DEPLOY_DIR"
    return
  fi

  if [[ -f "$COMPOSE_FILE" && -d "$DEPLOY_DIR/backend" && -d "$DEPLOY_DIR/frontend" ]]; then
    log "Using existing deployment source in $DEPLOY_DIR"
    return
  fi

  die "$DEPLOY_DIR is not empty and does not look like a PureBlog checkout"
}

run_extra_prepare_script() {
  local tmp_script

  if [[ -z "$PREPARE_SCRIPT_URL" ]]; then
    return 0
  fi

  require_command curl

  log "Downloading preparation script: $PREPARE_SCRIPT_URL"
  tmp_script="$(mktemp)"
  curl -fsSL "$PREPARE_SCRIPT_URL" -o "$tmp_script"
  chmod +x "$tmp_script"

  log "Running preparation script"
  PUREBLOG_DEPLOY_DIR="$DEPLOY_DIR" "$tmp_script" "$DEPLOY_DIR"
  rm -f "$tmp_script"
}

random_secret() {
  if have openssl; then
    openssl rand -hex 32
    return
  fi

  if [[ -r /dev/urandom ]]; then
    LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c 64
    return
  fi
}

replace_env_value() {
  local key="$1"
  local value="$2"
  local file="$3"
  local tmp_file="${file}.tmp.$$"

  awk -v key="$key" -v value="$value" '
    BEGIN { replaced = 0 }
    $0 ~ "^" key "=" {
      print key "=" value
      replaced = 1
      next
    }
    { print }
    END {
      if (!replaced) {
        print key "=" value
      }
    }
  ' "$file" >"$tmp_file"
  mv "$tmp_file" "$file"
}

ensure_env_file() {
  [[ -f "$COMPOSE_FILE" ]] || die "missing compose file: $COMPOSE_FILE"
  [[ -f "$ENV_EXAMPLE" ]] || die "missing env example: $ENV_EXAMPLE"

  if [[ -f "$ENV_FILE" ]]; then
    log "Using existing env file: $ENV_FILE"
  else
    log "Creating env file from example: $ENV_FILE"
    cp "$ENV_EXAMPLE" "$ENV_FILE"

    local secret
    secret="$(random_secret || true)"
    if [[ -n "$secret" ]]; then
      replace_env_value "JWT_SECRET" "$secret" "$ENV_FILE"
    fi
  fi

  if grep -Eq '^(JWT_SECRET=change-me|ADMIN_PASSWORD=change-me|PUBLIC_ORIGIN=http://localhost)' "$ENV_FILE"; then
    warn "review $ENV_FILE before exposing this service publicly"
  fi
}

prepare() {
  create_deployment_dir
  fetch_source
  run_extra_prepare_script
  ensure_env_file
}

require_docker() {
  require_command docker
  docker info >/dev/null 2>&1 || die "Docker daemon is not reachable"
}

start_services() {
  require_docker
  ensure_env_file
  log "Starting Docker Compose services"
  run_compose up -d --build
}

show_tail_logs() {
  require_docker
  ensure_env_file
  log "Showing last $LOG_LINES log lines"
  run_compose logs --tail "$LOG_LINES"
}

follow_logs() {
  require_docker
  ensure_env_file
  log "Following logs"
  run_compose logs -f --tail "$LOG_LINES"
}

show_status() {
  require_docker
  ensure_env_file
  run_compose ps
}

restart_services() {
  require_docker
  ensure_env_file
  log "Restarting services"
  run_compose restart
}

stop_services() {
  require_docker
  ensure_env_file
  log "Stopping services"
  run_compose down
}

print_summary() {
  local command_hint="$DEPLOY_DIR/deploy/docker-deploy.sh"

  if [[ ! -x "$command_hint" ]]; then
    command_hint="PUREBLOG_DEPLOY_DIR=$DEPLOY_DIR $0"
  fi

  cat <<EOF

PureBlog deployment is ready.

Deployment directory: $DEPLOY_DIR
Compose file:          $COMPOSE_FILE
Env file:              $ENV_FILE

Useful commands:
  $command_hint status
  $command_hint logs
  $command_hint restart
  $command_hint down
EOF
}

command="${1:-deploy}"
shift || true

case "$command" in
  deploy)
    prepare
    start_services
    show_tail_logs
    print_summary
    ;;
  prepare)
    prepare
    ;;
  up | start)
    start_services
    ;;
  logs)
    follow_logs
    ;;
  tail)
    show_tail_logs
    ;;
  status | ps)
    show_status
    ;;
  restart)
    restart_services
    ;;
  down | stop)
    stop_services
    ;;
  help | -h | --help)
    usage
    ;;
  *)
    usage
    die "unknown command: $command"
    ;;
esac
