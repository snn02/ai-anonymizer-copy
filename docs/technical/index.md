# technical docs

## Назначение

Технический раздел фиксирует реализацию trusted sidecar CLI и связанные инженерные правила.

Этот раздел отвечает на вопросы:
- из каких модулей состоит система;
- как идут данные от `scan` до `run`;
- какие проверки безопасности обязательны;
- как тестируем и выпускаем версии.

## Порядок чтения

1. `architecture.md`
2. `components.md`
3. `data-flow.md`
4. `security.md`
5. `configuration.md`
6. `error-model.md`
7. `testing.md`
8. `release.md`
9. `v2-mcp.md`

## Границы технического раздела

- Подробные продуктовые требования: `../scope.md`, `../roadmap.md`.
- Пользовательский процесс: `../user-guide.md`, `../user-scenarios.md`.
- Версионные обязательства: `../versions/<version>/plan.md`.
