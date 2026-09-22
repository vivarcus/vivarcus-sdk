package workflowactions

import "github.com/vivarcus/vivarcus-sdk/workflowaction"

// {{TYPE_NAME}} is a Record Workflow Action.
type {{TYPE_NAME}} struct{}

func ({{TYPE_NAME}}) Meta() workflowaction.Meta {
	return workflowaction.Meta{
		Label:     "{{LABEL}}",
		Object:    "{{OBJECT}}",
		StepTypes: []workflowaction.StepType{workflowaction.StepStart},
	}
}

func ({{TYPE_NAME}}) Execute(ctx workflowaction.RecordWorkflowActionContext) error {
	return nil
}

