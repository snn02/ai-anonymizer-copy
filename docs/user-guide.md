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

Для стандартного профиля `v1` секреты хранятся только в OS secret store.
Для временного профиля `v1-1-ws` (MVP, non-production) допустимы временные секреты в `config/env/CLI` при явном `insecure_no_secrets=true`.

### Как Сопоставить Данные От Разработчика API

Если вам передали:
1. `ключ` — это значение заголовка `Authorization` (кладется в keyring);
2. `партнер` — это `partner-id` и `api.partner_id` (UUID);
3. `пользователь` — это `user-id` (UUID-заголовок для `doctor` preflight, опционален).

### Что Важно Для Runtime V1

1. CLI отправляет `Authorization`, `partner-id`, `fields=anonymizer`.
2. Для `run`: `user-id`, `tags`, `tags-numeration` не параметризуются через CLI/env/config.
3. Для `doctor`: поддерживается `api.user_id` (`--api-user-id` / `ANON_API_USER_ID` / config).

## Настройка Секретов В OS Secret Store

### Что Является Секретом, А Что Нет

1. `api.partner_id` в `config` — это идентификатор (UUID), не секрет.
2. `audit.hmac_key_id` в `config` — это идентификатор записи HMAC, не секрет.
3. Секреты — это только значения:
   - `API_TOKEN` (для заголовка `Authorization`);
   - `HMAC_SECRET` (для подписи audit).

### Общий Шаблон Хранения

1. `service`: `ai-anonymizer/api`
2. API-token:
   - `account(user) = <api.partner_id из config>`
   - `value = <API_TOKEN>`
3. Audit HMAC:
   - `account(user) = <audit.hmac_key_id из config>`
   - `value = <HMAC_SECRET>`

Важно:
1. `<partner_uuid>` в командах ниже должен быть тем же значением, что и `api.partner_id` в `config`.
2. Если в `config` вы поменяли `audit.hmac_key_id`, используйте это же значение в секрете (не обязательно `audit-hmac-v1`).

### Генерация HMAC_SECRET (PowerShell)

Сгенерируйте случайное значение `HMAC_SECRET` (Base64, 32 байта):

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

Используйте выведенную строку как `<HMAC_SECRET>` в командах ниже.

### Windows (Credential Manager)

```powershell
cmdkey /generic:"ai-anonymizer/api:<partner_uuid>" /user:"<partner_uuid>" /pass:"<API_TOKEN>"
cmdkey /generic:"ai-anonymizer/api:<audit_hmac_key_id>" /user:"<audit_hmac_key_id>" /pass:"<HMAC_SECRET>"
```

Проверка:

```powershell
cmdkey /list | findstr "ai-anonymizer/api"
```

### Linux (Secret Service / Gnome Keyring)

```bash
echo -n "<API_TOKEN>" | secret-tool store --label="ai-anonymizer api partner" service "ai-anonymizer/api" username "<partner_uuid>"
echo -n "<HMAC_SECRET>" | secret-tool store --label="ai-anonymizer audit hmac" service "ai-anonymizer/api" username "<audit_hmac_key_id>"
```

### Macos (Keychain)

```bash
security add-generic-password -U -s "ai-anonymizer/api" -a "<partner_uuid>" -w "<API_TOKEN>"
security add-generic-password -U -s "ai-anonymizer/api" -a "<audit_hmac_key_id>" -w "<HMAC_SECRET>"
```

## Приоритетный Сценарий: Запуск Через Config

### Шаг 1. Подготовьте Config Файл

Рекомендуемый вариант:
1. возьмите `config.example.yaml` в корне проекта;
2. заполните значения под ваше окружение;
3. оставьте файл с именем `config.example.yaml`, если хотите запуск без `--config`.

Если хотите отдельный файл (`config.yaml`, `config.prod.yaml` и т.п.), это нормально, но тогда в командах указывайте `--config <путь>`.

Минимальный состав полей:

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

