## MODIFIED Requirements

### Requirement: auth token в v1-1-ws используется только при explicit enable
Runtime MUST использовать `api.auth_token` и `audit.hmac_secret` из env/config только в профиле `v1-1-ws` и только при явном insecure-включении.

#### Scenario: doctor/run в v1-1-ws с токеном
- Given профиль `v1-1-ws` активен
- And `api.auth_token` задан в env/config
- And `audit.hmac_secret` задан в env/config
- And `api.user_id` задан валидным UUID
- When пользователь запускает `anonym doctor` и `anonym run <id|path>`
- Then API auth выполняется через `api.auth_token`
- And audit HMAC выполняется через `audit.hmac_secret`
- And `run` отправляет header `user-id` в `POST /v1/tasks/file_anonymization`
