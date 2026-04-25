# doc map

## Порядок чтения

1. `../readme.md`
2. `scope.md`
3. `roadmap.md`
4. `user-guide.md`
5. `user-scenarios.md`
6. `test-scenarios.md`
7. `versions/<version>/plan.md`
8. `versions/<version>/openspec.md`

## Где источник истины

- Общий scope и границы: `scope.md`
- История и план версий: `roadmap.md`, `versions/<version>/plan.md`
- Пользовательские сценарии: `user-scenarios.md`
- Сценарии тестирования: `test-scenarios.md`
- Руководство пользователя: `user-guide.md`
- Ключевые решения: `decisions.md`
- Накопительные уроки: `lessons-learned.md`
- OpenSpec-рабочая зона: `../openspec/changes/`
- OpenSpec-архив: `../openspec/archive/`

## Правила связей

- `versions/<version>/plan.md` ссылается на релевантные сценарии.
- `versions/<version>/action-log.md` хранит изменения, ошибки и реакции версии.
- После закрытия версии выводы переносятся в `lessons-learned.md`.
- Feature-пакеты OpenSpec создаются только после согласования workflow разработки.
