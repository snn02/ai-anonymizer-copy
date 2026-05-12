# proposal: v1-i7 doctor userid pagination

## Why

Фактический контракт API для `GET /v1/tasks/anonymization_fields` требует дополнительные входы preflight-запроса `doctor`:
- заголовок `user-id`;
- query-параметры `page` и `per_page`.

Без этих параметров `doctor` может возвращать `401` даже при корректных `Authorization` и `partner-id`.

## What Changes

1. Добавить runtime-параметры `api.user_id`, `api.fields_page`, `api.fields_per_page` (config/env/CLI).
2. Передавать `user-id` в headers и `page/per_page` в query в `doctor`.
3. Обновить валидации, тесты и документацию `v1`.
