//go:build integration

package example_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestUpdateFieldIntegration(t *testing.T) {
	env := itest.Open(t)
	env.ApplyMDL(t, "mdl/01-object.mdl")
	env.SDKPut(t, "actions/set_title.go")

	rec := env.CreateRecord(t, "sdk_demo__c", map[string]any{"name__v": "update-field", "title__c": "before"})
	env.ExecuteAction(t, "sdk_demo__c", rec, "set_title__c")
	if got := env.Field(t, "sdk_demo__c", rec, "title__c"); got != "from-sdk" {
		t.Fatalf("title__c=%v", got)
	}
}
