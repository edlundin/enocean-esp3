package profiles

import "testing"

// TestFieldEnum verifies FieldEnum behavior.
func TestFieldEnum(t *testing.T) {
	f := Field{Enums: []EnumValue{{Raw: 1, Name: "one"}}}
	if e, ok := f.Enum(1); !ok || e.Name != "one" { t.Fatalf("enum hit = %#v %t", e, ok) }
	if _, ok := f.Enum(2); ok { t.Fatal("unexpected enum hit") }
}
