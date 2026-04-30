# Р СѓРєРѕРІРѕРґСЃС‚РІРѕ РџРѕР»СЊР·РѕРІР°С‚РµР»СЏ

## Р”Р»СЏ РљРѕРіРѕ Р­С‚Рѕ Р СѓРєРѕРІРѕРґСЃС‚РІРѕ

Р”РѕРєСѓРјРµРЅС‚ РґР»СЏ РїРѕР»СЊР·РѕРІР°С‚РµР»СЏ AI IDE, РєРѕС‚РѕСЂС‹Р№ С…РѕС‡РµС‚ Р·Р°РїСѓСЃС‚РёС‚СЊ `anonym` Рё РїРѕР»СѓС‡РёС‚СЊ СЂРµР·СѓР»СЊС‚Р°С‚ СЃ РїРµСЂРІРѕРіРѕ СЂР°Р·Р°.

РџРѕРґРґРµСЂР¶Р°РЅС‹ РґРІР° РІР°СЂРёР°РЅС‚Р° Р·Р°РїСѓСЃРєР°:
1. Сѓ РІР°СЃ РµСЃС‚СЊ РіРѕС‚РѕРІС‹Р№ `anonym.exe`;
2. РІС‹ Р·Р°РїСѓСЃРєР°РµС‚Рµ РёР· РёСЃС…РѕРґРЅРёРєРѕРІ С‡РµСЂРµР· `go run`.

## Р§С‚Рѕ Р”РµР»Р°РµС‚ РРЅСЃС‚СЂСѓРјРµРЅС‚

`anonym` вЂ” CLI РґР»СЏ Р±РµР·РѕРїР°СЃРЅРѕР№ Р°РЅРѕРЅРёРјРёР·Р°С†РёРё С„Р°Р№Р»РѕРІ:
1. `doctor` вЂ” РїСЂРѕРІРµСЂСЏРµС‚ РєРѕРЅС„РёРіСѓСЂР°С†РёСЋ, РґРѕСЃС‚СѓРї Рє API Рё СЃРµРєСЂРµС‚С‹;
2. `scan` вЂ” РёРЅРґРµРєСЃРёСЂСѓРµС‚ С„Р°Р№Р»С‹ РёР· `raw_path` РІ Р»РѕРєР°Р»СЊРЅС‹Р№ РєР°С‚Р°Р»РѕРі;
3. `list` вЂ” РїРѕРєР°Р·С‹РІР°РµС‚ РѕС‡РµСЂРµРґСЊ (`id`, РїСѓС‚СЊ, СЃС‚Р°С‚СѓСЃ);
4. `run <id|path>` вЂ” РѕС‚РїСЂР°РІР»СЏРµС‚ РІС‹Р±СЂР°РЅРЅС‹Р№ С„Р°Р№Р» РІ API Рё СЃРѕС…СЂР°РЅСЏРµС‚ СЂРµР·СѓР»СЊС‚Р°С‚ РІ `output_path`.

## РћР±СЏР·Р°С‚РµР»СЊРЅР°СЏ РњРѕРґРµР»СЊ РџР°РїРѕРє

РџРµСЂРµРґ РїРµСЂРІС‹Рј Р·Р°РїСѓСЃРєРѕРј РїРѕРґРіРѕС‚РѕРІСЊС‚Рµ:
1. `raw_path` вЂ” РїР°РїРєР° СЃ РёСЃС…РѕРґРЅС‹РјРё С„Р°Р№Р»Р°РјРё, РѕР±СЏР·Р°С‚РµР»СЊРЅРѕ РІРЅРµ IDE workspace;
2. `workspace_path` вЂ” РєРѕСЂРµРЅСЊ РїСЂРѕРµРєС‚Р° РІ AI IDE;
3. `output_path` вЂ” РїР°РїРєР° СЂРµР·СѓР»СЊС‚Р°С‚РѕРІ, РѕР±СЏР·Р°С‚РµР»СЊРЅРѕ РІРЅСѓС‚СЂРё `workspace_path`.

РџСЂРёРјРµСЂ:
1. `D:\secure-raw` (`raw_path`)
2. `C:\work\my-ai-project` (`workspace_path`)
3. `C:\work\my-ai-project\anonymized` (`output_path`)

