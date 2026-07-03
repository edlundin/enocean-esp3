# Smart Acknowledge + EEP 2.6.8 implementation plan

Scope audited: `/tmp/enocean-spec-text/SmartAcknowledge_Specification_v1.7.txt`, `eep268.xml`, and current Go packages. No project/source files were modified; this is an implementation plan only.

## Requirement summary

Implement this repo as a reusable Go library that can:

1. Parse/format Smart Acknowledge ERP1 telegram payloads from SmartAcknowledge v1.7.
2. Generate Go code from `eep268.xml` so library callers can parse and format normal EnOcean profile telegram payloads for the described EEP profiles.
3. Keep transport concerns separate: ESP3 frame handling, ERP1 wrapping, Smart Ack telegrams, and generated EEP profile codecs should be separate layers.

## Current support

- ESP3 frame serialize/parse exists in `pkg/esp3/esp3.go:59-143` with CRC8 header/data validation and packet type parsing.
- ERP1 wrapping exists in `pkg/erp1/erp1.go:11-94`: extracts `Rorg`, `UserData`, sender/destination IDs, status, subtelegram count, RSSI, security level; formats back to an ESP3 `RADIO_ERP1` telegram.
- RADIO_SUB_TEL-like parsing exists in `pkg/subtel/subtel.go`; note its packet-type check currently expects `PacketTypeRADIO_ERP1`, not `PacketTypeRADIO_SUB_TEL`.
- Basic EEP triplet type exists only: `pkg/eep/eep.go:21-99` has `EEP{Rorg, Func, Type}`, `FromString`, `FromTriplet`, and `String`; no XML metadata, profile validation, bit decoding, scaling, enums, status/condition support, or formatting.
- RORG enum already includes Smart Ack RORGs: `pkg/enums/rorg.go` has `SM_LRN_REQ=0xC6`, `SM_LRN_ANS=0xC7`, `SM_REC=0xA7`, `SIGNAL=0xD0`.
- ESP3 packet type enum already includes `SMART_ACK_COMMAND=0x06`: `pkg/enums/packettype.go:12-13`.
- Smart Ack command enum exists only as enum values: `pkg/enums/smartackcommand.go:5-78` defines `WR_LEARN_MODE` through `WR_WR_POSTMASTER`; no `pkg/smartackcommand` command structs or response parsers.
- Smart Ack events are partially parsed in `pkg/event/event.go:26-47` and `pkg/event/event.go:83-122` (`SA_RECLAIM_NOT_SUCCESSFUL`, `SA_CONFIRM_LEARN`, `SA_LEARN_ACK`). Current event parser uses `binary.LittleEndian` at `pkg/event/event.go:185-188`, while tests write big-endian fields; validation currently fails.
- Serial entrypoint is app-shaped, not library-shaped: `pkg/enocean.go:24-48` opens a serial port and starts a goroutine, but `parser` only prints telegrams at `pkg/enocean.go:206-210` and has a TODO for channels at `pkg/enocean.go:46`.
- `main.go` hard-codes a serial path and imports `pkg`, so keep codegen/library work out of `main.go` unless later converting it to an example.

## Spec requirements and gaps

### Smart Acknowledge telegram payloads

From SmartAcknowledge v1.7 section 3:

