//go:build integration

package example_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestRecordWorkflowActionIntegration(t *testing.T) {
	env := itest.Open(t)
	env.ApplyMDL(t, "mdl/01-object.mdl")
	env.ApplyMDL(t, "mdl/02-lifecycle.mdl")
	env.ApplyMDL(t, "mdl/03-bind-object-lifecycle.mdl")
	env.SDKPut(t, "workflowactions/capture_participants.go")
	env.ApplyMDL(t, "mdl/04-workflow.mdl")
	env.SDKGetContains(t, "acme.corp.sdkdemo.CaptureParticipants", "CaptureParticipants")

	rec := env.CreateRecord(t, "sdk_demo__c", map[string]any{"name__v": "capture-participants"})
	started := env.Transition(t, "sdk_demo__c", rec, "start_capture_wf__c")
	if started == nil {
		t.Fatal("workflow start returned nil")
	}
}
