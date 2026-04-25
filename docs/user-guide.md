# user guide

## Назначение

Руководство по работе с CLI `anonym` в ежедневном режиме.

## Термины

- `raw_path`: абсолютный путь к исходным файлам вне workspace IDE.
- `output_path`: папка в workspace для анонимизированных результатов.
- `id`: идентификатор файла из `anonym list`.

## Рекомендуемый поток [v1]

1. `anonym scan`
2. `anonym list`
3. `anonym run <id>`

Результат:
- в `output_path` появляется анонимизированный файл с безопасным техническим именем;
- исходный файл в `raw_path` не удаляется и не перемещается.

## Команды [v1]

- `anonym scan`
- `anonym list`
- `anonym run <id|path>`
- `anonym doctor`

Примечание:
- в `v1` команда `run` не поддерживает CLI override-флаги `--fields`, `--tags`, `--tags-numeration`; используются значения из конфигурации/окружения.

### `anonym doctor` (старт `v1-i1`)

Минимальный preflight в текущей итерации проверяет:
- `raw_path` находится вне `workspace_path`;
- `raw_path` существует, доступен и является директорией;
- `output_path` существует, доступен и является директорией;
- `output_path` проверяется на запись через write-probe;
- `api.base_url` использует `https`;
- host из `api.base_url` входит в `allowed_hosts`;
- наличие секрета для `api_partner_id` в OS secret store (service `ai-anonymizer/api`, account = `api_partner_id`);
- доступность `GET /v1/tasks/anonymization_fields` с timeout;
- `max_parallel_runs=1` для `v1`.

Пример запуска:
- `anonym doctor --raw-path D:/secure-raw --output-path C:/work/project/anonymized --workspace-path C:/work/project --api-base-url https://api.company.local --api-partner-id partner-1 --allowed-hosts api.company.local --request-timeout-sec 5 --max-parallel-runs 1`

### `anonym scan` и `anonym list` (промежуточный результат `v1-i1`)

- `scan` выполняет только локальное безопасное обнаружение в `raw_path` и сохраняет очередь в `<workspace>/.anonym/catalog.json`;
- `scan` применяет лимит `max_file_size_mb` и не индексирует oversized-файлы;
- `list` выводит элементы очереди в формате `id path status`;
- `scan/list` не отправляют файлы в API.

## Настройка конфигурации [v1]

1. Скопировать шаблон `config.example.yaml` в локальный рабочий конфиг.
2. Задать `raw_path` вне workspace IDE.
3. Задать `output_path` внутри workspace IDE.
4. Заполнить `api.base_url`, `api.partner_id`, `api.fields`.
5. Проверить `security.allowed_hosts` и лимиты.
6. Выполнить `anonym doctor`.

Важно:
- API-ключ не хранится в конфиг-файле и не передается в plaintext-переменных.
- Секреты читаются только из OS secret store.
- При конфликте источников приоритет: flags > env > config file.

## Параметры приложения [v1]

Ниже — полный каталог параметров из конфигурационного контракта.