- Telegram overview maps logical Smart Ack messages to telegram/RORGs: Learn Request `sm_lrn_req`, Learn Reply/Acknowledge `sm_lrn_ans`, Reclaim `sm_rec`, Signal `sig`, normal Data/Data Reply/Data Acknowledge as common EnOcean data telegrams (`SmartAcknowledge...txt:830-854`).
- Message indexes are required: Learn Reclaim `0b0`, Data Reclaim `0b1`, Learn Reply `0x01`, Learn Acknowledge `0x02`, Signal Mailbox Empty `0x01`, Mailbox Missing `0x02`, Reset `0x03` (`...txt:856-869`).
- Request codes are 5-bit values: default sensor `0b11111`, candidate/no-candidate and mailbox-space combinations `0b00000`-`0b00011` (`...txt:872-887`). Current repo has no enum/type for these.
- Ack codes include ranges, not only named constants: first LearnIn `0x00`, repeated LearnIn `0x01-0x0F`, failed LearnIn `0x10-0x1F`, complete LearnOut `0x20`, partial LearnOut `0x21-0x2F` (`...txt:888-894`). Current `LearnAckConfirmCode` only models selected values and `0xff`, so Smart Ack payload code should preserve unknown/application-specific range values.
- Learn Request: ERP1 RORG `0xC6`; data length 10 bytes; fields are Request Code 5 bits + Manufacturer ID 11 bits, EEP 3 bytes, RSSI 1 byte, Repeater ID 4 bytes; subtelegram count 3; sent by sensor, not repeated by sensor though Smart Ack devices alter/retransmit (`...txt:903-939`). Missing.
- Learn Reply: ERP1 RORG `0xC7`; message index `0x01`; 7-byte payload after index: response time 2 bytes, ack code 1 byte, sensor ID 4 bytes; subtelegram count 3; repeated; sent by controller to post master (`...txt:942-973`). Missing.
- Learn Acknowledge: ERP1 RORG `0xC7`; message index `0x02`; 4-byte payload after index: response time 2 bytes, ack code 1 byte, mailbox index 1 byte; subtelegram count 1; not repeated; sender ID is always controller ID even if repeater is promoted post master (`...txt:976-1005`). Missing.
- Learn Reclaim: ERP1 RORG `0xA7`; message index bit `0`; 1-bit data, rest unused; subtelegram count 1; not repeated (`...txt:1008-1029`). Missing.
- Data Reclaim: ERP1 RORG `0xA7`; high bit/index bit `1` plus 7-bit mailbox index (`...txt:1032-1052`). Missing.
- Signal telegrams: ERP1 RORG `0xD0`; 1-byte message indexes Mailbox Empty `0x01`, Mailbox Does Not Exist `0x02`, Reset `0x03`; subtelegram count 1, addressed to sensor (`...txt:1056-1129`). Missing.
- Normal Data/Data Reply/Data Acknowledge use common EnOcean telegrams with RORG defined by EEP profile; Data/Data Reply use subtelegram count 3/repeated, Data Acknowledge subtelegram count 1/not repeated (`...txt:1130-1143`). Current ERP1 wrapper can carry these but there is no generated EEP parser/formatter.

### EEP XML requirements

`eep268.xml` facts from local parsing:

- File is UTF-16LE XML (`file eep268.xml` reports UTF-16 little-endian).
- Parsed profile tree contains 4 RORG groups relevant for regular EEP profiles: `0xF6/RPS`, `0xD5/1BS`, `0xA5/4BS`, `0xD2/VLD`.
- Counts: 45 `func` nodes, 270 `type` nodes, 355 `case` nodes, 2576 `datafield` nodes.
- XML shape is nested `profile/rorg/func/type/case/datafield` with fields such as `number`, `title`, `status`, `condition`, `statusfield`, `datafield`, `bitoffs`, `bitsize`, `range`, `scale`, `unit`, `enum`, and `reserved`.
- Example datafield: `Push button / PB`, `bitoffs=3`, `bitsize=1`, enum `0 Released`, `1 Pressed & Hold`.
- Example conditional RPS case uses `condition` and `statusfield` (`T21`, `NU`) plus datafields.

Missing in repo:

- No XML reader or generator.
- No generated profile registry keyed by `(RORG,FUNC,TYPE)`.
- No bit extraction/insertion helper for non-byte-aligned fields.
- No support for EEP case selection by status fields/conditions.
- No support for enum labels, reserved fields, value ranges, scale/unit conversions, or raw fallback values.
- No VLD variable-length payload support beyond generic `UserData []byte` storage.

## High-value implementation constraints

