## Why

Сейчас `doctor` проходит успешно, но `run` получает `401` на `POST /v1/tasks/file_anonymization` в production API. По актуальной документации метода в headers присутствует `user-id`, и его отсутствие в `run` может блокировать авторизацию/маршрутизацию запроса.

## What Changes

- Добавить передачу заголовка `user-id` в `run`-запрос `POST /v1/tasks/file_anonymization` при наличии `api.user_id`.
- Уточнить поведение по умолчанию: если `api.user_id` не задан, используется текущий совместимый режим (без header или с согласованным fallback, определяется в реализации по контракту API).
- Добавить тесты на корректную отправку `user-id` и на отсутствие регрессий для `prod/mvp`.
- Синхронизировать документацию `run`-контракта, чтобы `doctor` и `run` использовали единые правила по `user-id`.

## Capabilities

### New Capabilities
- `v1-i10-run-userid-header`: поддержка `user-id` в `run`-вызове `file_anonymization` по API-контракту.

### Modified Capabilities
- `v1-1-ws-core`: уточнение контрактов заголовков `run` в режимах `prod|mvp` с учетом `user-id`.

## Impact

- Код: `internal/anonymizer/run.go`, `internal/config/config.go` (если потребуется уточнение defaults/валидации), `cmd/anonym`.
- Тесты: `internal/anonymizer/run_test.go`, `cmd/anonym/main_test.go`.
- Документация: `docs/technical/configuration.md`, `docs/user-guide.md`, `docs/test-scenarios.md`, `docs/versions/v1/plan.md`, `docs/versions/v1/action-log.md`.
- OpenSpec: change `v1-i10-run-userid-header` + delta spec для `v1-1-ws-core`.
