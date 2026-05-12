# tasks: v1-1-ws insecure mvp auth

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| IDE-изолированный запуск без доступа к host keyring | T1, T2 | U1, I1 | runtime поддерживает профиль `v1-1-ws` |
| auth и audit secret через env/config допускаются только при явном insecure-профиле | T2, T3 | U2, U3, I2 | `api.auth_token` и `audit.hmac_secret` читаются только в `v1-1-ws` и только при explicit enable |
| стандартный `v1` не ослабляется | T3 | U4, I3 | default-профиль продолжает требовать OS secret store |
| non-production ограничение и предупреждение пользователя | T4 | U5, I4 | runtime блокирует production-контур и выводит warning |

## security dod (для T1-T4)

- [x] runtime реализован в production-пути `doctor/run`;
- [x] минимум один негативный тест на каждый security-check;
- [x] отсутствует тихий insecure fallback;
- [x] evidence зафиксирован прогоном `go test ./...`.

## tasks

### T1. profile contract for v1-1-ws

- [x] добавить параметр профиля runtime (`v1` / `v1-1-ws`);
- [x] добавить явный переключатель insecure-режима;
- [x] сохранить обратную совместимость для текущего `v1`.

### T2. token source in v1-1-ws

- [x] добавить параметр `api.auth_token` (env/config) для `v1-1-ws`;
- [x] добавить параметр `audit.hmac_secret` (env/config) для `v1-1-ws`;
- [x] использовать токен и audit-secret в `doctor` и `run` вместо keyring только в `v1-1-ws`;
- [x] вернуть явную ошибку, если токен не задан.

### T3. no fallback in default v1

- [x] запретить использование `api.auth_token` в стандартном профиле `v1`;
- [x] гарантировать, что `v1` продолжает работать только через OS secret store;
- [x] покрыть негативными тестами кейс «токен задан, профиль не `v1-1-ws`».

### T4. policy guards and warnings

- [x] добавить проверку non-production ограничения для `v1-1-ws`;
- [x] добавить явное предупреждение в CLI при активном `v1-1-ws`;
- [x] синхронизировать `user-guide`, `security`, `configuration`, `test-scenarios`, `action-log`.

## tests

### unit (U*)

- [x] U1: профиль `v1-1-ws` корректно распознается в загрузке конфига.
- [x] U2: `api.auth_token` читается из env/config в `v1-1-ws`.
- [x] U2a: `audit.hmac_secret` читается из env/config в `v1-1-ws`.
- [x] U3: без insecure-включения запуск в `v1-1-ws` блокируется.
- [x] U4: стандартный `v1` игнорирует `api.auth_token` и требует keyring.
- [x] U5: `v1-1-ws` блокируется для production-контура.

### integration (I*)

- [x] I1: `doctor` проходит в изолированном окружении при корректном `v1-1-ws` токене.
- [x] I2: `run` выполняет запрос API в `v1-1-ws` с токеном из env/config.
- [x] I2a: `run` пишет audit в `v1-1-ws` с `audit.hmac_secret` из env/config.
- [x] I3: в стандартном `v1` запуск без keyring продолжает падать.
- [x] I4: CLI выводит предупреждение при активном `v1-1-ws`.

### regression

- [x] R1: `go test ./...` зеленый после изменений.
