# spec: v1-i1 core

## ADDED Requirements

### Requirement: doctor выполняет блокирующий preflight

CLI MUST выполнять preflight до рабочих операций и MUST блокировать выполнение при нарушении security-ограничений.

#### Scenario: корректное окружение

- Given заданы корректные `raw_path`, `output_path`, секреты и доступ к API
- When пользователь запускает `anonym doctor`
- Then CLI возвращает успешный результат preflight без предупреждений уровня block

#### Scenario: нарушение boundary/acl

- Given `raw_path` не проходит policy boundary или ACL-проверки
- When пользователь запускает `anonym doctor`
- Then CLI завершает команду с блокирующей ошибкой и рекомендацией по исправлению

### Requirement: scan выполняет только безопасное обнаружение

CLI MUST обрабатывать только файлы внутри `raw_path` и MUST отклонять небезопасные пути/объекты.

#### Scenario: безопасное обнаружение файлов

- Given в `raw_path` есть допустимые файлы поддерживаемых форматов
- When пользователь запускает `anonym scan`
- Then файлы добавляются в локальную очередь со статусом итерации

#### Scenario: попытка обхода path policy

- Given переданы пути с `..`, symlink/reparse или ADS-паттернами
- When выполняется обнаружение
- Then такие пути отклоняются и отражаются как ошибки в диагностике scan

### Requirement: list отображает локальную очередь без отправки в API

CLI MUST показывать локальное состояние и MUST NOT инициировать отправку в API при `scan/list`.

#### Scenario: просмотр очереди после scan

- Given пользователь выполнил `anonym scan`
- When пользователь выполняет `anonym list`
- Then CLI показывает элементы очереди с техническими идентификаторами и статусами
- And в процессе `scan/list` отсутствует вызов `POST /v1/tasks/file_anonymization`
