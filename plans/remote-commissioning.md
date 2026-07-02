# Remote Commissioning V1.5 / EEP XML code generation implementation plan

## Scope read

Task asks for a plan only. I did not edit project/source files. Goal to plan this repo as a Go library that can generate Go code from `eep268.xml` to parse and format ERP telegram user data for described EEP profiles, and cover RemoteCommissioning-V1.5 RPC parsing/formatting.

Sources inspected:
- `/tmp/enocean-spec-text/RemoteCommissioning-V1.5.txt`
- `eep268.xml` (UTF-16LE)
- current Go packages under `pkg/`, `internal/serializer/`, tests, and `go.mod`

## Current support in repo

- ESP3 frame parsing/serialization exists:
  - `pkg/esp3/esp3.go:61-68` serializes sync/header/CRC/data/optData.
  - `pkg/esp3/esp3.go:79-134` parses hex ESP3 frames and validates CRCs.
- ERP1 packet wrapper exists:
  - `pkg/erp1/erp1.go:11-20` models destination, RORG, RSSI, security level, status, subTelNum, sender ID, user data.
  - `pkg/erp1/erp1.go:22-67` extracts RORG/user data/sender/status from `RADIO_ERP1` ESP3.
  - `pkg/erp1/erp1.go:69-94` builds an ESP3 telegram from ERP1 packet fields.
- EEP identity support is only a triplet:
  - `pkg/eep/eep.go:21-25` has `EEP{Rorg, Func, Type}`.
  - `pkg/eep/eep.go:50-99` parses/formats strings like `A5-02-01`.
  - No generated profile metadata, no per-profile field parser, no formatter.
- Binary helper exists but is byte/field oriented, not bitfield/profile oriented:
  - `internal/serializer/serializer.go:75-133` serializes tagged struct fields into ESP3 common-command telegram data/optData.
  - `internal/serializer/deserializer.go:61-90` deserializes sequential fixed-size fields.
  - Defaults are big-endian (`serializer.go:27-39`, `deserializer.go:20-32`), matching the spec.
- Common command / TCM remote-management helpers exist, but they are ESP3 common commands, not Remote Commissioning RPCs over SYS_EX:
  - `pkg/commoncommand/reman.go:12-78` only covers WR_REMAN_CODE and Reman repeating.
  - `pkg/commoncommand/securedevice.go` covers TCM secure-device common commands, not Recom link table/security-profile RPC function codes.
- RORG enum already includes the major ERP RORGs:
  - `pkg/enums/rorg.go:7-22` includes RPS, 1BS, 4BS, VLD, SYS_EX, SEC_MAN, SIGNAL, UTE.
- There is currently a `main.go` program that opens a hard-coded serial port. Library goal should avoid forcing this into public API behavior; keep any CLI/demo out of core packages.

## Relevant spec requirements / missing pieces

### General Remote Commissioning constraints

- Big-endian everywhere: spec `RemoteCommissioning-V1.5.txt:724-729` says parameter values and all other telegram fields are big-endian, MSB first. Current serializer already defaults to big-endian, but generated EEP bit extraction must also define bit numbering consistently.
- Stateless command/response: spec `:731-740` says responses include metadata (for example index/length) even when request implies it. Do not drop metadata from parsed structs.
- Mandatory Remote Commissioning bundle: spec `:3355-3360` requires `Get Product ID Query & Response` and `Remote Commissioning Acknowledge`.

### Remote Commissioning RPC function coverage missing

All below use Manufacturer ID `0x7FF` unless otherwise stated and are Remote Management RPC/SYS_EX functions, not ESP3 common-command packet type.

Mandatory/common:
- `0x240` Remote Commissioning Acknowledge, spec `:767-783`.
- `0x227/0x827` Get Product ID query/response, Product ID is 2-byte manufacturer + 4-byte product ref, spec `:2873-2942`.
- `0x227/0x828` Get Product ID Selective, same request function with non-zero selection payload; selection types dBm/product/modulo, spec `:2943-3030`.
- `0x224` Reset Device Defaults, flags for config/inbound/outbound defaults, spec `:2817-2865`.
- `0x225` Radio Link Test Control, enable + 7-bit cluster count, spec `:2867-2900`.
- `0x226` Apply Changes, flags for link-table and configuration changes, spec `:2781-2811`.

Link table / teach:
- `0x210/0x810` Get Link Table Metadata, response 5 bytes with support flags and sizes, spec `:1324-1390`.
- `0x211/0x811` Get Link Table, direction + start/end; response direction + repeated 9-byte entries, spec `:1391-1470`.
- `0x212` Set Link Table Content, direction + repeated 9-byte entries, no paired response, spec `:1471-1524`.
- `0x213/0x813` Get Link Table GP Entry, direction + index + variable GP channel desc, spec `:1526-1604`.
- `0x214` Set Link Table GP Entry Content, direction + index + variable GP desc, spec `:1634-1671`.
- `0x220` Remote Set Learn Mode, 2-bit mode + inbound index, spec `:1673-1717` and repeated at `:2672-2712`.
- `0x221` Trigger Outbound Remote Teach Request, channel selection, paired response is the teach-in request itself, spec `:1719-1747`.

