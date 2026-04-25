# error model

## Цель

Сделать ошибки CLI одновременно понятными для пользователя и стабильными для автоматизации.

## Формат ошибки

- `code`: машинный код (`ANON_*`).
- `message`: короткое человеческое описание.
- `details`: безопасные технические детали (без PII/секретов).
- `next_step`: что сделать дальше.

## Классы ошибок

1. Конфигурация:
  - `ANON_CONFIG_INVALID`
  - `ANON_BOUNDARY_VIOLATION`
2. Безопасность пути:
  - `ANON_PATH_TRAVERSAL_BLOCKED`
  - `ANON_REPARSE_BLOCKED`
  - `ANON_ADS_BLOCKED`
3. Сеть и API:
  - `ANON_API_UNREACHABLE`
  - `ANON_API_AUTH_FAILED`
  - `ANON_API_CONTRACT_ERROR`
4. Лимиты:
  - `ANON_FILE_TOO_LARGE`
  - `ANON_TIMEOUT`
  - `ANON_PARALLEL_LIMIT`
5. Системные:
  - `ANON_IO_ERROR`
  - `ANON_INTERNAL_ERROR`

## Требования

- Ошибки детерминированы: одинаковая причина -> одинаковый код.
- Вывод не содержит чувствительных данных.
- Для `run` всегда показывается следующий практический шаг.
