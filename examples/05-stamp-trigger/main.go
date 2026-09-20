package main

import "github.com/vivarcus/vivarcus-sdk/trigger"

// StampName appends "-trig" to name__v on BEFORE_INSERT (demo Record Trigger).
type StampName struct{}

func (StampName) Meta() trigger.Meta {
	return trigger.Meta{
		Label:        "Stamp Name",
		Object:       "sdk_demo__c",
		Events:       []trigger.Event{trigger.BeforeInsert},
		EventSegment: trigger.PreCustom,
		Order:        3,
	}
}

func (StampName) Execute(ctx trigger.RecordTriggerContext) error {
	if len(ctx.RecordChanges) == 0 || ctx.RecordChanges[0].New == nil {
		return nil
	}
	name := ctx.RecordChanges[0].New.GetString("name__v")
	ctx.RecordChanges[0].New.SetValue("name__v", name+"-trig")
	return nil
}

func main() {}
