package actions

import (
	"github.com/acme.corp.sdkdemo/shared"
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// ClearTitle clears title__c via shared helpers.
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
	rec.SetValue(shared.TitleField, "")
	if err := platform.Update(rec.Object, rec.ID, shared.TitlePatch("")); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}