Security profile:
- `0x215/0x815` Get Security Profile, direction + index; response is direction, index, SLF, 4-byte RLC, 16-byte key, destination ID, source ID, spec `:1748-1835`.
- `0x216` Set Security Profile, same 32-byte content; spec table says data length 23 but field/data-structure totals 32 bytes. Treat table length as a spec typo and implement/validate 32 bytes with a comment, because rows `:1837-1905` enumerate 1+1+1+4+16+4+4.
- `0x234/0x834` Get Device Security Information, index + SLF/RLC/key, spec `:2574-2634`.
- `0x235` Set Device Security Information, index + SLF/RLC/key, spec `:2636-2670`.

Configuration:
- `0x230/0x830` Get Device Configuration, start/end index + optional length cap; response is repeated parameter records `{index uint16, length uint8, value []byte}`, spec `:2180-2261`.
- `0x231` Set Device Configuration, repeated parameter records, no paired response, spec `:2262-2327`.
- `0x232/0x832` Get Link Based Configuration, direction + link-table index + start/end + length; response includes direction/link index + parameter records, spec `:2328-2444`.
- `0x233` Set Link Based Configuration, direction/link index + parameter records, spec `:2445-2572`.
- V1.5 introduced limits: request parameter value length max 64 bytes and overall payload max 67 bytes for config responses/sets (`:2201-2207`, `:2258-2260`, `:2321-2326`). Text parameters limited to 64 bytes (`:2165-2169`). Enforce these in constructors/formatters.

Security process risk:
- Legacy Initial Secure Setup added in V1.5, spec `:3031-3205`. For a library parser/formatter, phase 1 should expose the telegram structs/constants only; do not implement session/security policy workflow until requested. It requires Signal 0x04/0x05 and secure telegram/RLC behavior outside current repo.

### EEP XML generation requirements

`eep268.xml` facts:
- File is UTF-16LE (`<?xml version="1.0" encoding="utf-16le"?>`). Generator must decode UTF-16LE before XML decoding; Go stdlib XML decoder does not handle UTF-16 unless supplied a charset reader. Smallest path: read bytes and use `unicode.UTF16` from `golang.org/x/text` only if already present? It is not in `go.mod`; avoid new dependency by detecting BOM/encoding and decoding with stdlib `unicode/utf16` manually.
- Profile hierarchy begins at `eep268.xml:470-509`: `<profile><rorg><number>...<func><number>...<type><number>...`.
- Data fields use `<datafield>`, `<bitoffs>`, `<bitsize>`, optional `<shortcut>`, `<data>`, `<enum><item><value>...`, e.g. `eep268.xml:546-562`.
- Counts from decoded XML: 4 RORG groups (`F6`, `D5`, `A5`, `D2`), 106 funcs, 370 cases/types, 2576 datafields. XML also has an `<rpc>` section around `eep268.xml:61061`; inspect before deciding whether it helps Remote Commissioning. It likely documents EEP RPCs, not necessarily Recom V1.5.
- Example generated field model must handle reserved fields, enums, units, ranges/scales, overlapping cases, and condition/case nodes. Do not flatten blindly into one struct per EEP without preserving cases/teach-in/status variants.

Minimum library API needed for EEP:
- Parse raw ERP1 user data by EEP triplet into named field values.
- Format named/raw field values back into ERP1 user data bytes.
- Provide generated profile metadata: title, RORG/FUNC/TYPE, fields with bit offset, bit size, shortcut/name, enum values, min/max/scale/unit when present, reserved flag.
- Keep raw-value access first. Physical unit conversion can be helper functions generated from min/max/scale later; do not block parser/formatter on full semantic conversion.

## Important current constraints and risks

