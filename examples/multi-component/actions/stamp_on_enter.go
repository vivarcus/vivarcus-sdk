package actions

import (
	"github.com/acme.corp.sdkdemo/shared"
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// StampOnEnter runs when a record enters in_review__c (lifecycle entry_action rule).
type StampOnEnter struct{}

func (StampOnEnter) Meta() action.Meta {
	return action.Meta{
		Label:  "Stamp On Enter",
		Object: "sdk_demo__c",
		Usages: []action.Usage{action.UsageLifecycleEntryAction},
		RunAs:  action.RunAsSystemUser,
	}
}

func (StampOnEnter) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (StampOnEnter) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	want := shared.EntryTitleValue
	platform.LogInfo("stamp_on_enter")
	rec.SetValue(shared.TitleField, want)
	if err := platform.Update(rec.Object, rec.ID, shared.TitlePatch(want)); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}
