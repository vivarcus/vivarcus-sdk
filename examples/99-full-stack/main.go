package main

import (
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// SetTitle writes title__c on the current record.
type SetTitle struct{}

func (SetTitle) Meta() action.Meta {
	return action.Meta{
		Label:        "Set Title",
		Object:       "sdk_demo__c",
		ObjectAction: "set_title__c",
		Usages:       []action.Usage{action.UsageUserAction},
		RunAs:        action.RunAsSystemUser,
	}
}

func (SetTitle) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (SetTitle) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	const want = "from-sdk"
	platform.LogInfo("set_title")
	rec.SetValue("title__c", want)
	if err := platform.Update(rec.Object, rec.ID, map[string]any{"title__c": want}); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}

// ClearTitle clears title__c on the current record.
type ClearTitle struct{}

func (ClearTitle) Meta() action.Meta {
	return action.Meta{
		Label:        "Clear Title",
		Object:       "sdk_demo__c",
		ObjectAction: "clear_title__c",
		Usages:       []action.Usage{action.UsageUserAction},
		RunAs:        action.RunAsSystemUser,
	}
}

func (ClearTitle) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (ClearTitle) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	platform.LogInfo("clear_title")
	rec.SetValue("title__c", "")
	if err := platform.Update(rec.Object, rec.ID, map[string]any{"title__c": ""}); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}

func main() {}
