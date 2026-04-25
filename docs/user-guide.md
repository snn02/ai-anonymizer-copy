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
- `anonym run <id|path> [--fields <value>] [--tags <csv>] [--tags-numeration 0|1]`
- `anonym doctor`

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