## РџСЂРёРѕСЂРёС‚РµС‚ РСЃС‚РѕС‡РЅРёРєРѕРІ РљРѕРЅС„РёРіСѓСЂР°С†РёРё

Runtime РїСЂРёРјРµРЅСЏРµС‚ Р·РЅР°С‡РµРЅРёСЏ РІ РїРѕСЂСЏРґРєРµ:
1. CLI-С„Р»Р°РіРё (СЃР°РјС‹Р№ РІС‹СЃРѕРєРёР№ РїСЂРёРѕСЂРёС‚РµС‚);
2. РџРµСЂРµРјРµРЅРЅС‹Рµ РѕРєСЂСѓР¶РµРЅРёСЏ;
3. `config.yaml`.

РџСЂР°РєС‚РёРєР° РґР»СЏ `v1`:
1. РѕСЃРЅРѕРІРЅРѕР№ СЃС†РµРЅР°СЂРёР№ вЂ” Р·Р°РїРѕР»РЅРёС‚СЊ `config.yaml` Рё Р·Р°РїСѓСЃРєР°С‚СЊ РєРѕРјРїР°РєС‚РЅС‹РјРё РєРѕРјР°РЅРґР°РјРё;
2. env/CLI РёСЃРїРѕР»СЊР·РѕРІР°С‚СЊ РєР°Рє С‚РѕС‡РµС‡РЅС‹Р№ override Р±РµР· РєРѕРїРёСЂРѕРІР°РЅРёСЏ РІСЃРµС… РїР°СЂР°РјРµС‚СЂРѕРІ.

## РћР±СЏР·Р°С‚РµР»СЊРЅС‹Рµ РџР°СЂР°РјРµС‚СЂС‹ Р РЎРµРєСЂРµС‚С‹

### РџР°СЂР°РјРµС‚СЂС‹ Runtime V1

РћР±СЏР·Р°С‚РµР»СЊРЅС‹Рµ РїРѕР»СЏ:
1. `paths.raw_path`
2. `paths.output_path`
3. `paths.workspace_path`
4. `api.base_url` (С‚РѕР»СЊРєРѕ `https`)
5. `api.partner_id` (СЃС‚СЂРѕРіРѕ UUID)
6. `security.allowed_hosts`
7. `limits.max_parallel_runs=1`

Р РµРєРѕРјРµРЅРґСѓРµРјС‹Рµ Р»РёРјРёС‚С‹:
1. `limits.max_file_size_mb=25`
2. `limits.max_pages=300`
3. `limits.request_timeout_sec=5`

### РЎРµРєСЂРµС‚С‹

РЎРµРєСЂРµС‚С‹ С…СЂР°РЅСЏС‚СЃСЏ С‚РѕР»СЊРєРѕ РІ OS secret store.
РќРµР»СЊР·СЏ С…СЂР°РЅРёС‚СЊ СЃРµРєСЂРµС‚С‹ РІ `config.yaml` Рё `.env`.

### РљР°Рє РЎРѕРїРѕСЃС‚Р°РІРёС‚СЊ Р”Р°РЅРЅС‹Рµ РћС‚ Р Р°Р·СЂР°Р±РѕС‚С‡РёРєР° API

Р•СЃР»Рё РІР°Рј РїРµСЂРµРґР°Р»Рё:
1. `РєР»СЋС‡` вЂ” СЌС‚Рѕ Р·РЅР°С‡РµРЅРёРµ Р·Р°РіРѕР»РѕРІРєР° `Authorization` (РєР»Р°РґРµС‚СЃСЏ РІ keyring);
2. `РїР°СЂС‚РЅРµСЂ` вЂ” СЌС‚Рѕ `partner-id` Рё `api.partner_id` (UUID);
3. `РїРѕР»СЊР·РѕРІР°С‚РµР»СЊ` вЂ” СЌС‚Рѕ `user-id` (РѕРїС†РёРѕРЅР°Р»СЊРЅС‹Р№ UUID-Р·Р°РіРѕР»РѕРІРѕРє, РЅРµ РёСЃРїРѕР»СЊР·СѓРµС‚СЃСЏ runtime `v1`).

### Р§С‚Рѕ Р’Р°Р¶РЅРѕ Р”Р»СЏ Runtime V1

