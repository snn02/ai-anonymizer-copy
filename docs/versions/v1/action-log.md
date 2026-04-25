# v1 action-log

## Формат записи

- Дата
- Тип: изменение | ошибка | реакция
- Описание
- Ссылка

## Записи

### 2026-04-25

- Тип: изменение
- Описание: перенесена структура документации в lowercase и заполнена на основе исходных документов плана и сценариев.
- Ссылка: `../../doc-map.md`

### 2026-04-25

- Тип: реакция
- Описание: зафиксировано, что feature-спеки в OpenSpec создаются только после согласования workflow.
- Ссылка: `openspec.md`

### 2026-04-25

- Тип: изменение
- Описание: требования v1 уточнены по безопасности и надежности: boundary enforcement, path hardening, PII-safe naming, anti-DoS лимиты, сетевая защита, lifecycle секретов (базовый), HMAC-аудит, правила fuzzy run.
- Ссылка: `plan.md`

### 2026-04-25

- Тип: изменение
- Описание: добавлен технический контракт конфигурации (`docs/technical/configuration.md`) и шаблон `config.example.yaml`; закреплены правила источников (flags > env > file) и запрет plaintext fallback для секретов.
- Ссылка: `../../technical/configuration.md`

### 2026-04-25

- Тип: изменение
- Описание: workflow для v1 согласован; фичи OpenSpec создаются по итерациям, каждая итерация должна быть законченной и тестируемой.
- Ссылка: `openspec.md`

### 2026-04-25

- Тип: изменение
- Описание: запланирована первая итерация `v1-i1` и создан feature-пакет OpenSpec с фокусом на `doctor` preflight и безопасный `scan/list` без API-отправки.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/proposal.md`

### 2026-04-25

- Тип: изменение
- Описание: начата реализация `v1-i1`: добавлен стартовый каркас trusted CLI на Go и реализован базовый `doctor preflight` с тестами на boundary/https/allowlist/max_parallel_runs.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: зафиксирован технический блокер окружения разработки: команда `go` отсутствует, из-за чего локальный прогон unit-тестов и форматирования `gofmt` отложен до установки Go toolchain.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: блокер по Go toolchain снят; подтвержден рабочий прогон `go test ./...` и `gofmt` для модулей `cmd/anonym`, `internal/config`, `internal/preflight`, `internal/scanner`, `internal/catalog`.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: изменение
- Описание: продолжена реализация `v1-i1`: добавлены команды `scan`/`list`, безопасный discovery файлов в `raw_path` и локальный каталог очереди `.anonym/catalog.json` без сетевой отправки.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: изменение
- Описание: по TDD расширен `doctor preflight`: добавлены проверки наличия partner secret (через `DoctorDeps.SecretChecker`) и доступности API endpoint `GET /v1/tasks/anonymization_fields` с timeout; добавлены новые unit-тесты `internal/preflight/doctor_test.go`.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: изменение
- Описание: по TDD добавлены лимиты discovery для `scan`: `max_file_size_mb` в конфигурации (flags/env) и фильтрация файлов по размеру через `scanner.DiscoverWithLimits`; добавлены/обновлены тесты `internal/config/config_test.go` и `internal/scanner/scanner_test.go`.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: результаты текущего состояния `v1-i1` зафиксированы после полного прогона `go test ./...`; подготовлено предложение на следующую итерацию `run + api adapter`.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/proposal.md`

### 2026-04-25

- Тип: изменение
- Описание: по TDD добавлены базовые preflight ACL-проверки: `raw_path` и `output_path` должны существовать и быть директориями, для `output_path` выполняется write-probe; тесты `internal/preflight/*` обновлены на реальные временные директории.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: подтверждено состояние итерации `v1-i1`: после внедрения лимитов discovery и ACL-checks полный прогон `go test ./...` проходит успешно.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`
