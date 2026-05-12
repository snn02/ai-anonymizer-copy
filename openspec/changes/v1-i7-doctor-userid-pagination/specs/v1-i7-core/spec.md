# spec: v1-i7 core

## Added Requirements

### Requirement: doctor preflight передает user-id и pagination параметры

Runtime MUST передавать `user-id` header и query-параметры `page/per_page` в `GET /v1/tasks/anonymization_fields` при `anonym doctor`.

#### Scenario: doctor request includes required inputs

- Given пользователь запускает `anonym doctor`
- When runtime формирует preflight-запрос
- Then заголовки включают `Authorization`, `partner-id`, `user-id`
- And query включает `page` и `per_page`

### Requirement: параметры доступны через config/env/CLI

Runtime MUST поддерживать параметры `api.user_id`, `api.fields_page`, `api.fields_per_page` с приоритетом `CLI > env > config`.

#### Scenario: source precedence for doctor request params

- Given значения заданы одновременно в config, env и CLI
- When выполняется `anonym doctor`
- Then runtime использует значения источника с более высоким приоритетом
