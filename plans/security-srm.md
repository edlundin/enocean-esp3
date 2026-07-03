# Security/SRM + EEP codegen implementation plan

## Scope read

- Local specs used:
  - `docs/specifications/SecureRemoteManagement-3.1.pdf`, extracted to `/tmp/SecureRemoteManagement-3.1.pdf.txt` for audit.
  - `docs/specifications/Security_of_EnOcean_Radio_Networks_v3.02.pdf`, extracted to `/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt` for audit.
  - No checked-in extracted `.txt` versions were present; extraction was temporary only.
- XML used: `eep268.xml` (UTF-16LE, 61,793 decoded lines).
- Go support reviewed: ESP3, ERP1, EEP triplet, RORG enums, common-command security/REMAN wrappers, serializer/deserializer, serial parser.

## Current support

- ESP3 frame serialization/parsing exists:
  - `pkg/esp3/esp3.go:51-79` builds ESP3 frames with sync/header CRC/data CRC.
  - `pkg/esp3/esp3.go:85-134` parses a hex ESP3 telegram and validates CRCs.
- ERP1 envelope exists but only as raw RORG/user-data mapping:
  - `pkg/erp1/erp1.go:11-20` models destination, RORG, RSSI, security level, status, sender, user data.
  - `pkg/erp1/erp1.go:22-67` parses ERP1 ESP3 data into raw `UserData`; no EEP profile parse.
  - `pkg/erp1/erp1.go:69-94` formats ERP1, but hard-codes RSSI `0xff` and security level `0x03` in opt-data at lines 79-83.
- EEP support is only a triplet parser/stringer:
  - `pkg/eep/eep.go:21-25` stores `Rorg`, `Func`, `Type`.
  - `pkg/eep/eep.go:27-40`, `50-100` validate and parse `RR-FF-TT` strings.
  - Bug/requirement gap: `maxFunc = 0x60` at `pkg/eep/eep.go:15-18`, but `eep268.xml:59658` has `FUNC 0xB0` (`Liquid Leakage Sensor`), so generated EEPs from 2.68 will not all validate.
- RORG enum is partial and security names are not spec-complete:
  - `pkg/enums/rorg.go:7-22` has RPS/1BS/4BS/VLD/SYS_EX/SEC/SEC_ENCAPS/SEC_MAN, but not `SEC_D 0x32`, `SEC_CDM 0x33`, `SEC_TI 0x35` names required by Security v3.02; `SEC_ENCAPS` appears to mean `SEC_R`.
- Common command wrappers cover device-table security setup, not over-the-air SRM SYS_EX/RMCC/RPC:
  - `pkg/commoncommand/securedevice.go:14-23` WR_SECUREDEVICE_ADD with 24-bit RLC.
  - `pkg/commoncommand/securedevice.go:308-315` WR_SECUREDEVICEV2_ADD with 32-bit RLC and direction.
  - `pkg/commoncommand/securedevice.go:393-430` REMAN key common commands.
  - `pkg/commoncommand/reman.go:12-25`, `28-59` old-style REMAN code/repetition commands.
- Reflection serializer is byte-oriented only:
  - `internal/serializer/serializer.go:75-140` serializes tagged fields to data/optdata.
  - `internal/serializer/serializer.go:173-216` writes whole numeric values; no bitfield packing, scale conversion, enum validation, or XML-driven field metadata.
  - `internal/serializer/deserializer.go:61-105`, `108-230` deserializes whole fields; slices consume remaining bytes (`deserializeSlice` below line 260), which is not enough for EEP bit ranges.
- Library API is not yet usable as a receive stream:
  - `pkg/enocean.go:24-49` opens a serial port and starts a goroutine.
  - `pkg/enocean.go:51-220` parser prints valid telegrams at lines 205-210 instead of returning them over a channel/callback.

## EEP XML facts relevant to code generation