1. CLI РѕС‚РїСЂР°РІР»СЏРµС‚ `Authorization`, `partner-id`, `fields=anonymizer`.
2. `user-id`, `tags`, `tags-numeration` РЅРµ РїР°СЂР°РјРµС‚СЂРёР·СѓСЋС‚СЃСЏ С‡РµСЂРµР· CLI/env/config.

## РќР°СЃС‚СЂРѕР№РєР° РЎРµРєСЂРµС‚РѕРІ Р’ OS Secret Store

### Р§С‚Рѕ РЇРІР»СЏРµС‚СЃСЏ РЎРµРєСЂРµС‚РѕРј, Рђ Р§С‚Рѕ РќРµС‚

1. `api.partner_id` РІ `config` вЂ” СЌС‚Рѕ РёРґРµРЅС‚РёС„РёРєР°С‚РѕСЂ (UUID), РЅРµ СЃРµРєСЂРµС‚.
2. `audit.hmac_key_id` РІ `config` вЂ” СЌС‚Рѕ РёРґРµРЅС‚РёС„РёРєР°С‚РѕСЂ Р·Р°РїРёСЃРё HMAC, РЅРµ СЃРµРєСЂРµС‚.
3. РЎРµРєСЂРµС‚С‹ вЂ” СЌС‚Рѕ С‚РѕР»СЊРєРѕ Р·РЅР°С‡РµРЅРёСЏ:
   - `API_TOKEN` (РґР»СЏ Р·Р°РіРѕР»РѕРІРєР° `Authorization`);
   - `HMAC_SECRET` (РґР»СЏ РїРѕРґРїРёСЃРё audit).

### РћР±С‰РёР№ РЁР°Р±Р»РѕРЅ РҐСЂР°РЅРµРЅРёСЏ

1. `service`: `ai-anonymizer/api`
2. API-token:
   - `account(user) = <api.partner_id РёР· config>`
   - `value = <API_TOKEN>`
3. Audit HMAC:
   - `account(user) = <audit.hmac_key_id РёР· config>`
   - `value = <HMAC_SECRET>`

Р’Р°Р¶РЅРѕ:
1. `<partner_uuid>` РІ РєРѕРјР°РЅРґР°С… РЅРёР¶Рµ РґРѕР»Р¶РµРЅ Р±С‹С‚СЊ С‚РµРј Р¶Рµ Р·РЅР°С‡РµРЅРёРµРј, С‡С‚Рѕ Рё `api.partner_id` РІ `config`.
2. Р•СЃР»Рё РІ `config` РІС‹ РїРѕРјРµРЅСЏР»Рё `audit.hmac_key_id`, РёСЃРїРѕР»СЊР·СѓР№С‚Рµ СЌС‚Рѕ Р¶Рµ Р·РЅР°С‡РµРЅРёРµ РІ СЃРµРєСЂРµС‚Рµ (РЅРµ РѕР±СЏР·Р°С‚РµР»СЊРЅРѕ `audit-hmac-v1`).

### Р“РµРЅРµСЂР°С†РёСЏ HMAC_SECRET (PowerShell)

РЎРіРµРЅРµСЂРёСЂСѓР№С‚Рµ СЃР»СѓС‡Р°Р№РЅРѕРµ Р·РЅР°С‡РµРЅРёРµ `HMAC_SECRET` (Base64, 32 Р±Р°Р№С‚Р°):

