# v1 openspec

## Статус

`v1` работает в post-release режиме: активна миграция к mode-first модели (`prod|mvp`).

## Где лежат артефакты

- Активные изменения: `../../../openspec/changes/`
- Закрытые изменения: `../../../openspec/archive/`

## Активные feature-пакеты v1

- `v1-i7-doctor-userid-pagination`:
  - `../../../openspec/changes/v1-i7-doctor-userid-pagination/proposal.md`
  - `../../../openspec/changes/v1-i7-doctor-userid-pagination/tasks.md`
  - `../../../openspec/changes/v1-i7-doctor-userid-pagination/specs/v1-i7-core/spec.md`
- `v1-i8-mvp-config-runtime`:
  - `../../../openspec/changes/v1-i8-mvp-config-runtime/proposal.md`
  - `../../../openspec/changes/v1-i8-mvp-config-runtime/design.md`
  - `../../../openspec/changes/v1-i8-mvp-config-runtime/tasks.md`
  - `../../../openspec/changes/v1-i8-mvp-config-runtime/specs/v1-i8-mvp-config-runtime/spec.md`
  - `../../../openspec/changes/v1-i8-mvp-config-runtime/specs/v1-1-ws-core/spec.md`

## Закрытые feature-пакеты v1

- `v1-i1-foundation-preflight-discovery`:
  - `../../../openspec/archive/v1-i1-foundation-preflight-discovery/proposal.md`
  - `../../../openspec/archive/v1-i1-foundation-preflight-discovery/tasks.md`
  - `../../../openspec/archive/v1-i1-foundation-preflight-discovery/specs/v1-i1-core/spec.md`
- `v1-i2-run-api-adapter`:
  - `../../../openspec/archive/v1-i2-run-api-adapter/proposal.md`
  - `../../../openspec/archive/v1-i2-run-api-adapter/tasks.md`
  - `../../../openspec/archive/v1-i2-run-api-adapter/specs/v1-i2-core/spec.md`
- `v1-i3-doctor-api-header-alignment`:
  - `../../../openspec/archive/v1-i3-doctor-api-header-alignment/proposal.md`
  - `../../../openspec/archive/v1-i3-doctor-api-header-alignment/tasks.md`
  - `../../../openspec/archive/v1-i3-doctor-api-header-alignment/specs/v1-i3-core/spec.md`
- `v1-i4-config-first-loading`:
  - `../../../openspec/archive/v1-i4-config-first-loading/proposal.md`
  - `../../../openspec/archive/v1-i4-config-first-loading/tasks.md`
  - `../../../openspec/archive/v1-i4-config-first-loading/specs/v1-i4-core/spec.md`
- `v1-1-ws-insecure-mvp-auth`:
  - `../../../openspec/changes/archive/2026-04-29-v1-1-ws-insecure-mvp-auth/proposal.md`
  - `../../../openspec/changes/archive/2026-04-29-v1-1-ws-insecure-mvp-auth/tasks.md`
  - `../../../openspec/specs/v1-1-ws-core/spec.md`
- `v1-i5-http-debug-log`:
  - `../../../openspec/changes/archive/2026-04-29-v1-i5-http-debug-log/proposal.md`
  - `../../../openspec/changes/archive/2026-04-29-v1-i5-http-debug-log/tasks.md`
  - `../../../openspec/specs/v1-i5-core/spec.md`
- `v1-i9-cli-mode-and-overrides-simplification`:
  - `../../../openspec/changes/archive/2026-05-01-v1-i9-cli-mode-and-overrides-simplification/proposal.md`
  - `../../../openspec/changes/archive/2026-05-01-v1-i9-cli-mode-and-overrides-simplification/design.md`
  - `../../../openspec/changes/archive/2026-05-01-v1-i9-cli-mode-and-overrides-simplification/tasks.md`
  - `../../../openspec/changes/archive/2026-05-01-v1-i9-cli-mode-and-overrides-simplification/specs/v1-i9-cli-mode-and-overrides-simplification/spec.md`
  - `../../../openspec/changes/archive/2026-05-01-v1-i9-cli-mode-and-overrides-simplification/specs/v1-1-ws-core/spec.md`

## Подготовленные следующие итерации

- после завершения `v1-i9` отдельные post-release итерации будут определяться по итогам стабилизации mode-first контракта; `v2` остается зарезервированной под MCP-этап.

## Правило

- Одна итерация `v1` = один feature-пакет OpenSpec.
- Feature-пакет должен описывать законченный и тестируемый результат итерации.
- После завершения итерации фиксируются результаты в `action-log.md` и планируется следующая итерация.

## Принципы планирования (обязательно для новых итераций)

1. Трассировка требований:
   - каждое требование из `docs/versions/v1/plan.md`, попадающее в итерацию, должно иметь явную связку в OpenSpec:
     - `requirement` -> `task` -> `test` -> `runtime evidence`.
2. Security-DoD:
   - для каждой security-задачи в `tasks.md` фиксируется DoD с проверяемыми пунктами:
     - runtime-реализация в production-пути;
     - минимум один негативный тест;
     - отсутствие небезопасного fallback.
3. Политика fallback:
   - default-заглушки вида `allow-all` допустимы только в тестах или под явным non-production feature-flag;
   - закрывать security-задачу при активном production fallback запрещено.
4. Гейт закрытия итерации:
   - перед переводом чекбокса в `[x]` выполняется проверка соответствия `plan.md` по всем пунктам, попавшим в scope итерации;
   - результаты проверки фиксируются в `docs/versions/v1/action-log.md`.

## Чек-лист планирования

- добавить в `tasks.md` явный трек security-подзадач с DoD;
- для каждой security-подзадачи предусмотреть отрицательные тест-кейсы;
- зафиксировать, какие пункты из `plan.md` остаются незакрытыми и где они будут закрыты;
- не отмечать security-задачу выполненной без runtime-подтверждения.
