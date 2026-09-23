//go:build integration

package example_test

import (
	"strings"
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestStampTriggerIntegration(t *testing.T) {
	env := itest.Open(t)
	env.ApplyMDL(t, "mdl/01-object.mdl")
	env.SDKPut(t, "triggers/stamp_name.go")

	rec := env.CreateRecord(t, "sdk_demo__c", map[string]any{"name__v": "stamp-trigger"})
	got, _ := env.Field(t, "sdk_demo__c", rec, "name__v").(string)
	if !strings.HasSuffix(got, "-trig") {
		t.Fatalf("name__v=%q", got)
	}
}
