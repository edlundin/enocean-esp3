# EEP 3.1 / eep268.xml implementation plan

## Scope read

User asked for an implementation plan only. I did not modify project/source files. This plan targets a library that can generate Go code from `eep268.xml` and use that generated code to parse and format ERP1 user-data telegram payloads for the profiles described by the XML, while preserving the existing ESP3/ERP1 layers.

## Current support in repo

- ESP3 packet framing exists in `pkg/esp3/esp3.go`:
  - `Telegram{PacketType, Data, OptData}` at lines 61-76 can serialize an ESP3 frame.
  - `NewEsp3TelegramFromHexString` at lines 79-137 validates sync byte and CRCs, then slices `Data` / `OptData`.
- ERP1 extraction exists in `pkg/erp1/erp1.go`:
  - `Packet` at lines 11-20 has `Rorg`, `UserData`, `SenderID`, `Status`, opt-data fields.
  - `NewPacketFromEsp3` at lines 22-67 extracts RORG from `Data[0]`, user data from between RORG and sender ID, sender ID and status.
  - `ToEsp3` at lines 69-94 writes RORG + `UserData` + sender ID + status.
  - This is the right boundary for EEP decoding/encoding: EEP logic should consume/produce `erp1.Packet.UserData`, not reimplement ESP3 framing.
- EEP type is currently only a triplet formatter/parser in `pkg/eep/eep.go`:
  - `EEP{Rorg, Func, Type}` at lines 21-25.
  - `FromString` and `String` format `RR-FF-TT` at lines 50-100.
  - Current bounds are stale for EEP 3.x: `maxFunc=0x60`, `maxType=0x7f` at lines 15-18, while the 3.1 spec says all fields are 8 bits.
- RORG enum coverage exists in `pkg/enums/rorg.go`:
  - Constants match the EEP 3.1 telegram list at lines 7-22.
  - Parser/string/valid switches at lines 25-120.
- Serial parser is not a reusable library API yet:
  - `OpenSerialPort` starts a goroutine and returns only the port at `pkg/enocean.go:24-49`.
  - Parsed telegrams are printed at `pkg/enocean.go:205-210` instead of returned on a channel/callback.
  - This is adjacent to the library goal but not required for code generation/parsing payloads; defer unless the public API must stream serial input.
- `main.go` is a hard-coded example executable (`/dev/cu.usbserial-EO7BBSFR` at line 26) in the module root. If this repo is to be imported cleanly as a library, move it under `cmd/...` or delete it in a later library-shaping phase.

## EEP 3.1 requirements that matter for code generation

Sources: `/tmp/enocean-spec-text/EnOcean-Equipment-Profiles-3-1.txt` and `eep268.xml`.

- EEP identity is triplet `(RORG, FUNC, TYPE)`:
  - Spec lines 194-198: profile elements are ERP radio telegram type (RORG), FUNC, TYPE.
  - Spec lines 209-215: from EEP 3.0 all fields are 8 bits; this affects 4BS teach-in if `FUNC > 0x3F` or `TYPE > 0x7F`.
  - Implementation impact: change `pkg/eep` validation to allow `Func` and `Type` 0x00..0xFF. Generated data should not inherit the old 0x60/0x7F caps.
- XML 2.6.8 remains the profile-description source:
  - Spec lines 24-25: profiles already defined are listed in former XML EEP-Specification 2.6.8 and still valid for descriptions.
  - `eep268.xml` is UTF-16LE (`<?xml ... encoding="utf-16le"?>`) and starts profiles at line 470.
  - Top-level profile RORGs in XML are only `F6`, `D5`, `A5`, `D2`:
    - `eep268.xml:471-...` RPS `0xF6`.
    - `eep268.xml:2442-...` 4BS `0xA5`.
    - `eep268.xml:21636-...` VLD `0xD2`.
    - Scripted count from XML: RPS 6 funcs / 14 types / 19 cases; 1BS 1 func / 2 types / 2 cases; 4BS 17 funcs / 133 types / 151 cases; VLD 21 funcs / 184 types / 187 cases.
