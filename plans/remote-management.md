# RemoteManagement_2.91 + eep268.xml implementation plan

## Scope audited

Request: keep repo as a Go library able to generate Go code from `eep268.xml`, then parse and format telegrams to/from described EnOcean Equipment Profiles and Remote Management messages. Plan only; no project/source edits.

Read:
- `/tmp/enocean-spec-text/RemoteManagement_2.91.txt`
- `eep268.xml` (repo root, UTF-16LE XML)
- Existing Go packages under `pkg/`, `internal/serializer/`, and `main.go`.

## Current support

- ESP3 framing exists: `pkg/esp3/esp3.go:41-68` has `Telegram{PacketType, Data, OptData}` and CRC-backed serialization; `:79-130` parses hex into an ESP3 telegram.
- ERP1 unpack/pack exists: `pkg/erp1/erp1.go:11-19` stores destination, sender, RORG, status, userdata; `:22-66` parses RADIO_ERP1; `:69-89` emits RADIO_ERP1.
- RORG enum already includes `SYS_EX` and `SEC_MAN`: `pkg/enums/rorg.go:17,20`.
- Packet type enum already includes `REMOTE_MAN_COMMAND` (`0x07`) and `RADIO_ERP2`: `pkg/enums/packettype.go:14,16`.
- `pkg/eep/eep.go:21-25` only models an EEP triplet. It does not parse EEP telegram payload fields. It also caps FUNC at `0x60` and TYPE at `0x7f` (`:15-18`), which conflicts with EEP 3.0+ notes in RemoteManagement_2.91 that FUNC and TYPE are 8-bit for newer EEPs.
- `pkg/commoncommand/reman.go:12-78` is ESP3 Common Command REMAN support (`WR_REMAN_CODE`, repeating config), not over-the-air Remote Management SYS_EX/RMCC/RPC.
- `internal/serializer` is byte-aligned struct serialization only. Generated EEP parsers will need bitfield extraction/scaling, not this reflection serializer as-is.
- `main.go` is a serial-port demo, not library API. If library-first is a real goal, move or ignore it later; not needed for the first parser/generator phases.

## Spec requirements that are missing

### Remote Management transport/message layer

From `RemoteManagement_2.91.txt`:
- Remote Management uses chained SYS_EX messages (`RORG 0xC5`) as containers (`4.1`, lines ~688-701).
- ERP1 SYS_EX telegram structure is 16 bytes: `RORG`, 1-byte `msg_id` containing `SEQ` and `IDX`, 8-byte data field, 4-byte sender id, 1-byte status, CRC8 (`4.1.2`, lines ~706-727).
- IDX 0 data field packs `data_length` (9 bits), `manufacturer ID` (11 bits), `fn_number` (12 bits), then 32 bits of payload (`4.1.2`, lines ~748-766). IDX > 0 carries 8 payload bytes (`lines ~768-769`).
- `SEQ = 0` is not allowed; `IDX` starts at 0 and increments (`4.1.3`, lines ~791-803).
- Merge grouping key: destination id, source id, SEQ; sort by IDX (`lines ~824-830`). Duplicate IDX discards message; chain period is 1 second (`4.2`, `4.2.1`, `4.2.2`, lines ~849-922).
- Maximum message parts: 64; max transferred data: 508 bytes (`6`, lines ~1944-1957).
- Broadcast responses need 0-2000 ms random delay (`3.1.4`, lines ~671-676). As a library, expose policy/hook; do not sleep inside parsers.
- Recommended REMAN send status is `0x0F` / no repeat (`4.3`, lines ~1028-1038). Current ERP1 `ToEsp3` hardcodes status from `Packet.Status`, but docs/comments mention repeated statuses differently; generator/remote API should default REMAN status deliberately.

### RMCC/RPC payloads and responses

