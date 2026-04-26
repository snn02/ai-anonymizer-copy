# v1 action-log

## Формат записи

- Дата
- Тип: изменение | ошибка | реакция
- Описание
- Ссылка

## Записи

### 2026-04-25

- Тип: изменение
- Описание: перенесена структура документации в lowercase и заполнена на основе исходных документов плана и сценариев.
- Ссылка: `../../doc-map.md`

### 2026-04-25

- Тип: реакция
- Описание: зафиксировано, что feature-спеки в OpenSpec создаются только после согласования workflow.
- Ссылка: `openspec.md`

### 2026-04-25

- Тип: изменение
- Описание: требования v1 уточнены по безопасности и надежности: boundary enforcement, path hardening, PII-safe naming, anti-DoS лимиты, сетевая защита, lifecycle секретов (базовый), HMAC-аудит, правила fuzzy run.
- Ссылка: `plan.md`

### 2026-04-25

- Тип: изменение
- Описание: добавлен технический контракт конфигурации (`docs/technical/configuration.md`) и шаблон `config.example.yaml`; закреплены правила источников (flags > env > file) и запрет plaintext fallback для секретов.
- Ссылка: `../../technical/configuration.md`

### 2026-04-25

- Тип: изменение
- Описание: workflow для v1 согласован; фичи OpenSpec создаются по итерациям, каждая итерация должна быть законченной и тестируемой.
- Ссылка: `openspec.md`

### 2026-04-25

- Тип: изменение
- Описание: запланирована первая итерация `v1-i1` и создан feature-пакет OpenSpec с фокусом на `doctor` preflight и безопасный `scan/list` без API-отправки.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/proposal.md`

### 2026-04-25

- Тип: изменение
- Описание: начата реализация `v1-i1`: добавлен стартовый каркас trusted CLI на Go и реализован базовый `doctor preflight` с тестами на boundary/https/allowlist/max_parallel_runs.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: зафиксирован технический блокер окружения разработки: команда `go` отсутствует, из-за чего локальный прогон unit-тестов и форматирования `gofmt` отложен до установки Go toolchain.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: блокер по Go toolchain снят; подтвержден рабочий прогон `go test ./...` и `gofmt` для модулей `cmd/anonym`, `internal/config`, `internal/preflight`, `internal/scanner`, `internal/catalog`.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: изменение
- Описание: продолжена реализация `v1-i1`: добавлены команды `scan`/`list`, безопасный discovery файлов в `raw_path` и локальный каталог очереди `.anonym/catalog.json` без сетевой отправки.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: изменение
- Описание: по TDD расширен `doctor preflight`: добавлены проверки наличия partner secret (через `DoctorDeps.SecretChecker`) и доступности API endpoint `GET /v1/tasks/anonymization_fields` с timeout; добавлены новые unit-тесты `internal/preflight/doctor_test.go`.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: изменение
- Описание: по TDD добавлены лимиты discovery для `scan`: `max_file_size_mb` в конфигурации (flags/env) и фильтрация файлов по размеру через `scanner.DiscoverWithLimits`; добавлены/обновлены тесты `internal/config/config_test.go` и `internal/scanner/scanner_test.go`.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: результаты текущего состояния `v1-i1` зафиксированы после полного прогона `go test ./...`; подготовлено предложение на следующую итерацию `run + api adapter`.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/proposal.md`

### 2026-04-25

- Тип: изменение
- Описание: по TDD добавлены базовые preflight ACL-проверки: `raw_path` и `output_path` должны существовать и быть директориями, для `output_path` выполняется write-probe; тесты `internal/preflight/*` обновлены на реальные временные директории.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: подтверждено состояние итерации `v1-i1`: после внедрения лимитов discovery и ACL-checks полный прогон `go test ./...` проходит успешно.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: ошибка
- Описание: выявлен разрыв между требованиями `v1` и планированием итерации `v1-i1`: в плане и proposal требование «секреты из OS secret store» присутствует, но в задачах итерации оно не было декомпозировано в проверяемую runtime-задачу (был закрыт интерфейсный слой `SecretChecker`, при этом runtime остался с `AllowAllSecretChecker` fallback).
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-25

