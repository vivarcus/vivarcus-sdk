//go:build integration

package example_test

import (
	"path/filepath"
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestLifecycleEntryIntegration(t *testing.T) {
	env := itest.Open(t)
	const (
		object = "demo_request__c"
		field  = "title__c"
		want   = "stamped-on-enter"
		fqn    = "acme.corp.sdkdemo.StampOnEnter"
	)
	env.ApplyMDL(t, "mdl/01-object.mdl")
	env.SDKPut(t, "actions/stamp_on_enter.go")

	dir := t.TempDir()
	env.RenderMDL(t, "mdl/02-lifecycle.mdl", filepath.Join(dir, "02-lifecycle.mdl"), []string{"ACTION_FQN=" + fqn})
	env.RenderMDL(t, "mdl/04-workflow-action.mdl", filepath.Join(dir, "04-workflow-action.mdl"), []string{"ACTION_FQN=" + fqn})
	env.RenderMDL(t, "mdl/05-workflow-cancel.mdl", filepath.Join(dir, "05-workflow-cancel.mdl"), []string{"ACTION_FQN=" + fqn})
	env.ApplyMDL(t, filepath.Join(dir, "02-lifecycle.mdl"))
	env.ApplyMDL(t, "mdl/03-bind-object-lifecycle.mdl")
	env.ApplyMDL(t, filepath.Join(dir, "04-workflow-action.mdl"))
	env.ApplyMDL(t, filepath.Join(dir, "05-workflow-cancel.mdl"))

	assertTitle := func(rec, label string) {
		t.Helper()
		if got := env.Field(t, object, rec, field); got != want {
			t.Fatalf("%s: title__c=%v", label, got)
		}
	}

	recEvent := env.CreateRecord(t, object, map[string]any{"name__v": "lc-event", field: "before-event"})
	assertTitle(recEvent, "event_action")

	recEntry := env.CreateRecord(t, object, map[string]any{"name__v": "lc-entry", field: "before-entry"})
	env.Transition(t, object, recEntry, "submit__c")
	assertTitle(recEntry, "entry_action")

	recWF := env.CreateRecord(t, object, map[string]any{"name__v": "lc-wf", field: "before-wf"})
	env.StartWorkflow(t, object, recWF, "stamp_on_enter_wf__c")
	assertTitle(recWF, "workflow action")

	recCancel := env.CreateRecord(t, object, map[string]any{"name__v": "lc-cancel", field: "before-cancel"})
	started := env.StartWorkflow(t, object, recCancel, "stamp_on_enter_wf_cancel__c")
	instanceID, _ := started["workflow_instance_id"].(string)
	if instanceID == "" {
		t.Fatalf("cancel workflow: missing workflow_instance_id: %v", started)
	}
	env.CancelWorkflow(t, instanceID)
	assertTitle(recCancel, "workflow cancel")
}
