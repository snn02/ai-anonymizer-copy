# v1 openspec

## Статус

`v1` завершена. Post-release итерация `v1-i3` закрыта и перенесена в архив.

## Где лежат артефакты

- Активные изменения: `../../../openspec/changes/`
- Закрытые изменения: `../../../openspec/archive/`

## Активные feature-пакеты v1

- нет.

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

## Подготовленные следующие итерации

- нет; планирование следующих итераций ведется в рамках `v2`.

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

## Чек-лист планирования v1-i2

- добавить в `tasks.md` явный трек security-подзадач с DoD;
- для каждой security-подзадачи предусмотреть отрицательные тест-кейсы;
- зафиксировать, какие пункты из `plan.md` остаются незакрытыми и где они будут закрыты;
- не отмечать security-задачу выполненной без runtime-подтверждения.