Важно для Windows-путей в YAML:
1. не используйте `\` внутри двойных кавычек (например, `"D:\secure-raw"`), это может вызвать ошибку `unknown escape character`;
2. используйте один из безопасных вариантов:
   - `"D:/secure-raw"` (рекомендуется);
   - `'D:\secure-raw'`;
   - `"D:\\secure-raw"`.

### Шаг 2. Запуск С Компактными Командами

Вариант A: файл в дефолтном месте и с дефолтным именем `./config.example.yaml` (можно без `--config`).

Если у вас `anonym.exe`:

```powershell
C:\tools\anonym\anonym.exe doctor
C:\tools\anonym\anonym.exe scan
C:\tools\anonym\anonym.exe list
C:\tools\anonym\anonym.exe run <id|path>
```

Если запуск через исходники:

```powershell
go run ./cmd/anonym doctor
go run ./cmd/anonym scan
go run ./cmd/anonym list
go run ./cmd/anonym run <id|path>
```

Вариант B: отдельный config-файл (нужен `--config`).

Если у вас `anonym.exe`:

```powershell
C:\tools\anonym\anonym.exe doctor --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe scan --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe list --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe run <id|path> --config C:\work\my-ai-project\config.yaml
```

Если запуск через исходники:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym scan --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym list --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym run <id|path> --config C:\work\my-ai-project\config.yaml
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
2. Если `--config` не указан, runtime ищет `./config.example.yaml` в текущей рабочей папке.
3. Если файл из п.2 отсутствует, runtime продолжит с env/CLI.
4. При битом YAML команда завершится явной ошибкой парсинга конфига.

## Важные Ограничения V1

1. `fields` фиксирован в `anonymizer`.
2. `user-id` параметризуется только для `doctor` preflight и не используется в `run`.
3. `max_parallel_runs` в `v1` должен быть `1`.
4. Секреты не читаются из plaintext-config/env fallback.
5. Audit пишется в `<workspace>/.anonym/audit.log` (JSONL), `file_id` хранится как HMAC.

## Временный Профиль V1-1-WS (MVP Для Изолированных IDE)

Когда использовать:
1. вы запускаете из изолированной среды AI IDE;
2. runtime не имеет доступа к OS secret store хоста.

Что меняется:
1. включается явный insecure-профиль `v1-1-ws`;
2. для API auth используется временный `api.auth_token` из env/config;
3. для audit используется временный `audit.hmac_secret` из env/config;
4. профиль допускается только для non-production контуров.

Что не меняется:
1. boundary и path hardening;
2. HTTPS + allowlist;
3. лимиты и audit.

Важно:
1. `v1-1-ws` — временный профиль MVP, не дефолтный режим;
2. для production используйте стандартный `v1` с OS secret store;
3. при `v1-1-ws` обязателен операционный контроль: rate-limit, быстрый revoke/rotation токена и HMAC-секрета.

## Быстрый Чек-Лист Перед Запуском

1. `raw_path` вне workspace, `output_path` внутри workspace.
2. `api.base_url` начинается с `https://`.
3. Host API присутствует в `allowed_hosts`.
4. Для `v1`: в OS secret store есть записи для `api.partner_id` и `audit.hmac_key_id`.
5. Для `v1-1-ws`: заданы `runtime.profile=v1-1-ws`, `security_flags.insecure_no_secrets=true`, `api.auth_token`, `audit.hmac_secret`.
6. `max_parallel_runs = 1`.
7. `doctor` возвращает `doctor: ok`.

## Уточнение По Секретам И Профилям (Актуально На 2026-04-29)

Этот раздел является приоритетным для интерпретации правил секретов.

1. Стандартный профиль `v1`:
- секреты берутся только из OS secret store;
- использование `api.auth_token` и `audit.hmac_secret` из `env/config` запрещено.

2. Временный профиль `v1-1-ws` (только non-production):
- секреты допускаются из `env/config` только при явном включении insecure-режима;
- обязательны оба значения: `api.auth_token` и `audit.hmac_secret`.

### Как Явно Включить `v1-1-ws`

Через CLI-флаги:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"

go run ./cmd/anonym scan --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"

go run ./cmd/anonym list --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"

go run ./cmd/anonym run <id|path> --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"
```

Через env-переменные:

```powershell
$env:ANON_RUNTIME_PROFILE="v1-1-ws"
$env:ANON_INSECURE_NO_SECRETS="true"
$env:ANON_API_AUTH_TOKEN="<API_TOKEN>"
$env:ANON_AUDIT_HMAC_SECRET="<HMAC_SECRET>"
```

Важно:
- при `v1-1-ws` CLI выводит предупреждение о временном insecure-профиле;
- для production-контура используйте только стандартный `v1`.

### Чек-Лист По Профилям

1. Для `v1`: проверьте записи в OS secret store (`api.partner_id`, `audit.hmac_key_id`).
2. Для `v1-1-ws`: вместо OS secret store задайте `api.auth_token` и `audit.hmac_secret` (CLI/env/config) и включите `insecure_no_secrets=true`.

## Диагностика HTTP-Вызовов (Опционально)

Для диагностики ошибок авторизации/доступности API можно включить debug-лог HTTP.

Как включить:

```powershell
$env:ANON_HTTP_DEBUG="true"
```

Куда пишется лог:
- `<workspace>/.anonym/http-debug.log`

Что пишется:
- время, команда (`doctor`/`run`), метод, URL, HTTP-статус;
- заголовки в безопасном виде (значение `Authorization` маскируется).

Важно:
- по умолчанию debug-лог выключен;
- при выключенном флаге файл лога не создается.


## ��������� ��������� Doctor (��������� �� 2026-04-30)

Preflight-�������� ������� `doctor` ����������:
- `GET /v1/tasks/anonymization_fields`
- headers: `Authorization`, `partner-id`, `user-id`
- query: `page`, `per_page`

����� runtime-��������� ��� `doctor`:
1. `api.user_id` (UUID)
2. `api.fields_page` (int, default `1`)
3. `api.fields_per_page` (int, default `10`)

CLI-�����:
- `--api-user-id`
- `--api-fields-page`
- `--api-fields-per-page`
