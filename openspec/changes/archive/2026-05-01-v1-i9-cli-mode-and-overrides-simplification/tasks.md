## 1. Mode Migration

- [x] 1.1 Ввести `mode` (`prod|mvp`) в runtime config loader с приоритетом `CLI > env > config > default`.
- [x] 1.2 Обеспечить backward-совместимость на миграционный период: старые поля профиля преобразуются в `mode` с явным предупреждением.

## 2. CLI Simplification

- [x] 2.1 Удалить лишние override-флаги из `doctor/scan/list/run`, оставить только `--config`, `--mode`, `--workspace-path`, `--output-path`.
- [x] 2.2 Обновить help/usage и ошибки для удаленных флагов с понятной подсказкой «перенесено в config/env».

## 3. Runtime Security Behavior

- [x] 3.1 Для `mode=prod` сохранить чтение API/audit-секретов только из OS Secret Store.
- [x] 3.2 Для `mode=mvp` использовать `api.auth_token` и `audit.hmac_secret` только из `config/env` и выводить предупреждение.
- [x] 3.3 Снять блокировку `mvp` для production API host, сохранив остальные preflight/security проверки.

## 4. Tests and Docs Sync (v1)

- [x] 4.1 Добавить/обновить unit+integration тесты на mode-переключение, минимальный CLI-контракт и удаленные флаги.
- [x] 4.2 Обновить `config.example.yaml` и `docs/user-guide.md` под новую модель `mode` и сокращенный CLI.
- [x] 4.3 Синхронизировать `docs/technical/configuration.md`, `docs/test-scenarios.md`, `docs/versions/v1/plan.md`, `docs/versions/v1/action-log.md`, `docs/versions/v1/openspec.md`.
