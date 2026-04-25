# proposal: v1-i2 run api adapter

## Контекст

Итерация `v1-i1` закрыла foundation для `doctor/scan/list` и security-baseline preflight.  
Итерация `v1-i2` должна дать законченный и тестируемый результат для команды `anonym run`.

## Цель итерации

Реализовать безопасный runtime-поток `run`:
- выбор цели по `id|path` и контроль неоднозначности;
- явная отправка в `POST /v1/tasks/file_anonymization`;
- адаптация ответа API в стабильный внутренний контракт;
- обновление статусов очереди (`sent/succeeded/failed`);
- сохранение результата в `output_path` с PII-safe именем.

## Обязательные условия планирования v1-i2

1. Для security-пунктов в scope фиксируется трассировка:
   - `requirement -> task -> test -> runtime evidence`.
2. Для каждой security-задачи задается DoD:
   - runtime реализован в production-пути;
   - есть минимум один негативный тест;
   - нет небезопасного production fallback.
3. Незакрытые security-пункты из `v1-i1` либо включаются в scope `v1-i2`, либо явно переносятся отдельной записью.

## In scope

1. `anonym run <id|path>` и resolver цели запуска.
2. API client + adapter для `file_anonymization`.
3. Обновление статусов каталога и кодов ошибок.
4. Runtime-проверки лимитов/timeout/сетевой политики в пути `run`.
5. Тесты:
   - unit: resolver/adapter/error mapping;
   - integration: API adapter + transitions catalog status;
   - e2e: `scan -> list -> run`.

## Out of scope

- MCP-обертка (`v2`).
- Расширенные operational-политики beyond `v1` (rotation/revocation playbooks).

## Критерии приемки

1. `run` выполняет отправку только по явной команде пользователя.
2. Ответ API стабильно нормализуется adapter-слоем.
3. Результат пишется в `output_path` с безопасным техническим именем.
4. Статусы очереди и ошибки предсказуемо отражают итог выполнения.
