## Why

В текущем MVP-профиле `v1-1-ws` команда `run` может игнорировать ожидаемые значения из `config` и уходить в чтение секретов из OS Secret Store. Из-за этого пользователь видит ошибку `audit key read failed: secret not found`, хотя `doctor` и `list` работают, а параметры MVP уже заданы в конфиге.

## What Changes

- Уточнить и зафиксировать единый runtime-контракт для MVP-профиля `v1-1-ws`: при включенном `security_flags.insecure_no_secrets=true` секреты и профиль берутся из config/env/CLI без обращения к OS Secret Store.
- Привести поведение `run` к тому же источнику конфигурации, который используется в `doctor`.
- Добавить проверки и тесты приоритета источников (`CLI > env > config`) для `runtime_profile`, `insecure_no_secrets`, `api.auth_token`, `audit.hmac_secret`.
- Обновить пользовательскую и техническую документацию с рабочими примерами запуска через config.

## Capabilities

### New Capabilities
- `v1-i8-mvp-config-runtime`: Гарантированная работа MVP-профиля `v1-1-ws` через config/env/CLI без обязательного OS Secret Store.

### Modified Capabilities
- `v1-1-ws-core`: Уточнение требований по источникам секретов и переключению профиля для команды `run`.

## Impact

- Код: `internal/config`, `internal/anonymizer`, `cmd/anonym`.
- Тесты: `internal/config/*_test.go`, `internal/anonymizer/*_test.go`, e2e smoke для `run`.
- Документация: `docs/user-guide.md`, `docs/technical/configuration.md`, `docs/versions/v1/plan.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md`.
- OpenSpec: новый change-пакет `openspec/changes/v1-i8-mvp-config-runtime/` и delta spec для `v1-1-ws-core`.
