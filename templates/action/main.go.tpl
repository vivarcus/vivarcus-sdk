package main

import (
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// {{TYPE_NAME}} is a Record Action for {{OBJECT}}.
type {{TYPE_NAME}} struct{}

func ({{TYPE_NAME}}) Meta() action.Meta {
	return action.Meta{
		Label:  "{{LABEL}}",
		Object: "{{OBJECT}}",
		Usages: []action.Usage{action.UsageUserAction},
		RunAs:  action.RunAsSystemUser,
	}
}

func ({{TYPE_NAME}}) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func ({{TYPE_NAME}}) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	platform.LogInfo("{{TYPE_NAME}} execute")
	// Example: update a field
	rec.SetValue("title__c", "updated-by-sdk")
	if err := platform.Update(rec.Object, rec.ID, map[string]any{
		"title__c": "updated-by-sdk",
	}); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}

func main() {}