- Telegram user-data lengths and bit order:
  - Spec lines 478-482: RPS and 1BS carry 1-byte user data; 1BS has an LRN bit and status-byte differences.
  - Spec lines 497-505: 4BS carries 4 bytes; radio order is DB_3 first, DB_0 last; offsets are in data-flow order (`DB_3.BIT_7` offset 0, `DB_0.BIT_3` offset 28), while bit value inside byte is normal MSB-to-LSB valuation.
  - Spec lines 506-510: VLD carries variable payload 1..14 bytes; example with 6 bytes has `DB_5.BIT_7` as offset 0.
  - Implementation impact: a single offset/size bit extractor over transmitted `UserData` bytes should work for RPS/1BS/4BS/VLD if offset 0 maps to bit 7 of byte 0. Do not reverse 4BS bytes before extracting; the XML/spec offsets are already in radio data-flow order.
- XML fields to model:
  - `eep268.xml:546-549` first `<datafield>` has `<bitoffs>` and `<bitsize>`.
  - `eep268.xml:2617`/`:2664` have `<scale>` and `<range>` for measurements.
  - XML datafield children include `data`, `shortcut`, `description`, `info`, `reserved`, `range`, `scale`, `unit`, `enum`, `value`.
  - XML case children include `condition`, `statusfield`, `datafield`, `direction` (`eep268.xml:12094`), title/description/status in some cases.
  - XML enum items can be individual values or ranges; spec lines 1148-1191 allow finite enums and enum ranges (e.g. `0...127`, `128...255`).
  - Spec lines 1239-1265 describe scale/range usage and prohibit extending scale limits via overflow.
  - Spec lines 1276-1305 define valid directions `FROM` and `TO` for commands; XML encodes these as direction values `1`/`2` in existing profile cases.
- Teach-in / security scope:
  - Spec lines 568-593: teach-in is a process; only 1BS/4BS reserve `DB_0.BIT_3` as LRN bit.
  - Spec lines 621-628: 1BS `DB_0.BIT_3` 0 teach-in, 1 data.
  - Spec lines 641-709: 4BS teach-in variations; if `FUNC > 0x3F` or `TYPE > 0x7F`, UTE must be used instead of 4BS teach-in.
  - Spec lines 758-887 and 888-915: UTE uses RORG `0xD4`, 8-bit RORG/FUNC/TYPE, query/response payload layouts and 500/700ms timing.
  - Spec lines 927-978: SEC/SEC_ENCAPS can wrap existing EEP data; EEPs do not require a particular security level. For encrypted SEC data, EEP parsing can only run after decryption/unencapsulation.
  - Implementation impact: phase 1 should parse operational RPS/1BS/4BS/VLD data telegrams. Teach-in helpers and SEC/UTE support should be separate small phases, not mixed into profile codegen.

## Missing spec requirements / gaps

1. No XML parser or generator exists.
2. No generated profile registry exists for `eep268.xml` profiles.
3. No EEP payload decoder/encoder exists; current `pkg/eep` only parses `RR-FF-TT` strings.
4. `pkg/eep` Func/Type bounds conflict with EEP 3.1 8-bit fields.
5. No generic bitfield read/write helpers exist for MSB-first radio payload offsets.
6. No generated support for:
   - datafield metadata (`bitoffs`, `bitsize`, descriptions, shortcuts),
   - reserved fields,
   - conditions and status fields,
   - enum values/ranges,
   - linear range-to-scale conversion and inverse formatting,
   - units,
   - direction/case selection.
7. No generated validation for user-data length by RORG/case (RPS/1BS 1 byte, 4BS 4 bytes, VLD variable up to 14 bytes).
8. No teach-in helpers for 1BS/4BS/UTE.
9. No SEC/SEC_ENCAPS unwrapping/decryption boundary; should be documented as out of initial generated parser scope unless security layer is added.
10. Public library API is thin: ESP3/ERP1 are importable, but root package serial parsing only prints frames.

## Smallest implementation phases

### Phase 0: make current EEP type spec-correct

Files:
- `pkg/eep/eep.go`
- `pkg/eep/eep_test.go`

Plan:
- Change `maxFunc` and `maxType` to `0xff`.
- Keep `EEP` triplet shape and `String()` unchanged.
- Add/adjust tests for `FF-FF-FF`, UTE-style `FUNC/TYPE` > old caps, and invalid values via string larger than byte.

