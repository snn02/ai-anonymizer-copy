# ai-anonymizer

Проект для безопасной анонимизации файлов перед использованием в ИИ-IDE, чтобы в рабочем пространстве IDE были только обезличенные данные.

## Что реализуем по версиям

- `v1`: sidecar CLI (`scan`, `list`, `run`, `doctor`).
- `v2`: MCP-обертка поверх trusted-core CLI.

## Базовый принцип безопасности

- raw-файлы находятся вне workspace IDE;
- доступ к raw и сети корпоративного API есть только у sidecar CLI;
- результат анонимизации сохраняется в workspace.

## Где читать документацию

- Карта: `docs/doc-map.md`
- Scope: `docs/scope.md`
- Roadmap: `docs/roadmap.md`
- Техническая документация: `docs/technical/index.md`
- Конфигурация trusted CLI: `docs/technical/configuration.md`
- Шаблон локальной конфигурации: `config.example.yaml`
- Руководство пользователя: `docs/user-guide.md`
- Источники требований: `docs/archive/2026-04-25-ai-ide-anonymization-plan.md`, `docs/archive/2026-04-25-ai-ide-anonymization-scenarios.md`
- OpenAPI: `docs/openapi.json`