```powershell
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

РСЃРїРѕР»СЊР·СѓР№С‚Рµ РІС‹РІРµРґРµРЅРЅСѓСЋ СЃС‚СЂРѕРєСѓ РєР°Рє `<HMAC_SECRET>` РІ РєРѕРјР°РЅРґР°С… РЅРёР¶Рµ.

### Windows (Credential Manager)

```powershell
cmdkey /generic:"ai-anonymizer/api:<partner_uuid>" /user:"<partner_uuid>" /pass:"<API_TOKEN>"
cmdkey /generic:"ai-anonymizer/api:<audit_hmac_key_id>" /user:"<audit_hmac_key_id>" /pass:"<HMAC_SECRET>"
```

РџСЂРѕРІРµСЂРєР°:

```powershell
cmdkey /list | findstr "ai-anonymizer/api"
```

### Linux (Secret Service / Gnome Keyring)

```bash
echo -n "<API_TOKEN>" | secret-tool store --label="ai-anonymizer api partner" service "ai-anonymizer/api" username "<partner_uuid>"
echo -n "<HMAC_SECRET>" | secret-tool store --label="ai-anonymizer audit hmac" service "ai-anonymizer/api" username "<audit_hmac_key_id>"
```

### Macos (Keychain)

```bash
security add-generic-password -U -s "ai-anonymizer/api" -a "<partner_uuid>" -w "<API_TOKEN>"
security add-generic-password -U -s "ai-anonymizer/api" -a "<audit_hmac_key_id>" -w "<HMAC_SECRET>"
```

## РџСЂРёРѕСЂРёС‚РµС‚РЅС‹Р№ РЎС†РµРЅР°СЂРёР№: Р—Р°РїСѓСЃРє Р§РµСЂРµР· Config

### РЁР°Рі 1. РџРѕРґРіРѕС‚РѕРІСЊС‚Рµ Config Р¤Р°Р№Р»

Р РµРєРѕРјРµРЅРґСѓРµРјС‹Р№ РІР°СЂРёР°РЅС‚:
1. РІРѕР·СЊРјРёС‚Рµ `config.example.yaml` РІ РєРѕСЂРЅРµ РїСЂРѕРµРєС‚Р°;
2. Р·Р°РїРѕР»РЅРёС‚Рµ Р·РЅР°С‡РµРЅРёСЏ РїРѕРґ РІР°С€Рµ РѕРєСЂСѓР¶РµРЅРёРµ;
3. РѕСЃС‚Р°РІСЊС‚Рµ С„Р°Р№Р» СЃ РёРјРµРЅРµРј `config.example.yaml`, РµСЃР»Рё С…РѕС‚РёС‚Рµ Р·Р°РїСѓСЃРє Р±РµР· `--config`.

Р•СЃР»Рё С…РѕС‚РёС‚Рµ РѕС‚РґРµР»СЊРЅС‹Р№ С„Р°Р№Р» (`config.yaml`, `config.prod.yaml` Рё С‚.Рї.), СЌС‚Рѕ РЅРѕСЂРјР°Р»СЊРЅРѕ, РЅРѕ С‚РѕРіРґР° РІ РєРѕРјР°РЅРґР°С… СѓРєР°Р·С‹РІР°Р№С‚Рµ `--config <РїСѓС‚СЊ>`.

РњРёРЅРёРјР°Р»СЊРЅС‹Р№ СЃРѕСЃС‚Р°РІ РїРѕР»РµР№:

```yaml
paths:
  raw_path: "D:/secure-raw"
  workspace_path: "C:/work/my-ai-project"
  output_path: "C:/work/my-ai-project/anonymized"
api:
  base_url: "https://production-retrievals.ai.rarus-cloud.ru"
  partner_id: "00000000-0000-0000-0000-000000000000"
limits:
  max_file_size_mb: 25
  max_pages: 300
  request_timeout_sec: 5
  max_parallel_runs: 1
security:
  allowed_hosts:
    - "production-retrievals.ai.rarus-cloud.ru"
audit:
  retention_days: 30
  hmac_key_id: "audit-hmac-v1"