- `eep268.xml` structure is regular enough for generator input:
  - `rorg -> func -> type -> case -> datafield` with `number`, `title`, `status`, `bitoffs`, `bitsize`, `range`, `scale`, `unit`, `enum`.
  - Counts from decoded XML: 65 `rorg` elements, 107 `func`, 342 `type`, 370 `case`, 2,578 `datafield`, 1,911 `enum`, 406 `range`, 664 `scale`, 606 `unit`.
- Example field layout:
  - `eep268.xml:2442-2443`: `0xA5` / `4BS Telegram`.
  - `eep268.xml:2602-2605`: FUNC `0x02` / `Temperature Sensors`.
  - `eep268.xml:2657-2667`: datafield `TMP`, bit offset 16, size 8, raw range `255..0`, scale `-40..0`, unit `°C`.
- Bit offsets are spec bit offsets, not Go struct byte offsets. Generator must emit bit extraction/packing helpers; current reflection serializer cannot represent these profiles correctly.
- Profiles include reserved fields and conditions/status fields; generator should preserve metadata but first phase can skip generating public fields for reserved bits.

## Spec requirements missing or partially missing

### Secure Remote Management 3.1

Evidence from extracted SRM text:

- SYS_EX payload is RORG `0xC5` and big-endian (`/tmp/SecureRemoteManagement-3.1.pdf.txt:267-270`).
- Alliance-defined RMCC/RPC header omits manufacturer ID `0x7FF`; format is 1 bit manufacturer-present `0`, 3 unused bits `0`, 12-bit function number, then payload (`/tmp/SecureRemoteManagement-3.1.pdf.txt:276-297`).
- Manufacturer-specific RPC header includes manufacturer-present `1`, 11-bit manufacturer ID, 12-bit function number, then payload (`/tmp/SecureRemoteManagement-3.1.pdf.txt:301-320`).
- Function ranges: reserved `0x000`; RMCC request `0x001..0x1FF`; RPC request `0x200..0x5FF`; responses `0x600..0xFFF` (`/tmp/SecureRemoteManagement-3.1.pdf.txt:327-345`).
- Required RMCCs: Action `0x005`, Ping `0x006`, Query Status `0x008` (`/tmp/SecureRemoteManagement-3.1.pdf.txt:350-366`).
- Required responses/payloads:
  - Ping response `0x606`, 1-byte RSSI (`/tmp/SecureRemoteManagement-3.1.pdf.txt:407-416`).
  - Query Status response `0x608`, 3-byte payload with 12-bit last function and 8-bit return code (`/tmp/SecureRemoteManagement-3.1.pdf.txt:419-450`).
  - Return codes include `0x00`, `0x01`, `0x04`, `0x05`, `0x07`, `0x08`, `0x0D`, `0x0E`, `0x0F` (`/tmp/SecureRemoteManagement-3.1.pdf.txt:450-463`).
- Required RPCs from SRM 3.1:
  - Remote Learn `0x201`, payload 1-byte flag (`/tmp/SecureRemoteManagement-3.1.pdf.txt:483-511`).
  - Remote Memory Write `0x203`, payload 32-bit address + 8-bit N + N bytes (`/tmp/SecureRemoteManagement-3.1.pdf.txt:516-537`).
  - Remote Memory Read `0x204`, request 5 bytes, response `0x804` N bytes (`/tmp/SecureRemoteManagement-3.1.pdf.txt:538-572`).
  - SMART ACK read `0x205`, responses `0x805` mailbox settings and `0x806` learned sensors (`/tmp/SecureRemoteManagement-3.1.pdf.txt:574-667`).
  - SMART ACK write `0x206`, operation-specific payloads (`/tmp/SecureRemoteManagement-3.1.pdf.txt:667-712`).
  - Remove Device `0x207`, 0-byte payload (`/tmp/SecureRemoteManagement-3.1.pdf.txt:713-733`).
