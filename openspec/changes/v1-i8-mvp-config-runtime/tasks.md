## 1. Runtime Behavior Alignment

- [x] 1.1 Зафиксировать в коде единый гейт для `v1-1-ws` в `run`: при `insecure_no_secrets=true` использовать секреты из config/env/CLI и не обращаться к OS Secret Store.
- [x] 1.2 Убедиться, что при отсутствии `insecure_no_secrets=true` в `v1-1-ws` команда `run` завершаетcя явной policy-ошибкой.

## 2. Config and Source Priority

- [x] 2.1 Проверить/доработать загрузку `runtime_profile`, `insecure_no_secrets`, `api.auth_token`, `audit.hmac_secret` с приоритетом `CLI > env > config > default`.
- [x] 2.2 Добавить/обновить тесты на приоритет источников для указанных параметров.

## 3. Command-Level Verification

- [x] 3.1 Добавить тест `run`-пути, подтверждающий работу в `v1-1-ws` без OS Secret Store при использовании полного стека источников (`config`, `env`, `CLI`) с приоритетом `CLI > env > config`.
- [x] 3.2 Добавить негативный тест на ошибку при `v1-1-ws` без `insecure_no_secrets=true`.

## 4. Documentation Sync (v1)

- [x] 4.1 Обновить `docs/user-guide.md` с рабочими примерами MVP-запуска через config.
- [x] 4.2 Обновить `docs/technical/configuration.md` по источникам параметров и ограничениям `v1-1-ws`.
- [x] 4.3 Синхронизировать `docs/versions/v1/plan.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md` по результатам изменения.
