# decisions

## Формат записи

- Дата
- Идентификатор решения
- Решение
- Причина
- Влияние

## Журнал решений

### 2026-04-25 / d-001

- Решение: документация проекта ведется на русском языке.
- Причина: единый язык коммуникации команды.
- Влияние: все документы в репозитории.

### 2026-04-25 / d-002

- Решение: имена файлов и папок документации ведем в lowercase.
- Причина: консистентность и переносимость между ОС.
- Влияние: корневые документы, `docs/*`, `openspec/*`.

### 2026-04-25 / d-003

- Решение: пользовательские и тестовые сценарии ведем общими файлами в `docs/` с маркировкой применимости по версиям.
- Причина: снизить дублирование и упростить сопровождение.
- Влияние: `docs/user-scenarios.md`, `docs/test-scenarios.md`.

### 2026-04-25 / d-004

- Решение: feature-спеки в `openspec/changes/` не создаем до согласования полного workflow разработки.
- Причина: сначала согласуем набор документов и процесс.
- Влияние: `docs/versions/*/openspec.md`, `openspec/changes/`.
- Статус: заменено решением `d-008` для `v1`; для `v2` ограничение действует до отдельного решения.

### 2026-04-25 / d-005

- Решение: исходные документы `2026-04-25-ai-ide-anonymization-plan.md` и `2026-04-25-ai-ide-anonymization-scenarios.md` перенесены в архив.
- Причина: рабочими документами стали структурированные файлы в `docs/` и `docs/versions/`.
- Влияние: `docs/archive/*`, обновленные ссылки в `readme.md` и планах версий.

### 2026-04-25 / d-006

- Решение: выявленные риски безопасности/надежности разложены по версиям как обязательные требования.
- Причина: зафиксировать, что обязательно закрываем в v1, а что переносим в v2.
- Влияние: `docs/versions/v1/plan.md`, `docs/versions/v2/plan.md`, `docs/test-scenarios.md`, `docs/roadmap.md`, `docs/user-guide.md`, `docs/user-scenarios.md`.

### 2026-04-25 / d-007

- Решение: закрепить единый контракт конфигурации trusted CLI (`docs/technical/configuration.md`) и шаблон `config.example.yaml`; секреты хранить только в OS secret store; приоритет источников конфигурации: flags > env > file.
- Причина: снизить риск небезопасной конфигурации и неоднозначного поведения между средами.
- Влияние: `docs/technical/configuration.md`, `config.example.yaml`, `docs/user-guide.md`, `docs/test-scenarios.md`, `docs/versions/v1/plan.md`, `docs/versions/v1/action-log.md`.

### 2026-04-25 / d-008

- Решение: workflow разработки для `v1` согласован и вводится в действие:
  1. формируется общий план разработки `v1`;
  2. работа идет итерациями, каждая итерация должна давать законченный и тестируемый результат;
  3. по итерации создается feature в `openspec/changes/`;
  4. после реализации итерации выполняются тесты, фиксируются выводы, планируется следующая итерация;
  5. после завершения всех итераций `v1` выполняются общее тестирование, релиз и запись в `lessons-learned`.
- Причина: перейти от подготовки документации к управляемому циклу поставки.
- Влияние: `docs/versions/v1/plan.md`, `docs/versions/v1/openspec.md`, `docs/versions/v1/action-log.md`, `docs/test-scenarios.md`, `docs/technical/release.md`, `agents.md`, `docs/doc-map.md`.

### 2026-04-25 / d-009

- Решение: для планирования итераций OpenSpec в `v1` вводятся обязательные quality-гейты:
  1. трассировка `requirement -> task -> test -> runtime evidence` для требований в scope;
  2. отдельный DoD для security-задач (runtime, негативный тест, отсутствие небезопасного fallback);
  3. запрет закрытия security-задач при active production fallback вида `allow-all`;
  4. обязательная запись результата release-checklist в `docs/versions/v1/action-log.md`.
- Причина: исключить повторение разрыва между требованиями `plan.md` и задачами OpenSpec.
- Влияние: `docs/versions/v1/openspec.md`, `openspec/changes/v1-i2-run-api-adapter/proposal.md`, `docs/versions/v1/action-log.md`.

### 2026-04-26 / d-010

- Решение: для `v1` закрепить контрактную синхронизацию API-вызова `file_anonymization` как обязательный критерий документации и тестов:
  1. `api.partner_id` и `partner-id` всегда указываются в UUID-формате;
  2. обязательные заголовки вызова фиксируются явно (`Authorization`, `partner-id`, `fields=anonymizer`);
  3. примеры в user-guide не должны использовать невалидные placeholder-значения вместо UUID.
- Причина: устранить расхождения между пользовательской документацией, техническим контрактом и runtime-проверками при вызове production API.
- Влияние: `docs/technical/configuration.md`, `docs/user-guide.md`, `docs/versions/v1/plan.md`, `docs/user-scenarios.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md`.

### 2026-04-27 / d-011

- Решение: ввести временный подэтап `v1-1-ws` для MVP-запуска в изолированных AI IDE без доступа к OS secret store; допустить auth-токен в env/config только при явном insecure-профиле и только для non-production контуров.
- Причина: основной сценарий ценности (запуск из Codex/Claude Code/Open Code) блокируется изоляцией окружения и недоступностью host keyring для runtime.
- Влияние: `docs/roadmap.md`, `docs/versions/v1/plan.md`, `docs/technical/security.md`, `docs/technical/configuration.md`, `docs/user-scenarios.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md`.

### 2026-05-01 / d-012

- Решение: упростить runtime-контракт `v1` до mode-first модели:
  1. единый режим `runtime.mode` (`prod|mvp`) вместо `runtime.profile + insecure flag`;
  2. минимальный CLI override: только `--config`, `--mode`, `--workspace-path`, `--output-path`;
  3. остальные параметры читаются из `config/env`;
  4. режим `mvp` допускается для production API host с явным warning в CLI.
- Причина: снизить операционные ошибки из-за перегруженного CLI и обеспечить практический MVP-запуск в AI IDE.
- Влияние: `cmd/anonym/main.go`, `internal/config/config.go`, `internal/preflight/preflight.go`, `internal/anonymizer/run.go`, `docs/user-guide.md`, `docs/technical/configuration.md`, `docs/test-scenarios.md`, `docs/user-scenarios.md`, `docs/versions/v1/plan.md`, `docs/versions/v1/action-log.md`.
