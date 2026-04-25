# components

## Модульная структура trusted-core

## `cmd/anonym`

- Точка входа CLI.
- Регистрация команд: `scan`, `list`, `run`, `doctor`.

## `internal/config`

- Загрузка конфигурации.
- Валидация обязательных параметров и лимитов.

## `internal/preflight`

- Проверка boundary `raw_path`/workspace.
- Проверка доступности секретов и API.
- Проверка обязательных security-политик.

## `internal/scanner`

- Обнаружение файлов в `raw_path`.
- Без сети и без отправки.

## `internal/catalog`

- Локальный индекс файлов и статусов:
  - `new`
  - `scanned`
  - `sent`
  - `succeeded`
  - `failed`

## `internal/resolver`

- Выбор цели `run` по `id`, относительному пути или fuzzy-строке.
- При нескольких fuzzy-совпадениях только список кандидатов.

## `internal/security/pathguard`

- Canonical path check.
- Запрет path traversal (`..`).
- Запрет symlink/reparse points.
- Запрет NTFS ADS.
- Разрешены только обычные файлы.

## `internal/anonymizer/client`

- Работа с API:
  - `POST /v1/tasks/file_anonymization`
  - `GET /v1/tasks/anonymization_file_types`
  - `GET /v1/tasks/anonymization_fields`
- Адаптация и нормализация ответов API.

## `internal/output`

- Генерация безопасных технических имен результата.
- Запись результата в `output_path`.

## `internal/audit`

- Структурный аудит без PII/секретов.
- HMAC-идентификатор файла в журнале.

## `internal/secrets`

- Доступ к API-ключам только через OS secret store.

## `internal/platform`

- Платформенные проверки для Windows/Linux/macOS.
