package actions

import (
	"fmt"

	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// QueryAndStamp runs a read-only VQL against the action object, then writes title__c on the context record.
type QueryAndStamp struct{}

func (QueryAndStamp) Meta() action.Meta {
	return action.Meta{
		Label:        "Stamp From Query",
		Object:       "sdk_demo__c",
		ObjectAction: "stamp_from_query__c",
		Usages:       []action.Usage{action.UsageUserAction},
		RunAs:        action.RunAsSystemUser,
	}
}

func (QueryAndStamp) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (QueryAndStamp) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	vql := fmt.Sprintf("select id from %s pagesize 50", rec.Object)
	res, err := platform.Query(vql)
	if err != nil {
		return action.ExecuteResult{}, err
	}
	title := fmt.Sprintf("query:%d", len(res.Records))
	platform.LogInfo("query rows=" + fmt.Sprint(len(res.Records)))
	if err := platform.Update(rec.Object, rec.ID, map[string]any{
		"title__c": title,
	}); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}

