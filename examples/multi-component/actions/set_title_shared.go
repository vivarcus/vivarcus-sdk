package actions

import (
	"github.com/acme.corp.sdkdemo/shared"
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// SetTitleShared writes title__c via shared helpers.
// Distinct from examples/02-update-field SetTitle so both can sit on one module.
type SetTitleShared struct{}

func (SetTitleShared) Meta() action.Meta {
	return action.Meta{
		Label:        "Set Title (shared)",
		Object:       "multi_demo__c",
		ObjectAction: "set_title_shared__c",
		Usages:       []action.Usage{action.UsageUserAction},
		RunAs:        action.RunAsSystemUser,
	}
}

func (SetTitleShared) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (SetTitleShared) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	want := shared.DemoTitleValue
	platform.LogInfo("set_title_shared")
	rec.SetValue(shared.TitleField, want)
	if err := platform.Update(rec.Object, rec.ID, shared.TitlePatch(want)); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}
