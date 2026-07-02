# Generic Profiles / eep268.xml implementation plan

## Scope read

- Spec: `/tmp/enocean-spec-text/GenericProfilesspecification_1.4.txt`.
- XML: `eep268.xml` is UTF-16LE; parsed as `/tmp/eep268.utf8.xml` for audit only.
- Go repo: ESP3, ERP1, EEP triplet, enums, serializer/deserializer, tests.

Goal is a library that can generate Go code from `eep268.xml` to parse and format telegram payloads for the profiles in that XML, while also supporting Generic Profiles v1.4 telegram/message rules where applicable. No source files were modified for this plan.

## Current support

- ESP3 frame parse/format exists:
  - `pkg/esp3/esp3.go:61-68` serializes ESP3 frames with CRC8H/CRC8D.
  - `pkg/esp3/esp3.go:79-130` parses hex ESP3 frames and validates sync + CRCs.
- ERP1 packet extraction/format exists:
  - `pkg/erp1/erp1.go:11-20` has RORG, sender, destination, status, user data fields.
  - `pkg/erp1/erp1.go:22-67` extracts ERP1 data into `Packet`.
  - `pkg/erp1/erp1.go:69-94` converts `Packet` back to ESP3.
- EEP support is only an identifier helper:
  - `pkg/eep/eep.go:21-25` defines `EEP{Rorg, Func, Type}`.
  - `pkg/eep/eep.go:27-40` and `:50-100` parse/format triplets.
  - There is no XML parsing, profile registry, field extraction, scaling, enum mapping, case selection, or generated code.
- RORG enum is missing Generic Profile RORGs:
  - `pkg/enums/rorg.go:7-23` has EEP/ERP RORGs but not GP `0xB0`..`0xB3`.
  - `pkg/enums/rorg.go:25-60`, `:62-97`, `:99-120` must be extended if GP telegrams are represented by `enums.Rorg`.
- Serial-port API is not yet a reusable telegram stream:
  - `pkg/enocean.go:24-48` opens a serial port and starts `parser`.
  - `pkg/enocean.go:51-190` is an internal state machine that prints parsed telegrams instead of returning them to callers; do not build generated profile parsing on this path first.
- Existing generic serializer/deserializer is byte-oriented reflection for common commands, not bitfield-oriented payload parsing:
  - `internal/serializer/serializer.go` and `internal/serializer/deserializer.go` can be reused as patterns for no-global config, but not for EEP bit offsets.

## Relevant spec requirements

### Generic Profiles v1.4

- Four GP message types, each max 512 bytes; each can be addressed with ADT: spec `:255-267`.
- GP API selects RORG and splits oversized messages into chained telegrams: spec `:274-287`.
- GP RORGs: `0xB0 GP_TI`, `0xB1 GP_TR`, `0xB2 GP_CD`, `0xB3 GP_SD`; all allow chaining, with broadcast/unicast restrictions: spec `:294-311`.
- Chaining is needed when payload exceeds ERP1 single-telegram payload; chained telegrams use `RORG CDM`, `SEQ`, `IDX`, first-fragment `LEN` and embedded GP RORG; `SEQ=00` is forbidden and `IDX` starts at 0: spec `:313-317`, `:340-348`.
- All signed OTA numbers are two's complement; all frames/bytes are big endian: spec `:503-507`.
- GP channels have channel type, signal type, and value type: spec `:508-549`.
- Channel types: `00 Teach-in information`, `01 Data`, `10 Flag`, `11 Enumeration`: spec `:523-529`.
- Resolution coding for Data/Enumeration supports 2,3,4,5,6,8,10,12,16,20,24,32 bits; Flag is implicit 1 bit: spec `:553-584`.
- Teach-in request header is 16 bits: manufacturer 11, data direction 1, purpose 2, unused 2: spec `:900-948`.
- Channel definitions are bit-packed with no byte alignment between definitions:
  - Data: 40 bits, fields `2+8+2+4+8+4+8+4`: spec `:949-979`.
  - Flag: 12 bits, fields `2+8+2`: spec `:986-995`.
  - Enumeration: 16 bits, fields `2+8+2+4`: spec `:996-1008`.
  - Teach-in information: `18+N` bits, fields `2+8+8+N`, does not affect operational indexing: spec `:1009-1022`.
