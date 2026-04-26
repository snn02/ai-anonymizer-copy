# spec: v1-i3 core

## ADDED Requirements

### Requirement: doctor использует контрактные заголовки при preflight API-check

`anonym doctor` MUST выполнять `GET /v1/tasks/anonymization_fields` с тем же auth/header-контекстом, который требуется контрактом API и используется runtime-интеграцией.

#### Scenario: preflight-запрос с обязательными заголовками

- Given корректная конфигурация `api.base_url`, `api.partner_id` и доступный секрет в OS secret store
- When пользователь запускает `anonym doctor`
- Then preflight отправляет `GET /v1/tasks/anonymization_fields` с заголовками `Authorization` и `partner-id`
- And при ответе API `2xx` `doctor` продолжает выполнение без контрактной ошибки

#### Scenario: отсутствие секрета для Authorization

- Given секрет `api.partner_id` отсутствует в OS secret store
- When пользователь запускает `anonym doctor`
- Then preflight не выполняет API-check
- And возвращает детерминированную ошибку об отсутствии секрета

### Requirement: doctor валидирует UUID-формат partner id

`anonym doctor` MUST валидировать, что `api.partner_id` соответствует UUID-формату до сетевого вызова.

#### Scenario: невалидный partner id

- Given `api.partner_id` задан значением, не соответствующим UUID
- When пользователь запускает `anonym doctor`
- Then команда завершается ошибкой валидации `api.partner_id`
- And сетевой вызов к `GET /v1/tasks/anonymization_fields` не выполняется

### Requirement: preflight не использует insecure fallback

Preflight path MUST NOT использовать plaintext fallback для auth-секрета и MUST использовать только OS secret store.

#### Scenario: попытка работы без OS secret store данных

- Given в `config/env` отсутствует поддерживаемый источник секрета и в OS secret store нет записи
- When пользователь запускает `anonym doctor`
- Then команда завершается ошибкой чтения/отсутствия секрета
- And preflight не отправляет неаутентифицированный запрос в API
