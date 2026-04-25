# user guide

## для кого это руководство

Этот документ для пользователя AI IDE, который хочет запустить `anonym` и получить рабочий результат с первого раза.

Поддержаны два сценария:
1. у вас уже есть готовый `anonym.exe`;
2. вы сделали `git clone` репозитория и запускаете через Go.

## что делает инструмент

`anonym` — это CLI для безопасной анонимизации файлов:
1. `scan` — находит файлы в `raw_path` и кладет их в локальную очередь;
2. `list` — показывает очередь (`id`, путь, статус);
3. `run <id|path>` — отправляет выбранный файл в API и сохраняет результат в `output_path`;
4. `doctor` — проверяет, что окружение настроено корректно.

## обязательная модель папок

Перед первым запуском создайте три папки:
1. `raw_path` — папка с исходными файлами, обязательно вне IDE workspace;
2. `workspace_path` — рабочая папка вашего проекта в AI IDE;
3. `output_path` — папка для результатов, обязательно внутри `workspace_path`.

Пример:
1. `D:\secure-raw` (`raw_path`)
2. `C:\work\my-ai-project` (`workspace_path`)
3. `C:\work\my-ai-project\anonymized` (`output_path`)

## обязательные параметры и секреты

### параметры (минимум для старта)

Обязательные параметры для `v1`:
1. `--raw-path` или `ANON_RAW_PATH`
2. `--output-path` или `ANON_OUTPUT_PATH`
3. `--workspace-path` или `ANON_WORKSPACE_PATH`
4. `--api-base-url` или `ANON_API_BASE_URL` (только `https`)
5. `--api-partner-id` или `ANON_API_PARTNER_ID`
6. `--allowed-hosts` или `ANON_ALLOWED_HOSTS`
7. `--max-parallel-runs 1` (или `ANON_MAX_PARALLEL_RUNS=1`)

Рекомендуемые лимиты:
1. `--max-file-size-mb 25`
2. `--max-pages 300`
3. `--request-timeout-sec 5`

### секреты (обязательно)

Секреты должны быть в OS secret store (не в `yaml`, не в `.env`):
1. `service`: `ai-anonymizer/api`
2. `account`: значение `api.partner_id` (например, `partner-1`)  
   `value`: API secret/token
3. `account`: значение `audit.hmac_key_id` (обычно `audit-hmac-v1`)  
   `value`: секрет для HMAC аудита

Проверка, что всё настроено:
1. запустите `doctor`;
2. если секретов нет, `doctor` вернет понятную ошибку.

## вариант 1: у вас уже есть `anonym.exe`

### что должно быть под рукой

1. файл `anonym.exe` (например, `C:\tools\anonym\anonym.exe`);
2. подготовленные папки `raw_path/workspace_path/output_path`;
3. настроенные секреты в OS secret store.

### первый запуск

```powershell
C:\tools\anonym\anonym.exe doctor `
  --raw-path D:\secure-raw `
  --output-path C:\work\my-ai-project\anonymized `
  --workspace-path C:\work\my-ai-project `
  --api-base-url https://api.company.local `
  --api-partner-id partner-1 `
  --allowed-hosts api.company.local `
  --max-file-size-mb 25 `
  --max-pages 300 `
  --request-timeout-sec 5 `
  --max-parallel-runs 1
```

Если в ответе `doctor: ok`, можно работать.

### рабочий цикл

```powershell
C:\tools\anonym\anonym.exe scan --raw-path D:\secure-raw --output-path C:\work\my-ai-project\anonymized --workspace-path C:\work\my-ai-project --api-base-url https://api.company.local --api-partner-id partner-1 --allowed-hosts api.company.local --max-file-size-mb 25 --max-pages 300 --request-timeout-sec 5 --max-parallel-runs 1

C:\tools\anonym\anonym.exe list --raw-path D:\secure-raw --output-path C:\work\my-ai-project\anonymized --workspace-path C:\work\my-ai-project --api-base-url https://api.company.local --api-partner-id partner-1 --allowed-hosts api.company.local --max-file-size-mb 25 --max-pages 300 --request-timeout-sec 5 --max-parallel-runs 1

C:\tools\anonym\anonym.exe run <id> --raw-path D:\secure-raw --output-path C:\work\my-ai-project\anonymized --workspace-path C:\work\my-ai-project --api-base-url https://api.company.local --api-partner-id partner-1 --allowed-hosts api.company.local --max-file-size-mb 25 --max-pages 300 --request-timeout-sec 5 --max-parallel-runs 1
```

