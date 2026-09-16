package action

import (
	"strings"
	"testing"
	"time"

	"github.com/vivarcus/vivarcus-sdk/wire"
)

func TestMetaToWireDefaults(t *testing.T) {
	t.Parallel()
	wm := ToWireMeta(Meta{Label: "Hello"})
	if wm.RunAs != string(RunAsSystemUser) {
		t.Fatalf("run_as=%q", wm.RunAs)
	}
	if err := (Meta{Label: "Hello"}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestContextFromWireUserInputAndTime(t *testing.T) {
	t.Parallel()
	wctx := wire.Context{
		Records: []wire.Record{{
			ID:     "1",
			Object: "study__v",
			Fields: []wire.Field{
				{Name: "name__v", Kind: "string", S: "S1"},
				{Name: "when__v", Kind: "time", T: time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC).Format(time.RFC3339)},
			},
		}},
		UserInputRecord: &wire.Record{
			ID:     "IN-1",
			Object: "input__c",
			Fields: []wire.Field{{Name: "note__c", Kind: "string", S: "ok"}},
		},
		VaultID:          "v",
		CurrentUserID:    "u1",
		InitiatingUserID: "u1",
		Usage:            "UserAction",
	}
	ctx := ContextFromWire(wctx, nil)
	if len(ctx.Records) != 1 || ctx.Records[0].GetString("name__v") != "S1" {
		t.Fatalf("records=%+v", ctx.Records)
	}
	if ctx.Records[0].GetTime("when__v").IsZero() {
		t.Fatal("expected GetTime")
	}
	if ctx.UserInputRecord == nil || ctx.UserInputRecord.GetString("note__c") != "ok" {
		t.Fatal("user input")
	}
	if ctx.Usage != UsageUserAction {
		t.Fatalf("usage=%q", ctx.Usage)
	}
}

func TestSupportsUsageGuest(t *testing.T) {
	t.Parallel()
	if SupportsUsage(nil, UsageUserBulkAction) {
		t.Fatal("unspecified excludes bulk")
	}
	if !strings.Contains((Meta{}).Validate().Error(), "label") {
		t.Fatal("validate")
	}
}
