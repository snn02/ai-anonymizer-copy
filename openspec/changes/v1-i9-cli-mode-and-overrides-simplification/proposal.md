## Why

CLI перегружен override-флагами, из-за чего пользователи регулярно ошибаются в запуске и попадают в неочевидные runtime-ветки. Для MVP нужно упростить управление режимами и сделать `config/env` единственным источником большинства параметров.

## What Changes

- Ввести единый режим запуска `mode` (`prod|mvp`) вместо связки `runtime_profile` + `insecure_no_secrets`.
- **BREAKING**: сократить CLI override до минимального набора: `--config`, `--mode`, `--workspace-path`, `--output-path`.
- **BREAKING**: убрать CLI override для остальных runtime-параметров (`allowed_hosts`, `max_pages`, `max_file_size_mb`, API/paging и т.д.); они читаются из `config/env`.
- Для `mvp` разрешить запуск на production API при явном `mode=mvp`, с предупреждением в CLI.
- Упростить структуру конфигурации и синхронизировать документацию/тесты под новый контракт.

## Capabilities

### New Capabilities
- `v1-i9-cli-mode-and-overrides-simplification`: упрощенный CLI-контракт с mode-first запуском и config-first параметризацией.

### Modified Capabilities
- `v1-1-ws-core`: заменить профильную модель `v1-1-ws + insecure flag` на `mode=mvp`, сохранить security-warning и требования к mvp-секретам.

## Impact

- Код: `cmd/anonym/main.go`, `internal/config/config.go`, `internal/preflight/preflight.go`, `internal/anonymizer/run.go`.
- Тесты: `cmd/anonym/main_test.go`, `internal/config/config_test.go`, `internal/preflight/doctor_test.go`, `internal/anonymizer/run_test.go`.
- Документация: `config.example.yaml`, `docs/user-guide.md`, `docs/technical/configuration.md`, `docs/test-scenarios.md`, `docs/versions/v1/plan.md`, `docs/versions/v1/action-log.md`.
- OpenSpec: новый feature-пакет `v1-i9-cli-mode-and-overrides-simplification` и delta spec для `v1-1-ws-core`.
