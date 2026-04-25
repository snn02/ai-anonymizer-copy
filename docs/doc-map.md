# doc map

## Порядок чтения

1. `../readme.md`
2. `scope.md`
3. `roadmap.md`
4. `technical/index.md`
5. `user-guide.md`
6. `user-scenarios.md`
7. `test-scenarios.md`
8. `versions/<version>/plan.md`
9. `versions/<version>/openspec.md`

## Где источник истины

- Общий scope и границы: `scope.md`
- История и план версий: `roadmap.md`, `versions/<version>/plan.md`
- Пользовательские сценарии: `user-scenarios.md`
- Сценарии тестирования: `test-scenarios.md`
- Руководство пользователя: `user-guide.md`
- Техническая архитектура и инженерные контракты: `technical/*`
- Ключевые решения: `decisions.md`
- Накопительные уроки: `lessons-learned.md`
- OpenSpec-рабочая зона: `../openspec/changes/`
- OpenSpec-архив: `../openspec/archive/`

## Правила связей

- `versions/<version>/plan.md` ссылается на релевантные сценарии.
- `versions/<version>/action-log.md` хранит изменения, ошибки и реакции версии.
- `technical/index.md` ссылается на обязательные технические документы реализации.
- После закрытия версии выводы переносятся в `lessons-learned.md`.
- Для `v1` feature-пакеты OpenSpec создаются по итерациям в `../openspec/changes/` согласно согласованному workflow.
- Для `v2` создание feature-пакетов открывается отдельным решением.
