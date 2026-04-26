# proposal: v1-i3 doctor api header alignment

## Контекст

После закрытия `v1` выявлен контрактный разрыв между документацией и runtime preflight:
1. в `v1` требованиях закреплено, что preflight-проверка API должна использовать тот же auth/header-контекст, что и runtime-вызовы;
2. фактический `doctor` при проверке `GET /v1/tasks/anonymization_fields` не передает обязательные заголовки контракта (`Authorization`, `partner-id`);
3. в `v1` требуется UUID-формат `api.partner_id`, но это должно быть гарантировано runtime-валидацией.

## Цель итерации

Закрыть контрактный разрыв preflight без изменения поведенческих границ `v1`:
1. `doctor` проверяет API с корректными заголовками и секретом из OS secret store;
2. `api.partner_id` валидируется как UUID в runtime-пути;
3. покрытие тестами подтверждает детерминированные ошибки и отсутствие insecure fallback.

## In scope

1. Runtime-правка preflight для `GET /v1/tasks/anonymization_fields`:
   - `Authorization` из OS secret store;
   - `partner-id` из `api.partner_id`.
2. Runtime-валидация UUID-формата `api.partner_id` для `doctor`.
3. Обновление unit/integration тестов для preflight-контракта.
4. Синхронизация OpenSpec-артефактов по трассировке `requirement -> task -> test -> runtime evidence`.

## Out of scope

1. Изменения поведения `run`/`scan`/`list`.
2. Добавление `user-id` в runtime `v1`.
3. Параметризация `fields/tags/tags-numeration` через CLI/env/config.
4. Любые `v2`-задачи MCP/операционного расширения.

## Критерии приемки

1. `doctor` выполняет API-проверку с заголовками `Authorization` и `partner-id`.
2. При не-UUID значении `api.partner_id` `doctor` завершается явной ошибкой валидации до сетевого вызова.
3. Все новые проверки покрыты unit/integration тестами и входят в `p0`-гейт `v1`.
