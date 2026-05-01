# roadmap

## Правило

`roadmap.md` хранит только статус версий и ссылки. Детали требований и шагов живут в `docs/versions/*` и `docs/versions/*/action-log.md`.

## v1

- Цель: безопасный sidecar CLI для анонимизации.
- Статус: in progress (`mode=prod|mvp` migration).
- Итерации: `v1-i1`, `v1-i2`, `v1-i3`, `v1-i4` завершены; в работе post-release итерации `v1-i7`, `v1-i8`, `v1-i9`.
- Режим `mvp`: временный MVP-режим для запуска в изолированных AI IDE без OS secret store (с компенсирующими ограничениями и явным warning).
- План версии: `versions/v1/plan.md`.
- OpenSpec версии: `versions/v1/openspec.md`.
- Лог прогресса: `versions/v1/action-log.md`.

## v2

- Цель: MCP-обертка над trusted-core CLI.
- Статус: reserved (не стартовала).
- План версии: `versions/v2/plan.md`.
- OpenSpec версии: `versions/v2/openspec.md`.
- Лог прогресса: `versions/v2/action-log.md`.

## Примечание по OpenSpec

- Для `v1` feature-пакеты создаются по итерациям в `../openspec/changes/`.
- Для `v2` создание feature-пакетов открывается отдельным решением.
