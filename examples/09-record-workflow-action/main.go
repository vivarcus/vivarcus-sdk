package main

import "github.com/vivarcus/vivarcus-sdk/workflowaction"

// CaptureParticipants is a start-step Record Workflow Action that pre-fills
// the current participant group with the workflow initiating user.
type CaptureParticipants struct{}

func (CaptureParticipants) Meta() workflowaction.Meta {
	return workflowaction.Meta{
		Label:     "Capture Participants",
		Object:    "sdk_demo__c",
		StepTypes: []workflowaction.StepType{workflowaction.StepStart},
	}
}

func (CaptureParticipants) Execute(ctx workflowaction.RecordWorkflowActionContext) error {
	if ctx.Event != workflowaction.EventGetParticipants {
		return nil
	}
	group := "approver"
	if ctx.ParticipantGroup != nil && ctx.ParticipantGroup.Name != "" {
		group = ctx.ParticipantGroup.Name
	}
	if ctx.InitiatingUserID == "" {
		return nil
	}
	ctx.SetParticipantUsers(group, []string{ctx.InitiatingUserID})
	return nil
}

func main() {}
