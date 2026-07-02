# eep268.xml generator context and plan

## Scope

Plan only. Do not modify project/source files. The useful first step is a generator that reads the checked-in `eep268.xml` and emits a small, idiomatic Go profile package for parsing/formatting ERP1 profile user payloads.

## Source facts from `eep268.xml`

- Encoding/header: UTF-16LE XML declaration, stylesheet and DTD refs at file start. Local `eep.dtd`/`eep.xsl` are referenced but not present in this repo, so generator should infer from the XML and not require the DTD.
- Profile tree starts at `eep/profile` around line 470.
- Profile RORGs present:
  - `0xF6` RPS, line 471; 6 funcs, 14 types, 19 cases.
  - `0xD5` 1BS, line 2355; 1 func, 1 type, 1 case.
  - `0xA5` 4BS, line 2441; 17 funcs, 132 types, 148 cases.
  - `0xD2` VLD, later in file; 21 funcs, 123 types, 187 cases.
- Counts observed under `profile`: 45 profile funcs, 270 profile types, 355 profile cases, 2435 profile case datafields, 26 profile case statusfields, 126 case conditions, 61 type refs.

### Inferred XML schema subset for profiles

```text
eep
  profile
    rorg
      number        // hex, e.g. 0xA5
      title
      oldnumber?
      fullname?
      telegram?     // RPS, 1BS, 4BS, VLD
      description?
      teachin?      // rorg-level pseudo type/case/datafield, not normal profile payload
      func*
        number      // hex, e.g. 0x02
        title
        description?
        type*
          number    // hex, e.g. 0x05; empty in teachin blocks
          title
          status?
          description?
          ref*      // rorg/func/type references to another profile definition
          case*
            title?
            description?
            direction?
            status?
            condition?       // one or more statusfield/datafield value predicates
            statusfield*     // fields from ERP1 Status byte, mostly RPS T21/NU
            datafield*       // fields from ERP1 UserData bytes
```

`datafield`/`statusfield` common shape:

```text
reserved?             // field should be skipped but bit range retained
value?                // for condition predicates
shortcut?             // good Go identifier seed, e.g. TMP, LRN, CO
data?                 // human label, e.g. Temperature
info?/description?
bitoffs               // integer, sometimes has whitespace
bitsize               // integer
range? min/max        // raw integer range; min can be greater than max
scale? min/max        // physical range; also appears as scalar in a few item nodes
unit?
enum? item*
```

Enum item forms are not uniform:

- Most common: `<value>...</value><description>...</description>`; e.g. F6-01-01 Push Button line 562.
- Ranges: `<min>...</min><max>...</max><description>...</description>`.
- Scaled enum/range items: item-local `scale` and `unit` also occur.
- Some malformed-for-codegen-ish items have only `description`; generator should preserve comments/metadata but not generate strict constants for those.

### Concrete profile snippets to anchor tests

- `F6-01-01` Push Button:
  - RORG `0xF6` at lines 471-476; func `0x01` line 500; type `0x01` line 509.
  - Case line 545 has two reserved datafields (`bitoffs` 0 size 3, `bitoffs` 4 size 4) and `PB` at bit 3 size 1 with enum `0=Released`, `1=Pressed & Hold` lines 556-571.
- `F6-02-01` Rocker Switch has multiple cases selected by status bits T21/NU:
  - Conditions use `statusfield bitoffs 2/3 bitsize 1 value ...` around lines 600-611 and repeated in later cases.
  - Datafields include R1 bits 0..2, EB bit 3, R2 bits 4..6, SA bit 7 around lines 1485-1632.
- `D5-00-01` Single Input Contact:
  - RORG `0xD5` line 2356; func `0x00` line 2385; type `0x01` line 2388.
  - Contact `CO` bit 7 size 1 enum `0=open`, `1=closed` lines 2393-2414.
  - Learn Button `LRN` bit 4 size 1 enum `0=pressed`, `1=not pressed` lines 2415-2436.
- `A5-02-01` Temperature Sensor -40..0 C:
  - Func `0x02` Temperature Sensors lines 2602-2604; type `0x01` lines 2605-2608.
  - Reserved fields at bitoffs 0 size 16, 24 size 4, 29 size 3 lines 2610-2639.
  - LRN bit at bitoffs 28 size 1 enum data/teach-in lines 2640-2656.
  - Temperature `TMP` bitoffs 16 size 8, raw range `255..0`, scale `-40..0`, unit `°C` lines 2657-2667.
- `A5-02-02`/`03` show same shape with different scale ranges (`-30..+10`, `-20..+20`) lines 2670-2810.

## Existing Go code constraints

