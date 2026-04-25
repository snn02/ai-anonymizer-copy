# tasks: v1-i2 run api adapter

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| `run <id|path>` запускает отправку явно | T1, T2 | U1, I1, E1 | CLI `anonym run` вызывает отправку только в `run` |
| адаптер `file_anonymization` обязателен | T3 | U2, I2 | adapter-слой в `internal/anonymizer/client` |
| статусы очереди `sent/succeeded/failed` | T4 | U3, I3, E2 | переходы статусов в каталоге |
| PII-safe output naming | T5 | U4, E3 | запись результата в `output_path` по technical id |
| security run path (https/tls/allowlist/timeout/limits/no-fallback) | T6 | U5-U8, I4 | runtime guardrails в production пути `run` |

## security DoD (для T6)

- [ ] runtime реализован в production-пути `run`;
- [ ] есть минимум один негативный тест на каждый security-check;
- [ ] нет небезопасного production fallback;
- [ ] evidence подтвержден прогоном `go test ./...`.

## tasks

### T1. run command + target resolver

- [ ] реализовать `anonym run <id|path>` с явным запуском отправки;
- [ ] реализовать выбор цели по `id|path` и детерминированные ошибки.

### T2. run ambiguity policy (fuzzy)

- [ ] добавить политику: автозапуск только при одном fuzzy-совпадении;
- [ ] при нескольких совпадениях возвращать список кандидатов и требовать явный `id`.

### T3. API client + response adapter

- [ ] реализовать client для `POST /v1/tasks/file_anonymization`;
- [ ] добавить adapter нормализации ответа API во внутренний контракт.

### T4. catalog status transitions

- [ ] реализовать переходы статусов `scanned -> sent -> succeeded/failed`;
- [ ] стандартизировать error mapping для CLI.

### T5. output writer (PII-safe)

- [ ] сохранять результат в `output_path` с безопасным техническим именем;
- [ ] исключить утечку исходного raw-имени в output filename.

### T6. security guardrails in run path

- [ ] enforce `https` + tls verify + host allowlist для `run`;
- [ ] enforce timeout и обязательные лимиты runtime;
- [ ] запретить небезопасные production fallback.

## tests

### unit (U*)

- [ ] U1: resolver выбирает цель по `id|path`, ошибки на not-found.
- [ ] U2: adapter нормализует варианты `file_anonymization`.
- [ ] U3: transitions статусов каталога корректны.
- [ ] U4: output naming использует technical id (без PII).
- [ ] U5: run rejects non-https endpoint.
- [ ] U6: run rejects host вне allowlist.
- [ ] U7: run enforces timeout/limit violations.
- [ ] U8: run path не имеет production `allow-all` fallback.

### integration/e2e (I*/E*)

- [ ] I1: `run` вызывает API только по явной команде.
- [ ] I2: adapter + API mock совместимы с контрактом.
- [ ] I3: catalog transitions отражают outcome `run`.
- [ ] I4: security guardrails срабатывают в runtime path.
- [ ] E1: `scan -> list -> run <id>` завершает отправку.
- [ ] E2: ошибки API отображаются и фиксируют `failed`.
- [ ] E3: output создан в `output_path` с безопасным именем.