Missing all SYS_EX RMCC/RPC constructors and parsers:
- RMCC function codes: `0x001` Unlock, `0x002` Lock, `0x003` Set Code, `0x004` Query ID, `0x005` Action, `0x006` Ping, `0x007` Query Function, `0x008` Query Status (`5.1`, lines ~991-1005).
- Deprecated Query ID Answer `0x604`, preferred Query ID Answer Extended `0x704` with `locked by other manager` bit (`5.1.4.1-5.1.4.2`, lines ~1269-1299).
- Ping answer `0x606`, Query Function answer `0x607`, Query Status answer `0x608` (`5.1.6.1-5.1.8.1`, lines ~1322-1435).
- RPCs: Remote Learn `0x201`, Flash Write `0x203`, Flash Read `0x204`, Flash Read Answer `0x804`, Smart Ack read/write settings and answers `0x205/0x805/0x806/0x206` (`5.2`, lines ~1436-1940).
- Status/return codes `0x00-0x11`, including v2.91 additions `Session is closed 0x10`, `Insufficient rights 0x11` (`7.4`, lines ~2423-2444).
- Security/session behavior is mostly device-side policy, but library needs types and constructors. Non-secure unlock/lock/set-code uses 32-bit codes; reserved codes are `0x00000000` and `0xFFFFFFFF` (`2.1`, lines ~285-413; `6`, lines ~1958-1964).

### SEC_MAN compatibility layer

RemoteManagement_2.91 marks chapter 7 as deprecated in favor of Secure Remote Management, but still specifies compatibility:
- `SEC_MAN` RORG `0x34` already exists in enum.
- SEC_MAN has key/type nibble, payload, RLC, CMAC; types: `0x00` single data, `0x01` chained, `0x02` SYS_EX encapsulated (`7.2.1`, lines ~1996-2040).
- Secure ReMan adds mandatory RMCCs: Start Session `0x009`, Close Session `0x00A`; replies `0x609` (`7.3.1-7.3.4`, lines ~2344-2415).
- Actual VAES/CMAC needs Security spec, not just this spec. Keep this as a later phase unless user explicitly wants secure encryption now.

### EEP XML generation

- `eep268.xml` is UTF-16LE. Go `encoding/xml` will not decode UTF-16 by itself without a `CharsetReader` (`golang.org/x/net/html/charset`) or a tiny local UTF-16 decoder. Since no new dependency is needed, use stdlib `unicode/utf16` + BOM/LE handling.
- XML shape: `<profile><rorg><number>0xF6</number>...<func><number>0x01</number>...<type><number>0x01</number>...<case><datafield>...`.
- Data fields include `<bitoffs>`, `<bitsize>`, optional `<range><min/max>`, `<scale><min/max>`, `<unit>`, and `<enum><item><value><description>>`.
- Some `datafield`s are reserved or empty; generator should preserve enough metadata but only generate public fields for non-reserved fields with a usable `shortcut`/`data` name.
- RemoteManagement notes Query ID/Ping/Remote Learn legacy EEP encoding only works with 21-bit EEP; for EEP 3.0 FUNC/TYPE may be 8-bit and FUNC > `0x3f` or TYPE > `0x7f` will not fit (`5.1.4`, `5.1.6.1`, `5.2.1`). The current `pkg/eep` limits need separation: full EEP triplet vs legacy 21-bit REMAN EEP mask format.

## Implementation plan: smallest phases

### Phase 1: Pure bit/field primitives and EEP triplet fix

Files:
- Add `pkg/bitfield/bitfield.go` or keep private in `pkg/eep/internal` if only generated code needs it.
- Update `pkg/eep/eep.go` and tests.

Work:
- Introduce tiny MSB/LSB bit extraction/insertion helpers with table tests. Avoid reflection.
- Split concepts:
  - `eep.EEP{Rorg, Func, Type}` allows full byte FUNC/TYPE.
  - `reman.LegacyEEP21` (or methods in `pkg/reman`) validates/remaps only the 21-bit REMAN mask use case.
- Keep `EEP.String()` stable.

Validation:
- New tests for 1-bit, cross-byte, and 32-bit field extraction/insert round trips.
- `go test ./pkg/eep ./pkg/bitfield`.