- Current `EEP.FromTriplet` does not validate RORG (`pkg/eep/eep.go:27-40`) and allows any byte. That is useful for raw XML-generated triplets, but generated lookup should only return known generated profiles.
- `pkg/enums/rorg.go` has only hardcoded known RORGs and `ParseRorgFromByte` rejects unknown bytes. Since XML only includes F6/D5/A5/D2 this is OK for EEP 2.6.8, but avoid requiring enum validation for future generated EEP XML.
- Existing reflection serializer/deserializer cannot parse bitfields or variable-length parameter lists in the middle of structs. Reuse its big-endian behavior where it fits, but EEP and Recom payloads need explicit byte/bit functions.
- Existing `erp1.Packet.ToEsp3()` hard-codes optional data RSSI/security bytes to `0xff, 0x03` (`pkg/erp1/erp1.go:79-83`) instead of preserving packet fields. Not central to EEP generation, but round-trip formatting at ESP3 level will surprise users.
- Existing test suite currently fails unrelated to this plan: `go test ./...` fails in `pkg/enums`, `pkg/event`, and `pkg/subtel`. Record this as baseline before implementation; do not treat new Recom/EEP work as done until targeted tests pass and existing failures are either fixed or isolated.
- License/IP note: spec text says EnOcean Alliance spec has non-commercial/IP restrictions. Do not check generated copies of spec text into repo; generated code from `eep268.xml` should preserve only necessary metadata and source attribution.

## Smallest implementation phases

### Phase 1 — Add explicit bit/payload primitives, no generator yet

Files to add/change:
- Add `pkg/eep/bitfield.go` + tests.
- Add `pkg/recom/payload.go` + tests.

Work:
- Implement `GetBitsBE(data []byte, bitOffset, bitSize int) (uint64, error)` and `SetBitsBE(data []byte, bitOffset, bitSize int, value uint64) error`.
- Define bit numbering from XML with tests using one known XML example (F6-02-01 fields at bit offsets 0/3/4/7) and one cross-byte field (A5 temperature bit offset 16 size 8).
- Implement simple Recom helpers: direction flag byte, capped parameter records, product ID type.

Validation:
- `go test ./pkg/eep ./pkg/recom`

### Phase 2 — Hand-written Recom RPC library surface

Files to add/change:
- Add `pkg/recom/recom.go`, `pkg/recom/linktable.go`, `pkg/recom/config.go`, `pkg/recom/security.go`, `pkg/recom/common.go`.
- Possibly add `pkg/sysex/sysex.go` if current repo lacks SYS_EX RPC framing. Keep it minimal: manufacturer ID, function code, payload, addressed/broadcast metadata if present in existing Remote Management framing.

Work:
- Define function code constants for all RPCs above.
- Define structs and `MarshalPayload`/`Parse...Payload` functions for query/response/content payloads.
- Enforce V1.5 length limits for config params (64-byte value/request length, 67-byte payload unless parsing response bytes from the wild should return a descriptive error).
- Represent paired/no-paired response and broadcast/addressed support as metadata if useful; do not implement workflows/state machine.

Validation:
- Table tests from spec byte layouts for each RPC payload.
- `go test ./pkg/recom`

### Phase 3 — XML reader and generator CLI/tool

Files to add/change:
- Add `internal/eepxml/eepxml.go` parser + tests.
- Add `cmd/eepgen/main.go` generator or `internal/cmd/eepgen` plus `go:generate` in `pkg/eep/generated.go` header. Keep generated output in `pkg/eep/generated_profiles.go`.

Work:
- Decode UTF-16LE manually using stdlib `unicode/utf16`.
- Parse only needed XML nodes first: rorg/func/type/case/datafield/shortcut/data/description/bitoffs/bitsize/enum/range/scale/unit/reserved.
- Generate metadata tables, not custom code per field initially. This avoids 370 fragile generated structs.
- Generate deterministic sorted slices/maps keyed by `EEP.String()`.

Validation:
- Parser test checks counts: 4 RORGs, 106 funcs, 370 case/type nodes, and known profile `F6-02-01` has `R1`, `EB`, `R2`, `SA` fields.
- Generator golden test on a tiny XML fixture, not the full copyrighted XML.
- `go test ./internal/eepxml ./pkg/eep`

### Phase 4 — EEP parse/format runtime API over generated metadata

Files to add/change:
- Extend `pkg/eep/eep.go` or add `pkg/eep/profile.go`, `pkg/eep/parse.go`, `pkg/eep/format.go`.

Work:
- API shape:
  - `Lookup(profile EEP) (*Profile, bool)`
  - `Parse(profile EEP, userData []byte) (TelegramValues, error)`
  - `Format(profile EEP, values map[string]uint64) ([]byte, error)`
  - Include raw reserved handling policy: parser may ignore reserved by default; formatter zeroes reserved unless raw override explicitly provided.
- Start with raw integer/enumeration values. Add optional enum description lookup.
- Use generated field metadata and phase-1 bitfield helpers.

Validation:
- Round-trip tests on F6-02-01, D5-00-01, A5-02-01, D2-01-00 samples.
- Property-ish small test: for each generated field with bitSize <= 64, set max value and parse back on a zero buffer for non-overlapping simple cases.

### Phase 5 — Integrate ERP1 convenience without coupling

Files to add/change:
- Add `pkg/erp1/eep.go` or `pkg/eep/erp1.go`.

