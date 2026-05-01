## ADDED Requirements

### Requirement: MVP runtime MUST работать из config без OS Secret Store
В профиле `v1-1-ws` при `insecure_no_secrets=true` система MUST использовать `api.auth_token` и `audit.hmac_secret` из источников конфигурации (`CLI > env > config`) и MUST NOT обращаться к OS Secret Store.

#### Scenario: run использует секреты из config
- **WHEN** пользователь запускает `anonym run <id|path>` с `runtime_profile=v1-1-ws`, `insecure_no_secrets=true`, `api.auth_token` и `audit.hmac_secret` из config
- **THEN** команда выполняет API-вызов и аудит без чтения секретов из OS Secret Store

#### Scenario: приоритет CLI/env/config применяется единообразно
- **WHEN** значения `runtime_profile`, `insecure_no_secrets`, `api.auth_token`, `audit.hmac_secret` заданы одновременно в CLI, env и config
- **THEN** runtime использует приоритет `CLI > env > config > default` для всех перечисленных параметров
