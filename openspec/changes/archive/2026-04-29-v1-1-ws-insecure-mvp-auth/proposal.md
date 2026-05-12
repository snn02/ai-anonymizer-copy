# proposal: v1-1-ws insecure mvp auth

## Why

В изолированных средах AI IDE (Codex/Claude Code/Open Code) runtime CLI не всегда имеет доступ к OS secret store хоста.
Это блокирует базовую ценность продукта: запуск анонимизации прямо из IDE-агента.

Для MVP требуется временный профиль `v1-1-ws`, который позволяет выполнять auth без keyring, но с явными ограничениями безопасности.

## What Changes

1. Добавляется явный insecure-профиль `v1-1-ws`.
2. Разрешается `api.auth_token` через env/config только в профиле `v1-1-ws`.
3. Разрешается `audit.hmac_secret` через env/config только в профиле `v1-1-ws`.
4. Сохраняются базовые guardrails `v1` (boundary/path hardening/https/allowlist/limits/audit).
5. Фиксируются non-production ограничения и компенсирующие меры.

## Контекст

В изолированных средах AI IDE (Codex/Claude Code/Open Code) runtime CLI не всегда имеет доступ к OS secret store хоста.
Это блокирует базовую ценность продукта: запуск анонимизации прямо из IDE-агента.

Для MVP требуется временный профиль `v1-1-ws`, который позволяет выполнять auth без keyring, но с явными ограничениями безопасности.

## Цель подэтапа

Снять блокер запуска в изолированных IDE, не затрагивая `v2` и MCP:
1. добавить явный insecure-профиль `v1-1-ws`;
2. разрешить auth-токен через env/config только в этом профиле;
3. разрешить `audit.hmac_secret` через env/config только в этом профиле;
4. сохранить базовые guardrails `v1` (boundary/path hardening/https/allowlist/limits/audit);
5. зафиксировать non-production ограничение и компенсирующие меры.

## In Scope

1. Контракт профиля `v1-1-ws` в runtime-конфигурации.
2. Явный флаг/параметр insecure-режима без скрытого fallback.
3. Использование `api.auth_token` и `audit.hmac_secret` из env/config в `doctor/run` только в `v1-1-ws`.
4. Проверки блокировки для production-контура.
5. Документация и тестовые сценарии с трассировкой `requirement -> task -> test -> runtime evidence`.

## Out Of Scope

1. Любые изменения `v2` и MCP-архитектуры.
2. Отмена/удаление модели OS secret store в стандартном `v1`.
3. Новые функциональные API-возможности вне auth-профиля.

## Критерии приемки

1. В стандартном профиле `v1` поведение секретов не меняется (только OS secret store).
2. В `v1-1-ws` auth-токен из env/config работает только при явном insecure-включении.
3. Без insecure-включения запуск в `v1-1-ws` блокируется с понятной ошибкой.
4. Для production-контура `v1-1-ws` блокируется policy-проверкой.
5. Тесты `p0` для `v1-1-ws` проходят.
