package job

import "testing"

func TestMetaValidate(t *testing.T) {
	t.Parallel()
	if err := (Meta{}).Validate(); err == nil {
		t.Fatal("expected label required")
	}
	if err := (Meta{Label: "Evaluate"}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestParamAccess(t *testing.T) {
	t.Parallel()
	ctx := InitContext{Params: map[string]any{
		"id":  "r1",
		"ids": []any{"a", "b"},
	}}
	if ctx.ParamString("id") != "r1" {
		t.Fatalf("ParamString=%q", ctx.ParamString("id"))
	}
	got := ctx.ParamStrings("ids")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("ParamStrings=%v", got)
	}
}
