# tasks: v1-i3 doctor api header alignment

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| `doctor` проверяет доступность API с auth/header-инвариантами runtime | T1, T2 | U1, I1 | `internal/preflight` формирует `GET /v1/tasks/anonymization_fields` c `Authorization` + `partner-id` |
| `api.partner_id` должен быть UUID | T3 | U2, I2 | `doctor` блокирует невалидный `partner_id` до сетевого вызова |
| запрет небезопасного fallback | T2, T4 | U3, I3 | preflight использует только секрет из OS secret store, без plaintext fallback |

## security DoD (для T1-T4)

- [x] runtime реализован в production-пути `doctor`/preflight;
- [x] минимум один негативный тест на каждый security-check;
- [x] отсутствует insecure production fallback;
- [x] evidence зафиксирован прогоном `go test ./...`.

## tasks

### T1. preflight api request header alignment

- [x] добавить формирование заголовка `Authorization` в `doctor` при проверке `GET /v1/tasks/anonymization_fields`;
- [x] добавить формирование заголовка `partner-id` из `api.partner_id`;
- [x] сохранить текущий timeout/allowlist/https контроль без ослаблений.

### T2. secret-source hardening in doctor path

- [x] убедиться, что `doctor` использует секрет только из OS secret store;
- [x] исключить любой plaintext fallback в preflight-пути;
- [x] стандартизировать сообщение ошибки при отсутствии секрета для preflight API-check.

### T3. partner id uuid validation

- [x] добавить runtime-валидацию UUID-формата `api.partner_id` в `doctor` path;
- [x] обеспечить явную детерминированную ошибку при невалидном формате.

### T4. tests and regression protection

- [x] добавить/обновить unit-тесты на обязательные заголовки preflight-запроса;
- [x] добавить/обновить unit-тесты на UUID-валидацию `api.partner_id`;
- [x] добавить integration-тест сценария `doctor` с API mock, проверяющего headers;
- [x] подтвердить, что существующий runtime `run` не регрессировал.

## tests

### unit (U*)

- [x] U1: `doctor` отправляет `Authorization` и `partner-id` в `GET /v1/tasks/anonymization_fields`.
- [x] U2: `doctor` возвращает ошибку валидации при не-UUID `api.partner_id`.
- [x] U3: `doctor` возвращает явную ошибку при отсутствии секрета в OS secret store.

### integration (I*)

- [x] I1: mock API принимает preflight-запрос только при корректных заголовках, `doctor` проходит.
- [x] I2: при невалидном `partner_id` сетевой вызов не выполняется, `doctor` падает на валидации.
- [x] I3: при пустом/отсутствующем секрете preflight API-check не выполняется и возвращается deterministic error.

### regression

- [x] R1: `go test ./...` зеленый после изменений.
