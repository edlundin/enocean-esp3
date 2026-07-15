# enocean-esp3

Go library for EnOcean ESP3/ERP1 telegrams, EEP profile payloads generated from `eep268.xml`, Smart Ack, Remote Management, Remote Commissioning, and security helpers.

## Install

```sh
go get github.com/edlundin/enocean-esp3
```

Import the packages you need:

```go
import (
    "github.com/edlundin/enocean-esp3/pkg/eep"
    "github.com/edlundin/enocean-esp3/pkg/eep/profiles"
    "github.com/edlundin/enocean-esp3/pkg/erp1"
    "github.com/edlundin/enocean-esp3/pkg/esp3"
)
```

## Parse ESP3 and ERP1

```go
telegram, err := esp3.NewEsp3TelegramFromHexString("55000707017af6...")
if err != nil {
    panic(err)
}

packet, err := erp1.NewPacketFromEsp3(telegram)
if err != nil {
    panic(err)
}

fmt.Println(packet.Rorg, packet.UserData, packet.SenderID)
```

## Decode an EEP profile payload

EEP profile metadata is generated from the checked-in `eep268.xml` into `pkg/eep/profiles`.

```go
profile, err := eep.FromString("A5-02-01")
if err != nil {
    panic(err)
}

decoded, err := profiles.ParsePacket(packet, profile)
if err != nil {
    panic(err)
}

fmt.Println(decoded.Format())
```

For generated metadata-driven profiles, type assert to `profiles.Decoded` if you need raw fields:

```go
d := decoded.(profiles.Decoded)
for key, value := range d.Values {
    fmt.Println(key, value.Raw, value.Text, value.Scaled, value.Unit)
}
```

Enum values are stable typed metadata:

```go
p, ok := profiles.Lookup(profile)
if ok {
    for _, f := range p.Fields {
        for _, ev := range f.Enums {
            fmt.Println(f.Shortcut, ev.Raw, ev.Name, ev.Description)
        }
    }
}
```

## Encode an EEP profile payload

```go
profile, _ := eep.FromString("F6-02-01")
userData, status, err := profiles.Encode(profile, map[string]uint64{
    "EB": 1,
    "SA": 1,
})
if err != nil {
    panic(err)
}

packet := erp1.Packet{
    Rorg:     profile.Rorg,
    UserData: userData,
    Status:   status,
}
```

Some common profiles also have small concrete types, e.g.:

```go
t := profiles.D50001{ContactClosed: true, LearnButton: true}
userData, status, err := t.MarshalERP1UserData()
```

## Regenerate EEP profiles

When `eep268.xml` changes:

```sh
go run ./cmd/eepgen -xml eep268.xml -out pkg/eep/profiles
go test ./...
```

The generator decodes UTF-16LE XML and writes `pkg/eep/profiles/profiles_gen.go`.

## Smart Ack

```go
msg, err := smartack.Parse(packet)
if err != nil {
    panic(err)
}

switch m := msg.(type) {
case smartack.LearnRequest:
    fmt.Println(m.EEP, m.RequestCode, m.ManufacturerID)
case smartack.DataReclaim:
    fmt.Println(m.MailboxIndex)
}
```

## Remote Management / Remote Commissioning / Security

Packages are split by layer:

- `pkg/reman`: Remote Management SYS_EX message split/merge and RMCC basics
- `pkg/recom`: Remote Commissioning payload helpers
- `pkg/srm`: Secure Remote Management SYS_EX/RPC headers
- `pkg/security`: SEC_R encode/decode, CMAC, SEC_CDM chain helpers
- `pkg/ddf`: minimal DDF V2 metadata parser

Plain ReMan/ReCom SYS_EX messages do not contain a rolling code. Secure envelopes
(legacy SEC_MAN or modern SEC_R/SEC_CDM) use the RLC defined by their SLF; RLC
fields in ReCom security-profile RPCs are configuration state, not an extra envelope.

Example Remote Management message:

```go
msg := reman.Message{
    Seq:            1,
    ManufacturerID: reman.ManufacturerID,
    Function:       reman.FuncQueryID,
    Payload:        nil,
}
packets, err := msg.Packets()
```

## Development

```sh
go test ./...
```

No runtime dependency on `eep268.xml`; generated profile metadata is committed as Go code.