## вариант 2: вы сделали `git clone`

### какие папки нужны в репозитории

Для запуска через `go run` должны быть:
1. `cmd/anonym/` (точка входа CLI)
2. `internal/` (runtime-модули)
3. `go.mod`
4. `go.sum`
5. `config.example.yaml` (шаблон параметров)

### запуск из исходников

Из корня репозитория:

```powershell
go run ./cmd/anonym doctor `
  --raw-path D:\secure-raw `
  --output-path C:\work\my-ai-project\anonymized `
  --workspace-path C:\work\my-ai-project `
  --api-base-url https://api.company.local `
  --api-partner-id partner-1 `
  --allowed-hosts api.company.local `
  --max-file-size-mb 25 `
  --max-pages 300 `
  --request-timeout-sec 5 `
  --max-parallel-runs 1
```

Дальше используйте те же команды `scan`, `list`, `run`, заменив `anonym.exe` на `go run ./cmd/anonym`.

## как задавать параметры без ошибок

### безопасный способ для первого запуска

Используйте CLI-флаги в каждой команде. Так проще диагностировать ошибки.

### через переменные окружения (когда команды запускаются часто)

Пример для PowerShell:

```powershell
$env:ANON_RAW_PATH="D:\secure-raw"
$env:ANON_OUTPUT_PATH="C:\work\my-ai-project\anonymized"
$env:ANON_WORKSPACE_PATH="C:\work\my-ai-project"
$env:ANON_API_BASE_URL="https://api.company.local"
$env:ANON_API_PARTNER_ID="partner-1"
$env:ANON_ALLOWED_HOSTS="api.company.local"
$env:ANON_MAX_FILE_SIZE_MB="25"
$env:ANON_MAX_PAGES="300"
$env:ANON_REQUEST_TIMEOUT_SEC="5"
$env:ANON_MAX_PARALLEL_RUNS="1"
$env:ANON_AUDIT_RETENTION_DAYS="30"
$env:ANON_AUDIT_HMAC_KEY_ID="audit-hmac-v1"
```

После этого можно запускать короче:

```powershell
go run ./cmd/anonym doctor
go run ./cmd/anonym scan
go run ./cmd/anonym list
go run ./cmd/anonym run <id>
```

## важные замечания по v1

1. `run` в `v1` не поддерживает CLI override-флаги `--fields`, `--tags`, `--tags-numeration`.
2. Секреты из `config.yaml` и plaintext env не используются как fallback.
3. Аудит пишется в `<workspace>/.anonym/audit.log` в формате JSONL:
   - `file_id` хранится только как HMAC;
   - исходный raw-файл и секреты в лог не пишутся.
4. `config.example.yaml` используйте как шаблон и чек-лист параметров; для надежного запуска в `v1` задавайте параметры через `flags/env`.

## быстрый чек-лист перед началом

1. `raw_path` вне workspace, `output_path` внутри workspace.
2. `api.base_url` начинается с `https://`.
3. Host API входит в `allowed_hosts`.
4. В OS secret store есть два секрета (`api.partner_id` и `audit.hmac_key_id`).
5. `max_parallel_runs = 1`.
6. `doctor` возвращает `doctor: ok`.

## известные ограничения v1

1. Подсчет страниц PDF в `run` эвристический (по структуре PDF), в редких сложных файлах возможна неточность.
2. Для DOCX основной источник страниц — `docProps/app.xml`; при отсутствии метаданных применяется fallback по page-break в `word/document.xml`.
