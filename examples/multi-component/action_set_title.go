package main

import (
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
	"github.com/acme.corp.sdkdemo/shared"
)

// SetTitle writes title__c via shared helpers.
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
	want := shared.DemoTitleValue
	platform.LogInfo("set_title")
	rec.SetValue(shared.TitleField, want)
	if err := platform.Update(rec.Object, rec.ID, shared.TitlePatch(want)); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}