| Параметр | Назначение | Где задается | Кто меняет |
|---|---|---|---|
| `paths.raw_path` | путь к исходным файлам вне workspace | `--raw-path`, `ANON_RAW_PATH`, `config.yaml` | пользователь/администратор рабочего места |
| `paths.output_path` | путь для анонимизированных файлов в workspace | `--output-path`, `ANON_OUTPUT_PATH`, `config.yaml` | пользователь/администратор рабочего места |
| `api.base_url` | базовый URL API анонимизации | `--api-base-url`, `ANON_API_BASE_URL`, `config.yaml` | администратор/техлид |
| `api.partner_id` | идентификатор партнера для вызовов API | `--api-partner-id`, `ANON_API_PARTNER_ID`, `config.yaml` | администратор/техлид |
| `api.fields` | профиль/поля анонимизации | `ANON_API_FIELDS`, `config.yaml` | администратор/техлид |
| `api.tags` | дополнительные теги запроса | `ANON_API_TAGS`, `config.yaml` | администратор/техлид |
| `api.tags_numeration` | режим нумерации тегов (`0/1`) | `ANON_API_TAGS_NUMERATION`, `config.yaml` | администратор/техлид |
| `limits.max_file_size_mb` | лимит размера входного файла | `--max-file-size-mb`, `ANON_MAX_FILE_SIZE_MB`, `config.yaml` | администратор/техлид |
| `limits.max_pages` | лимит страниц для поддерживаемых форматов | `ANON_MAX_PAGES`, `config.yaml` | администратор/техлид |
| `limits.request_timeout_sec` | timeout API-запроса | `--request-timeout-sec`, `ANON_REQUEST_TIMEOUT_SEC`, `config.yaml` | администратор/техлид |
| `limits.max_parallel_runs` | лимит параллельных запусков (для `v1` = `1`) | `--max-parallel-runs`, `ANON_MAX_PARALLEL_RUNS`, `config.yaml` | администратор/техлид |
| `security.allowed_hosts` | allowlist допустимых API host | `--allowed-hosts`, `ANON_ALLOWED_HOSTS`, `config.yaml` | администратор/техлид |
| `security.require_tls_verify` | обязательная проверка TLS-сертификата | `ANON_REQUIRE_TLS_VERIFY`, `config.yaml` | администратор/техлид |
| `audit.retention_days` | срок хранения аудита | `ANON_AUDIT_RETENTION_DAYS`, `config.yaml` | администратор/техлид |
| `audit.hmac_key_id` | ключ HMAC для `file-id` в аудит-логах | `ANON_AUDIT_HMAC_KEY_ID`, `config.yaml` | администратор/техлид |
| `catalog.db_path` | путь хранения локального каталога | `ANON_CATALOG_DB_PATH`, `config.yaml` | администратор/техлид |
| `logging.level` | уровень логирования | `ANON_LOG_LEVEL`, `config.yaml` | администратор/техлид |
| `logging.format` | формат логов | `ANON_LOG_FORMAT`, `config.yaml` | администратор/техлид |

## Статус параметров в текущей реализации v1

- Уже применяются в runtime: `raw_path`, `output_path`, `api.base_url`, `api.partner_id`, `security.allowed_hosts`, `limits.max_file_size_mb`, `limits.max_pages`, `limits.request_timeout_sec`, `limits.max_parallel_runs`, `audit.retention_days`, `audit.hmac_key_id`.
- В `run` дополнительно применяются проверки path hardening для Windows (`ADS` и reparse-point deny).
- Параметры целевого конфигурационного контракта, которые должны быть синхронизированы с runtime отдельными задачами: `security.require_tls_verify`, `catalog.db_path`, `logging.*`, `api.tags*`.

## Аудит [v1]

- журнал хранится в `<workspace>/.anonym/audit.log` (формат JSONL);
- `file_id` в журнале хранится только в виде HMAC-идентификатора;
- raw-путь, исходное имя файла и секреты в журнал не записываются;
- при каждой новой записи применяется retention по `audit.retention_days`.

## Варианты запуска [v1+v2]

- По `id` из списка.
- По точному относительному пути внутри `raw_path`.
- По примерному имени (с обработкой неоднозначности).

## Требования к сообщениям CLI [v1+v2]

- четкий итог операции;
- `id` и путь результата;
- код и причина ошибки без утечки PII;
- понятный следующий шаг.

## Что добавляется в v2

- MCP-обертка над тем же trusted-core CLI.
- Операционные регламенты секретов и эксплуатации.

## Известные ограничения v1 (risks)

- Текущий page-counter для PDF в `run` реализован эвристически (по структуре PDF), поэтому для отдельных сложных PDF возможна неточная оценка числа страниц.
- Для `DOCX` первичный источник количества страниц — `docProps/app.xml`; при отсутствии метаданных используется fallback по page-break в `word/document.xml`.