```

Р’Р°Р¶РЅРѕ РґР»СЏ Windows-РїСѓС‚РµР№ РІ YAML:
1. РЅРµ РёСЃРїРѕР»СЊР·СѓР№С‚Рµ `\` РІРЅСѓС‚СЂРё РґРІРѕР№РЅС‹С… РєР°РІС‹С‡РµРє (РЅР°РїСЂРёРјРµСЂ, `"D:\secure-raw"`), СЌС‚Рѕ РјРѕР¶РµС‚ РІС‹Р·РІР°С‚СЊ РѕС€РёР±РєСѓ `unknown escape character`;
2. РёСЃРїРѕР»СЊР·СѓР№С‚Рµ РѕРґРёРЅ РёР· Р±РµР·РѕРїР°СЃРЅС‹С… РІР°СЂРёР°РЅС‚РѕРІ:
   - `"D:/secure-raw"` (СЂРµРєРѕРјРµРЅРґСѓРµС‚СЃСЏ);
   - `'D:\secure-raw'`;
   - `"D:\\secure-raw"`.

### РЁР°Рі 2. Р—Р°РїСѓСЃРє РЎ РљРѕРјРїР°РєС‚РЅС‹РјРё РљРѕРјР°РЅРґР°РјРё

Р’Р°СЂРёР°РЅС‚ A: С„Р°Р№Р» РІ РґРµС„РѕР»С‚РЅРѕРј РјРµСЃС‚Рµ Рё СЃ РґРµС„РѕР»С‚РЅС‹Рј РёРјРµРЅРµРј `./config.example.yaml` (РјРѕР¶РЅРѕ Р±РµР· `--config`).

Р•СЃР»Рё Сѓ РІР°СЃ `anonym.exe`:

```powershell
C:\tools\anonym\anonym.exe doctor
C:\tools\anonym\anonym.exe scan
C:\tools\anonym\anonym.exe list
C:\tools\anonym\anonym.exe run <id|path>
```

Р•СЃР»Рё Р·Р°РїСѓСЃРє С‡РµСЂРµР· РёСЃС…РѕРґРЅРёРєРё:

```powershell
go run ./cmd/anonym doctor
go run ./cmd/anonym scan
go run ./cmd/anonym list
go run ./cmd/anonym run <id|path>
```

Р’Р°СЂРёР°РЅС‚ B: РѕС‚РґРµР»СЊРЅС‹Р№ config-С„Р°Р№Р» (РЅСѓР¶РµРЅ `--config`).

Р•СЃР»Рё Сѓ РІР°СЃ `anonym.exe`:

```powershell
C:\tools\anonym\anonym.exe doctor --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe scan --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe list --config C:\work\my-ai-project\config.yaml
C:\tools\anonym\anonym.exe run <id|path> --config C:\work\my-ai-project\config.yaml
```

Р•СЃР»Рё Р·Р°РїСѓСЃРє С‡РµСЂРµР· РёСЃС…РѕРґРЅРёРєРё:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym scan --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym list --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym run <id|path> --config C:\work\my-ai-project\config.yaml
```

## Override РџР°СЂР°РјРµС‚СЂРѕРІ: РљРѕРіРґР° Р РљР°Рє Р”РµР»Р°С‚СЊ

### РљРµР№СЃ 1. Override Р§РµСЂРµР· Env (Р”Р»СЏ РЎРµСЂРёРё Р—Р°РїСѓСЃРєРѕРІ)

РљРѕРіРґР° РёСЃРїРѕР»СЊР·РѕРІР°С‚СЊ:
1. РІ СЌС‚РѕР№ СЃРµСЃСЃРёРё РЅСѓР¶РЅРѕ РІСЂРµРјРµРЅРЅРѕ СЃРјРµРЅРёС‚СЊ 1-2 РїР°СЂР°РјРµС‚СЂР°;
2. РєРѕРјР°РЅРґС‹ Р·Р°РїСѓСЃРєР°СЋС‚СЃСЏ РјРЅРѕРіРѕ СЂР°Р·, Рё РЅРµ С…РѕС‡РµС‚СЃСЏ РїРѕРІС‚РѕСЂСЏС‚СЊ С„Р»Р°Рі.

РџСЂРёРјРµСЂ:

```powershell
$env:ANON_REQUEST_TIMEOUT_SEC="15"
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml
go run ./cmd/anonym run <id> --config C:\work\my-ai-project\config.yaml
```

### РљРµР№СЃ 2. Override Р§РµСЂРµР· CLI (Р Р°Р·РѕРІС‹Р№ Р—Р°РїСѓСЃРє)

РљРѕРіРґР° РёСЃРїРѕР»СЊР·РѕРІР°С‚СЊ:
1. СЂР°Р·РѕРІРѕ РЅСѓР¶РЅРѕ РїСЂРѕРІРµСЂРёС‚СЊ РґСЂСѓРіРѕР№ С…РѕСЃС‚/Р»РёРјРёС‚;
2. РІР°Р¶РЅРѕ СЏРІРЅРѕ Р·Р°С„РёРєСЃРёСЂРѕРІР°С‚СЊ override РІ РёСЃС‚РѕСЂРёРё РєРѕРјР°РЅРґС‹.

