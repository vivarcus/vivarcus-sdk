package workflowaction_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/workflowaction"
)

func TestMetaValidate(t *testing.T) {
	t.Parallel()
	m := workflowaction.Meta{
		Label:     "Capture",
		StepTypes: []workflowaction.StepType{workflowaction.StepStart},
	}
	if err := m.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (workflowaction.Meta{StepTypes: []workflowaction.StepType{workflowaction.StepStart}}).Validate(); err == nil {
		t.Fatal("label required")
	}
	if err := (workflowaction.Meta{Label: "X"}).Validate(); err == nil {
		t.Fatal("step_types required")
	}
}

func TestSetParticipantUsers(t *testing.T) {
	t.Parallel()
	ctx := workflowaction.RecordWorkflowActionContext{
		Participants: map[string][]string{},
		ParticipantGroup: &workflowaction.ParticipantGroup{
			Name: "approver",
		},
	}
	ctx.SetParticipantUsers("approver", []string{"user-1"})
	if got := ctx.Participants["approver"]; len(got) != 1 || got[0] != "user-1" {
		t.Fatalf("%v", ctx.Participants)
	}
	if len(ctx.ParticipantGroup.Users) != 1 || ctx.ParticipantGroup.Users[0] != "user-1" {
		t.Fatalf("%+v", ctx.ParticipantGroup)
	}
}
