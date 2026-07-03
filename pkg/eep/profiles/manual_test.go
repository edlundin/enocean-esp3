package profiles

import "testing"

func TestManualProfilesRoundTrip(t *testing.T) {
	cases := []Telegram{D50001{ContactClosed: true, LearnButton: false}, F60101{Pressed: true}, A50201{TemperatureC: -20, LearnButton: false}}
	for _, in := range cases {
		data, status, err := in.MarshalERP1UserData()
		if err != nil { t.Fatal(err) }
		var out Telegram
		switch in.(type) {
		case D50001:
			out, err = parseD50001(data, status)
		case F60101:
			out, err = parseF60101(data, status)
		case A50201:
			out, err = parseA50201(data, status)
		}
		if err != nil { t.Fatal(err) }
		if out.EEP() != in.EEP() || out.Format() == "" { t.Fatalf("bad round trip %#v", out) }
	}
	if _, err := parseD50001(nil, 0); err == nil { t.Fatal("D5 short accepted") }
	if _, err := parseF60101(nil, 0); err == nil { t.Fatal("F6 short accepted") }
	if _, err := parseA50201(nil, 0); err == nil { t.Fatal("A5 short accepted") }
}