- Teach-in response header is 16 bits: manufacturer 11, result 2, undefined 3; result values include general reject, success, teach-out, rejected channel list: spec `:1036-1045`.
- If channels are rejected, response includes one acknowledgement bit per indexed channel, same order as definitions; only sent for result `11`: spec `:1077-1100`.
- Channel indexes start at first outbound channel index 0 and continue through inbound channels; Teach-in information channels are not indexed: spec `:1101-1106`.
- Timing values matter for a full teach-in session state machine, but can be postponed for a payload-codegen library: transmitter timeout 750 ms, receiver response time 500 ms, receiver timeout 750 ms, transmitter response time 500 ms: spec `:1122-1135`.
- Complete data message concatenates every channel value from channel 0, preserves MSB order, pads zeros to next byte; no channel IDs in the payload: spec `:1605-1625`.
- Selective data message starts with a 4-bit channel count; each selected entry is 6-bit channel index + channel value bits; preserves MSB order and pads to byte: spec `:1636-1663`.
- Remote Management minimum support is required only for continuously powered GP devices; for this library plan, model it as future integration, not phase-1 profile payload parsing.

### `eep268.xml`

- Encoding is UTF-16LE with XML stylesheet and DTD; Go generator should use `encoding/xml` over an `encoding/unicode.UTF16` reader or decode first. Do not assume UTF-8.
- Parsed profile counts from XML:
  - 4 RORGs under `<profile>`: `0xF6`, `0xD5`, `0xA5`, `0xD2`.
  - 45 funcs, 270 profile types, 359 cases, 2544 datafields.
  - RORG breakdown: `F6/RPS` 6 funcs / 14 types, `D5/1BS` 1 / 1, `A5/4BS` 17 / 132, `D2/VLD` 21 / 123.
- XML hierarchy and fields:
  - RORG/function/type/case/datafield example starts at `/tmp/eep268.utf8.xml:470-565`.
  - Datafields carry `bitoffs`, `bitsize`, optional `reserved`, `enum`, `range`, `scale`, `unit`, `value`, `info`: examples `/tmp/eep268.utf8.xml:555-565`, `:2656-2666`.
  - 4BS teach-in variations are embedded in `<teachin>` before normal funcs: `/tmp/eep268.utf8.xml:2447-2600`.
  - XML can include profile-level metadata like manufacturer, datalength, broadcast, addressable, answer in RPC section: `/tmp/eep268.utf8.xml:61060-61090`. Treat RPC as separate/out of phase unless explicitly required.
- Important XML modeling constraints:
  - `bitoffs` are EEP DB bit offsets, not Go struct byte offsets.
  - Some values include whitespace or formulas (`"0x04 + N"`, `"N*9"`) in RPC; keep generator tolerant.
  - Cases may have conditions on statusfields or direction. Generated parsers must choose a matching case or report ambiguity/no match.
  - Reserved fields should participate in formatting validation where XML gives fixed values, but should not appear as public semantic fields unless useful.

## Missing requirements / gaps

1. No generated profile model or generator from `eep268.xml`.
2. No bit-level reader/writer that handles arbitrary offset/size, MSB order, padding, and signed two's-complement values.
3. No EEP datafield scale/range conversion. XML fields need raw integer parse/format and physical value conversion using `range` -> `scale` for EEP, while GP spec needs quantization formulas for GP channels.
4. No enum support for datafield values.
5. No case selection based on statusfield/datafield conditions.
6. No RORG-specific payload length/status handling beyond simple ERP1 slicing.
7. No GP RORG constants, GP message structs, GP teach-in headers, GP channel definitions, complete/selective data message codecs, or GP chaining helpers.
8. No library API shape for consumers to parse an `erp1.Packet` into a typed/generated profile result or format values back to an ERP1/ESP3 telegram.
9. Existing serial parser prints instead of returning telegrams; usable library API should stay payload-level first.
10. Existing `go test ./...` currently fails in unrelated packages (`pkg/enums`, `pkg/event`, `pkg/subtel`), so new validation needs targeted package tests until baseline is fixed.

## Smallest implementation phases

### Phase 1 — profile XML model and generator smoke test

Files to add/change:
- Add `internal/eepxml` or `internal/eepgen/xmlmodel.go` for the minimal XML structs needed: profile/rorg/func/type/case/statusfield/datafield/enum/range/scale.
- Add `cmd/eepgen/main.go` with flags like `-xml eep268.xml -out pkg/eep/generated`.
- Add `internal/eepgen` tests with a tiny UTF-16 fixture and one real-XML smoke test behind normal unit test if file exists.

Minimum output:
- Generator can read `eep268.xml`, count profiles, and emit deterministic `profiles_gen.go` metadata for RORG/FUNC/TYPE/case/fields.
- No parsing/formatting yet.

Validation:
- `go test ./internal/eepgen ./internal/eepxml`.
- `go run ./cmd/eepgen -xml ./eep268.xml -out /tmp/eepgen-check`.

### Phase 2 — bitfield core, no generated code dependency

