# TODO / real-device validation gaps

`go test ./...` passes. These are the remaining checks most likely to surface during WireBOS integration and real EnOcean device testing.

## Runtime/session API

- [ ] `pkg/enocean.go`: finish channel handling for ESP3 telegram reads.
- [ ] `pkg/enocean.go`: add context cancellation for long-running read/write loops.
- [ ] `pkg/enocean.go`: validate packet type before accepting parsed telegrams.
- [ ] Add an integration-style test with a fake serial port/reader for read loop cancellation and packet dispatch.

## ReMan / RMCC

- [ ] Confirm `reman.QueryIDPayload()` is spec-correct as an empty payload; add a test and comment.
- [ ] Confirm `reman.PingPayload()` is spec-correct as an empty payload; add a test and comment.
- [ ] Add real-device or captured-frame tests for ReMan Query ID, Ping, and Code flows.

## GP / Generic Profiles

- [ ] `pkg/gp`: implement or explicitly document unsupported GP channel definitions.
- [ ] Add captured telegram tests for every GP channel definition WireBOS devices emit.

## EEP coverage

- [ ] When WireBOS hits `unsupported EEP`, add the profile under `pkg/eep/profiles` with captured telegram tests.
- [ ] Track WireBOS-required EEPs separately from full-spec EEPs so integration can stay focused.

## WireBOS integration checks

- [ ] In WireBOS `enoceaninterpreter`, use this module for ESP3/ERP1/ReCom/ReMan parsing/building; do not duplicate protocol code.
- [ ] Add a WireBOS adapter test: ESP3/ERP1 telegram → parsed value → async event → datagateway write.
- [ ] Add a WireBOS command test: domain command → package telegram builder → MQTT/WireGet command payload.

## Release hygiene

- [ ] Tag a version once WireBOS imports it without a local `replace`.
- [ ] Keep `go test ./...` green before updating WireBOS dependency.
