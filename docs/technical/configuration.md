# configuration

## Назначение

Единая конфигурация trusted CLI для `scan/list/run/doctor` без хранения секретов в проектных файлах.

## Принципы

1. Секреты не хранятся в `yaml/json/env` и читаются только из OS secret store.
2. Конфиг валидируется на старте и повторно проверяется в `anonym doctor`.
3. Небезопасные значения блокируют рабочие команды.

## Приоритет источников

1. CLI-флаги.
2. Переменные окружения.
3. `config.yaml`.

Если есть конфликт, применяется источник с более высоким приоритетом.

## Загрузка Config Файла

1. Runtime читает путь, переданный в `--config`.
2. Если `--config` указан явно и файл недоступен, команда завершается ошибкой загрузки конфига.
3. Если `--config` не указан явно, используется дефолтный путь команды; при отсутствии файла runtime продолжает работу с env/CLI.
4. При невалидном YAML runtime возвращает явную ошибку парсинга.

## Таблица параметров runtime v1

| Ключ | Тип | Обяз. | Default | Пример | Env |
|---|---|---|---|---|---|
| `paths.raw_path` | string (abs path) | да | - | `D:/secure-raw` | `ANON_RAW_PATH` |
| `paths.output_path` | string (abs path) | да | - | `C:/work/project/anonymized` | `ANON_OUTPUT_PATH` |
| `paths.workspace_path` | string (abs path) | да | - | `C:/work/project` | `ANON_WORKSPACE_PATH` |
| `api.base_url` | string (https url) | да | `https://production-retrievals.ai.rarus-cloud.ru` | `https://production-retrievals.ai.rarus-cloud.ru` | `ANON_API_BASE_URL` |
| `api.partner_id` | string (uuid) | да | - | `00000000-0000-0000-0000-000000000000` | `ANON_API_PARTNER_ID` |
| `limits.max_file_size_mb` | int | да | `25` | `50` | `ANON_MAX_FILE_SIZE_MB` |
| `limits.max_pages` | int | да | `300` | `500` | `ANON_MAX_PAGES` |
| `limits.request_timeout_sec` | int | да | `5` | `15` | `ANON_REQUEST_TIMEOUT_SEC` |
| `limits.max_parallel_runs` | int | да | `1` | `1` | `ANON_MAX_PARALLEL_RUNS` |
| `security.allowed_hosts` | string list | да | `["production-retrievals.ai.rarus-cloud.ru"]` | `["production-retrievals.ai.rarus-cloud.ru"]` | `ANON_ALLOWED_HOSTS` |
| `security.require_tls_verify` | bool | да | `true` | `true` | `ANON_REQUIRE_TLS_VERIFY` |
| `audit.retention_days` | int | да | `30` | `90` | `ANON_AUDIT_RETENTION_DAYS` |
| `audit.hmac_key_id` | string | да | `audit-hmac-v1` | `audit-hmac-v1` | `ANON_AUDIT_HMAC_KEY_ID` |
| `audit.path` | string (abs path) | нет | `<workspace>/.anonym/audit.log` | `C:/work/project/.anonym/audit.log` | - |

`catalog.path` в runtime `v1` не задается через `config/env/CLI` и вычисляется автоматически как `<workspace>/.anonym/catalog.json`.

## Контракт вызова API для `run` (v1)

Обязательные заголовки при `POST /v1/tasks/file_anonymization`:
1. `Authorization` (значение берется из OS secret store по `api.partner_id`).
2. `partner-id` (`api.partner_id`, строго UUID).
3. `fields` (в runtime v1 фиксировано значение `anonymizer`).

Опциональные заголовки:
1. `tags-numeration` (`0|1`).
2. `user-id` (UUID).

Важно для v1 runtime:
1. `tags`, `tags-numeration` и `user-id` не настраиваются через CLI/config/env.
2. `fields` не настраивается через CLI/config/env и отправляется как фиксированное `anonymizer`.

## Валидация и security-ограничения

1. `paths.raw_path` должен быть вне workspace IDE.
2. `paths.output_path` должен быть внутри workspace IDE.
3. `api.base_url` только `https`.
4. `security.allowed_hosts` не может быть пустым в production-профиле.
5. `limits.max_parallel_runs` в `v1` фиксируется в `1`.
6. Секреты для API (например, `Authorization` ключ) читаются только из keyring по `partner_id`/профилю.
7. Ключ для audit HMAC читается только из OS secret store по `audit.hmac_key_id`; plaintext fallback запрещен.
8. Audit-журнал хранится в `<workspace>/.anonym/audit.log` и применяет retention по `audit.retention_days`.

## Формат переменных окружения

- Для list-полей (`allowed_hosts`) использовать `,`:
  - `ANON_ALLOWED_HOSTS=production-retrievals.ai.rarus-cloud.ru,api.backup.local`

## Профили окружения

Рекомендуемые профили:
1. `dev` - для локальной разработки.
2. `pilot` - для ограниченного теста с коллегами.
3. `prod` - для боевого контура.

Профиль влияет только на дефолты и строгость валидации, но не отключает базовые security-проверки.