Files to add/change:
- Add `pkg/bitfield` or unexported `internal/bitfield` if only generated code uses it.
- Tests for extract/insert across byte boundaries, MSB order, signed two's complement, padding.

Minimum output:
- `ReadUnsigned(data []byte, bitOffset, bitSize int)`, `ReadSigned`, `WriteUnsigned`, `WriteSigned`.
- Keep this tiny; do not use reflection.

Validation:
- Unit tests with examples matching XML offsets (`bitoffs=16`, `bitsize=8`) and GP complete/selective packing examples from spec.

### Phase 3 — generated EEP metadata + dynamic parser/formatter

Files to add/change:
- Add `pkg/eep/profile.go` for public types: `ProfileID`, `Profile`, `Case`, `Field`, `EnumItem`, `Scale`, `Value`.
- Add `pkg/eep/parse.go` and `format.go` using generated metadata + bitfield core.
- Generator emits `pkg/eep/generated/profiles_gen.go` or `pkg/eep/profiles_gen.go`.

Minimum public API:
- `eep.Lookup(id EEP) (*Profile, bool)`.
- `eep.Parse(id EEP, userData []byte, status byte) (Message, error)`.
- `eep.Format(id EEP, fields map[string]any) ([]byte, byte, error)` or a typed generated wrapper later.
- `Message` should include selected case, raw fields, enum labels, and scaled values where present.

Smallest correct semantics:
- Parse/format EEP field raw values by `bitoffs/bitsize`.
- Apply enum descriptions.
- Apply linear range/scale conversion for fields that have both.
- Enforce fixed `value` and case conditions.
- Ignore RPC section initially; document as unsupported.

Validation:
- Golden tests for `F6-01-01` push button (`/tmp/eep268.utf8.xml:508-565`).
- Golden tests for `A5-02-01` temperature (`/tmp/eep268.utf8.xml:2602-2666`).
- Round-trip parse -> format -> parse for these profiles.

### Phase 4 — typed generated wrappers, only after dynamic parser works

Files to add/change:
- Extend generator to emit one Go type per profile/case or compact typed accessors if full type-per-profile gets noisy.
- Generated code should call the same `pkg/eep` runtime, not duplicate bit logic.

Minimum output:
- Stable Go identifiers from XML titles/shortcuts.
- Compile-time docs from XML descriptions where concise.
- Avoid generating hundreds of hand-maintained files; one generated package is fine.

Validation:
- `go generate ./...` or `go run ./cmd/eepgen ...` produces no diff after committed generation.
- Compile generated package with `go test ./pkg/eep/...`.

### Phase 5 — Generic Profiles runtime messages

Files to add/change:
- Extend `pkg/enums/rorg.go` with GP RORGs `B0`..`B3`.
- Add `pkg/gp` for GP-specific runtime: message types, teach-in request/response headers, channel definitions, complete/selective data payload codecs.
- Reuse bitfield core.
- Add `pkg/gp/chaining.go` only if actual ERP1 chained telegram assembly/disassembly is required at this layer.

Minimum output:
- Encode/decode teach-in request/response payloads exactly per spec.
- Encode/decode complete/selective data messages using GP channel definitions.
- Validate 512-byte max message length and channel index rules.

Validation:
- Tests using spec examples:
  - teach-in request header `0xFFF0` / toggle `0xFFF8` from examples.
  - data channel example `0x4195001051` and flag example `0x826` from spec pages around `:1360-1370`.
  - complete/selective bit packing examples from spec `:1614-1679`.

### Phase 6 — ERP1/ESP3 integration

Files to add/change:
- Add helpers that accept/return `erp1.Packet`, not serial ports.
- Optional: add a channel/callback API to `pkg/enocean.go` later; not needed for codegen payload library.

Minimum output:
- `eep.ParsePacket(packet erp1.Packet, id EEP)`.
- `gp.ParsePacket(packet erp1.Packet)` that routes by RORG `B0`..`B3`.
- `ToPacket` helpers for formatting.

Validation:
- Unit tests build `erp1.Packet` directly and verify RORG/userData/status; avoid flaky serial tests.

## Hard constraints for implementer

- Do not hand-code individual profiles from `eep268.xml`; generate metadata/code.
- Do not add a new dependency unless UTF-16 decoding cannot be done with already available stdlib/x packages. Prefer `encoding/xml` plus standard or `golang.org/x/text/encoding/unicode` if already acceptable; currently `golang.org/x/sys` only is installed, so adding x/text is a decision point.
- Do not start with the serial-port parser; it is not necessary for payload parsing/formatting.
- Do not implement Remote Management/RPC in phase 1-5 unless explicitly requested; XML RPC shape is different and would widen scope.
- Keep generated code deterministic and checked by tests.

