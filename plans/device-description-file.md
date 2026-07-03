# Device Description File V2.0 / eep268.xml implementation plan

## Scope read

- User asked for an audit and implementation plan only. No project/source files were modified.
- Sources inspected:
  - `/tmp/enocean-spec-text/Device-Description-File-V2.0.txt`
  - `eep268.xml` (UTF-16LE; decoded locally to inspect only)
  - Go packages under `pkg/`, especially `pkg/eep`, `pkg/erp1`, `pkg/esp3`, `pkg/enums`, `pkg/event`, `pkg/subtel`, `pkg/commoncommand`.

## Spec requirements that matter for the requested library

The goal is specifically: generate Go code from `eep268.xml` to parse and format telegrams from/to described profiles. DDF V2.0 is broader than EEP telegram parsing, so keep the first implementation to the smallest useful subset.

DDF V2.0 high-value requirements/evidence:

- DDF root/device model: `Enocean_Devices` root with fixed `schemaVersion`; `Device` children carry `Product_ID` and metadata. Spec lines 209-216 and 1016-1032.
- Product IDs are strict 6-byte hex strings: `Hex6Bytes` pattern `0x[0-9A-Fa-f]{12}`. Spec lines 994-1000; changelog says V2.0 enforces 6-byte Product_ID, lines 1682-1683.
- TX/RX structure distinguishes transmitted and received profiles:
  - `TX` contains `EURID` and `BaseID`; `RX` contains `EEP`, `GP`, `MSC`; `EURID`/`BaseID` may contain `EEP`, `GP`, `Signals`, `MSC`. Spec lines 1040-1048.
  - `BaseID` may repeat `0..127` and has required `IDOffset`. Spec lines 1043-1048.
- DDF EEP references are only profile identifiers, not field layouts: `EEP` has `LinkEntry` attr and `Rorg`, optional `Func`, optional `Type`; `Func`/`Type` skipped for MSC (`Rorg=0xD1`). Spec lines 1049-1052.
- GP support is structurally separate and includes teach-in payload, channels, channel/value/signal types, enum lists, and restrictions. Spec lines 1053-1066 and 1066-1115.
- MSC support requires field descriptions when `Rorg 0xD1` is referenced; `MSC` contains `Bitfield`, whose `Field` has `name`, `offset`, `bitSize`, and child `Enum` or `Scaled`. Spec lines 1123-1130.
- `SubsequentPayload` handles dynamic MSC/VLD-like payloads where an earlier type/CMD field selects later interpretation; can define fields directly or reference another `SubsequentPayload` by id/refId. Spec lines 1138-1141 and 1403-1412.
- Device parameters/ReCom are large and separate from EEP telegram data: recom params have `index`, `accessLevel`, `recommendedUserLevel`, optional RPC write address/length, `BitField`, `Scaled`, `Text`, `Enum`, `Private`, apps, linked params. Spec lines 1206-1267 and 1358-1367.
- Optional command/RPC metadata exists but is not telegram profile parsing. Spec lines 1368-1390.
- DDF server/index retrieval exists (`index.xml`, `Product Path`, `Product_ID`), but is outside code generation from local `eep268.xml`. Spec lines 1614-1624.

## eep268.xml findings

`eep268.xml` is not a DDF V2.0 file. It is the EEP profile catalog in XML, encoded as UTF-16LE with declaration `encoding="utf-16le"` and DOCTYPE `eep.dtd`.

Observed shape after UTF-16LE decode:

- Root: `<eep>`.
- Counts from parsed XML: 47,068 elements; 65 `<rorg>`, 106 `<func>`, 341 `<type>`, 2,576 `<datafield>`, 1,910 `<enum>`, 663 `<scale>`, 370 `<case>`, 130 `<condition>`, 48 `<statusfield>`.
- First real profiles:
  - `/tmp/eep268.utf8.xml:470-471`: `<rorg><number>0xF6</number>` RPS Telegram.
  - Similar top-level RORGs: `0xD5` 1BS, `0xA5` 4BS, `0xD2` VLD.
- Field layout examples:
  - `/tmp/eep268.utf8.xml:545-552`: `<datafield>` with `<bitoffs>` and `<bitsize>`.
  - `/tmp/eep268.utf8.xml:559-562`: field with `<enum>`.
  - `/tmp/eep268.utf8.xml:599-619`: `<condition>` entries select fields/cases.
- XPath pattern for generator input: `eep/rorg[number]/func[number]/type[number]/case/datafield` plus adjacent `statusfield`, `condition`, `direction`, `scale`, `range`, `enum`, `reserved`, `unit`, `shortcut`, `description`, `status`.