РџСЂРёРјРµСЂ:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml --api-base-url https://pilot-retrievals.company.local --allowed-hosts pilot-retrievals.company.local
```

## Р”РёР°РіРЅРѕСЃС‚РёРєР° Config Р РµР¶РёРјР°

1. Р•СЃР»Рё `--config` СѓРєР°Р·Р°РЅ СЏРІРЅРѕ Рё С„Р°Р№Р» РЅРµ РЅР°Р№РґРµРЅ, РєРѕРјР°РЅРґР° Р·Р°РІРµСЂС€РёС‚СЃСЏ РѕС€РёР±РєРѕР№ Р·Р°РіСЂСѓР·РєРё С„Р°Р№Р»Р°.
2. Р•СЃР»Рё `--config` РЅРµ СѓРєР°Р·Р°РЅ, runtime РёС‰РµС‚ `./config.example.yaml` РІ С‚РµРєСѓС‰РµР№ СЂР°Р±РѕС‡РµР№ РїР°РїРєРµ.
3. Р•СЃР»Рё С„Р°Р№Р» РёР· Рї.2 РѕС‚СЃСѓС‚СЃС‚РІСѓРµС‚, runtime РїСЂРѕРґРѕР»Р¶РёС‚ СЃ env/CLI.
4. РџСЂРё Р±РёС‚РѕРј YAML РєРѕРјР°РЅРґР° Р·Р°РІРµСЂС€РёС‚СЃСЏ СЏРІРЅРѕР№ РѕС€РёР±РєРѕР№ РїР°СЂСЃРёРЅРіР° РєРѕРЅС„РёРіР°.

## Р’Р°Р¶РЅС‹Рµ РћРіСЂР°РЅРёС‡РµРЅРёСЏ V1

1. `fields` С„РёРєСЃРёСЂРѕРІР°РЅ РІ `anonymizer`.
2. `user-id` РЅРµ РѕС‚РїСЂР°РІР»СЏРµС‚СЃСЏ РёР· CLI runtime.
3. `max_parallel_runs` РІ `v1` РґРѕР»Р¶РµРЅ Р±С‹С‚СЊ `1`.
4. РЎРµРєСЂРµС‚С‹ РЅРµ С‡РёС‚Р°СЋС‚СЃСЏ РёР· plaintext-config/env fallback.
5. Audit РїРёС€РµС‚СЃСЏ РІ `<workspace>/.anonym/audit.log` (JSONL), `file_id` С…СЂР°РЅРёС‚СЃСЏ РєР°Рє HMAC.

## Р’СЂРµРјРµРЅРЅС‹Р№ РџСЂРѕС„РёР»СЊ V1-1-WS (MVP Р”Р»СЏ РР·РѕР»РёСЂРѕРІР°РЅРЅС‹С… IDE)

РљРѕРіРґР° РёСЃРїРѕР»СЊР·РѕРІР°С‚СЊ:
1. РІС‹ Р·Р°РїСѓСЃРєР°РµС‚Рµ РёР· РёР·РѕР»РёСЂРѕРІР°РЅРЅРѕР№ СЃСЂРµРґС‹ AI IDE;
2. runtime РЅРµ РёРјРµРµС‚ РґРѕСЃС‚СѓРїР° Рє OS secret store С…РѕСЃС‚Р°.

Р§С‚Рѕ РјРµРЅСЏРµС‚СЃСЏ:
1. РІРєР»СЋС‡Р°РµС‚СЃСЏ СЏРІРЅС‹Р№ insecure-РїСЂРѕС„РёР»СЊ `v1-1-ws`;
2. РґР»СЏ API auth РёСЃРїРѕР»СЊР·СѓРµС‚СЃСЏ РІСЂРµРјРµРЅРЅС‹Р№ `api.auth_token` РёР· env/config;
3. РґР»СЏ audit РёСЃРїРѕР»СЊР·СѓРµС‚СЃСЏ РІСЂРµРјРµРЅРЅС‹Р№ `audit.hmac_secret` РёР· env/config;
4. РїСЂРѕС„РёР»СЊ РґРѕРїСѓСЃРєР°РµС‚СЃСЏ С‚РѕР»СЊРєРѕ РґР»СЏ non-production РєРѕРЅС‚СѓСЂРѕРІ.

Р§С‚Рѕ РЅРµ РјРµРЅСЏРµС‚СЃСЏ:
1. boundary Рё path hardening;
2. HTTPS + allowlist;
3. Р»РёРјРёС‚С‹ Рё audit.

Р’Р°Р¶РЅРѕ:
1. `v1-1-ws` вЂ” РІСЂРµРјРµРЅРЅС‹Р№ РїСЂРѕС„РёР»СЊ MVP, РЅРµ РґРµС„РѕР»С‚РЅС‹Р№ СЂРµР¶РёРј;
2. РґР»СЏ production РёСЃРїРѕР»СЊР·СѓР№С‚Рµ СЃС‚Р°РЅРґР°СЂС‚РЅС‹Р№ `v1` СЃ OS secret store;
3. РїСЂРё `v1-1-ws` РѕР±СЏР·Р°С‚РµР»РµРЅ РѕРїРµСЂР°С†РёРѕРЅРЅС‹Р№ РєРѕРЅС‚СЂРѕР»СЊ: rate-limit, Р±С‹СЃС‚СЂС‹Р№ revoke/rotation С‚РѕРєРµРЅР° Рё HMAC-СЃРµРєСЂРµС‚Р°.

## Р‘С‹СЃС‚СЂС‹Р№ Р§РµРє-Р›РёСЃС‚ РџРµСЂРµРґ Р—Р°РїСѓСЃРєРѕРј

1. `raw_path` РІРЅРµ workspace, `output_path` РІРЅСѓС‚СЂРё workspace.
2. `api.base_url` РЅР°С‡РёРЅР°РµС‚СЃСЏ СЃ `https://`.
3. Host API РїСЂРёСЃСѓС‚СЃС‚РІСѓРµС‚ РІ `allowed_hosts`.
4. Р’ OS secret store РµСЃС‚СЊ Р·Р°РїРёСЃРё РґР»СЏ `api.partner_id` Рё `audit.hmac_key_id`.
5. `max_parallel_runs = 1`.
6. `doctor` РІРѕР·РІСЂР°С‰Р°РµС‚ `doctor: ok`.