Risk:
- Bit numbering must match XML. Confirm with one F6 1-bit enum and one A5 scaled field fixture before generating broad code.

### Phase 2: Generate metadata and code from `eep268.xml`

Files:
- Add `cmd/eepgen/main.go` or `internal/eepgen` plus `go generate` entry. Smallest lazy path: one command that reads XML and writes generated files.
- Add generated package under `pkg/eep/profiles` or `pkg/eep/generated`.
- Add `pkg/eep/profile.go` for shared interfaces/types.

Work:
- Parse UTF-16LE XML into a compact internal model: RORG, FUNC, TYPE, title, cases, fields, enums, ranges/scales/units.
- Generate a registry keyed by `eep.EEP` plus per-profile structs only for parse/format fields.
- Public API shape should stay boring:
  - `profiles.Parse(packet erp1.Packet, profile eep.EEP) (any, error)`
  - `profiles.Format(profile eep.EEP, value any) ([]byte, error)` or profile-specific `Encode()`/`Decode()`.
- Do not generate a bespoke abstraction for every XML detail. Store descriptions/ranges in metadata; generate code only for bits/scaling/enums needed to parse/format.
- Initial generated support can target ERP1 user data and status; ERP2 can be metadata-compatible but separate packet parsing later.

Validation:
- Golden tests generated from 2-3 representative profiles: F6-01-01 enum bit, one A5 4BS scaled sensor, one D2 VLD variable length profile.
- Test generator determinism: run generator twice, `git diff --exit-code`.

Risk:
- XML contains inconsistent empty tags/whitespace. Generator must be tolerant and skip invalid fields with comments in generated metadata rather than failing entire generation.

### Phase 3: Remote Management SYS_EX message encode/decode

Files:
- Add `pkg/reman/message.go`, `pkg/reman/sys_ex.go`, `pkg/reman/status.go`, tests.
- Reuse `pkg/erp1.Packet` and `pkg/deviceid`.

Work:
- Define `reman.Message{Seq, ManufacturerID, Function, Payload, SourceID, DestinationID}`.
- Implement encode to one or more ERP1 SYS_EX packets:
  - RORG `0xC5`
  - msg id split into SEQ/IDX per spec (confirm exact SEQ/IDX bit widths from table; likely 2 bits SEQ and 6 bits IDX for SEC_MAN examples, but SYS_EX table text is ambiguous in OCR)
  - IDX 0 packs 9-bit data length, 11-bit manufacturer ID, 12-bit function number, first 4 payload bytes.
  - IDX > 0 packs next 8 payload bytes.
  - enforce seq nonzero, max 64 parts, max payload 508.
- Implement decoder/merger that accepts `erp1.Packet`, validates RORG, extracts message part, merges by destination/source/seq, sorts IDX, rejects duplicate IDX, and supports chain timeout as caller-provided timestamp.
- Library should return state (`NeedMore`, `Complete`, `Discarded`) rather than sleeping/timers internally.

Validation:
- Unit tests for exact Query ID example payload from spec section 5.1.4.
- Split payload sizes: 0, 4, 5, 12, 508.
- Duplicate IDX and chain-timeout tests.

Risk:
- Spec text extraction makes bit layout hard to see. Use the original PDF if possible before coding msg_id layout and IDX 0 packing. Escalate if PDF unavailable.

### Phase 4: RMCC/RPC typed constructors/parsers

Files:
- Add `pkg/reman/rmcc.go`, `pkg/reman/rpc.go`, `pkg/reman/response.go`, tests.

Work:
- Function constants and return codes.
- Constructors/parsers for mandatory RMCCs first: unlock, lock, set code, query ID, action, ping, query function, query status, start/close session constants.
- Parse paired answers: `0x704`, `0x606`, `0x607`, `0x608`, `0x609`.
- Add RPC payload structs only for spec-defined RPCs in 5.2; do not implement manufacturer-specific RPC generation.
- Validate reserved security codes for constructors where relevant.