Implementation risk: Go `encoding/xml` cannot read UTF-16 without a charset reader. The minimal robust path is to use `golang.org/x/text/encoding/unicode` + `transform` if already allowed/added, or a tiny UTF-16LE decoder for this generator only. Since this is a generator tool, adding `x/text` is acceptable only if the next agent wants less custom code; otherwise decode bytes via `unicode/utf16` after BOM/encoding detection.

## Current Go support

What exists now:

- ESP3 framing and CRC:
  - `pkg/esp3/esp3.go:61-68` serializes ESP3 telegrams.
  - `pkg/esp3/esp3.go:79-137` parses ESP3 hex strings and verifies CRC.
- ERP1 packet extraction/building:
  - `pkg/erp1/erp1.go:11-20` has `Packet` with `Rorg`, `UserData`, sender/destination IDs, status, RSSI/security.
  - `pkg/erp1/erp1.go:22-67` extracts ERP1 from ESP3; `UserData` is the EEP payload (`telegram.Data[1:senderIdOffset]`).
  - `pkg/erp1/erp1.go:69-94` builds ERP1 back to ESP3 bytes.
- EEP identifier only:
  - `pkg/eep/eep.go:21-25` defines `EEP{Rorg, Func, Type}`.
  - `pkg/eep/eep.go:27-100` validates/parses/formats EEP triplets.
- RORG constants:
  - `pkg/enums/rorg.go:7-23` includes RPS, 1BS, 4BS, VLD, MSC, SIGNAL, UTE, etc.
  - `pkg/enums/rorg.go:25-120` parse/string/valid helpers.
- Serial side is not yet a clean library API:
  - `pkg/enocean.go:24-49` opens serial port and starts parser goroutine.
  - `pkg/enocean.go:46` has TODO: handle channel for ESP3 telegrams and cancellation.
  - `main.go` is a runnable example/hard-coded serial path, not needed for profile codegen.

What is missing for the requested goal:

- No parser for `eep268.xml`.
- No intermediate profile model for RORG/FUNC/TYPE/cases/datafields.
- No code generator (`go:generate`, `cmd/...`, templates) for generated profile code.
- No bit extraction/packing helpers for EEP data fields.
- No generated per-profile parse/format functions.
- No dynamic dispatch from `(RORG,FUNC,TYPE)` to profile codecs.
- No support for scale/range conversions, enums, reserved fields, condition-selected cases, status fields, directions, learn-in/teach-in metadata, or VLD command/sub-payload selection.
- DDF V2.0 device metadata (`Product_ID`, TX/RX/BaseID/EURID, GP, MSC, ReCom params, OptionalCommands, SupportedRPC, index.xml retrieval) is entirely absent.

## Smallest useful architecture

Do not try to implement full DDF V2.0 first. The user asked for generated Go from `eep268.xml` to parse/format described telegrams. Start with EEP payloads, then bridge DDF later.

Proposed packages/files:

- `internal/eepxml/` or `internal/eepgen/model.go`: XML structs + normalization for `eep268.xml`.
- `cmd/eepgen/main.go`: generator executable; input `eep268.xml`, output generated Go.
- `pkg/eep/bit.go`: tiny shared bit get/set helpers. Use MSB/LSB behavior verified against EEP examples before locking API.
- `pkg/eep/profile.go`: public minimal runtime types:
  - `type ProfileID struct { Rorg enums.Rorg; Func, Type byte }` or reuse existing `EEP`.
  - `type FieldValue struct { Raw uint64; Value float64/string; Enum string; Unit string }` (keep lean; exact shape can be refined).
  - `type Parsed struct { EEP EEP; Case string; Fields map[string]FieldValue }`.
  - `type Codec interface { Parse(userData []byte, status byte) (Parsed, error); Format(Parsed) ([]byte, byte, error) }`.
- `pkg/eep/generated_profiles.go`: generated registry and per-profile codecs.
- `pkg/eep/generated_profiles_test.go`: generated or fixture tests for a small profile subset.
- Optional later: `pkg/ddf/` for actual DDF V2.0 Product_ID/TX/RX metadata parsing. This should not block EEP XML codegen.

Reuse existing code:

- Keep `pkg/eep.EEP` and its string parsing/formatting. It is the natural key.
- Keep `pkg/erp1.Packet` as the boundary from radio telegram to EEP payload; generated codecs should accept `Packet.UserData` plus `Packet.Status`, not ESP3 bytes.
- Keep `pkg/enums.Rorg` constants.

## Phase plan

### Phase 0 - Decide the exact first target

Outcome: generator produces code for a very small subset (`F6-02-01`, `D5-00-01`, one simple `A5-02-xx`) and a shared runtime API.