Validation:
- `go test ./pkg/eep ./pkg/enums`

### Phase 1: add tiny runtime payload codec primitives

Files:
- new `pkg/eep/bitfield.go`
- new `pkg/eep/value.go` or keep in same small file if concise
- tests in `pkg/eep/bitfield_test.go`

Plan:
- Add MSB-first bit extraction/insertion over `[]byte` where offset 0 is bit 7 of byte 0.
- Support up to 64-bit raw fields initially. XML `bitsize` values fit typical sensor/control fields; generator should reject wider fields explicitly if any appear.
- Add linear conversion helpers:
  - decode physical value from raw using XML `range min/max` and `scale min/max`.
  - encode raw from physical with range checks and rounding policy documented/tested.
- Keep it generic; no generated code yet.

Validation:
- Unit test using spec 4BS mapping: 4-byte payload, offset 0 maps to `DB_3.BIT_7`, offset 28 maps to `DB_0.BIT_3`.
- Test linear decode with XML temperature example (`range 255..0`, `scale -40..0`).

### Phase 2: parse `eep268.xml` into an internal IR

Files:
- new `internal/eepxml/parser.go`
- new `internal/eepxml/model.go`
- tests/fixtures in `internal/eepxml`

Plan:
- Use stdlib `encoding/xml`; no new dependency. Since the file is UTF-16LE, wrap input with `unicode.UTF16(unicode.LittleEndian, unicode.ExpectBOM).NewDecoder()` from `golang.org/x/text/encoding/unicode` only if already available indirectly? It is not in `go.mod`; to avoid a dependency, decode UTF-16LE with `unicode/utf16` + byte pairs in ~15 lines before `xml.Unmarshal`/`Decoder`.
- Model only what codegen needs:
  - RORG number/title/telegram/fullname,
  - FUNC number/title/description,
  - TYPE number/title/status/description,
  - Case title/description/status/direction, conditions/statusfields/datafields,
  - DataField: name/shortcut/description/info/bitoffs/bitsize/value/reserved/range/scale/unit/enums.
- Normalize hex numbers to byte, directions to enum/string, ranges to inclusive numeric ranges.
- Generate a count/report test that asserts current XML counts: 4 profile RORGs, 45 funcs, 270 profile types, 359 cases (from profile RORGs; exact case count should be measured in the test after parser normalization).

Validation:
- `go test ./internal/eepxml`
- A test that reads repository `eep268.xml` and validates representative profiles: `F6-01-01`, `A5-02-05`, one VLD profile.

### Phase 3: generator emits metadata-only Go registry

Files:
- new `cmd/eepgen/main.go`
- new `internal/eepgen/generate.go`
- new generated package, suggested `pkg/eep/generated` or `pkg/eep/profiles`
- add `//go:generate go run ./cmd/eepgen -in ../../eep268.xml -out profiles_gen.go` near generated package

Plan:
- Generate one boring Go file containing compact structs/slices/maps, not one Go type per profile yet.
- Public lookup API in `pkg/eep`:
  - `Lookup(EEP) (*Profile, bool)`
  - `Profiles() []Profile`
- Include metadata enough for parse/format: case fields, conditions, enums, ranges/scales/units, directions.
- Ensure generated output is deterministic: stable profile/case/field ordering and gofmt.

Validation:
- Golden test or generated-file smoke test.
- `go generate ./...` then `go test ./pkg/eep ./internal/eepxml ./internal/eepgen`.

### Phase 4: generic parse/format by generated metadata

Files:
- `pkg/eep/decoder.go`
- `pkg/eep/encoder.go`
- generated metadata package from phase 3
- tests in `pkg/eep`

Plan:
- Add API with minimal shape:
  - `func Decode(profile EEP, userData []byte, status byte) (Telegram, error)`
  - `func Encode(profile EEP, values map[string]Value) ([]byte, error)` or a small typed `FieldValue` slice if maps are too loose.
