#!/bin/sh
set -eu
# apt packages used by this script on Debian/Ubuntu:
#   golang-go
#   git
#   iproute2
#   postgresql-client
#   redis-tools
#
# external tools expected separately:
#   goose
#   docker
BUILD_SIGNATURE_FILE='./.runtime-cache/build.hash'
die() {
  printf '%s\n' "$*" >&2
  exit 1
}
show_help() {
  cat <<'EOF'
Usage:
  ./dev.sh
  ./dev.sh help
  ./dev.sh -h
  ./dev.sh --help
  ./dev.sh migrate
  ./dev.sh reset
  ./dev.sh reset <seed_name> [seed_name...]
  ./dev.sh seed [seed_name...]
  ./dev.sh mail-test
  ./dev.sh goose-*
Commands:
  ./dev.sh
    Запускает backend.
    Перед запуском применяет миграции и выполняет проверку сидов:
    если количество строк в CSV и в таблице не совпадает, выводится сообщение;
    если в таблице 0 строк, сид выполняется автоматически.
  ./dev.sh migrate
    Выполняет goose-up.
  ./dev.sh reset
    Полный reset базы через goose-reset.
  ./dev.sh reset <seed_name> [seed_name...]
    Точечно очищает таблицы сидов через seeds-import reset.
    Пример:
      ./dev.sh reset minerals
  ./dev.sh seed [seed_name...]
    Импортирует все сиды или только указанные.
    Примеры:
      ./dev.sh seed
      ./dev.sh seed minerals
      ./dev.sh seed kitchen_categories kitchen_filter_options
  ./dev.sh mail-test
    Запускает тест отправки почты.
Available seed names:
  minerals
  kitchen_categories
  kitchen_filter_options
  kitchen_ingredients_catalog
Notes:
  Имена сидов формируются из имён CSV-файлов в db/seeds.
  Для прямого использования importer CLI:
    go run ./cmd/seeds-import list
    go run ./cmd/seeds-import verify --auto-seed-zero
    go run ./cmd/seeds-import import minerals
    go run ./cmd/seeds-import reset minerals
EOF
}
need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing required command: $1"
}
load_env_file() {
  [ -f "./.env" ] || return 0
  set -a
  . "./.env"
  set +a
}
require_env() {
  var_name="$1"
  eval "var_value=\${$var_name-}"
  [ -n "$var_value" ] || die "missing required env: $var_name"
}
run_or_die() {
  if ! "$@"; then
    die "command failed: $*"
  fi
}
ensure_runtime_cache_dir() {
  mkdir -p "./.runtime-cache"
}
collect_hash_inputs() {
  for path in "$@"; do
    [ -e "$path" ] || continue
    if [ -d "$path" ]; then
      find "$path" -type f | LC_ALL=C sort
    else
      printf '%s\n' "$path"
    fi
  done
}
hash_paths() {
  need_cmd sha256sum
  ensure_runtime_cache_dir
  files_file="$(mktemp "./.runtime-cache/hash-inputs.XXXXXX")"
  collect_hash_inputs "$@" >"$files_file"
  if [ ! -s "$files_file" ]; then
    rm -f "$files_file"
    die "no input files found for hashing"
  fi
  digest="$(
    while IFS= read -r file; do
      sha256sum "$file"
    done < "$files_file" | sha256sum | awk '{print $1}'
  )"
  rm -f "$files_file"
  printf '%s\n' "$digest"
}
signature_is_current() {
  signature_file="$1"
  current_signature="$2"
  [ -f "$signature_file" ] || return 1
  saved_signature="$(tr -d '\r\n' < "$signature_file")"
  [ -n "$saved_signature" ] && [ "$saved_signature" = "$current_signature" ]
}
write_signature() {
  signature_file="$1"
  current_signature="$2"
  ensure_runtime_cache_dir
  printf '%s\n' "$current_signature" > "$signature_file"
}
current_app_signature() {
  hash_paths \
    ./go.mod \
    ./go.sum \
    ./cmd/app \
    ./internal \
    ./db \
    ./utils
}
build_binary() {
  need_cmd go
  if command -v git >/dev/null 2>&1 && [ -d .git ]; then
    rev="$(git rev-parse --short=12 HEAD 2>/dev/null || printf '%s' unknown)"
    build_time="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    dirty="false"
    if ! git diff --quiet --ignore-submodules HEAD -- 2>/dev/null; then
      dirty="true"
    fi
    if ! git diff --quiet --ignore-submodules --cached 2>/dev/null; then
      dirty="true"
    fi
    ldflags="-s -w -X 'main.BuildRevision=${rev}' -X 'main.BuildTime=${build_time}' -X 'main.BuildDirty=${dirty}'"
    run_or_die go build -buildvcs=true -trimpath -ldflags "$ldflags" -o ./main ./cmd/app
    return 0
  fi
  run_or_die go build -trimpath -ldflags "-s -w" -o ./main ./cmd/app
}
ensure_binary() {
  current_signature="$(current_app_signature)"
  if [ ! -x ./main ]; then
    build_binary
    write_signature "$BUILD_SIGNATURE_FILE" "$current_signature"
    return 0
  fi
  if ! signature_is_current "$BUILD_SIGNATURE_FILE" "$current_signature"; then
    build_binary
    write_signature "$BUILD_SIGNATURE_FILE" "$current_signature"
  fi
}
wait_postgres() {
  need_cmd psql
  timeout_seconds="${POSTGRES_WAIT_TIMEOUT:-60}"
  started_at="$(date +%s)"
  while ! PGPASSWORD="$POSTGRES_PASSWORD" psql \
    -h "$POSTGRES_HOST" \
    -p "$POSTGRES_PORT" \
    -U "$POSTGRES_USER" \
    -d "$POSTGRES_DB" \
    -c 'select 1' >/dev/null 2>&1
  do
    now="$(date +%s)"
    elapsed=$((now - started_at))
    [ "$elapsed" -lt "$timeout_seconds" ] || die "postgres is not ready after ${timeout_seconds}s"
    sleep 1
  done
}
run_make_target() {
  target="$1"
  need_cmd make
  [ -n "$target" ] || die "missing make target"
  ensure_runtime_cache_dir
  log_file="$(mktemp "./.runtime-cache/make.XXXXXX")"
  pipe_file="$(mktemp -u "./.runtime-cache/make.pipe.XXXXXX")"
  mkfifo "$pipe_file"
  tee "$log_file" <"$pipe_file" &
  tee_pid=$!
  if make "$target" >"$pipe_file" 2>&1; then
    wait "$tee_pid"
    rm -f "$pipe_file"
    rm -f "$log_file"
    return 0
  fi
  wait "$tee_pid"
  rm -f "$pipe_file"
  if grep -qi 'no next version found' "$log_file"; then
    rm -f "$log_file"
    return 0
  fi
  cat "$log_file" >&2
  rm -f "$log_file"
  die "make target failed: ${target}"
}
run_go_command() {
  cmd_path="$1"
  shift
  [ -d "$cmd_path" ] || die "missing command directory: ${cmd_path}"
  run_or_die go run "$cmd_path" "$@"
}
run_seed_check() {
  run_go_command ./cmd/seeds-import verify --auto-seed-zero
}
ensure_schema_ready() {
  run_make_target goose-up
}
run_action() {
  action="${1-}"
  case "$action" in
    help|-h|--help)
      show_help
      ;;
    "")
      exec ./main
      ;;
    migrate)
      run_make_target goose-up
      ;;
    reset)
      if [ $# -gt 1 ]; then
        shift
        run_go_command ./cmd/seeds-import reset "$@"
        return 0
      fi
      run_make_target goose-reset
      ;;
    seed)
      shift
      ensure_schema_ready
      run_go_command ./cmd/seeds-import import "$@"
      ;;
    mail-test)
      run_go_command ./cmd/mail-test
      ;;
    goose-*)
      run_make_target "$action"
      ;;
    *)
      return 0
      ;;
  esac
}
ensure_app_port_available() {
  need_cmd ss
  line="$(ss -ltnp "( sport = :${GO_PORT} )" 2>/dev/null | awk 'NR>1 {print; exit}')"
  [ -z "$line" ] && return 0
  die "application port is already in use: ${GO_PORT}; ${line}"
}
main() {
  load_env_file
  case "${1-}" in
    help|-h|--help)
      show_help
      return 0
      ;;
  esac
  require_env POSTGRES_HOST
  require_env POSTGRES_PORT
  require_env POSTGRES_USER
  require_env POSTGRES_PASSWORD
  require_env POSTGRES_DB
  require_env GO_PORT
  ensure_app_port_available
  ensure_binary
  wait_postgres
  if [ $# -eq 0 ]; then
    ensure_schema_ready
    run_seed_check
  fi
  run_action "$@"
}
main "$@"