## РЈС‚РѕС‡РЅРµРЅРёРµ РџРѕ РЎРµРєСЂРµС‚Р°Рј Р РџСЂРѕС„РёР»СЏРј (РђРєС‚СѓР°Р»СЊРЅРѕ РќР° 2026-04-29)

Р­С‚РѕС‚ СЂР°Р·РґРµР» СЏРІР»СЏРµС‚СЃСЏ РїСЂРёРѕСЂРёС‚РµС‚РЅС‹Рј РґР»СЏ РёРЅС‚РµСЂРїСЂРµС‚Р°С†РёРё РїСЂР°РІРёР» СЃРµРєСЂРµС‚РѕРІ.

1. РЎС‚Р°РЅРґР°СЂС‚РЅС‹Р№ РїСЂРѕС„РёР»СЊ `v1`:
- СЃРµРєСЂРµС‚С‹ Р±РµСЂСѓС‚СЃСЏ С‚РѕР»СЊРєРѕ РёР· OS secret store;
- РёСЃРїРѕР»СЊР·РѕРІР°РЅРёРµ `api.auth_token` Рё `audit.hmac_secret` РёР· `env/config` Р·Р°РїСЂРµС‰РµРЅРѕ.

2. Р’СЂРµРјРµРЅРЅС‹Р№ РїСЂРѕС„РёР»СЊ `v1-1-ws` (С‚РѕР»СЊРєРѕ non-production):
- СЃРµРєСЂРµС‚С‹ РґРѕРїСѓСЃРєР°СЋС‚СЃСЏ РёР· `env/config` С‚РѕР»СЊРєРѕ РїСЂРё СЏРІРЅРѕРј РІРєР»СЋС‡РµРЅРёРё insecure-СЂРµР¶РёРјР°;
- РѕР±СЏР·Р°С‚РµР»СЊРЅС‹ РѕР±Р° Р·РЅР°С‡РµРЅРёСЏ: `api.auth_token` Рё `audit.hmac_secret`.

### РљР°Рє РЇРІРЅРѕ Р’РєР»СЋС‡РёС‚СЊ `v1-1-ws`

Р§РµСЂРµР· CLI-С„Р»Р°РіРё:

```powershell
go run ./cmd/anonym doctor --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"

go run ./cmd/anonym scan --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"

go run ./cmd/anonym list --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"

go run ./cmd/anonym run <id|path> --config C:\work\my-ai-project\config.yaml --runtime-profile v1-1-ws --insecure-no-secrets --api-auth-token "<API_TOKEN>" --audit-hmac-secret "<HMAC_SECRET>"
```

