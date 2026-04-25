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

## Таблица параметров

| Ключ | Тип | Обяз. | Default | Пример | Env |
|---|---|---|---|---|---|
| `paths.raw_path` | string (abs path) | да | - | `D:/secure-raw` | `ANON_RAW_PATH` |
| `paths.output_path` | string (abs path) | да | - | `C:/work/project/anonymized` | `ANON_OUTPUT_PATH` |
| `api.base_url` | string (https url) | да | - | `https://api.company.local` | `ANON_API_BASE_URL` |
| `api.partner_id` | string (uuid) | да | - | `00000000-0000-0000-0000-000000000000` | `ANON_PARTNER_ID` |
| `api.fields` | string | да | `anonymizer` | `anonymizer` | `ANON_API_FIELDS` |
| `api.tags` | string list | нет | `[]` | `["fio","phone"]` | `ANON_API_TAGS` |
| `api.tags_numeration` | int (`0|1`) | нет | `1` | `1` | `ANON_API_TAGS_NUMERATION` |
| `limits.max_file_size_mb` | int | да | `25` | `50` | `ANON_MAX_FILE_SIZE_MB` |
| `limits.max_pages` | int | да | `300` | `500` | `ANON_MAX_PAGES` |
| `limits.request_timeout_sec` | int | да | `120` | `180` | `ANON_REQUEST_TIMEOUT_SEC` |
| `limits.max_parallel_runs` | int | да | `1` | `1` | `ANON_MAX_PARALLEL_RUNS` |
| `security.allowed_hosts` | string list | да | `[]` | `["api.company.local"]` | `ANON_ALLOWED_HOSTS` |
| `security.require_tls_verify` | bool | да | `true` | `true` | `ANON_REQUIRE_TLS_VERIFY` |
| `audit.retention_days` | int | да | `30` | `90` | `ANON_AUDIT_RETENTION_DAYS` |
| `audit.hmac_key_id` | string | да | `default` | `audit-hmac-v1` | `ANON_AUDIT_HMAC_KEY_ID` |
| `catalog.db_path` | string (abs path) | нет | `<workspace>/.anonym/catalog.db` | `C:/work/project/.anonym/catalog.db` | `ANON_CATALOG_DB_PATH` |
| `logging.level` | enum | нет | `info` | `debug` | `ANON_LOG_LEVEL` |
| `logging.format` | enum | нет | `json` | `json` | `ANON_LOG_FORMAT` |

## Валидация и security-ограничения

1. `paths.raw_path` должен быть вне workspace IDE.
2. `paths.output_path` должен быть внутри workspace IDE.
3. `api.base_url` только `https`.
4. `security.allowed_hosts` не может быть пустым в production-профиле.
5. `limits.max_parallel_runs` в `v1` фиксируется в `1`.
6. Секреты для API (например, `Authorization` ключ) читаются только из keyring по `partner_id`/профилю.

## Формат переменных окружения

- Для list-полей (`allowed_hosts`, `tags`) использовать `,`:
  - `ANON_ALLOWED_HOSTS=api.company.local,api.backup.local`
  - `ANON_API_TAGS=fio,phone,email`

## Профили окружения

Рекомендуемые профили:
1. `dev` - для локальной разработки.
2. `pilot` - для ограниченного теста с коллегами.
3. `prod` - для боевого контура.

Профиль влияет только на дефолты и строгость валидации, но не отключает базовые security-проверки.