- SRM must be carried as secure messages. For ERP1, use `SEC` if entire SYS_EX fits one telegram; otherwise split into `SEC_CDM`; ERP2 can carry whole SEC telegram (`/tmp/SecureRemoteManagement-3.1.pdf.txt:225-233`).
- Session commands from older REMAN are no longer needed; authentication is possession of security key (`/tmp/SecureRemoteManagement-3.1.pdf.txt:234-246`). Current `WrRemanCode` is therefore legacy/common-command support, not SRM 3.1 OTA compliance.

Missing in repo:

- No `pkg/srm` or `pkg/sysex` package for RMCC/RPC header packing/parsing.
- No 12-bit function number packing, manufacturer-present header handling, return-code enum, or SRM request/response parser.
- No SRM-over-secure-ERP integration.
- No broadcast random response delay logic is needed for a library formatter/parser unless simulating a remote device; document it, do not implement in initial library unless device-side support is requested.

### Security of EnOcean Radio Networks v3.02

Evidence from extracted security text:

- PK is 128-bit; unique per device; RLC resets to 0 when PK changes (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:418-472`).
- RLC is replay protection; receiver accepts only higher RLC; explicit RLC recommended; implicit RLC discouraged; 32-bit recommended, 16-bit deprecated (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:473-595`).
- CMAC authenticates telegram data, RLC, PK; 3-byte and 4-byte CMAC supported, 2-byte deprecated (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:596-613`).
- SLF layout is `RLC_TYPE` bits 7..5, `CMAC_TYPE` bits 4..3, `ENCRYPTION_TYPE` bits 2..0 (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:624-646`).
- Supported SLFs include standard `0xF3`, energy-reduced `0xAB`/`0xCB`, ultra-low-power `0x8B` (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:647-704`).
- Security RORG-S values: `SEC 0x30`, `SEC_R 0x31`, `SEC_D 0x32`, `SEC_CDM 0x33`, `SEC_TI 0x35` (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:782-813`).
- Secure telegrams with original RORG encrypt/authenticate RORG+DATA and transmit `SEC_R 0x31` + encrypted RORG/DATA + plain RLC + CMAC (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:866-886`).
- Chaining required in ERP1 when payload+RLC+CMAC exceeds 14 bytes; with 4-byte RLC and 4-byte CMAC, payload >6 bytes requires SEC_CDM (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:889-907`).
- Secure chain message uses 2-bit nonzero `SEQ` and 6-bit `IDX`; first message includes 16-bit length and 10 bytes content, subsequent messages carry up to 13 bytes (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:908-970`).
- SEC-TI includes teach-in info, SLF, current RLC, PK; ERP1 SEC-TI must be split into two messages (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:997-1150`, `1428-1470`).
- PSK teach-in encrypts `RLC||PK` using VAES with RLC `0x0000` and PSK (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:1168-1183`).
- VAES IV constant is `3410DE8F1ABA3EFF9F5A117172EACABD` (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:1204-1279`).
- CMAC input is `RORG-S || encrypted data || RLC`; RFC4493-style AES-CMAC with K1/K2, use 3 or 4 MSB (`/tmp/Security_of_EnOcean_Radio_Networks_v3.02.pdf.txt:1280-1439`).

Missing in repo:

- No AES/VAES/CMAC implementation. Use Go stdlib `crypto/aes`, `crypto/cipher` primitives; do not add a dependency.
- No `SecurityLevelFormat` type to parse RLC/CMAC/encryption lengths.
- No RLC state/window manager; for a stateless parser library, expose explicit-RLC decrypt/verify first and leave implicit-RLC windowing behind an optional stateful helper.
- No SEC/SEC_R/SEC_CDM/SEC_TI parser/formatter.
- No chain assembler/splitter.
- RORG enum lacks `SEC_D`, `SEC_CDM`, `SEC_TI` and should rename/alias `SEC_ENCAPS` to spec `SEC_R`.

## Smallest implementation phases

### Phase 0 — tighten library boundary, no feature breadth

Files likely touched:

- `main.go`: move demo CLI behind `cmd/...` or leave as example; module currently has a root `main`, so library consumers import `github.com/edlundin/enocean-esp3/pkg` but root package is not library-shaped.
- `pkg/enocean.go`: stop printing from parser; expose a channel/callback of `esp3.Telegram` or just leave serial I/O out of the EEP/security work.

