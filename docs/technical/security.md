# security

## Security-модель `v1`

- `raw_path` расположен вне workspace IDE.
- Только trusted sidecar CLI имеет доступ к raw и сети API.
- IDE не читает raw-файлы напрямую.
- Отправка только по явной команде пользователя.

## Обязательные проверки

1. Boundary:
  - `raw_path` вне workspace;
  - отказ при небезопасной конфигурации.
2. Path hardening:
  - canonical path;
  - anti-traversal;
  - anti-symlink/reparse;
  - anti-ADS;
  - только обычные файлы.
3. Сеть:
  - только HTTPS;
  - TLS verification обязательно;
  - allowlist host/base URL.
4. Секреты:
  - только OS secret store;
  - plaintext fallback запрещен.
5. Аудит:
  - без PII;
  - без секретов;
  - file-id в HMAC-виде.

## Лимиты устойчивости

- `max_file_size`
- `max_pages` (для поддерживаемых форматов)
- `request_timeout`
- `max_parallel_runs=1` по умолчанию

## Политика логирования

- Структурные события с техническими кодами.
- Ошибки пригодны для диагностики, но без чувствительных данных.
