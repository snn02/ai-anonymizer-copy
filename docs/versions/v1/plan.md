# v1 plan

## Цель версии

Реализовать безопасный sidecar CLI для анонимизации файлов в контуре ИИ-IDE.

## Общий план разработки v1

1. Зафиксировать обязательные функциональные, security и reliability требования.
2. Выполнять работу итерациями с закончиваемым и тестируемым результатом.
3. На каждую итерацию создавать отдельный feature-пакет в `openspec/changes/`.
4. После каждой итерации выполнять тесты, фиксировать выводы и планировать следующую.
5. После завершения всех итераций выполнить общее тестирование, релиз и обновить `docs/lessons-learned.md`.

## План итераций v1

### Итерация 1 (v1-i1): foundation preflight + safe discovery

Цель итерации:
- получить законченный и тестируемый baseline для безопасного запуска CLI без отправки файлов в API.

Фичи итерации:
1. `anonym doctor` выполняет preflight-проверки конфигурации, boundary/ACL, доступности API и наличия секретов.
2. `anonym scan` обнаруживает файлы только в `raw_path`, применяет path hardening и формирует локальную очередь.
3. `anonym list` показывает локальную очередь с техническими идентификаторами и статусами без утечки PII.

Границы итерации:
- отправка в API (`anonym run`) и обработка результата анонимизации переносятся на следующую итерацию;
- в рамках итерации подтверждается отсутствие неявной отправки в API из `scan/list`.

Артефакты OpenSpec:
- feature-пакет: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/`.

### Итерация 2 (v1-i2): run + api adapter

Цель итерации:
- получить законченный и тестируемый runtime-поток явного запуска анонимизации через `anonym run <id|path>`.

Фичи итерации:
1. `anonym run <id|path>` выполняет явную отправку в API только по команде пользователя.
2. API-ответ `file_anonymization` обрабатывается через adapter внутреннего контракта.
3. Каталог очереди поддерживает переходы `scanned -> sent -> succeeded/failed`.
4. Результат сохраняется в `output_path` с PII-safe именем по technical id.
5. В runtime-пути `run` применяются guardrails: `https`, TLS verify, allowlist, timeout, лимиты, запрет production fallback.

Границы итерации:
- MCP-слой и интеграции `v2` не входят в scope;
- расширенные operational-политики секретов beyond `v1` не входят в scope.

Артефакты OpenSpec:
- feature-пакет: `../../../openspec/changes/v1-i2-run-api-adapter/`.

### Итерация 4 (v1-i4): config-first runtime configuration

Цель итерации:
- обеспечить полноценную загрузку параметров из `config.yaml` как базового сценария запуска, с предсказуемым переопределением через env и CLI.

Фичи итерации:
1. Runtime реально читает `config.yaml` (включая путь через `--config`).
2. Приоритет источников строго соблюдается: `CLI > env > config file`.
3. Все runtime-команды (`doctor/scan/list/run`) поддерживают компактный вызов с минимальным набором флагов при заполненном конфиге.
4. Ошибки чтения/парсинга/валидации конфига возвращаются в явном, диагностируемом виде.

Границы итерации:
- новые функциональные возможности API не добавляются;
- контракт заголовков `run` не меняется (`Authorization`, `partner-id`, `fields=anonymizer`).

Артефакты OpenSpec:
- feature-пакет: `../../../openspec/changes/v1-i4-config-first-loading/`.

### Подверсия (v1-1-ws): workspace-isolated mvp profile

Цель подверсии:
- обеспечить запуск в изолированных AI IDE (Codex/Claude Code/Open Code), где runtime не видит OS secret store хоста.

Фичи подверсии:
1. Вводится явный временный профиль `v1-1-ws` без OS secret store для API auth.
2. Авторизация выполняется через runtime-параметр токена (env/config) только при включенном insecure-режиме.
3. Audit HMAC ключ (`HMAC_SECRET`) также задается через runtime-параметр (env/config) только при включенном insecure-режиме.
4. Добавляются компенсирующие ограничения: non-production scope, явные предупреждения, операционный регламент revoke/rotation.
4. Документация и тестовые сценарии фиксируют, что профиль временный и должен быть заменен защищенной моделью в следующем этапе.

Границы подверсии:
- MCP-архитектура и `v2` не затрагиваются;
- базовые security-проверки `v1` (boundary/path hardening/https/allowlist/audit/limits) сохраняются;
- изменения не должны автоматически ослаблять default-профиль `v1` с OS secret store.

Артефакты OpenSpec:
- feature-пакет: `../../../openspec/changes/v1-1-ws-insecure-mvp-auth/` (планируется).

### Итерация 5 (v1-i5): http debug log для диагностики

Цель итерации:
- добавить управляемый диагностический http-лог для `doctor/run` без ослабления security-модели.

Фичи итерации:
1. Debug-лог HTTP-вызовов включается только через `ANON_HTTP_DEBUG=true`.
2. Лог пишется в `<workspace>/.anonym/http-debug.log`.
3. Чувствительные заголовки (включая `Authorization`) маскируются.

Границы итерации:
- default-поведение команд и policy секретов не меняются;
- `v2` и MCP не затрагиваются.

Артефакты OpenSpec:
- feature-пакет: `../../../openspec/changes/v1-i5-http-debug-log/`.

## Функциональные требования v1

- Команды CLI: `anonym scan`, `anonym list`, `anonym run`, `anonym doctor`.
- Локальное состояние со статусами: `new/scanned/sent/succeeded/failed`.
- Контракт конфигурации зафиксирован в `docs/technical/configuration.md` и шаблоне `config.example.yaml`.
- Базовый сценарий запуска `v1`: параметры задаются в `config.yaml`, override при необходимости выполняется через env или CLI.
- Дефолты production-контура для API в `v1`: `api.base_url=https://production-retrievals.ai.rarus-cloud.ru`, `security.allowed_hosts=["production-retrievals.ai.rarus-cloud.ru"]`.
- Интеграция с API:
  - `POST /v1/tasks/file_anonymization`
  - `GET /v1/tasks/anonymization_file_types`
  - `GET /v1/tasks/anonymization_fields`
  - preflight `doctor` для `anonymization_fields` использует headers `Authorization`, `partner-id`, `user-id` и query `page`, `per_page`.
