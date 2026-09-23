package triggers

import (
	"github.com/acme.corp.sdkdemo/shared"
	"github.com/vivarcus/vivarcus-sdk/trigger"
)

// StampDemoName runs BEFORE_INSERT and stamps name__v.
// Distinct from examples/05-stamp-trigger StampName; both skip when the suffix is already present.
type StampDemoName struct{}

func (StampDemoName) Meta() trigger.Meta {
	return trigger.Meta{
		Label:        "Stamp Demo Name",
		Object:       "multi_demo__c",
		Events:       []trigger.Event{trigger.BeforeInsert},
		EventSegment: trigger.PreCustom,
		Order:        4,
	}
}

func (StampDemoName) Execute(ctx trigger.RecordTriggerContext) error {
	if len(ctx.RecordChanges) == 0 || ctx.RecordChanges[0].New == nil {
		return nil
	}
	name := ctx.RecordChanges[0].New.GetString(shared.NameField)
	ctx.RecordChanges[0].New.SetValue(shared.NameField, shared.StampName(name))
	return nil
}