- Use stdlib `encoding/xml` plus UTF-16 handling. Do not add a dependency just to parse XML; if needed, decode UTF-16 with `unicode/utf16` and strip BOM before `xml.Decoder`.
- Generate compact metadata tables plus one generic codec rather than thousands of handwritten/profile-specific parsing functions. That still satisfies “generate Go code” and keeps maintenance sane.
- Keep Smart Ack protocol telegrams separate from EEP profile codecs. Smart Ack controls learning/reclaim/signals; normal data payload interpretation belongs to generated EEP codec.
- Preserve raw values when spec gives application-specific ranges (ack code ranges, unknown enum values, reserved bits) instead of rejecting too early.
- ERP1 byte order for multi-byte protocol fields should be spec/test-confirmed and implemented manually with `encoding/binary.BigEndian` unless ESP3 spec evidence says otherwise. Avoid reflection/binary over structs for radio telegrams; bit-level layouts need explicit code.
- Do not solve full post-master state machines/mailbox management in the first pass unless explicitly requested. The smallest library requirement is parse/format telegrams, not operate a controller.

## Files to change/add in implementation

### Likely new files

- `pkg/smartack/smartack.go`: public Smart Ack payload types and parse/format entrypoints.
- `pkg/smartack/telegram.go`: explicit payload encoders/decoders for Learn Request, Learn Reply, Learn Ack, Reclaim, Signal.
- `pkg/smartack/smartack_test.go`: spec vector tests for every Smart Ack payload shape.
- `internal/eepxml/parse.go`: parse UTF-16 XML into a small IR (`RORG`, `Func`, `Type`, `Case`, `Field`, `Enum`, `Range`, `Scale`, `Condition`).
- `cmd/eepgen/main.go`: code generator invoked by `go generate`; input `eep268.xml`, output generated Go.
- `pkg/eep/generated_profiles.go`: generated metadata table. Commit generated file; do not require users to have XML at runtime.
- `pkg/eep/codec.go`: generic EEP parse/format over generated metadata.
- `pkg/eep/bit.go`: bit extraction/insertion helpers; test heavily.
- `pkg/eep/codec_test.go` and `internal/eepxml/parse_test.go`: minimal generator and round-trip tests.

### Likely existing files to touch

- `pkg/eep/eep.go`: add registry lookup/types; keep existing `EEP` API backward compatible.
- `pkg/erp1/erp1.go`: possibly add convenience `ParseProfile(eep.EEP)` / `WithProfileData` helpers or leave as wrapper and call generated codec externally. Do not bake profile logic into ERP1 wrapper.
- `pkg/enums/event.go`: consider widening `LearnAckConfirmCode` handling or add Smart Ack-specific ack code type in `pkg/smartack` so application-specific ranges are preserved.
- `pkg/event/event.go`: fix endianness only if implementation touches event validation; currently tests show endian failures.
- `pkg/subtel/subtel.go`: verify/fix packet type if RADIO_SUB_TEL support is needed for validation; not required for first Smart Ack/EEP payload phase.
- `pkg/enocean.go`: later expose parsed telegram channel/callback; not needed for codegen/profile parsing plan except to make repo library-shaped.

## Smallest implementation phases

### Phase 1 — Smart Ack payload library, no XML generator

Goal: add `pkg/smartack` that parses/formats Smart Ack ERP1 payloads using existing `erp1.Packet`.

Minimum API:

- `Parse(p erp1.Packet) (Message, error)` dispatches by `p.Rorg` and message index bits/bytes.
- `func (m LearnRequest) ERP1(senderID deviceid.DeviceID) erp1.Packet` or `AppendToERP1` style helpers for formatting.
- Types: `LearnRequest`, `LearnReply`, `LearnAcknowledge`, `LearnReclaim`, `DataReclaim`, `MailboxEmpty`, `MailboxDoesNotExist`, `Reset`.
- Enums/constants: request codes, ack code classifier, message indexes, expected subtelegram counts.