- Контракт `POST /v1/tasks/file_anonymization` в `v1`:
  - обязательные заголовки: `Authorization`, `partner-id` (UUID), `fields=anonymizer`;
  - optional заголовки: `tags-numeration`, `user-id` (UUID);
  - формат тела: `multipart/form-data` с обязательным полем `file`.
- `api.partner_id` валидируется в UUID-формате единообразно в runtime-путях `scan/list/run/doctor`.
- `scan` только обнаруживает, не отправляет в API.
- `run <id|path>` запускает отправку явно по команде пользователя.

## Требования безопасности и надежности v1

- Граница доступа должна быть проверяемой:
  - `raw_path` обязательно вне workspace IDE;
  - запуск блокируется, если путь не проходит preflight-проверки;
  - требования к ОС-доступам (ACL) фиксируются как обязательные для окружения.
- Выбор файла должен быть защищен от обхода пути:
  - canonical path check;
  - запрет выхода за `raw_path`;
  - запрет `..`;
  - запрет симлинков/reparse points и NTFS ADS при выборе файла;
  - только файлы, не директории.
- Имена результатов не должны раскрывать PII:
  - в `output_path` сохраняется безопасное имя по `id`/техническому идентификатору, а не исходное имя raw-файла.
- Лимиты anti-DoS обязательны в v1:
  - `max_file_size`, `max_pages` (для поддерживаемых форматов), `request_timeout`, `max_parallel_runs=1` по умолчанию;
  - явная ошибка при превышении лимитов.
- Сетевая защита обязательна в v1:
  - только HTTPS;
  - проверка TLS-сертификата;
  - allowlist допустимых host/base-url для API.
  - `doctor` проверяет доступность API с теми же auth/header-инвариантами, которые используются в runtime-вызовах.
- Секреты в v1:
  - хранятся в OS secret store;
  - запрещен fallback в plaintext-конфиг;
  - в логи/ошибки секреты не попадают.
- Временное исключение для `v1-1-ws`:
  - допускаются `api.auth_token` и `audit.hmac_secret` из env/config только при явном insecure-профиле;
  - профиль ограничен non-production использованием;
  - обязательны компенсирующие меры (ротация/revoke, предупреждения в CLI и документации).
- Аудит в v1:
  - без PII и секретов;
  - `file-id` в журнале хранится в виде HMAC-хэша;
  - задается базовый retention-период аудита.
- Контракт результата анонимизации в v1:
  - результат `file_anonymization` обрабатывается через адаптер форматов ответа;
  - адаптер покрывается интеграционными тестами на реальных/эталонных ответах.
- Ограничения параметризации заголовков в `v1`:
  - `fields` фиксируется в `anonymizer` (без CLI/env override);
- Fuzzy-поиск в `run`:
  - автозапуск только при одном совпадении;
  - при нескольких совпадениях требуется явный выбор `id`.

## Критерии приемки v1

1. IDE не имеет доступа к raw-папке и не читает raw-файлы.
2. Пользователь может обнаружить файлы (`list`) и запустить обработку (`run`).
3. Результат создается в `output_path` с безопасным именем без PII.
4. Неподдерживаемые форматы, лимиты и API-ошибки возвращают понятный статус.
5. Логи не содержат PII/секретов и используют HMAC-идентификатор файла.

## Что переносим на v2

- Расширенная операционная политика секретов: ротация/отзыв и регламент реакции на компрометацию.
- Усиленная сетевая политика (включая требования корпоративного proxy при необходимости).
- Расширение диагностик и контрактов MCP-слоя без ослабления требований v1.

## Связанные документы

- Источник требований: `../../archive/2026-04-25-ai-ide-anonymization-plan.md`
- Источник сценариев: `../../archive/2026-04-25-ai-ide-anonymization-scenarios.md`
- Технический контракт конфигурации: `../../technical/configuration.md`
- Общие пользовательские сценарии: `../../user-scenarios.md`
- Общие тестовые сценарии: `../../test-scenarios.md`
- OpenSpec индекс версии: `openspec.md`
- Журнал хода версии: `action-log.md`