- Тип: реакция
- Описание: зафиксирован postmortem для переноса в `lessons-learned`: причина — отсутствие трассировки «требование -> подзадача -> тест -> runtime-проверка» при декомпозиции preflight; corrective actions — 1) для security-требований вводить отдельные task-чекбоксы с явным DoD «runtime реализован, fallback запрещен, есть негативный тест», 2) перед закрытием итерации проходить release-checklist на соответствие `plan.md` по каждому security-пункту, 3) запрещать для production default-заглушки вида `allow-all` без отдельного feature-flag и явной пометки non-production.
- Ссылка: `plan.md`

### 2026-04-25

- Тип: изменение
- Описание: обновлены принципы planning workflow OpenSpec для `v1`: в `docs/versions/v1/openspec.md` добавлены обязательные quality-гейты (трассировка requirement->task->test->runtime evidence, security-DoD, запрет закрытия задач при production `allow-all` fallback, release-checklist перед закрытием итерации); эти гейты распространены на планирование `v1-i2`.
- Ссылка: `openspec.md`

### 2026-04-26

- Тип: изменение
- Описание: доработана `v1-i1` по security-требованию OS secret store: реализован runtime `internal/secrets` (keyring), `anonym doctor` использует `OSSecretChecker` в production-пути, default fallback `AllowAllSecretChecker` удален.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-26

- Тип: реакция
- Описание: по актуализированным правилам OpenSpec зафиксирован DoD/evidence для security-задачи `doctor` (runtime, негативные тесты, отсутствие allow-all fallback); подтверждено прогоном `go test ./...` с зеленым статусом всех пакетов.
- Ссылка: `../../../openspec/changes/v1-i1-foundation-preflight-discovery/tasks.md`

### 2026-04-26

- Тип: изменение
- Описание: обновлен OpenSpec-пакет итерации `v1-i2-run-api-adapter`: зафиксированы `proposal`, декомпозиция `tasks` с трассировкой `requirement -> task -> test -> runtime evidence`, security DoD и `spec` с проверяемыми сценариями `run/api-adapter/status/security`.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/proposal.md`

### 2026-04-26

- Тип: изменение
- Описание: начата и реализована итерация `v1-i2` по TDD: добавлены команда `anonym run`, сервисный слой `internal/anonymizer` (resolver `id|path|fuzzy`, API client/adapter, переходы статусов `sent/succeeded/failed`, PII-safe output naming), а также чтение секрета через OS secret store без production fallback.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/tasks.md`

### 2026-04-26

- Тип: реакция
- Описание: подтверждено покрытие ключевых проверок `v1-i2` (явный запуск run, ambiguous fuzzy policy, API adapter, статусы каталога, security guardrails `https/allowlist/limits/no-fallback`) и зеленый прогон `go test ./...`.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/specs/v1-i2-core/spec.md`

### 2026-04-26

- Тип: изменение
- Описание: в пользовательской документации добавлен полный каталог параметров приложения с источниками (`flags/env/config`) и зонами ответственности по изменению лимитов/безопасности; отдельно зафиксирован статус параметров, уже применяемых в runtime `v1`.
- Ссылка: `../../user-guide.md`

### 2026-04-26

- Тип: изменение
- Описание: по `v1-i2` закрыты незавершенные security/reliability пункты: в runtime `run` добавлены `max_pages` (для поддерживаемых форматов) и усиленный Windows path hardening (запрет ADS и reparse points); добавлены/обновлены unit и integration проверки.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/tasks.md`

### 2026-04-26