- `pkg/eep/eep.go` lines 21-25 defines only the EEP triplet `{Rorg enums.Rorg, Func byte, Type byte}`.
- `pkg/eep/eep.go` lines 27-40 and 50-100 parse/format triplets (`RR-FF-TT`) and enforce FUNC <= `0x60`, TYPE <= `0x7f`.
- `pkg/enums/rorg.go` lines 7-22 already defines RORG constants for `RPS`, `1BS`, `4BS`, `VLD`, etc.; reuse these, do not generate duplicate RORG constants.
- `pkg/erp1/erp1.go` lines 11-20 models an ERP1 radio packet. Profile parsing should consume `erp1.Packet.Rorg`, `erp1.Packet.UserData`, and `erp1.Packet.Status`, not full ESP3 bytes.
- `pkg/erp1/erp1.go` lines 22-66 extracts `UserData` as bytes between RORG and sender ID. `ToEsp3` writes `Rorg + UserData + SenderID + Status` at lines 69-89.
- Internal reflection serializer/deserializer is byte-sequential and tag-based for common commands; it does not handle bitfields, scaling, ERP1 status conditions, or dynamic profile selection. Reusing it for EEP payloads would be the wrong abstraction.

## Bit/scaling model to generate

- Treat `bitoffs` as an offset into the ERP1 profile payload byte slice for `datafield`, and into `erp1.Packet.Status` for `statusfield`.
- Field bit index is least-significant-bit based within each byte: `byteIndex := bitoffs / 8`, `shift := bitoffs % 8`. This matches examples such as D5 contact at bit 7 and A5 LRN at bit 28 (`DB0` bit 4 in a 4-byte payload ordered DB3, DB2, DB1, DB0).
- Multi-byte fields should be extracted by a tiny generated/runtime helper that walks bits by global bit offset. Do not use `encoding/binary` structs for profile fields; XML fields are not byte-aligned in general.
- Scaled numeric conversion should support reversed ranges:
  - raw-to-physical: `scaled = scaleMin + (raw-rangeMin)*(scaleMax-scaleMin)/(rangeMax-rangeMin)`.
  - This naturally handles A5 temperature raw `255..0` to `-40..0`.
  - Keep raw value too, or expose encode/decode with rounding, because round-trip formatting needs exact bits.
- Enums should become typed Go constants only when all items have concrete `value` (or clean min/max ranges for a `String` classifier). Preserve descriptions as comments; skip unusable description-only items.
- Reserved fields should be generated as masks/validation comments, not exported struct fields.
- Conditions select one of multiple cases. Minimal first generator can return an error on ambiguous/no matching case; do not invent defaults beyond exact XML predicates.
- `ref` nodes should be resolved or explicitly skipped in first pass. There are 61 under profile types, so full generation needs a resolver eventually.

## Proposed package/generator layout

Minimal layout, few files:

```text
cmd/eepgen/main.go              // reads eep268.xml, emits generated Go
eep268.xml                      // input, UTF-16LE
pkg/eep                         // keep existing triplet code
pkg/eep/profiles                // generated and tiny hand-written runtime helpers
  bit.go                        // extract/insert bits, scale helper, Parse entrypoint
  profiles_gen.go               // registry: map[eep.EEP] Profile metadata/parser
  f6_01_01_gen.go               // generated profile code
  d5_00_01_gen.go
  a5_02_01_gen.go
```

Do not add dependencies. Go stdlib `encoding/xml`, `unicode/utf16`/BOM handling, `text/template`, and `go/format` are enough. The lazy safe generator reads the XML, normalizes strings/numbers, builds an intermediate model, then templates Go.

Generated API shape:

```go
package profiles

type Telegram interface {
    EEP() eep.EEP
    MarshalERP1UserData() ([]byte, byte, error)
    Format() string
}

func ParsePacket(p erp1.Packet, prof eep.EEP) (Telegram, error)
func ParseUserData(prof eep.EEP, userData []byte, status byte) (Telegram, error)
```

For each profile, generate concrete structs with exported fields named from shortcut/data, e.g.:

```go
type D50001 struct {
    Contact D50001Contact
    LearnButton D50001LearnButton
}

type A50201 struct {
    TemperatureC float64
    TemperatureRaw uint8
    LearnButton A50201LearnButton
}
```

Formatting should be boring: `String()`/`Format()` returns stable key/value text (`A5-02-01 Temperature=...°C LRN=...`). Avoid JSON/custom pretty format unless asked.

## Minimal first generated profile set

Start with three profiles because they cover the required mechanics without dragging in VLD complexity:

1. `D5-00-01` Single Input Contact: 1-byte payload, simple enum bits, learn bit.
2. `F6-01-01` Push Button: RPS payload + reserved fields, simple enum.
3. `A5-02-01` Temperature Sensor: 4-byte payload, reserved masks, reversed range scaling, LRN bit.

Optional fourth after those pass: `F6-02-01`, because it proves statusfield conditions/case selection (T21/NU). Do not start with D2/VLD; variable-length/message profiles will widen scope.

## Validation tests to add when implementing

- XML parser/generator tests:
  - Decode `eep268.xml` as UTF-16LE and load all four RORGs.
  - Assert counts: RORGs=4, profile funcs=45, profile types=270, profile cases=355.
  - Assert parsed facts for `D5-00-01`, `F6-01-01`, `A5-02-01` exactly match bit offsets/ranges above.