Validation:
- Golden byte tests for each RMCC payload.
- Round-trip RMCC -> Message -> ERP1 SYS_EX parts -> merge -> RMCC parse.

Risk:
- Query ID uses legacy 21-bit EEP/mask; do not reuse full EEP byte packing without a dedicated conversion test.

### Phase 5: Integrate EEP generated profiles with telegram parse/format API

Files:
- Add `pkg/telegram` or extend `pkg/eep/profiles` with high-level helpers. Keep package count low; prefer `pkg/eep/profiles` if sufficient.

Work:
- API for ERP1 telegram to typed EEP value:
  - ESP3 parse already exists -> ERP1 parse -> profile decode.
- API for typed EEP value to ERP1/ESP3 telegram with sender/destination/status options.
- Document that Remote Management is separate from normal EEP payload parsing, except where RMCC/RPC carries EEP-like fields.

Validation:
- End-to-end test: ESP3 hex -> ERP1 -> EEP profile value -> format back user data.
- End-to-end REMAN test: build Query ID broadcast, serialize to ESP3, parse back, merge message, parse RMCC.

Risk:
- Current repo tests already fail unrelated packages. Before implementation, decide whether to fix baseline or run targeted package tests in CI.

### Phase 6: Optional SEC_MAN compatibility

Files:
- Add `pkg/reman/secman.go`, tests.

Work:
- First implement framing only: key/type nibbles, SEC_SYS_EX clear header + encrypted payload placeholders.
- Do not implement VAES/CMAC without reading Security spec and approval; this is crypto and scope expansion.

Validation:
- Use v2.91 examples in `7.2.2` as golden vectors after encryption primitives exist.

Risk:
- Crypto correctness. Escalate before adding dependencies or implementing VAES/CMAC manually.

## Files likely to change later

- `pkg/eep/eep.go`, `pkg/eep/eep_test.go`: widen full EEP FUNC/TYPE and add legacy REMAN packing boundary tests.
- `pkg/enums/rorg.go`: maybe no change; already has SYS_EX/SEC_MAN.
- `pkg/erp1/erp1.go`: possibly add option-preserving `ToEsp3` because current method hardcodes RSSI/security level in opt data (`:79-83`). For outgoing radio, maybe okay; for round-trip formatting, preserve options separately.
- `pkg/esp3/esp3.go`: maybe add byte-slice parser rather than hex-only parser. Current parser works but hex-only API is awkward for library use.
- New `pkg/reman/*`: SYS_EX, RMCC/RPC, responses, merge state.
- New generator: `cmd/eepgen/main.go` and generated `pkg/eep/profiles/*_generated.go`.
- New tests alongside each package.

## Validation baseline observed

Command run: `go test ./...`.

Result: failed before any source changes:
- `pkg/enums`: `TestParseCommonCommandFromByte` expects `0x34` invalid, but code has `CommonCommandSET_CRCMode = 0x34`; test panics after nil error.
- `pkg/event`: byte-order expectation failures for several event fields.
- `pkg/subtel`: sender/userdata/SubTel count expectation failures.

Targeted validation for implementation should start with new package tests (`go test ./pkg/reman ./pkg/eep ./pkg/eep/profiles ./cmd/eepgen`) until baseline is repaired.

## Meta-prompt handoff for next implementation agent

Goal: implement library support, in small phases, for generated EEP profile parsing/formatting from `eep268.xml` plus RemoteManagement_2.91 SYS_EX/RMCC/RPC message construction/parsing. Do not widen into secure crypto unless explicitly approved.

