package actions

import (
	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
)

// StampOnEnter runs on lifecycle entry, event_action, workflow action steps, or
// workflow cancel rules that reference Recordaction.<FQN> (or Objectaction.*).
type StampOnEnter struct{}

func (StampOnEnter) Meta() action.Meta {
	return action.Meta{
		Label:  "Stamp On Enter",
		Object: "demo_request__c",
		// System-invoked paths (see examples/04-lifecycle-entry README Steps 4–8).
		// Add UsageLifecycleUserAction only if you also expose an Objectaction button.
		Usages: []action.Usage{
			action.UsageLifecycleEntryAction,
			action.UsageEventAction,
			action.UsageWorkflowStep,
			action.UsageWorkflowCancel,
		},
		RunAs: action.RunAsSystemUser,
	}
}

func (StampOnEnter) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (StampOnEnter) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	rec := ctx.Records[0]
	platform.LogInfo("stamp_on_enter")
	rec.SetValue("title__c", "stamped-on-enter")
	if err := platform.Update(rec.Object, rec.ID, map[string]any{
		"title__c": "stamped-on-enter",
	}); err != nil {
		return action.ExecuteResult{}, err
	}
	return action.ExecuteResult{}, nil
}