Р§РµСЂРµР· env-РїРµСЂРµРјРµРЅРЅС‹Рµ:

```powershell
$env:ANON_RUNTIME_PROFILE="v1-1-ws"
$env:ANON_INSECURE_NO_SECRETS="true"
$env:ANON_API_AUTH_TOKEN="<API_TOKEN>"
$env:ANON_AUDIT_HMAC_SECRET="<HMAC_SECRET>"
```

Р’Р°Р¶РЅРѕ:
- РїСЂРё `v1-1-ws` CLI РІС‹РІРѕРґРёС‚ РїСЂРµРґСѓРїСЂРµР¶РґРµРЅРёРµ Рѕ РІСЂРµРјРµРЅРЅРѕРј insecure-РїСЂРѕС„РёР»Рµ;
- РґР»СЏ production-РєРѕРЅС‚СѓСЂР° РёСЃРїРѕР»СЊР·СѓР№С‚Рµ С‚РѕР»СЊРєРѕ СЃС‚Р°РЅРґР°СЂС‚РЅС‹Р№ `v1`.

### Р§РµРє-Р›РёСЃС‚ РџРѕ РџСЂРѕС„РёР»СЏРј

1. Р”Р»СЏ `v1`: РїСЂРѕРІРµСЂСЊС‚Рµ Р·Р°РїРёСЃРё РІ OS secret store (`api.partner_id`, `audit.hmac_key_id`).
2. Р”Р»СЏ `v1-1-ws`: РІРјРµСЃС‚Рѕ OS secret store Р·Р°РґР°Р№С‚Рµ `api.auth_token` Рё `audit.hmac_secret` (CLI/env/config) Рё РІРєР»СЋС‡РёС‚Рµ `insecure_no_secrets=true`.

## Р”РёР°РіРЅРѕСЃС‚РёРєР° HTTP-Р’С‹Р·РѕРІРѕРІ (РћРїС†РёРѕРЅР°Р»СЊРЅРѕ)

Р”Р»СЏ РґРёР°РіРЅРѕСЃС‚РёРєРё РѕС€РёР±РѕРє Р°РІС‚РѕСЂРёР·Р°С†РёРё/РґРѕСЃС‚СѓРїРЅРѕСЃС‚Рё API РјРѕР¶РЅРѕ РІРєР»СЋС‡РёС‚СЊ debug-Р»РѕРі HTTP.

РљР°Рє РІРєР»СЋС‡РёС‚СЊ:

```powershell
$env:ANON_HTTP_DEBUG="true"
```

РљСѓРґР° РїРёС€РµС‚СЃСЏ Р»РѕРі:
- `<workspace>/.anonym/http-debug.log`

Р§С‚Рѕ РїРёС€РµС‚СЃСЏ:
- РІСЂРµРјСЏ, РєРѕРјР°РЅРґР° (`doctor`/`run`), РјРµС‚РѕРґ, URL, HTTP-СЃС‚Р°С‚СѓСЃ;
- Р·Р°РіРѕР»РѕРІРєРё РІ Р±РµР·РѕРїР°СЃРЅРѕРј РІРёРґРµ (Р·РЅР°С‡РµРЅРёРµ `Authorization` РјР°СЃРєРёСЂСѓРµС‚СЃСЏ).

Р’Р°Р¶РЅРѕ:
- РїРѕ СѓРјРѕР»С‡Р°РЅРёСЋ debug-Р»РѕРі РІС‹РєР»СЋС‡РµРЅ;
- РїСЂРё РІС‹РєР»СЋС‡РµРЅРЅРѕРј С„Р»Р°РіРµ С„Р°Р№Р» Р»РѕРіР° РЅРµ СЃРѕР·РґР°РµС‚СЃСЏ.


## Уточнение Контракта Doctor (Актуально На 2026-04-30)

Preflight-проверка команды `doctor` использует:
- `GET /v1/tasks/anonymization_fields`
- headers: `Authorization`, `partner-id`, `user-id`
- query: `page`, `per_page`

Новые runtime-параметры для `doctor`:
1. `api.user_id` (UUID)
2. `api.fields_page` (int, default `1`)
3. `api.fields_per_page` (int, default `10`)

CLI-флаги:
- `--api-user-id`
- `--api-fields-page`
- `--api-fields-per-page`
