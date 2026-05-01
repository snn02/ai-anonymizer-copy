## ADDED Requirements

### Requirement: CLI override минимизирован и унифицирован
Система MUST поддерживать единый минимальный набор CLI override для команд `doctor`, `scan`, `list`, `run`: `--config`, `--mode`, `--workspace-path`, `--output-path`.

#### Scenario: запуск с минимальным override
- **WHEN** пользователь запускает любую из команд с `--config` и без других override-флагов
- **THEN** runtime успешно читает остальные параметры из `config/env`

#### Scenario: override путей для отдельного запуска
- **WHEN** пользователь передает `--workspace-path` и/или `--output-path`
- **THEN** runtime применяет эти значения как временный override поверх `config/env`

### Requirement: остальные параметры runtime берутся из config/env
Параметры API, лимитов и security (включая `allowed_hosts`, `max_file_size_mb`, `max_pages`, `request_timeout_sec`, `api.user_id`, `api.fields_page`, `api.fields_per_page`) MUST NOT иметь CLI override и MUST читаться из `config/env`.

#### Scenario: попытка использовать удаленный override-флаг
- **WHEN** пользователь передает удаленный флаг (например, `--allowed-hosts`)
- **THEN** CLI возвращает явную ошибку неизвестного флага

### Requirement: режим запуска определяется одним параметром mode
Система MUST определять runtime-поведение только по `mode` со значениями `prod|mvp`.

#### Scenario: mode=prod
- **WHEN** `mode=prod`
- **THEN** секреты API и audit читаются только из OS Secret Store

#### Scenario: mode=mvp
- **WHEN** `mode=mvp`
- **THEN** секреты API и audit читаются из `config/env`, и CLI выводит предупреждение о сниженных security-гарантиях
