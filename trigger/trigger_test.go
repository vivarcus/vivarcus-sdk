package trigger_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/trigger"
)

func TestMetaValidate(t *testing.T) {
	t.Parallel()
	m := trigger.Meta{Label: "T", Events: []trigger.Event{trigger.BeforeInsert}}
	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (trigger.Meta{Events: []trigger.Event{trigger.BeforeInsert}}).Validate(); err == nil {
		t.Fatal("label required")
	}
}

func TestSetError(t *testing.T) {
	t.Parallel()
	var c trigger.RecordChange
	c.SetError("MISSING_PARAMETER", "name required")
	if c.Err() == nil || c.Err().Error() != "name required" {
		t.Fatalf("%v", c.Err())
	}
}
