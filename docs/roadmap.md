# roadmap

## v1

- Цель: запустить безопасный sidecar CLI для анонимизации файлов.
- Шаги:
  - реализовать `scan/list/run/doctor`;
  - подключить API и preflight-валидации;
  - обеспечить проверяемую границу доступа (`raw_path`/ACL/preflight);
  - внедрить hardening выбора файла (canonical path, anti-traversal, anti-symlink/reparse);
  - ввести безопасные имена результирующих файлов без утечки PII;
  - включить лимиты anti-DoS (`max_file_size`, timeout, параллельность);
  - включить сетевую защиту (HTTPS + TLS verify + host allowlist);
  - хранить секреты в OS secret store без plaintext fallback;
  - добавить аудит без PII/секретов с HMAC-идентификатором;
  - закрепить адаптер результата `file_anonymization` тестами.
- Статус: in progress (текущая итерация `v1-i2`).
- План версии: `versions/v1/plan.md`.

## v2

- Цель: добавить MCP-обертку поверх существующего trusted-core CLI.
- Шаги:
  - подключить MCP-интерфейс;
  - сохранить security-модель v1;
  - подключить OpenCode и Claude через MCP tool-интерфейс;
  - добавить операционный lifecycle секретов (rotation/revocation/compromise response);
  - расширить эксплуатационные политики сети/прокси в корпоративном контуре.
- Статус: planned.
- План версии: `versions/v2/plan.md`.

## Примечание

- Для `v1` feature-пакеты в `../openspec/changes/` создаются по итерациям (см. `docs/versions/v1/openspec.md`).
- Для `v2` создание feature-пакетов открывается отдельным решением.
