# tasks: v1-i4 config first loading

## requirement trace

| requirement (v1 plan) | task | test | runtime evidence |
|---|---|---|---|
| `config.yaml` используется как базовый источник параметров runtime | T1, T2 | U1, U2, I1 | `internal/config.Load` читает YAML и применяет значения при отсутствии override |
| приоритет источников `CLI > env > config file` | T2, T3 | U3, U4, I2 | runtime берет значения из флагов, затем env, затем файла |
| компактные вызовы `doctor/scan/list/run` с `--config` | T1, T4 | I1, I3 | команды работают с минимальным набором флагов при заполненном конфиге |
| явные ошибки чтения/парсинга конфига | T1, T3 | U5, I4 | runtime возвращает диагностируемую ошибку загрузки конфига |

## security dod (для T1-T4)

- [ ] runtime реализован в production-пути `doctor/scan/list/run`;
- [ ] минимум один негативный тест на каждый security-check;
- [ ] отсутствует insecure fallback для секретов и API-заголовков;
- [ ] evidence зафиксирован прогоном `go test ./...`.

## tasks

### T1. yaml loading in internal config

- [ ] реализовать чтение `config.yaml` в `internal/config` с поддержкой пути из `--config`;
- [ ] добавить маппинг полей `paths/api/limits/security/audit` в runtime-структуру;
- [ ] сохранить текущую валидацию обязательных полей и безопасные дефолты `v1`.

### T2. precedence and merge rules

- [ ] реализовать merge-логику, где приоритет строго `CLI > env > config file`;
- [ ] зафиксировать единые правила для scalar/list параметров (`allowed_hosts`);
- [ ] убедиться, что источники секретов не меняются (только OS secret store).

### T3. tests for config-first runtime behavior

- [ ] добавить/обновить unit-тесты на чтение YAML и fallback-поведение;
- [ ] добавить/обновить unit-тесты на приоритет `CLI > env > config file`;
- [ ] добавить/обновить негативные тесты на битый YAML и отсутствующий конфиг.

### T4. docs synchronization for config-first

- [ ] обновить `docs/user-guide.md`: config-first как основной сценарий, override через env/CLI как отдельные кейсы;
- [ ] привести заголовки и подзаголовки user-guide к написанию с заглавной буквы;
- [ ] синхронизировать `docs/technical/configuration.md`, `docs/user-scenarios.md`, `docs/test-scenarios.md`, `docs/versions/v1/action-log.md`.

## tests

### unit (U*)

- [ ] U1: значения из валидного `config.yaml` попадают в `Config`.
- [ ] U2: при отсутствии `env/CLI` runtime использует значения из файла.
- [ ] U3: env переопределяет значения из `config.yaml`.
- [ ] U4: CLI переопределяет значения из env и `config.yaml`.
- [ ] U5: невалидный YAML возвращает явную ошибку загрузки.

### integration (I*)

- [ ] I1: `anonym doctor --config <path>` проходит preflight при заполненном конфиге и секретах.
- [ ] I2: при конфликте источников runtime берет значения строго по приоритету `CLI > env > config file`.
- [ ] I3: компактные вызовы `scan/list/run` с `--config` работают без дублирования полного набора флагов.
- [ ] I4: при отсутствии/битом конфиге команда завершается диагностируемой ошибкой до выполнения runtime-операции.

### regression

- [ ] R1: `go test ./...` зеленый после изменений.
