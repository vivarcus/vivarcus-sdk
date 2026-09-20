package main

import (
	"github.com/vivarcus/vivarcus-sdk/trigger"
	"github.com/acme.corp.sdkdemo/shared"
)

// StampName runs BEFORE_INSERT and stamps name__v.
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
	name := ctx.RecordChanges[0].New.GetString(shared.NameField)
	ctx.RecordChanges[0].New.SetValue(shared.NameField, shared.StampName(name))
	return nil
}