Why: confirms bit numbering, enum/scale model, and round-trip shape before generating thousands of fields.

Validation:

- `go test ./pkg/eep ./pkg/erp1 ./pkg/esp3`.
- Golden tests with hand-picked `UserData` bytes from simple RPS/1BS/4BS fixtures.

### Phase 1 - Parse eep268.xml into a normalized model

Implement generator-only parser:

- Decode UTF-16LE input.
- Ignore DOCTYPE/external DTD.
- Parse `rorg[number,title,telegram] -> func[number,title] -> type[number,title,status] -> case -> datafield/statusfield`.
- Capture: name/data, shortcut, description, bit offset, bit size, enum items, scale min/max/range/unit, reserved, conditions, direction.
- Emit a JSON debug dump or generator test so future changes can inspect model without reading XML manually.

Validation:

- Parser test asserts counts near current observed counts: 65 rorg nodes, 106 funcs, 341 types, 2576 datafields.
- Assert known profile `F6-02-01` and `D5-00-01` fields exist.

### Phase 2 - Runtime bit helpers and manual codecs for 2-3 profiles

Implement smallest runtime before generating code:

- Bit get/set with explicit bit numbering tests. EEP uses `bitoffs`/`bitsize`; sample XML shows RPS and 1BS fields with bit offsets such as `Contact` at offset 7 size 1.
- Parse raw fields from payload/status into `uint64`.
- Format raw fields back into payload/status.
- Enum lookup and scaled conversion only where the sample profiles require it.

Validation:

- Unit tests for bit extraction/packing crossing byte boundaries.
- Round-trip parse -> format for chosen RPS/1BS/4BS payloads.

### Phase 3 - Generate codecs for broad EEP subset

Generator emits registry and codec definitions for profiles with straightforward cases:

- RPS (`0xF6`), 1BS (`0xD5`), 4BS (`0xA5`) first.
- Generate field metadata tables and generic table-driven parse/format code instead of thousands of custom functions. Lazy win: less generated code, easier review.
- Skip unsupported constructs with explicit generator warnings and generated metadata flag (`UnsupportedReason`) rather than silently wrong codecs.

Validation:

- Generator test on full `eep268.xml`.
- Generated compile test: `go generate ./... && go test ./pkg/eep`.
- Coverage stats in test output: profiles generated, profiles skipped by reason.

### Phase 4 - Conditions/cases and VLD

Add selection logic:

- Support `<case>` alternatives and `<condition>` fields.
- Add VLD (`0xD2`) command/sub-payload style dispatch. This is analogous to DDF `SubsequentPayload` requirements at spec lines 1403-1412.
- Keep parse safe: if a case cannot be selected, return structured unknown/ambiguous error with raw fields.

Validation:

- Tests for at least one conditional RPS case and one VLD profile.

### Phase 5 - DDF V2.0 parser as metadata layer (only after EEP codecs work)

Implement `pkg/ddf` if the library must consume real DDF files:

- Parse `Enocean_Devices`, `Device Product_ID`, metadata, TX/RX/EURID/BaseID EEP lists.
- Link DDF EEP declarations to generated EEP codecs by `(Rorg, Func, Type)`.
- Add `MSC` bitfield support using same bit/enum/scaled runtime model.
- Add `GP`, device params, optional commands, RPC metadata only as metadata structs unless there is a concrete parse/format requirement.
- Add index.xml retrieval only if remote server lookup is explicitly requested.

Validation:

- XML unmarshal tests with minimal DDF snippets covering required Product_ID, TX/RX/EURID/BaseID, MSC bitfield, SupportedRPC.
- Product_ID validation test for exactly `0x` + 12 hex chars.

## Validation commands

Current baseline command run:

- `go test ./...` downloads dependencies and does not pass today. Failures are unrelated to this plan but must be tracked before claiming green:
  - `pkg/enums`: `TestParseCommonCommandFromByte` expects `0x34` invalid, parser accepts it, then test panics on nil error.
  - `pkg/event`: endian expectations fail for ManufacturerID/DeviceID-like fields.
  - `pkg/subtel`: packet parsing expectations fail (sender ID/user data/SubTel count).

For the next implementation agent, use targeted checks until baseline is repaired:

- `go test ./pkg/eep ./pkg/erp1 ./pkg/esp3`
- `go test ./internal/eepgen ./cmd/eepgen ./pkg/eep` once generator exists.
- `go generate ./...` if generator is wired with `//go:generate`.

## Risks and decisions for planner

