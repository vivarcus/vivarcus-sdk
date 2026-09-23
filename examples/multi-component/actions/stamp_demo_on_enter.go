package actions

import (
	"github.com/acme.corp.sdkdemo/shared"
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// StampDemoOnEnter runs when a record enters in_review__c (lifecycle entry_action rule).
// Distinct from examples/04-lifecycle-entry StampOnEnter (that one binds demo_request__c).
type StampDemoOnEnter struct{}

func (StampDemoOnEnter) Meta() action.Meta {
	return action.Meta{
		Label:  "Stamp Demo On Enter",
		Object: "sdk_demo__c",
		Usages: []action.Usage{action.UsageLifecycleEntryAction},
		RunAs:  action.RunAsSystemUser,
	}
}

func (StampDemoOnEnter) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (StampDemoOnEnter) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	want := shared.EntryTitleValue
	platform.LogInfo("stamp_demo_on_enter")
	rec.SetValue(shared.TitleField, want)
	if err := platform.Update(rec.Object, rec.ID, shared.TitlePatch(want)); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}