## Implementation risks

- `eep268.xml` may not be the Generic Profiles appendix referenced by GP spec for signal types; it is the EEP 2.6.8 XML profile database. GP channel signal-type lists may need an additional appendix source if true GP teach-in signal names are required.
- Existing repo baseline tests fail unrelated to this work; full `go test ./...` is not a clean validation gate yet.
- EEP bit numbering can be easy to invert. Lock this down with golden tests from known EEP examples before generating all profiles.
- Generated typed wrappers can explode API surface; dynamic parser first is the lazy/safe base.
- Chaining references Remote Management SYS_EX and ERP2 details not implemented locally. Treat chaining as a separate helper after single-telegram payload codecs work.

## Suggested next-agent contract

Goal: implement Phase 1 and Phase 2 only unless told otherwise: XML loading/generator smoke output plus bitfield core tests. This creates the foundation without widening into hundreds of generated profile APIs.

Evidence to use:
- Current EEP helper only parses IDs: `pkg/eep/eep.go:21-100`.
- Current RORG enum lacks GP RORGs: `pkg/enums/rorg.go:7-120`.
- ESP3/ERP1 payload layers already exist: `pkg/esp3/esp3.go:61-130`, `pkg/erp1/erp1.go:22-94`.
- XML structure and counts above; use `/tmp/eep268.utf8.xml` line references only as audit evidence, not as a source file.

Success criteria:
- Generator reads UTF-16LE `eep268.xml` and emits deterministic metadata or a JSON/Go summary in a temp/generated target.
- Bitfield read/write handles arbitrary bit offsets/sizes and signed values with tests.
- No source-wide serial API refactor.

Validation:
- Targeted `go test` for new packages.
- `go run ./cmd/eepgen -xml ./eep268.xml -out /tmp/eepgen-check` succeeds.
- Run `git diff --exit-code` for committed files if generated output is expected to be checked in.

Stop/escalation:
- Ask before adding dependencies for UTF-16 decoding.
- Ask if the real Generic Profiles appendix is required and not represented by `eep268.xml`.
- Stop after Phase 2 if the goal was a small safe increment.

## Commands run during audit

- `python3` decoded `eep268.xml` as UTF-16LE to `/tmp/eep268.utf8.xml` for inspection and counted XML elements.
- `go test ./...` — failed in existing unrelated packages: `pkg/enums`, `pkg/event`, `pkg/subtel`; many packages passed.
- `grep -R "RorgGP\|0xB0\|Generic Profiles\|GP_TI\|Parse.*EEP\|eep268" ...` — no existing GP/codegen support found.
- `git status --porcelain=v1 && git diff --cached --name-only` — untracked `docs/`, `eep268.xml`, `main.go`; no staged files listed.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only at plans/generic-profiles.md; no project/source files were modified. Scope is phased around XML codegen, bitfield core, EEP dynamic parser, GP runtime, then integration."
    }
  ],
  "changedFiles": [
    "plans/generic-profiles.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "python3 decode/count script for eep268.xml",
      "result": "passed",
      "summary": "Confirmed UTF-16LE XML and counted 4 RORGs, 45 funcs, 270 types, 359 cases, 2544 datafields."
    },
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Existing baseline failures in pkg/enums, pkg/event, pkg/subtel; unrelated packages passed."
    },
    {
      "command": "grep -R \"RorgGP\\|0xB0\\|Generic Profiles\\|GP_TI\\|Parse.*EEP\\|eep268\" -n -- *.go pkg internal docs .github",
      "result": "passed",
      "summary": "No existing Generic Profiles/codegen support found."
    },
    {
      "command": "git status --porcelain=v1 && git diff --cached --name-only",
      "result": "passed",
      "summary": "Showed untracked docs/, eep268.xml, main.go and no staged files."
    }
  ],
  "validationOutput": [
    "Plan file written to /Users/edlundin/work/e/enocean-esp3/plans/generic-profiles.md.",
    "go test ./... currently fails before implementation work: pkg/enums commoncommand test panic, pkg/event endian expectations, pkg/subtel parsing expectations."
  ],
  "residualRisks": [
    "eep268.xml appears to be the EEP XML database, not necessarily the Generic Profiles appendix containing GP signal-type lists.",
    "Full repo test suite is not currently green, so implementation should use targeted tests until baseline is fixed."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added audit/implementation plan only; no source changes.",
  "reviewFindings": [
    "no blockers in the plan output"
  ],
  "manualNotes": "Repository already had untracked docs/, eep268.xml, and main.go before/after this audit; they were not staged."
}
```
