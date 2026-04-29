# spec: v1-i4 core

## Added Requirements

### Requirement: runtime загружает параметры из config файла

`anonym doctor/scan/list/run` MUST использовать значения из `config.yaml` как базовый источник runtime-параметров, если для поля не задан override в env или CLI.

#### Scenario: запуск doctor с заполненным config

- Given в `config.yaml` заданы валидные обязательные параметры runtime
- And в OS secret store доступны необходимые секреты
- When пользователь запускает `anonym doctor --config <path>`
- Then команда использует параметры из указанного config файла
- And preflight проходит при корректной конфигурации

#### Scenario: config файл отсутствует или недоступен

- Given путь `--config` указывает на отсутствующий или недоступный файл
- When пользователь запускает `anonym doctor --config <path>`
- Then команда завершается явной ошибкой загрузки конфигурации
- And runtime-операции не выполняются

### Requirement: приоритет источников конфигурации детерминирован

Runtime MUST применять значения по приоритету `CLI > env > config file`.

#### Scenario: env переопределяет значение из config файла

- Given в `config.yaml` задан `api.base_url=A`
- And в env задан `ANON_API_BASE_URL=B`
- When пользователь запускает команду без CLI-override `api-base-url`
- Then runtime использует `api.base_url=B`

#### Scenario: cli переопределяет env и config файл

- Given в `config.yaml` задан `api.base_url=A`
- And в env задан `ANON_API_BASE_URL=B`
- When пользователь запускает команду с `--api-base-url=C`
- Then runtime использует `api.base_url=C`

### Requirement: compact runtime calls with config first

При заполненном `config.yaml` runtime MUST поддерживать компактные вызовы без дублирования полного набора параметров.

#### Scenario: compact scan list run

- Given `config.yaml` содержит обязательные runtime-параметры
- And preflight успешно выполнен
- When пользователь запускает `anonym scan --config <path>`, затем `anonym list --config <path>`, затем `anonym run <id> --config <path>`
- Then все команды выполняются с параметрами из config файла
- And override применяется только к явно переопределенным полям
