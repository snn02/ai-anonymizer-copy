## MODIFIED Requirements

### Requirement: runtime поддерживает режим mvp для изолированных IDE
Runtime MUST поддерживать режим `mvp` как эквивалент временного insecure-профиля для запуска без доступа к host OS secret store.

#### Scenario: режим mvp включен явно
- Given пользователь запускает CLI в изолированной AI IDE
- And в конфигурации выбран `mode=mvp`
- When пользователь выполняет `anonym doctor`
- Then runtime использует контракт `mvp` и не требует host keyring

### Requirement: auth token в mvp используется только при explicit enable
Runtime MUST использовать `api.auth_token` и `audit.hmac_secret` из env/config только в режиме `mvp` и MUST NOT использовать эти значения в режиме `prod`.

#### Scenario: doctor/run в mvp с токеном
- Given режим `mvp` активен
- And `api.auth_token` задан в env/config
- And `audit.hmac_secret` задан в env/config
- When пользователь запускает `anonym doctor` и `anonym run <id|path>`
- Then API auth выполняется через `api.auth_token`
- And audit HMAC выполняется через `audit.hmac_secret`

#### Scenario: режим prod не принимает token fallback
- Given активен режим `prod`
- And в env/config задан `api.auth_token`
- When пользователь запускает `anonym run <id|path>`
- Then runtime игнорирует token fallback и требует секреты из OS secret store

### Requirement: mvp допускается в production API-контуре с предупреждением
Режим `mvp` MUST допускать запуск в production API-контуре при явном включении режима и MUST предупреждать пользователя о снижении security-гарантий.

#### Scenario: mvp с production endpoint
- Given режим `mvp` активен
- And `api.base_url` указывает на production-контур
- When пользователь запускает `anonym doctor`
- Then preflight не блокируется только по признаку production-host
- And CLI выводит предупреждение о временном insecure-режиме