Validation:

- Existing `go test ./...` currently fails for unrelated enum/event/subtel tests; record baseline before changes and run targeted new tests.

### Phase 1 — generated EEP parse/format without security

Smallest useful deliverable: generated code from `eep268.xml` for profile metadata and bit-level raw/scaled field access.

Files to add/change:

- Add `internal/eepxml` or `internal/eepgen` parser used only by generator.
- Add `cmd/eepgen/main.go` or `internal/cmd/eepgen` generator.
- Add `pkg/eep/generated/...` for generated profile registry and field codecs.
- Change `pkg/eep/eep.go` validation: `maxFunc` must allow at least `0xB0` from XML; safest is `0xff` unless spec text elsewhere still constrains it.
- Add `pkg/eep/bitfield.go` with tiny stdlib helpers:
  - `GetBitsBE(data []byte, bitOffset, bitSize int) (uint64, error)`
  - `SetBitsBE(data []byte, bitOffset, bitSize int, value uint64) error`
  - scaling helper for raw range reversed or normal: `(raw-rawMin)*(scaleMax-scaleMin)/(rawMax-rawMin)+scaleMin`.
- Add `pkg/eep/registry.go` with `Profile` metadata and generated lookup by `EEP`.

Design constraints:

- Do not generate reflection structs first. The XML is bitfield-oriented; a metadata-driven codec is smaller and covers all profiles.
- Generate reserved fields in metadata but not public values.
- First pass can support RPS/1BS/4BS/VLD data payloads and defer teach-in condition semantics; parser should still expose raw fields if a condition cannot be evaluated.

Validation:

- Generator unit test with a tiny inline XML fixture.
- Bitfield tests with A5-02-05 `TMP` field: raw 255 => -40°C, raw 0 => 0°C using `eep268.xml:2657-2667`.
- Generated registry test: `A5-02-05`, `D2-xx`, and `FUNC 0xB0` profiles exist.

### Phase 2 — ERP1 profile integration

Files to change:

- `pkg/erp1/erp1.go`: add methods, not a rewrite:
  - `func (p Packet) Payload() []byte` returning user data.
  - `func (p Packet) ParseEEP(profile eep.EEP) (eep.Values, error)` using generated registry.
  - `func NewPacketFromEEP(profile eep.EEP, values eep.Values, sender, dest deviceid.DeviceID, status byte) (Packet, error)`.
- Add tests under `pkg/erp1` and/or `pkg/eep`.

Validation:

- Round-trip A5-02-05 values -> ERP1 user data -> values.
- Existing ERP1 tests still pass.

### Phase 3 — SRM SYS_EX packing/parsing, unsecured payload only

Files to add:

- `pkg/srm/sysex.go`: header pack/unpack for Alliance and manufacturer-specific messages.
- `pkg/srm/rmcc.go`: Action, Ping, QueryStatus requests/responses and return codes.
- `pkg/srm/rpc.go`: RemoteLearn, MemoryRead/Write, SmartAck read/write, RemoveDevice payloads.

Key implementation:

- Use plain byte packing for 12-bit fields; no reflection serializer.
- API shape:
  - `type Message struct { ManufacturerID *uint16; Function uint16; Payload []byte }`
  - `func ParseSYSEx([]byte) (Message, error)`
  - `func (m Message) MarshalBinary() ([]byte, error)`
- Validate function ranges and payload lengths where fixed by spec.

Validation:

- Header tests:
  - Alliance function `0x006` encodes two bytes with manufacturer-present bit 0, unused bits 0, function 0x006.
  - Manufacturer-specific RPC encodes present bit 1, 11-bit manufacturer, 12-bit function.
- Payload tests for Ping response and Query Status 12-bit function + 8-bit return.

### Phase 4 — security primitives and SEC/SEC_R single telegrams

Files to add/change:

- `pkg/security/slf.go`: SLF parse/format with constants `0xF3`, `0xAB`, `0xCB`, `0x8B`.
- `pkg/security/vaes.go`: VAES using `crypto/aes`; IV constant from spec.
- `pkg/security/cmac.go`: AES-CMAC RFC4493 with truncated 3/4 MSB.
- `pkg/security/telegram.go`: secure/unsecure single ERP message.
- `pkg/enums/rorg.go`: add `RorgSEC_D=0x32`, `RorgSEC_CDM=0x33`, `RorgSEC_TI=0x35`; alias `RorgSEC_R=0x31` while keeping old name for compatibility.

Smallest API:

- `func EncryptERP1(rorg enums.Rorg, data []byte, key [16]byte, rlc uint32, slf SLF) (rorgS enums.Rorg, secureData []byte, err error)`
- `func DecryptERP1(rorgS enums.Rorg, secureData []byte, key [16]byte, slf SLF) (rorg enums.Rorg, data []byte, rlc uint32, err error)`
- First pass supports explicit RLC only (`0xF3`, `0xAB`, `0xCB`) and returns `ErrImplicitRLCNeedsState` for `0x8B`.

Validation:

- Use test vectors/examples from Security Annex A.4 if extracted text is readable enough; otherwise create round-trip tests plus independent CMAC RFC4493 known vectors.
- Verify CMAC fails on tampered byte and stale RLC.

### Phase 5 — SEC_CDM chain splitting/assembly and SRM over security

Files to add/change:

- `pkg/security/chain.go`: split/reassemble SEC_CDM chains.
- `pkg/srm/secure.go`: helpers to wrap SYS_EX SRM messages into secure ERP1 packets and unwrap them.

Validation:

- ERP1 standard SLF `0xF3`: any SRM SYS_EX payload longer than 6 bytes uses SEC_CDM, per spec.
- Chain tests: first chunk has nonzero `SEQ`, `IDX=0`, 16-bit length; middle chunks max 13 bytes; reassembly rejects mixed SEQ, missing chunks, duplicate conflict, length mismatch.
- SRM MemoryWrite with payload >6 bytes round-trips through secure chain.

### Phase 6 — SEC-TI teach-in support

Files to add/change:

- `pkg/security/teachin.go`: SEC_TI parse/format for two ERP1 telegrams.
- Optional PSK encryption for `RLC||PK` using VAES with RLC zero.

Validation:

- Parse/format non-PTM unidirectional and bidirectional teach-in info.
- PSK round-trip encrypt/decrypt of RLC+PK.
- Reject invalid IDX/CNT/reserved bits.

## Validation strategy

- Baseline command run: `go test ./...` currently fails before this plan's implementation:
  - `pkg/enums` test expects `0x34` invalid though code includes `SET_CRCMode`.
  - `pkg/event` has endian expectation failures.
  - `pkg/subtel` has sender/user-data/SubTel count failures.
- For implementation phases, require targeted passing tests first:
  - `go test ./pkg/eep ./pkg/erp1 ./pkg/srm ./pkg/security ./internal/eepxml`
  - `go test ./pkg/esp3 ./pkg/erp1` for envelope compatibility.
- Full `go test ./...` remains a final gate only after deciding whether unrelated existing failures are fixed or accepted as baseline debt.
- Add `go generate ./...` or `go generate ./pkg/eep` only after generator path is stable.

## Risks and open decisions

- Spec text extraction was temporary; if reviewer requires exact source text files, add checked-in extracted text or cite PDFs only. No local extracted `.txt` existed at start.
- EEP XML conditions/teach-in semantics are richer than simple field extraction. First generated codec should expose raw conditional cases rather than pretending full semantic support.
- Security implementation is easy to get subtly wrong. Use stdlib AES, known CMAC vectors, and spec annex vectors before using against hardware.
- RLC state persistence is application-specific. Library should not silently store counters globally; expose state hooks/types and keep stateless explicit-RLC helpers.
- Existing common-command REMAN/security functions configure a TCM/module, not the same as OTA SRM 3.1 RMCC/RPC. Do not mix APIs.
- Current repo has untracked `docs/`, `eep268.xml`, `main.go`, and `plans/`; no source files were modified for this audit besides writing this plan.

