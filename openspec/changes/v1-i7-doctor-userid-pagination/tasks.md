# tasks: v1-i7 doctor userid pagination

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| `doctor` выполняет API-check по актуальному контракту tasks endpoint | T1 | U1, U2 | `user-id` header и `page/per_page` query отправляются в preflight |
| параметры доступны пользователю в config/env/CLI | T2 | U3 | флаги и env читаются с приоритетом `CLI > env > config` |

## tasks

### T1. doctor request contract

- [x] добавить `user-id` header в preflight `doctor`;
- [x] добавить query `page` и `per_page` в preflight `doctor`;
- [x] добавить валидацию `api.user_id` (UUID) и `api.fields_page/per_page` (>=1).

### T2. config + cli + docs

- [x] добавить параметры `api.user_id`, `api.fields_page`, `api.fields_per_page` в config/env;
- [x] добавить CLI-флаги `--api-user-id`, `--api-fields-page`, `--api-fields-per-page`;
- [x] синхронизировать `docs/versions/v1/plan.md`, `docs/user-guide.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md`.

## tests

- [x] U1: `doctor` отправляет `user-id` header.
- [x] U2: `doctor` отправляет query `page` и `per_page`.
- [x] U3: параметры читаются из config/env/CLI с ожидаемым приоритетом.
- [x] U4: `go test ./...` зеленый.
