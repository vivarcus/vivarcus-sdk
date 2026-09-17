package main

import (
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// ConfirmDialog demonstrates OnPreExecute confirm and OnPostExecute banner.
type ConfirmDialog struct{}

func (ConfirmDialog) Meta() action.Meta {
	return action.Meta{
		Label:         "Confirm Update",
		Object:        "sdk_demo__c",
		ObjectAction:  "confirm_update__c",
		Usages: []action.Usage{action.UsageUserAction},
		RunAs:  action.RunAsSystemUser,
	}
}

func (ConfirmDialog) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (ConfirmDialog) OnPreExecute(ctx action.RecordActionContext) (action.PreExecuteResult, error) {
	return action.PreExecuteResult{
		Title:          "Confirm",
		ConfirmMessage: "Update title__c on this record?",
	}, nil
}

func (ConfirmDialog) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	rec.SetValue("title__c", "confirmed-by-sdk")
	return action.ExecuteResult{}, platform.Update(rec.Object, rec.ID, map[string]any{
		"title__c": "confirmed-by-sdk",
	})
}

func (ConfirmDialog) OnPostExecute(ctx action.RecordActionContext) (action.PostExecuteResult, error) {
	return action.PostExecuteResult{
		Message: "Title updated successfully.",
	}, nil
}

func main() {}
