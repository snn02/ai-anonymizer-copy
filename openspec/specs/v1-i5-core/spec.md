# v1-i5-core Specification

## Purpose
TBD - created by archiving change v1-i5-http-debug-log. Update Purpose after archive.
## Requirements
### Requirement: опциональный http debug-лог для doctor и run

Runtime MUST поддерживать опциональный debug-лог HTTP-вызовов для команд `doctor` и `run`.

#### Scenario: debug включен

- Given задан `ANON_HTTP_DEBUG=true`
- When пользователь запускает `anonym doctor` или `anonym run <id>`
- Then runtime записывает запись в `<workspace>/.anonym/http-debug.log`

### Requirement: debug-лог не раскрывает секреты и не влияет на default-путь

Runtime MUST маскировать чувствительные заголовки в debug-логе и MUST NOT изменять default-поведение команд при выключенном debug.

#### Scenario: security of debug log

- Given `Authorization` содержит секрет
- When runtime пишет debug-лог
- Then значение `Authorization` маскировано
- And при выключенном debug-флаге лог не создается