- Runtime bit helper tests:
  - Extract one-bit MSB: payload `{0x80}` at bitoffs 7 size 1 => 1.
  - Extract A5 temp byte: payload `{0x00,0x00,0x80,0x10}` at bitoffs 16 size 8 => 0x80; bitoffs 28 size 1 => 1.
  - Insert/extract round trip for fields crossing a byte boundary.
- Generated profile tests:
  - `D5-00-01`: `{0x80}` parses contact closed and LRN pressed/not pressed per bit 4; marshal round-trips.
  - `F6-01-01`: `{0x08}` parses PB pressed; `{0x00}` released.
  - `A5-02-01`: raw 255 -> -40°C, raw 0 -> 0°C, raw 128 -> about -19.92°C; marshal round-trips raw values.
  - `ParsePacket` rejects RORG/EEP mismatch and short `UserData`.
- Existing package validation: run `go test ./...` after implementation.

## Implementation risks / decisions for next agent

- XML is UTF-16LE; plain `os.ReadFile` into `encoding/xml.Unmarshal` will not be enough unless decoded first.
- The XML contains lots of narrative/doc nodes; generator should parse only structural profile nodes and ignore rich text inside descriptions except flattened comments.
- `scale` appears both as child min/max and occasionally scalar text. First pass only needs datafield range+scale min/max for A5-02-01.
- `range min > max` is intentional; do not normalize away direction.
- RPS profiles use status byte predicates. A first profile set can skip multi-case RPS except F6-01-01; add F6-02-01 when condition logic is implemented.
- VLD profiles likely need variable length and message/function fields; defer until bit helpers and simple ERP1 profile dispatch are stable.
- Existing `pkg/eep` should remain the canonical triplet package. Put generated payload parsers below it (`pkg/eep/profiles`) rather than bloating `pkg/eep/eep.go`.

## Meta-prompt handoff for implementation agent

Goal: implement a minimal EEP profile generator/runtime for `eep268.xml` that emits idiomatic Go parsers/formatters for `D5-00-01`, `F6-01-01`, and `A5-02-01`, with tests proving XML parsing, bit extraction, scaling, and ERP1 integration.

Context/evidence: use facts and line anchors above. Reuse `pkg/eep.EEP` and `pkg/enums.Rorg*`; consume `pkg/erp1.Packet.UserData` and `Status`. Do not use the internal reflection serializer for bitfields.

Success criteria: generator decodes UTF-16LE XML; generated code compiles; selected profiles parse, format, and marshal ERP1 user data; tests cover XML facts and generated profile behavior; `go test ./...` passes.

Hard constraints: no new dependencies unless stdlib cannot decode the XML; no duplicate RORG/triplet types; no D2/VLD in first cut; preserve exact raw bit round trips.

Suggested approach: write a tiny XML model for just `profile/rorg/func/type/case/(condition|datafield|statusfield)`, normalize names/numbers, generate three profile files plus shared bit/scale helpers, then add tests.

Validation: targeted `go test ./pkg/eep/... ./pkg/erp1/...` during work, final `go test ./...`.

Stop/escalation: ask before expanding beyond the three-profile first set or changing public `pkg/eep.EEP`; stop once the three generated profiles and tests pass.

Resolved assumptions: bit offsets are LSB-indexed within payload bytes; A5 payload byte order is ERP1 user data order DB3, DB2, DB1, DB0; reversed raw ranges are valid and used for scaling.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Only the requested generator-focused plan was written to plans/eep268-generator.md; no project/source files were modified."
    }
  ],
  "changedFiles": [
    "plans/eep268-generator.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "python3 XML inspection scripts over eep268.xml",
      "result": "passed",
      "summary": "Parsed UTF-16LE XML, counted schema nodes/RORGs, and sampled profile datafields/enums/scales."
    },
    {
      "command": "iconv -f UTF-16LE -t UTF-8 eep268.xml | nl -ba | sed/grep ...",
      "result": "passed",
      "summary": "Collected line-numbered XML evidence for profile and field examples."
    },
    {
      "command": "nl -ba pkg/eep/eep.go pkg/erp1/erp1.go pkg/enums/rorg.go",
      "result": "passed",
      "summary": "Collected existing Go line anchors for EEP triplets, ERP1 payload shape, and RORG constants."
    },
    {
      "command": "git status --short",
      "result": "passed",
      "summary": "Confirmed working tree had pre-existing untracked paths; no source edits were made by this task."
    },
    {
      "command": "git diff --cached --stat && git status --short plans/eep268-generator.md",
      "result": "passed",
      "summary": "No staged diff; plan file is untracked/created as requested."
    }
  ],
  "validationOutput": [
    "Plan file written to /Users/edlundin/work/e/enocean-esp3/plans/eep268-generator.md"
  ],
  "residualRisks": [
    "No code/tests were run because the task requested planning only and source files were not modified.",
    "Bit ordering is inferred from EEP field examples and should be locked down with implementation tests."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added generator-focused handoff plan only.",
  "reviewFindings": [
    "no blockers"
  ],
  "manualNotes": "Repository already showed untracked docs/, eep268.xml, main.go, and plans/ before writing this plan."
}
```