Validation:

- Unit tests for exact payload bytes from spec sections 3.1.2-3.1.9.
- Round-trip parse/format for every Smart Ack message.
- Negative tests for wrong RORG, wrong lengths, invalid message index, mailbox index >127 for Data Reclaim.

Why first: It is independent of EEP XML and exercises existing ERP1/ESP3 layers.

### Phase 2 — EEP XML reader + generated metadata table

Goal: parse `eep268.xml` and generate a committed metadata file; no full value formatting yet.

Minimum API:

- `internal/eepxml.Parse(io.Reader) (Model, error)`.
- `cmd/eepgen -in eep268.xml -out pkg/eep/generated_profiles.go`.
- `pkg/eep.Lookup(eep.EEP) (Profile, bool)`.
- Generated profile metadata includes RORG/FUNC/TYPE title/status, cases, fields, bit offsets/sizes, enum labels, ranges/scales/units where present.

Validation:

- Golden/count test: generator sees 4 RORG groups, 45 funcs, 270 types, 355 cases from current `eep268.xml`.
- Generated file compiles under `go test ./pkg/eep`.
- Lookup tests for known profiles like `F6-01-01`, `D5-00-01`, one `A5-*`, one `D2-*`.

### Phase 3 — Generic EEP parse/format codec

Goal: parse/format payload bits using generated metadata.

Minimum API:

- `type Values map[string]Value` or similar. Prefer stable field shortcut keys (`PB`, `CO`, etc.) with raw key fallback when duplicate shortcuts exist.
- `func Decode(profile EEP, rorg enums.Rorg, userData []byte, status byte) (DecodedTelegram, error)`.
- `func Encode(profile EEP, values Values) (userData []byte, status byte, error)`.
- Bit helper tested independently: extract/insert across byte boundaries, preserve reserved/unset bits where possible.

Validation:

- Round-trip tests on simple profiles: RPS `F6-01-01`, 1BS `D5-00-01`, one simple 4BS scalar profile, one VLD profile with variable length if metadata is clear.
- Enum label tests and raw numeric fallback tests.
- Condition/status-field test for an RPS case.

### Phase 4 — Integrate normal Data/Data Reply/Data Acknowledge with Smart Ack

Goal: let callers parse Smart Ack control telegrams and common data telegrams in one library flow.

Minimum API:

- `smartack.ParseERP1(p erp1.Packet, profile *eep.EEP)` returns either Smart Ack message or decoded EEP data for normal data RORGs.
- Formatting helpers for Data/Data Reply/Data Acknowledge that set correct subtelegram count defaults from spec (3 for repeated Data/Data Reply, 1 for Data Acknowledge) while allowing caller override.

Validation:

- Tests combining `erp1.Packet -> smartack/control` and `erp1.Packet -> eep/data` dispatch.
- ESP3 serialization round-trip with existing `esp3.Serialize`/`NewEsp3TelegramFromHexString`.

### Phase 5 — Library transport cleanup (only if desired)

Goal: make serial parser useful as a library.

Minimum:

- Add a channel/callback API for parsed `esp3.Telegram` instead of `fmt.Println` in `pkg/enocean.go:206-210`.
- Leave `main.go` as example or move to `cmd/...` later.

Validation:

- Parser unit test using a fake `io.Reader`/serial-like interface, or leave untouched if out of scope.

## Existing validation status

Command run: `go test ./...`

Result: failed before any modifications.

Failures observed:

- `pkg/enums`: `TestParseCommonCommandFromByte` expected `0x34` invalid but parser returned nil, then test panicked on nil error.
- `pkg/event`: endian mismatches for `SA_CONFIRM_LEARN`, `SA_LEARN_ACK`, and secure device ID; parser uses little-endian while tests expect big-endian.
- `pkg/subtel`: sender/user-data/subtel count failures; likely test/implementation disagreement on packet type/layout and trailing optional subtelegram parsing.

