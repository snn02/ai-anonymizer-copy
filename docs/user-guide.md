# Руководство Пользователя

## Для Кого Это Руководство

Документ для пользователя AI IDE, который хочет запустить `anonym` и получить результат с первого раза.

Поддержаны два варианта запуска:
1. у вас есть готовый `anonym.exe`;
2. вы запускаете из исходников через `go run`.

## Что Делает Инструмент

`anonym` — CLI для безопасной анонимизации файлов:
1. `doctor` — проверяет конфигурацию, доступ к API и секреты;
2. `scan` — индексирует файлы из `raw_path` в локальный каталог;
3. `list` — показывает очередь (`id`, путь, статус);
4. `run <id|path>` — отправляет выбранный файл в API и сохраняет результат в `output_path`.

## Обязательная Модель Папок

Перед первым запуском подготовьте:
1. `raw_path` — папка с исходными файлами, обязательно вне IDE workspace;
2. `workspace_path` — корень проекта в AI IDE;
3. `output_path` — папка результатов, обязательно внутри `workspace_path`.

Пример:
1. `D:\secure-raw` (`raw_path`)
2. `C:\work\my-ai-project` (`workspace_path`)
3. `C:\work\my-ai-project\anonymized` (`output_path`)

## Приоритет Источников Конфигурации

Runtime применяет значения в порядке:
1. CLI-флаги (самый высокий приоритет);
2. Переменные окружения;
3. `config.yaml`.

Практика для `v1`:
1. основной сценарий — заполнить `config.yaml` и запускать компактными командами;
2. env/CLI использовать как точечный override без копирования всех параметров.

## Обязательные Параметры И Секреты

### Параметры Runtime V1

Обязательные поля:
1. `paths.raw_path`
2. `paths.output_path`
3. `paths.workspace_path`
4. `api.base_url` (только `https`)
5. `api.partner_id` (строго UUID)
6. `security.allowed_hosts`
7. `limits.max_parallel_runs=1`

Рекомендуемые лимиты:
1. `limits.max_file_size_mb=25`
2. `limits.max_pages=300`
3. `limits.request_timeout_sec=5`

### Секреты

Секреты хранятся только в OS secret store.
Нельзя хранить секреты в `config.yaml` и `.env`.

### Как Сопоставить Данные От Разработчика API

Если вам передали:
1. `ключ` — это значение заголовка `Authorization` (кладется в keyring);
2. `партнер` — это `partner-id` и `api.partner_id` (UUID);
3. `пользователь` — это `user-id` (опциональный UUID-заголовок, не используется runtime `v1`).

### Что Важно Для Runtime V1

1. CLI отправляет `Authorization`, `partner-id`, `fields=anonymizer`.
2. `user-id`, `tags`, `tags-numeration` не параметризуются через CLI/env/config.

## Настройка Секретов В OS Secret Store

Общий шаблон:
1. `service`: `ai-anonymizer/api`
2. `account=<api.partner_id>`, `value=<API_TOKEN>`
3. `account=<audit.hmac_key_id>`, `value=<HMAC_SECRET>`

### Windows (Credential Manager)

```powershell
cmdkey /generic:"ai-anonymizer/api:<partner_uuid>" /user:"<partner_uuid>" /pass:"<API_TOKEN>"
cmdkey /generic:"ai-anonymizer/api:audit-hmac-v1" /user:"audit-hmac-v1" /pass:"<HMAC_SECRET>"
```

Проверка:

```powershell
cmdkey /list | findstr "ai-anonymizer/api"
```

### Linux (Secret Service / Gnome Keyring)

```bash
echo -n "<API_TOKEN>" | secret-tool store --label="ai-anonymizer api partner" service "ai-anonymizer/api" username "<partner_uuid>"
echo -n "<HMAC_SECRET>" | secret-tool store --label="ai-anonymizer audit hmac" service "ai-anonymizer/api" username "audit-hmac-v1"
```

### Macos (Keychain)

```bash
security add-generic-password -U -s "ai-anonymizer/api" -a "<partner_uuid>" -w "<API_TOKEN>"
security add-generic-password -U -s "ai-anonymizer/api" -a "audit-hmac-v1" -w "<HMAC_SECRET>"
```

## Приоритетный Сценарий: Запуск Через Config

### Шаг 1. Подготовьте Config Файл

Создайте `config.yaml` по образцу `config.example.yaml` и заполните минимум:

```yaml
paths:
  raw_path: "D:/secure-raw"
  workspace_path: "C:/work/my-ai-project"
  output_path: "C:/work/my-ai-project/anonymized"
api:
  base_url: "https://production-retrievals.ai.rarus-cloud.ru"
  partner_id: "00000000-0000-0000-0000-000000000000"
limits:
  max_file_size_mb: 25
  max_pages: 300
  request_timeout_sec: 5
  max_parallel_runs: 1
security:
  allowed_hosts:
    - "production-retrievals.ai.rarus-cloud.ru"
audit:
  retention_days: 30
  hmac_key_id: "audit-hmac-v1"
```

### Шаг 2. Запуск С Компактными Командами

Если у вас `anonym.exe`:

```powershell
C:\tools\anonym\anonym.exe doctor --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe scan --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe list --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe run <id> --config C:\work\my-ai-project\config.yaml
```

Если запуск через исходники:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym scan --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym list --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym run <id> --config C:\work\my-ai-project\config.yaml
```

## Override Параметров: Когда И Как Делать

### Кейс 1. Override Через Env (Для Серии Запусков)

Когда использовать:
1. в этой сессии нужно временно сменить 1-2 параметра;
2. команды запускаются много раз, и не хочется повторять флаг.

Пример:

```powershell
$env:ANON_REQUEST_TIMEOUT_SEC="15"
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym run <id> --config C:\work\my-ai-project\config.yaml
```

### Кейс 2. Override Через CLI (Разовый Запуск)

Когда использовать:
1. разово нужно проверить другой хост/лимит;
2. важно явно зафиксировать override в истории команды.

Пример:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml --api-base-url https://pilot-retrievals.company.local --allowed-hosts pilot-retrievals.company.local
```

## Диагностика Config Режима

1. Если `--config` указан явно и файл не найден, команда завершится ошибкой загрузки файла.
2. Если файл не указан явно и дефолтный путь отсутствует, runtime продолжит с env/CLI.
3. При битом YAML команда завершится явной ошибкой парсинга конфига.

## Важные Ограничения V1

1. `fields` фиксирован в `anonymizer`.
2. `user-id` не отправляется из CLI runtime.
3. `max_parallel_runs` в `v1` должен быть `1`.
4. Секреты не читаются из plaintext-config/env fallback.
5. Audit пишется в `<workspace>/.anonym/audit.log` (JSONL), `file_id` хранится как HMAC.

## Быстрый Чек-Лист Перед Запуском

1. `raw_path` вне workspace, `output_path` внутри workspace.
2. `api.base_url` начинается с `https://`.
3. Host API присутствует в `allowed_hosts`.
4. В OS secret store есть записи для `api.partner_id` и `audit.hmac_key_id`.
5. `max_parallel_runs = 1`.
6. `doctor` возвращает `doctor: ok`.
