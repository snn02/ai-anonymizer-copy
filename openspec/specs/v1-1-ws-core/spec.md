# v1-1-ws-core Specification

## Purpose
TBD - created by archiving change v1-1-ws-insecure-mvp-auth. Update Purpose after archive.
## Requirements
### Requirement: runtime поддерживает профиль v1-1-ws для изолированных IDE

Runtime MUST поддерживать профиль `v1-1-ws`, который разрешает запуск без доступа к host OS secret store.

#### Scenario: профиль v1-1-ws включен явно

- Given пользователь запускает CLI в изолированной AI IDE
- And в конфигурации выбран профиль `v1-1-ws`
- And включен явный insecure-переключатель
- When пользователь выполняет `anonym doctor`
- Then runtime использует контракт `v1-1-ws` и не требует host keyring

### Requirement: auth token в v1-1-ws используется только при explicit enable

Runtime MUST использовать `api.auth_token` и `audit.hmac_secret` из env/config только в профиле `v1-1-ws` и только при явном insecure-включении.

#### Scenario: doctor/run в v1-1-ws с токеном

- Given профиль `v1-1-ws` активен
- And `api.auth_token` задан в env/config
- And `audit.hmac_secret` задан в env/config
- When пользователь запускает `anonym doctor` и `anonym run <id>`
- Then API auth выполняется через `api.auth_token`
- And audit HMAC выполняется через `audit.hmac_secret`

#### Scenario: токен задан, но insecure не включен

- Given профиль `v1-1-ws` активен
- And `api.auth_token` задан
- And insecure-включение отсутствует
- When пользователь запускает `anonym run <id>`
- Then команда блокируется с явной ошибкой policy violation

### Requirement: стандартный v1 не ослабляется

Стандартный профиль `v1` MUST продолжать использовать только OS secret store и MUST NOT использовать `api.auth_token` из env/config.

#### Scenario: token fallback в standard v1 запрещен

- Given активен профиль `v1`
- And в env задан `api.auth_token`
- When пользователь запускает `anonym run <id>`
- Then runtime игнорирует token fallback и требует секрет из OS secret store

### Requirement: v1-1-ws ограничен non-production контуром

Профиль `v1-1-ws` MUST быть ограничен non-production использованием и MUST предупреждать пользователя о снижении security-гарантий.

#### Scenario: попытка использовать v1-1-ws для production endpoint

- Given профиль `v1-1-ws` активен
- And `api.base_url` указывает на production-контур
- When пользователь запускает `anonym doctor`
- Then команда блокируется policy-проверкой

#### Scenario: предупреждение пользователя при v1-1-ws

- Given профиль `v1-1-ws` активен
- When пользователь запускает `anonym doctor` или `anonym run`
- Then CLI выводит явное предупреждение о временном insecure-режиме