This matters because future Smart Ack work should add targeted package tests and not rely on full `go test ./...` being green until baseline failures are fixed or quarantined.

## Residual risks / open questions

- EEP bit numbering must be confirmed against the EEP PDF/examples before final codec implementation. XML `bitoffs` appears to use bit positions, but byte/bit ordering for multi-byte fields needs explicit tests.
- VLD profile formatting may need per-profile length/rule handling; start with metadata-driven variable length and add profile exceptions only where tests prove needed.
- XML contains conditions/status fields; generic case selection may be ambiguous for some profiles. Decode can return all matching cases or require caller-supplied case selection for encode.
- Smart Ack spec describes behavior/state machines (post-master election, mailboxes, timing), but the requested library goal is parse/format. Do not implement controller/post-master runtime state until a separate requirement asks for it.
- Current repo has untracked files (`docs/`, `eep268.xml`, `main.go`, `plans/`) already present; do not assume a clean worktree.

## Handoff contract for next implementation agent

Goal: implement Phase 1 first, then Phase 2 only if Phase 1 is accepted. Produce minimal, library-shaped code with tests.

Context/evidence:

- Smart Ack payload specs are in `/tmp/enocean-spec-text/SmartAcknowledge_Specification_v1.7.txt` lines around 830-1143.
- Existing RORG and packet type constants already cover Smart Ack: `pkg/enums/rorg.go`, `pkg/enums/packettype.go`.
- Existing ERP1 wrapper is the right transport-level reuse point: `pkg/erp1/erp1.go:11-94`.
- Existing EEP package is only triplet parsing: `pkg/eep/eep.go:21-99`.
- `eep268.xml` is UTF-16LE and parsed locally as 4 RORG groups / 45 funcs / 270 types / 355 cases / 2576 datafields.

Success criteria:

- Phase 1: every Smart Ack telegram in spec sections 3.1.2-3.1.9 parses/formats with exact byte tests.
- Phase 2: generator parses XML, generated metadata compiles, known profile lookups work.
- Do not widen into runtime controller/mailbox state machine unless explicitly approved.

Validation:

- Run targeted tests for packages touched, e.g. `go test ./pkg/smartack ./pkg/erp1 ./pkg/eep`.
- Run `go test ./...` and report existing unrelated failures separately if still present.

Stop/escalation rules:

- Stop after Phase 1 if scope/time is tight; it is a coherent deliverable.
- Ask before adding dependencies, changing public module structure, or implementing post-master/mailbox state.
- If XML semantics are ambiguous, add raw metadata generation first and defer lossy scaling/case-selection behavior.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only for the requested Smart Acknowledge + eep268.xml audit; no project/source files were modified."
    }
  ],
  "changedFiles": [
    "plans/smart-ack.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "file eep268.xml && python3 ...",
      "result": "passed",
      "summary": "Confirmed eep268.xml is UTF-16LE XML."
    },
    {
      "command": "python3 XML parse/count scripts",
      "result": "passed",
      "summary": "Parsed eep268.xml and counted 4 RORG groups, 45 funcs, 270 types, 355 cases, 2576 datafields."
    },
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Baseline failures in pkg/enums, pkg/event, and pkg/subtel; no source changes made."
    },
    {
      "command": "git diff --cached --name-only",
      "result": "passed",
      "summary": "No staged files."
    }
  ],
  "validationOutput": [
    "go test ./... failed before modifications: commoncommand enum invalid-value test panic, event endian mismatches, subtel packet parsing mismatches."
  ],
  "residualRisks": [
    "EEP bit numbering and VLD edge cases need confirmation during implementation.",
    "Existing baseline tests are not green, so full-suite validation will need separate cleanup or documented failures."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added implementation plan/audit document only.",
  "reviewFindings": [
    "no blockers for plan-only deliverable"
  ],
  "manualNotes": "Worktree already contains untracked docs/, eep268.xml, main.go, and plans/ entries; this run intentionally changed only the requested plan file."
}
```