- **DDF vs EEP confusion:** `eep268.xml` describes standard EEP profiles; DDF V2.0 describes device metadata and references EEP/GP/MSC profiles. Do EEP codegen first; DDF parser later.
- **Bit numbering:** Must be proven with XML examples and EEP PDF semantics before broad generation. This is the highest correctness risk.
- **Conditions/cases:** Many profiles have conditional layouts; unsupported conditions must fail closed, not parse wrong.
- **VLD/subsequent payload:** VLD and MSC dynamic payloads need command/type-based secondary dispatch; this maps to DDF `SubsequentPayload` but should be phased after fixed layouts.
- **Generated API shape:** Keep table-driven metadata + generic parser unless benchmarks prove per-profile custom code is needed.
- **Existing failing tests:** Full repo is not green before this work. Do not mask unrelated failures.
- **Licensing/spec use:** The DDF spec text contains non-commercial/IP restrictions. Product use may need EnOcean Alliance membership/legal review.

## Handoff meta-prompt for next agent

Goal: Implement the smallest library slice that generates Go EEP profile codecs from `eep268.xml`: parse XML into a normalized model, generate a registry and table-driven parse/format support for a small validated subset, and leave unsupported constructs explicit.

Context/evidence:

- Current EEP identifier is in `pkg/eep/eep.go:21-100`; reuse it.
- ERP1 payload boundary is `pkg/erp1/erp1.go:22-67`, especially `UserData` at lines 50-55; generated codecs should operate on `UserData` plus `Status`.
- ESP3 framing is already handled in `pkg/esp3/esp3.go:61-137`; do not duplicate it.
- RORG constants are in `pkg/enums/rorg.go:7-23`.
- DDF V2.0 EEP declarations are metadata only (`Rorg`, optional `Func`, optional `Type`) per spec lines 1049-1052; field layouts come from `eep268.xml`.
- `eep268.xml` is UTF-16LE and has 65 rorg / 106 func / 341 type / 2576 datafield elements.

Success criteria:

- Generator parses full `eep268.xml` without external DTD access.
- Generated code compiles.
- At least 2-3 simple profiles parse and format raw telegram payloads round-trip.
- Unsupported profile constructs are reported/skipped explicitly.
- Existing ESP3/ERP1 behavior remains unchanged.

Hard constraints:

- Do not broaden into serial port orchestration, remote DDF server fetch, or full ReCom config unless explicitly requested.
- Do not silently parse fields with unsupported conditions/scales; return an error or skip with reason.
- Do not duplicate ESP3/ERP1 packet parsing.

Suggested approach:

- Build `cmd/eepgen` + internal parser first.
- Add minimal `pkg/eep` bit helpers and table-driven codec runtime.
- Generate metadata tables for a chosen profile subset, then broaden.
- Keep generated code boring and deterministic; stable sort profiles/fields.

Validation:

- Start with `go test ./pkg/eep ./pkg/erp1 ./pkg/esp3`.
- Add generator/parser tests with count assertions and known profiles.
- Run full `go test ./...` but document existing unrelated failures until fixed.

Stop/escalation rules:

- Stop and ask if expected public API shape must be stable before implementation.
- Stop and ask if full DDF V2.0 device metadata is required in the first deliverable; it is larger than EEP codegen.
- Stop broad generation if bit numbering cannot be confirmed by tests.

Resolved assumptions:

- First useful deliverable is EEP codegen from `eep268.xml`, not complete DDF V2.0 server/index/ReCom support.
- Generated codecs should integrate at ERP1 `UserData` layer, not raw serial or ESP3 layer.

## Structured acceptance report

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only and did not modify project/source files; only the required plan file was written."
    }
  ],
  "changedFiles": [
    "plans/device-description-file.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Baseline repo is not green: failures in pkg/enums, pkg/event, and pkg/subtel unrelated to this planning-only audit."
    },
    {
      "command": "git status --short",
      "result": "passed",
      "summary": "Checked working tree before writing plan; existing untracked docs/, eep268.xml, and main.go were present."
    }
  ],
  "validationOutput": [
    "Spec/DDF requirements and eep268.xml structure inspected.",
    "Current Go support audited: ESP3, ERP1, EEP identifier, RORG enums exist; EEP XML codegen and profile codecs are missing."
  ],
  "residualRisks": [
    "Full repo tests currently fail before implementation work.",
    "Bit numbering and conditional profile semantics need confirmation with focused fixtures before broad generation.",
    "DDF V2.0 licensing/IP terms may need legal review for commercial use."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added planning handoff only; no source implementation changes.",
  "reviewFindings": [
    "no blockers for planning-only task"
  ],
  "manualNotes": "This plan intentionally scopes first implementation to EEP codegen from eep268.xml; full DDF V2.0 device metadata is a later phase."
}
```
