# tasks: v1-i5 http debug log

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| диагностируемость API/check вызовов без утечки секретов | T1, T2 | U1, U2, I1 | debug-лог пишется в `<workspace>/.anonym/http-debug.log` по флагу |
| default security-модель не ослабляется | T2 | U3 | по умолчанию debug выключен, секреты маскируются |

## tasks

### T1. debug лог для doctor/run

- [x] добавить опциональный http debug-лог для preflight `doctor` и runtime `run`;
- [x] включение только по `ANON_HTTP_DEBUG=true`;
- [x] путь лога: `<workspace>/.anonym/http-debug.log`.

### T2. безопасный формат и документация

- [x] маскировать `Authorization` в debug-логе;
- [x] не менять поведение команд при выключенном debug;
- [x] синхронизировать `docs/user-guide.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md`.

## tests

- [x] U1: при `ANON_HTTP_DEBUG=true` создается `http-debug.log` для `doctor`.
- [x] U2: при `ANON_HTTP_DEBUG=true` создается `http-debug.log` для `run`.
- [x] U3: `Authorization` в логе маскируется, а при выключенном debug лог не создается.
