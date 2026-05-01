## 1. Run Header Implementation

- [x] 1.1 Добавить передачу header `user-id` в `POST /v1/tasks/file_anonymization` при непустом `api.user_id`.
- [x] 1.2 Сохранить backward-compatible поведение: если `api.user_id` пуст, `run` выполняется без `user-id` header.

## 2. Validation and Diagnostics

- [x] 2.1 Подтвердить, что невалидный `api.user_id` (не UUID) блокируется preflight/runtime до API-запроса.
- [x] 2.2 Проверить, что ошибки авторизации `401` в `run` остаются диагностируемыми и не маскируются.

## 3. Test Coverage

- [x] 3.1 Добавить/обновить тесты `internal/anonymizer/run_test.go` на отправку `user-id` и сценарий без `user-id`.
- [x] 3.2 Обновить CLI/integration тесты (`cmd/anonym/main_test.go`) для config-first сценария с `api.user_id`.

## 4. Documentation Sync (v1)

- [x] 4.1 Обновить `docs/technical/configuration.md` и `docs/user-guide.md` по контракту `run` (`user-id` в header при наличии).
- [x] 4.2 Синхронизировать `docs/test-scenarios.md`, `docs/versions/v1/plan.md`, `docs/versions/v1/action-log.md`, `docs/versions/v1/openspec.md`.