- Тип: реакция
- Описание: подтвержден зеленый прогон `go test ./...` после внедрения `max_pages` и Windows hardening в `run`; OpenSpec `v1-i2` синхронизирован с фактическим runtime evidence.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/tasks.md`

### 2026-04-26

- Тип: изменение
- Описание: устранены замечания техлид-ревью по `v1-i2`: `max_pages` в `run` расширен на целевые форматы PDF и DOCX (добавлены тесты), а в технической документации синхронизировано имя env-переменной `api.partner_id` с фактическим кодом (`ANON_API_PARTNER_ID`).
- Ссылка: `../../technical/configuration.md`

### 2026-04-26

- Тип: изменение
- Описание: скорректирована пользовательская документация `v1`: удалены неподдерживаемые CLI override-флаги `run` (`--fields/--tags/--tags-numeration`) и добавлен блок известных ограничений/рисков по текущему page-counter для PDF и DOCX.
- Ссылка: `../../user-guide.md`

### 2026-04-26

- Тип: изменение
- Описание: по TDD закрыт `T7` для `v1-i2`: добавлен runtime audit trail (`internal/audit`) с безопасным payload, HMAC `file_id`, retention по `audit.retention_days`; интегрировано в `anonym run` с обязательным чтением audit-ключа из OS secret store (`audit.hmac_key_id`) без plaintext fallback.
- Ссылка: `../../../openspec/changes/v1-i2-run-api-adapter/tasks.md`

### 2026-04-26

- Тип: реакция
- Описание: подтверждено `p0`-покрытие `T7` (безопасные поля журнала, HMAC file-id, retention, совместимость adapter) и зеленый прогон `go test ./...` после синхронизации OpenSpec и документации.
- Ссылка: `../../test-scenarios.md`

### 2026-04-26

- Тип: реакция
- Описание: выполнен финальный release gate `v1` по `p0`-сценариям (`T1`, `T2`, `T3`, `T5`, `T6`, `T7`, `T9`, `T10`) с подтверждением зеленого прогона `go test ./...`; технических блокеров в runtime v1 не выявлено.
- Ссылка: `../../test-scenarios.md`

### 2026-04-26

- Тип: изменение
- Описание: версия `v1` закрыта по процессу: roadmap переведен в статус `completed`, feature-пакеты `v1-i1`/`v1-i2` перенесены в `openspec/archive`, добавлены выводы в `docs/lessons-learned.md`.
- Ссылка: `../../roadmap.md`

### 2026-04-26

- Тип: изменение
- Описание: полностью переписано пользовательское руководство `v1` понятным пошаговым форматом: два сценария запуска (готовый `exe` и `git clone`), обязательная структура папок, детальная настройка параметров и секретов OS secret store, чек-лист запуска без релиза.
- Ссылка: `../../user-guide.md`

### 2026-04-26

- Тип: изменение
- Описание: расширен раздел секретов в user-guide: добавлены пошаговые инструкции для Windows/Linux/macOS, пояснение по `audit.hmac_key_id` и источник `HMAC`-секрета, а также команда генерации значения через PowerShell.
- Ссылка: `../../user-guide.md`

### 2026-04-26

- Тип: изменение
- Описание: выполнена контрактная синхронизация документации `v1` для API anonymization: в `plan.md` зафиксированы обязательные заголовки `POST /v1/tasks/file_anonymization` (`Authorization`, `partner-id`, `fields=anonymizer`), формат `multipart/form-data` и ограничения параметризации заголовков в runtime `v1`.
- Ссылка: `plan.md`

### 2026-04-26

- Тип: изменение
- Описание: обновлены пользовательские и тестовые сценарии `v1`: добавлены проверяемые критерии UUID-формата `api.partner_id`, явные проверки обязательных заголовков в `run` и auth/header-контекста preflight-проверки `anonymization_fields`.
- Ссылка: `../../user-scenarios.md`

### 2026-04-26

- Тип: изменение
- Описание: исправлен user-guide под production-контракт API: примеры `api.partner_id` переведены на UUID-формат (вместо `partner-1`), добавлено сопоставление входных данных `ключ/партнер/пользователь` к заголовкам `Authorization/partner-id/user-id` и зафиксированы ограничения runtime `v1` по `fields/user-id/tags-numeration`.
- Ссылка: `../../user-guide.md`

### 2026-04-26

- Тип: изменение
- Описание: переписан технический контракт конфигурации `v1` (`docs/technical/configuration.md`) в синхронизации с runtime: обновлены фактические defaults/пути, выделен раздел обязательных заголовков API-вызова и явно отмечены заголовки, не параметризуемые через CLI/env/config в `v1`.
- Ссылка: `../../technical/configuration.md`

### 2026-04-26

- Тип: реакция
- Описание: проверена актуальность task-контрактов относительно production OpenAPI (`/v1/tasks/file_anonymization`, `/v1/tasks/anonymization_fields`, `/v1/tasks/anonymization_file_types`): структура обязательных заголовков и форматы тела/авторизации совпадают с локальной спецификацией.
- Ссылка: `../../openapi.json`

### 2026-04-26

- Тип: изменение
- Описание: сформирован post-release feature-пакет OpenSpec `v1-i3-doctor-api-header-alignment` для кодовых правок preflight: выравнивание заголовков `doctor` с контрактом API (`Authorization`, `partner-id`), UUID-валидация `api.partner_id`, security DoD и тестовый контур.
- Ссылка: `../../../openspec/archive/v1-i3-doctor-api-header-alignment/proposal.md`

### 2026-04-26

- Тип: изменение
- Описание: реализованы кодовые правки `v1-i3` в preflight/doctor: API-check `GET /v1/tasks/anonymization_fields` теперь отправляет `Authorization` и `partner-id`; добавлена runtime-валидация UUID для `api.partner_id`; секрет для preflight берется через `GetSecret` из OS secret store без plaintext fallback.
- Ссылка: `../../../openspec/archive/v1-i3-doctor-api-header-alignment/tasks.md`

### 2026-04-26

- Тип: реакция
- Описание: подтверждено тестами `v1-i3`: обновлены проверки `internal/preflight/doctor_test.go` (headers, UUID validation, missing secret behavior), адаптирован `cmd/anonym/main_test.go` под новый контракт secret checker; полный прогон `go test ./...` зеленый.
- Ссылка: `../../test-scenarios.md`

### 2026-04-26

- Тип: изменение
- Описание: feature-пакет `v1-i3-doctor-api-header-alignment` закрыт и перенесен в `openspec/archive/`; индекс `docs/versions/v1/openspec.md` синхронизирован.
- Ссылка: `openspec.md`

### 2026-04-26

- Тип: изменение
- Описание: по техлид-ревью доработан `v1-i3`: в preflight добавлен детерминированный маппинг `secret not found` в стабильную бизнес-ошибку, усилен тест на источник значения `Authorization`, и UUID-валидация `api.partner_id` унифицирована в общем runtime-пути (`scan/list/run/doctor`).
- Ссылка: `../../test-scenarios.md`

### 2026-04-26

- Тип: изменение
- Описание: обновлены дефолты API-конфигурации под production endpoint: в `config.example.yaml` и документации (`technical/configuration.md`, `user-guide.md`) зафиксированы значения по умолчанию `api.base_url=https://production-retrievals.ai.rarus-cloud.ru` и `security.allowed_hosts=["production-retrievals.ai.rarus-cloud.ru"]`.
- Ссылка: `../../technical/configuration.md`

### 2026-04-26

- Тип: изменение
- Описание: запланирована итерация `v1-i4` для доведения runtime до config-first сценария: чтение `config.yaml`/`--config`, гарантированный приоритет `CLI > env > config file`, и компактные вызовы `doctor/scan/list/run` с точечными override.
- Ссылка: `plan.md`

### 2026-04-26

- Тип: изменение
- Описание: создан OpenSpec feature-пакет `v1-i4-config-first-loading` (`proposal/tasks/spec`) для реализации config-first runtime-конфигурации и синхронизирован индекс `docs/versions/v1/openspec.md` с переводом `v1-i4` в активные итерации.
- Ссылка: `openspec.md`
