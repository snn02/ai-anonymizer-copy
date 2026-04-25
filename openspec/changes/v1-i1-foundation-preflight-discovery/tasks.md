# tasks: v1-i1 foundation preflight discovery

## progress

- 2026-04-25: стартована реализация итерации, добавлен базовый каркас trusted CLI на Go (`go.mod`, `cmd/anonym`, `internal/config`, `internal/preflight`) и тесты preflight.
- 2026-04-25: в preflight закрыт baseline (`boundary`, `https`, host allowlist, `max_parallel_runs=1`); проверки секретов и активной доступности API остаются в работе.
- 2026-04-25: реализованы `scan/list` с локальным каталогом очереди (`.anonym/catalog.json`) и безопасным discovery; добавлены unit/integration тесты CLI и модулей `catalog/scanner`.

## 1. doctor preflight

- [ ] Реализовать обязательные preflight-проверки конфигурации, boundary/ACL, секретов и доступности API.
- [x] Добавить негативные кейсы для блокирующих ошибок preflight.
- [x] Обновить пользовательскую документацию по диагностике `doctor`.

## 2. scan safe discovery

- [x] Реализовать безопасное обнаружение файлов в пределах `raw_path`.
- [x] Добавить проверки path hardening: traversal, symlink/reparse, ADS, директории вместо файлов.
- [ ] Ввести лимиты, влияющие на discovery в рамках v1-требований.

## 3. list локальная очередь

- [x] Реализовать локальное состояние очереди и вывод `list` с техническим `id`.
- [x] Подтвердить, что `scan/list` не выполняют отправку файлов в API.
- [x] Добавить тесты на стабильность формата вывода и статусов итерации.

## 4. проверка и фиксация

- [x] Прогнать `p0`-тесты, относящиеся к `doctor`, boundary и discovery.
- [ ] Зафиксировать результаты в `docs/versions/v1/action-log.md`.
- [ ] Подготовить предложение на итерацию 2 (`run` + адаптер ответа API).
