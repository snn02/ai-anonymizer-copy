# tasks: v1-i2 run api adapter

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| `run <id|path>` запускает отправку явно | T1, T2 | U1, I1, E1 | CLI `anonym run` вызывает отправку только в `run` |
| адаптер `file_anonymization` обязателен | T3 | U2, I2 | adapter-слой в `internal/anonymizer/run.go` |
| статусы очереди `sent/succeeded/failed` | T4 | U3, I3, E2 | переходы статусов в каталоге |
| PII-safe output naming | T5 | U4, E3 | запись результата в `output_path` по technical id |
| security run path (https/tls/allowlist/timeout/limits/no-fallback) | T6 | U5-U8, I4 | runtime guardrails в `internal/anonymizer/run.go` production пути `run` |
| audit logging (safe fields + HMAC file-id + retention) | T7 | U9-U11, I5 | runtime audit recorder в `internal/audit/recorder.go` и вызов из `run` |

## security DoD (для T6)

- [x] runtime реализован в production-пути `run`;
- [x] есть минимум один негативный тест на каждый security-check;
- [x] нет небезопасного production fallback;
- [x] evidence подтвержден прогоном `go test ./...`.

## tasks

### T1. run command + target resolver

- [x] реализовать `anonym run <id|path>` с явным запуском отправки;
- [x] реализовать выбор цели по `id|path` и детерминированные ошибки.

### T2. run ambiguity policy (fuzzy)

- [x] добавить политику: автозапуск только при одном fuzzy-совпадении;
- [x] при нескольких совпадениях возвращать список кандидатов и требовать явный `id`.

### T3. API client + response adapter

- [x] реализовать client для `POST /v1/tasks/file_anonymization`;
- [x] добавить adapter нормализации ответа API во внутренний контракт.

### T4. catalog status transitions

- [x] реализовать переходы статусов `scanned -> sent -> succeeded/failed`;
- [x] стандартизировать error mapping для CLI.

### T5. output writer (PII-safe)

- [x] сохранять результат в `output_path` с безопасным техническим именем;
- [x] исключить утечку исходного raw-имени в output filename.

### T6. security guardrails in run path

- [x] enforce `https` + tls verify + host allowlist для `run`;
- [x] enforce timeout и обязательные лимиты runtime (`max_file_size`, `max_pages` для PDF/DOCX/текстовых форматов);
- [x] enforce Windows path hardening для `run` (ADS/reparse-point deny);
- [x] запретить небезопасные production fallback.

### T7. audit trail and retention

- [x] реализовать runtime-аудит в `run` с записью только безопасных полей;
- [x] записывать `file-id` в журнал в виде HMAC-хэша (без raw/path/PII);
- [x] применять retention-политику по `audit.retention_days` при каждой записи;
- [x] добавить обязательную проверку `audit.hmac_key_id` и чтение audit-ключа из OS secret store.

## tests

### unit (U*)

- [x] U1: resolver выбирает цель по `id|path`, ошибки на not-found.
- [x] U2: adapter нормализует варианты `file_anonymization`.
- [x] U3: transitions статусов каталога корректны.
- [x] U4: output naming использует technical id (без PII).
- [x] U5: run rejects non-https endpoint.
- [x] U6: run rejects host вне allowlist.
- [x] U7: run enforces timeout/limit violations.
- [x] U7a: run rejects PDF/DOCX/text documents that exceed `max_pages`.
- [x] U8: run path не имеет production `allow-all` fallback.
- [x] U8a: run rejects Windows ADS/reparse paths.
- [x] U9: run rejects missing audit hmac key in OS secret store.
- [x] U10: run writes audit entry with HMAC `file_id` and safe payload.
- [x] U11: audit recorder prunes entries older than retention window.

### integration/e2e (I*/E*)

- [x] I1: `run` вызывает API только по явной команде.
- [x] I2: adapter + API mock совместимы с контрактом.
- [x] I3: catalog transitions отражают outcome `run`.
- [x] I4: security guardrails срабатывают в runtime path.
- [x] I5: audit trail создается при `run`, содержит только безопасные поля и соблюдает retention.
- [x] E1: `scan -> list -> run <id>` завершает отправку.
- [x] E2: ошибки API отображаются и фиксируют `failed`.
- [x] E3: output создан в `output_path` с безопасным именем.