Evidence/context:
- ESP3 framing: `pkg/esp3/esp3.go:41-68`, parse at `:79-130`.
- ERP1 packet model: `pkg/erp1/erp1.go:11-89`.
- EEP triplet only: `pkg/eep/eep.go:21-99`; current FUNC/TYPE caps are too low for EEP 3.0 full profiles.
- RORG enum has `SYS_EX 0xC5` and `SEC_MAN 0x34`: `pkg/enums/rorg.go:17,20`.
- Packet type enum has `REMOTE_MAN_COMMAND` and `RADIO_ERP2`: `pkg/enums/packettype.go:14,16`.
- Remote spec requires SYS_EX chained messages, IDX 0 header packing, 1s chain period, 64 parts, 508-byte max payload, RMCC/RPC function constants, and v2.91 session/status additions.
- `eep268.xml` is UTF-16LE and contains profile hierarchy plus bit offsets, sizes, enum/range/scale metadata.

Success criteria:
- Generator reads `eep268.xml` deterministically and generates Go metadata/code for profile parse/format.
- Public API can parse ESP3/ERP1 telegrams into generated EEP profile values and format values back to ERP1 userdata/ESP3.
- Public API can build and parse Remote Management SYS_EX RMCC/RPC messages for the mandatory spec functions.
- Unit/golden tests cover bit packing, XML generator determinism, at least 3 EEP profiles, and REMAN Query ID/Ping/Status round trips.

Hard constraints:
- Do not add a crypto dependency or implement VAES/CMAC without explicit approval.
- Do not make parser logic depend on serial ports or `main.go`.
- Keep generated code deterministic.
- Keep stateful time behavior outside pure parsers; accept timestamps/deadlines from caller.

Suggested approach:
1. Add bitfield helpers and fix full-vs-legacy EEP modeling.
2. Add XML parser/generator for metadata and minimal parse/format methods.
3. Add `pkg/reman` SYS_EX message split/merge.
4. Add typed RMCC/RPC constructors/parsers.
5. Add high-level examples/tests only after low-level round trips pass.

Validation:
- Run focused tests for changed packages.
- Run generator twice and verify no diff.
- Run `go test ./...` last, but note existing unrelated failures.

Stop/escalation:
- Escalate if exact SYS_EX `SEQ/IDX` bit layout cannot be confirmed from source PDF/OCR.
- Escalate before adding dependencies beyond stdlib/current deps.
- Stop after non-secure REMAN + generated EEP profiles unless user explicitly asks for SEC_MAN encryption.

Resolved assumptions:
- `eep268.xml` in repo root is the intended XML source.
- `RemoteManagement_2.91.txt` is source-backed enough for plan; coding SYS_EX bit layout should verify original PDF diagrams if available.
- Acceptance is plan-only, so no source files changed.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only at plans/remote-management.md; did not modify project/source files."
    }
  ],
  "changedFiles": [
    "plans/remote-management.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "wc -l /tmp/enocean-spec-text/RemoteManagement_2.91.txt /tmp/enocean-spec-text/eep268.xml 2>/dev/null || true",
      "result": "passed",
      "summary": "Confirmed RemoteManagement_2.91.txt exists; /tmp eep268.xml was not present."
    },
    {
      "command": "find /tmp/enocean-spec-text -maxdepth 2 -type f -print",
      "result": "passed",
      "summary": "Listed available spec text files."
    },
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Baseline failures in pkg/enums, pkg/event, and pkg/subtel unrelated to this plan-only audit."
    },
    {
      "command": "git diff --cached --name-only | wc -l && git diff --name-only | wc -l",
      "result": "passed",
      "summary": "No staged files and no tracked source diffs before writing the plan."
    }
  ],
  "validationOutput": [
    "go test ./... failed at baseline: pkg/enums invalid-common-command test expects 0x34 invalid; pkg/event byte-order failures; pkg/subtel parsing/count failures.",
    "Only generated plan file was written."
  ],
  "residualRisks": [
    "SYS_EX SEQ/IDX bit layout should be verified against the original PDF diagrams before implementation because text extraction is ambiguous.",
    "SEC_MAN encryption requires Security spec work and was intentionally deferred."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added implementation plan document only; no source changes.",
  "reviewFindings": [
    "no blockers for plan-only handoff"
  ],
  "manualNotes": "Repo has pre-existing untracked files (docs/, eep268.xml, main.go) and baseline test failures; source was not modified."
}
```
