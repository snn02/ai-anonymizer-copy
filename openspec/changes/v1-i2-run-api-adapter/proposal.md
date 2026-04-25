# proposal: v1-i2 run api adapter

## Контекст

Итерация `v1-i1` закрывает foundation для `doctor/scan/list`: preflight, безопасный discovery, локальная очередь и базовые лимиты discovery.

Следующий завершенный и тестируемый результат для `v1` — управляемый запуск отправки по `anonym run` с обработкой ответа API через адаптер контракта.

## Цель итерации

Реализовать безопасный поток `run`:
- выбор элемента очереди по `id|path`;
- явная отправка в `POST /v1/tasks/file_anonymization`;
- адаптация ответа API в стабильный внутренний контракт;
- запись результата в `output_path` с безопасным техническим именем без PII.

## In scope

1. Команда `anonym run <id|path>` и разрешение цели запуска.
2. API client + adapter для контракта `file_anonymization`.
3. Обновление статусов очереди (`sent/succeeded/failed`) и ошибок.
4. Базовые проверки лимитов и timeout в runtime-path `run`.
5. Тесты:
   - unit для resolver/adapter/error mapping;
   - integration для API adapter и catalog status transitions;
   - e2e сценарий `scan -> list -> run`.

## Out of scope

- MCP-обертка (`v2`).
- Расширенные операционные политики beyond `v1` (rotation/revocation playbooks).

## Критерии приемки

1. `run` выполняет отправку только по явной команде пользователя.
2. Контракт ответа API стабильно нормализуется адаптером.
3. Результат в `output_path` сохраняется с PII-safe именем.
4. Статусы очереди и ошибки предсказуемо отражают outcome выполнения.
