package wire_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/wire"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	in := wire.Meta{Label: "X", Usages: []string{"UserAction"}, RunAs: "SystemUser"}
	raw, err := wire.Encode(in)
	if err != nil {
		t.Fatal(err)
	}
	var out wire.Meta
	if err := wire.Decode(raw, &out); err != nil {
		t.Fatal(err)
	}
	if out.Label != "X" || out.RunAs != "SystemUser" {
		t.Fatalf("%+v", out)
	}
}