Work:
- Convenience wrapper: `eep.ParseERP1(profile EEP, packet erp1.Packet)` validates `packet.Rorg == profile.Rorg` and parses `packet.UserData`.
- Convenience formatter: build `erp1.Packet` user data from profile values, leaving sender/destination/status to caller.
- Do not hide ESP3/ERP1 packet metadata.

Validation:
- Round-trip from `erp1.Packet` to values and back for one profile.

### Phase 6 — Documentation and examples

Files to add/change:
- `docs/eep-generation.md`
- `docs/recom.md`
- examples under `examples/` only if repo convention accepts examples; otherwise package tests with examples.

Work:
- Document generator command and that full XML is an input artifact, not fetched at runtime.
- Document Recom RPC payload support vs. not implemented workflows/security sessions.

Validation:
- `go test ./...` after deciding what to do with existing baseline failures.

## Validation plan

Baseline command run during audit:
- `go test ./...` failed before any source changes:
  - `pkg/enums`: `TestParseCommonCommandFromByte` expects `0x34` invalid but code accepts `CommonCommandSET_CRCMode`; test then panics on nil error.
  - `pkg/event`: endian expectations fail (`0x1234` parsed as `0x3412`, IDs reversed).
  - `pkg/subtel`: sender/user data/subtelegram count expectations fail.
- `git diff --cached --name-only` produced no output; no staged files.

Targeted checks after implementation:
1. `gofmt` on changed Go files.
2. `go test ./pkg/eep ./internal/eepxml ./pkg/recom ./pkg/erp1`
3. Generator determinism: run generator twice and assert no diff.
4. Full `go test ./...`; if still failing, list pre-existing failures separately.

## Suggested next-agent contract

Goal:
- Implement phases 1 and 2 first: minimal reusable bitfield primitives and hand-written Remote Commissioning V1.5 RPC payload parser/formatter package. Do not start full XML codegen until these foundations exist.

Context/evidence:
- Current repo has ESP3/ERP1 wrappers and triplet-only EEP support (`pkg/esp3/esp3.go`, `pkg/erp1/erp1.go`, `pkg/eep/eep.go`).
- Remote Commissioning V1.5 required RPCs and byte layouts are listed above with spec line references.
- `eep268.xml` is UTF-16LE and hierarchical rorg/func/type/case/datafield metadata starts at `eep268.xml:470`.

Success criteria:
- Public Go package can marshal/parse all listed Recom payloads with table tests.
- Bitfield helpers can parse/format known EEP field examples.
- No source generator yet unless phase 1/2 are complete and tested.

Hard constraints:
- Do not modify source files when only producing a plan/review.
- Avoid new dependencies unless stdlib cannot reasonably do it; UTF-16LE can be decoded with stdlib.
- Preserve big-endian behavior.
- Do not implement legacy secure setup workflow/session policy without explicit scope; expose only payloads/constants if needed.

Stop/escalation:
- Ask for decision before changing public API shape that breaks existing packages.
- Stop after phase 2 if existing repo tests unrelated to new work block full `go test ./...`; report baseline failures.

## Residual risks

- Remote Management SYS_EX framing details are in another spec, not fully present in RemoteCommissioning-V1.5 text. Payload structs can be correct, but full over-the-air SYS_EX wrapping may require reading Remote Management / Secure Remote Management specs.
- EEP XML contains complex `case`, `condition`, teach-in/status structures. Metadata-table generation is safer than generated structs initially, but semantic profile-specific formatting may need iterative fixes.
- Existing test failures suggest current endian/device ID handling may already be inconsistent in some packages.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only at plans/remote-commissioning.md; no project/source files were modified."
    }
  ],
  "changedFiles": [
    "plans/remote-commissioning.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Baseline failures in pkg/enums, pkg/event, and pkg/subtel before source changes."
    },
    {
      "command": "git diff --cached --name-only",
      "result": "passed",
      "summary": "No staged files."
    }
  ],
  "validationOutput": [
    "Read RemoteCommissioning-V1.5.txt, eep268.xml, and relevant Go files; wrote plan only.",
    "go test ./... baseline: pkg/enums invalid-command test/panic; pkg/event endian mismatches; pkg/subtel parsing/count mismatches."
  ],
  "residualRisks": [
    "SYS_EX Remote Management framing may require additional Remote Management/Secure Remote Management specs.",
    "EEP XML case/condition semantics are complex; initial generator should preserve metadata and avoid over-flattening.",
    "Existing unrelated test failures need triage before full-suite validation can pass."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added audit and phased implementation plan document only.",
  "reviewFindings": [
    "no blockers for plan-only task"
  ],
  "manualNotes": "Repository already has untracked docs/, eep268.xml, main.go, and plans/ in git status; only the requested plan file was written by this task."
}
```