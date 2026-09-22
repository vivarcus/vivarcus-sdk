package actions

import (
	"github.com/vivarcus/vivarcus-sdk/action"
)

// NoopAction is a minimal UserAction used as the vivarcus-sdk example.
type NoopAction struct{}

func (NoopAction) Meta() action.Meta {
	return action.Meta{
		Label:  "Wasm Noop",
		Usages: []action.Usage{action.UsageUserAction},
		RunAs:  action.RunAsSystemUser,
	}
}

func (NoopAction) IsExecutable(ctx action.RecordActionContext) bool {
	return len(ctx.Records) > 0
}

func (NoopAction) Execute(ctx action.RecordActionContext) (action.ExecuteResult, error) {
	return action.ExecuteResult{}, nil
}

