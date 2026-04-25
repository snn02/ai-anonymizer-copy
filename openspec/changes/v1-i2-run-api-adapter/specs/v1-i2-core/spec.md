# spec: v1-i2 core

## ADDED Requirements

### Requirement: run выполняет только явный запуск отправки

CLI MUST отправлять файл в API только при явной команде пользователя `anonym run <id|path>` и MUST NOT отправлять файлы в API в командах `scan/list`.

#### Scenario: успешный явный запуск по id

- Given пользователь выполнил `anonym scan` и в каталоге есть элемент со статусом `scanned`
- When пользователь запускает `anonym run <id>`
- Then CLI вызывает `POST /v1/tasks/file_anonymization` для выбранного элемента и обновляет статус элемента на `sent`

#### Scenario: неоднозначный fuzzy-поиск

- Given аргумент `run` совпадает более чем с одним элементом каталога
- When пользователь запускает `anonym run <path-fragment>`
- Then CLI не выполняет отправку и возвращает список кандидатов с требованием явного `id`

### Requirement: ответ API нормализуется adapter-слоем

Runtime MUST обрабатывать ответ `file_anonymization` через адаптер внутреннего контракта и MUST возвращать детерминированный результат для статусов и ошибок.

#### Scenario: валидный ответ анонимизации

- Given API возвращает успешный ответ по контракту `file_anonymization`
- When adapter обрабатывает ответ
- Then CLI получает нормализованный результат и продолжает запись выходного файла

#### Scenario: несовместимый формат ответа

- Given API возвращает ответ с отсутствующими обязательными полями
- When adapter обрабатывает ответ
- Then возвращается контролируемая ошибка адаптера и статус элемента фиксируется как `failed`

### Requirement: run обновляет статусы каталога предсказуемо

CLI MUST выполнять переходы статусов `scanned -> sent -> succeeded/failed` и MUST фиксировать причину ошибки при неуспехе.

#### Scenario: успешный переход до succeeded

- Given элемент каталога имеет статус `scanned`
- When `run` завершен успешно
- Then элемент проходит переходы `sent -> succeeded`

#### Scenario: ошибка API

- Given элемент каталога имеет статус `scanned`
- When API возвращает ошибку выполнения
- Then элемент получает статус `failed` и CLI отображает стандартизированную причину

### Requirement: результат сохраняется с PII-safe именем

Runtime MUST сохранять анонимизированный результат в `output_path` по техническому идентификатору и MUST NOT использовать исходное имя raw-файла в выходном имени.

#### Scenario: запись результата без утечки исходного имени

- Given анонимизация завершилась успешно
- When результат сохраняется в `output_path`
- Then имя выходного файла строится по technical id и не содержит raw filename

### Requirement: run применяет security guardrails runtime-пути

Runtime MUST применять сетевые и ресурсные ограничения (`https`, `tls verify`, allowlist host, timeout, лимиты) и MUST NOT иметь production fallback, ослабляющий проверки.

#### Scenario: попытка использовать non-https endpoint

- Given в конфигурации указан endpoint с `http://`
- When пользователь запускает `anonym run`
- Then выполнение блокируется с явной ошибкой policy violation

#### Scenario: endpoint вне allowlist

- Given endpoint не входит в allowlist
- When пользователь запускает `anonym run`
- Then выполнение блокируется до отправки запроса

#### Scenario: отключенный security fallback недопустим

- Given runtime сконфигурирован без корректных security-параметров
- When пользователь запускает `anonym run`
- Then команда завершается ошибкой и не использует `allow-all` fallback в production-пути
