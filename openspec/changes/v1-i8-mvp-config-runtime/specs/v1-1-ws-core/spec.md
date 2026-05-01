## MODIFIED Requirements

### Requirement: auth token в v1-1-ws используется только при explicit enable
Runtime MUST использовать `api.auth_token` и `audit.hmac_secret` из env/config/CLI только в профиле `v1-1-ws` и только при явном insecure-включении.

#### Scenario: doctor/run в v1-1-ws с токеном
- Given профиль `v1-1-ws` активен
- And `api.auth_token` задан в env/config/CLI
- And `audit.hmac_secret` задан в env/config/CLI
- And `insecure_no_secrets=true`
- When пользователь запускает `anonym doctor` и `anonym run <id|path>`
- Then API auth выполняется через `api.auth_token`
- And audit HMAC выполняется через `audit.hmac_secret`
- And OS Secret Store не используется

#### Scenario: токен задан, но insecure не включен
- Given профиль `v1-1-ws` активен
- And `api.auth_token` задан
- And insecure-включение отсутствует
- When пользователь запускает `anonym run <id|path>`
- Then команда блокируется с явной ошибкой policy violation