## Handoff contract for implementation agent

Goal: implement a Go library that generates EEP codecs from `eep268.xml`, parses/formats ERP1 telegrams to/from those EEP profiles, and then adds SRM 3.1 messages over Security v3.02 secure ERP telegrams in the smallest phases above.

Context/evidence: use the line/file references in this plan. Reuse existing `esp3.Telegram`, `erp1.Packet`, `eep.EEP`, `deviceid.DeviceID`, and stdlib `crypto/aes`; do not add dependencies for XML, AES, or CMAC.

Success criteria:

- Generated EEP registry/code covers profiles from `eep268.xml`, including FUNC `0xB0`.
- At least one generated profile round-trips raw and scaled values.
- SRM SYS_EX RMCC/RPC headers parse/format per 12-bit function/manufacturer rules.
- Security package supports SLF parsing, VAES, CMAC, SEC_R single telegrams, SEC_CDM chain split/reassemble, and SEC_TI after the earlier phases.
- Targeted tests pass; unrelated baseline full-suite failures are documented until fixed.

Hard constraints:

- Keep generator/runtime simple; no new dependency unless stdlib cannot do it.
- No global mutable security key/RLC registry in the library.
- Do not claim full Security/SRM compliance until spec annex vectors or hardware captures pass.

Stop/escalation rules:

- Stop after each phase with runnable tests.
- Ask for decision before changing public module layout (`main.go`/root package), fixing unrelated baseline tests, or choosing persistent RLC storage semantics.

## Acceptance report

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Produced an implementation plan only; no project/source files were modified beyond the requested plan output."
    },
    {
      "id": "criterion-2",
      "status": "satisfied",
      "evidence": "Plan includes spec-backed requirements, current support with file/line references, missing requirements, phased files to change, validation, residual risks, and command output summary."
    }
  ],
  "changedFiles": [
    "plans/security-srm.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "python3 -m venv /tmp/enocean_pdf_venv && /tmp/enocean_pdf_venv/bin/pip install --quiet pypdf && /tmp/enocean_pdf_venv/bin/python ...",
      "result": "passed",
      "summary": "Extracted SecureRemoteManagement-3.1.pdf (24 pages, 34,251 chars) and Security_of_EnOcean_Radio_Networks_v3.02.pdf (58 pages, 92,830 chars) to /tmp for reading."
    },
    {
      "command": "python3 ... parse eep268.xml",
      "result": "passed",
      "summary": "Decoded UTF-16LE XML; counted profile elements and confirmed FUNC max 0xB0."
    },
    {
      "command": "go test ./...",
      "result": "failed",
      "summary": "Existing baseline failures in pkg/enums, pkg/event, and pkg/subtel; unrelated to this plan-only audit."
    },
    {
      "command": "git diff --cached --name-only",
      "result": "passed",
      "summary": "No staged files."
    }
  ],
  "validationOutput": [
    "go test ./...: pkg/enums TestParseCommonCommandFromByte fails on 0x34 expected invalid but parser accepts SET_CRCMode; pkg/event endian expectations fail; pkg/subtel sender/user-data/SubTel count expectations fail.",
    "git diff --cached --name-only: empty output."
  ],
  "residualRisks": [
    "Local extracted spec .txt files were not present; PDFs were extracted to /tmp for audit.",
    "Full test suite is already failing before implementation work.",
    "Security algorithms require spec annex vectors or hardware captures before claiming compliance."
  ],
  "noStagedFiles": true,
  "diffSummary": "Added/updated only the requested implementation plan at plans/security-srm.md.",
  "reviewFindings": [
    "no blockers for plan-only handoff"
  ],
  "manualNotes": "Repository has pre-existing untracked docs/, eep268.xml, main.go, and plans/ entries; no staging was performed."
}
```
