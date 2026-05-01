## ADDED Requirements

### Requirement: run поддерживает header user-id для file_anonymization
Команда `run` SHALL отправлять `user-id` в `POST /v1/tasks/file_anonymization`, если `api.user_id` задан в runtime-конфигурации.

#### Scenario: user-id задан и передан в run-запросе
- **WHEN** пользователь запускает `anonym run <id|path>` с заданным `api.user_id`
- **THEN** исходящий `POST /v1/tasks/file_anonymization` содержит header `user-id` с этим значением

#### Scenario: user-id не задан
- **WHEN** пользователь запускает `anonym run <id|path>` без `api.user_id`
- **THEN** запрос выполняется без header `user-id` (backward-compatible)

### Requirement: preflight/runtime валидируют формат user-id
Система SHALL отклонять запуск при невалидном `api.user_id` (не UUID) до отправки API-запроса.

#### Scenario: невалидный user-id
- **WHEN** в конфигурации указан `api.user_id`, не соответствующий UUID
- **THEN** preflight/runtime возвращает явную ошибку валидации и не отправляет запрос в API