- `Decode` selects a case by evaluating statusfields/conditions; if ambiguous, return all matching cases or an explicit ambiguity error. Lazy option: return first exact match and include case ID/name; only add multi-case API if tests reveal ambiguity.
- Output field values by shortcut/IP-key/name with raw value, enum label if matched, scaled value/unit if scale exists.
- `Encode` sets only non-reserved fields, validates enums/ranges/scales, and returns userData bytes. The caller can pass returned bytes into `erp1.Packet{Rorg: profile.Rorg, UserData: ...}.ToEsp3()`.

Validation:
- Decode known examples from XML profiles:
  - RPS `F6-01-01` push button bit.
  - 1BS `D5-00-01` LRN bit/data field.
  - 4BS temperature profile (e.g. `A5-02-05`) with scaled values.
  - One VLD profile with variable length.
- Encode/decode round trip for those profiles.

### Phase 5: teach-in helpers, separately

Files:
- `pkg/eep/teachin.go`
- tests in `pkg/eep/teachin_test.go`

Plan:
- Add helpers to detect data vs teach-in:
  - 1BS/4BS LRN bit (`DB_0.BIT_3` offset 4 for 1BS, offset 28 for 4BS per spec).
  - 4BS LRN type / bidirectional bits (`DB_0.BIT_7..3`, offsets 24..28).
- Add UTE query/response structs for RORG `D4` exactly as spec lines 835-915.
- Enforce 4BS teach-in limitation: if FUNC > `0x3F` or TYPE > `0x7F`, use UTE instead.

Validation:
- Unit tests for 1BS, 4BS unidirectional/profile/bidirectional bit layouts, and UTE query/response round trips.

### Phase 6: library API cleanup, only if desired

Files:
- `pkg/enocean.go`
- `main.go` -> `cmd/enocean-esp3/main.go` or remove example

Plan:
- Make serial parser return frames over a channel or accept a callback instead of `fmt.Println`.
- Keep generated EEP decode independent: serial API should deliver `esp3.Telegram` or `erp1.Packet`; callers opt into EEP decode by profile.

Validation:
- Add parser tests using an in-memory reader/fake serial port if current `serial.Port` interface allows it.

## Suggested public API shape

Keep it small; generated structs should be data, runtime code should be generic:

```go
profile, ok := eep.Lookup(eep.MustParse("A5-02-05"))
msg, err := eep.Decode(profile.ID, erp.UserData, erp.Status)

userData, err := eep.Encode(profile.ID, map[string]eep.Value{
    "TMP": eep.Float(21.5),
})
packet := erp1.Packet{Rorg: profile.ID.Rorg, UserData: userData, SenderID: id, Status: status}
```

Do not generate hundreds of profile-specific structs first. A metadata-driven decoder covers the XML with far less code and keeps generated output stable.

## Implementation risks

- XML irregularity: old XML has mixed HTML-ish descriptions and optional/missing tags. Parser must ignore unknown markup and preserve text only where useful.
- UTF-16LE input: stdlib XML decoder does not transparently decode UTF-16 without a charset reader; handle this explicitly in `internal/eepxml`.
- Case selection: some profiles use status fields/conditions/directions. Ambiguous matching must be reported clearly rather than silently decoding the wrong case.
- Scaling direction: XML ranges can be descending (`255..0`) while scales ascend/descend. Tests must cover both.
- Value keys: `shortcut` is not guaranteed unique across every case/profile. Internally key by profile+case+field index; expose shortcut/name as metadata, not as the only identity.
- Security: encrypted SEC payloads cannot be decoded as EEP data until security handling decrypts/unwraps them. SEC_ENCAPS needs original RORG restored before EEP profile lookup.
- Current full test suite is already red in unrelated packages; use targeted tests during this work and record full-suite residual failures.

## Validation baseline

Commands run during audit:

- `go test ./...` failed before any source changes. Failures observed:
  - `pkg/enums`: `TestParseCommonCommandFromByte` expects error for `0x34`, got nil, then nil pointer panic in test.
  - `pkg/event`: endian/order mismatches in SmartAck/security event parsing tests.
  - `pkg/subtel`: sender/user-data/subtel count mismatches.
- `git status --short` before writing this plan showed pre-existing untracked `docs/`, `eep268.xml`, `main.go`. No staged files (`git diff --cached --stat` empty).

Targeted validation once implementation starts:

1. `go test ./pkg/eep ./internal/eepxml ./internal/eepgen`
2. `go generate ./... && git diff --exit-code -- pkg/eep/profiles_gen.go` in CI after generated file is checked in.
3. `go test ./pkg/esp3 ./pkg/erp1 ./pkg/eep` for ESP3/ERP1/EEP integration.
4. Full `go test ./...` after unrelated existing failures are fixed or accepted as known failures.

## Meta-prompt handoff for next agent

Goal: Implement the library/codegen plan above in the smallest useful phases. Start with EEP 3.1 8-bit triplet correctness, bitfield helpers, XML parser, generated metadata registry, then generic decode/encode. Do not begin with hundreds of hand-written profile structs.

Context/evidence:
- Existing reusable layers: `pkg/esp3/esp3.go:61-137`, `pkg/erp1/erp1.go:11-94`, `pkg/eep/eep.go:21-100`, `pkg/enums/rorg.go:7-120`.
- Spec requires 8-bit RORG/FUNC/TYPE (`EnOcean-Equipment-Profiles-3-1.txt:194-215`), RPS/1BS 1 byte (`:478-482`), 4BS 4 bytes with DB_3 first and offset 0 at `DB_3.BIT_7` (`:497-505`), VLD 1..14 bytes (`:506-510`), enum ranges (`:1148-1191`), scale rules (`:1239-1265`), directions (`:1276-1305`), teach-in bits (`:568-709`), UTE (`:758-915`), and SEC caveats (`:927-978`).
- XML is UTF-16LE, profile data begins at `eep268.xml:470`; key tags are `rorg/func/type/case/datafield/bitoffs/bitsize/range/scale/enum/direction`.

Success criteria:
- Generated Go metadata from `eep268.xml` can be looked up by EEP triplet.
- Runtime can decode raw `erp1.Packet.UserData` + status into field values for representative RPS/1BS/4BS/VLD profiles.
- Runtime can encode field values back to userData for the same representative profiles.
- No ESP3 framing duplication; all frame serialization/parsing remains in existing `esp3`/`erp1` packages.

Hard constraints:
- Keep dependencies unchanged unless stdlib cannot reasonably solve UTF-16/XML parsing. Prefer a tiny UTF-16LE decoder over adding a module.
- Generated output must be deterministic and gofmt'd.
- Do not mix teach-in/security implementation into initial profile decode/encode unless explicitly requested.

Suggested approach:
- Phase the work exactly enough to keep tests green in touched packages.
- Build a generic metadata-driven decoder before considering generated typed structs.
- Reuse existing `eep.EEP` and `erp1.Packet` as boundaries.

Stop/escalation rules:
- Ask for a decision before changing public API names if there are multiple reasonable options.
- Stop after metadata registry + representative decode/encode works; broad profile-specific ergonomics can come later.
- If XML contains fields wider than the chosen raw integer size, fail generation with a clear error and report the profiles.

Resolved assumptions:
- `eep268.xml` is authoritative for concrete profiles; EEP 3.1 text supplies rules and corrections.
- Initial library does not decrypt SEC telegrams.
- Initial serial streaming cleanup is optional and separate from EEP generation.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only at plans/eep-3-1.md; no project/source files were modified."
    }
  ],
  "changedFiles": [
    "plans/eep-3-1.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Baseline suite already fails in pkg/enums, pkg/event, and pkg/subtel; failures recorded in plan."
    },
    {
      "command": "git status --short && git diff --stat && git diff --cached --stat",
      "result": "passed",
      "summary": "Confirmed no staged files; pre-existing untracked docs/, eep268.xml, main.go were visible before writing the plan."
    }
  ],
  "validationOutput": [
    "Plan file written to /Users/edlundin/work/e/enocean-esp3/plans/eep-3-1.md",
    "No source files intentionally changed."
  ],
  "residualRisks": [
    "Full test suite has unrelated pre-existing failures; implementation validation should use targeted package tests until those are fixed."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added implementation plan only.",
  "reviewFindings": [
    "no blockers for planning task"
  ],
  "manualNotes": "Acceptance criterion says implement requested change; requested change was an implementation plan only, not source changes."
}
```
