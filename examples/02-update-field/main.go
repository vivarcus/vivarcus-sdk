package main

import (
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// SetTitle writes title__c on the current record via SetValue and host update.
type SetTitle struct{}

func (SetTitle) Meta() action.Meta {
	return action.Meta{
		Label:  "Set Title",
		Object: "sdk_demo__c",
		Usages: []action.Usage{action.UsageUserAction},
		RunAs:  action.RunAsSystemUser,
	}
}

func (SetTitle) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (SetTitle) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	platform.LogInfo("set_title")
	rec.SetValue("title__c", "from-sdk")
	if err := platform.Update(rec.Object, rec.ID, map[string]any{
		"title__c": "from-sdk",
	}); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}

func main() {}
