//go:build integration

package example_test

import (
	"strings"
	"testing"

	"github.com/vivarcus/vivarcus-sdk/examples/itest"
)

func TestQueryFieldIntegration(t *testing.T) {
	env := itest.Open(t)
	env.ApplyMDL(t, "mdl/01-object.mdl")
	env.SDKPut(t, "actions/query_and_stamp.go")

	rec := env.CreateRecord(t, "sdk_demo__c", map[string]any{"name__v": "query-field", "title__c": "before"})
	env.ExecuteAction(t, "sdk_demo__c", rec, "stamp_from_query__c")
	got, _ := env.Field(t, "sdk_demo__c", rec, "title__c").(string)
	if !strings.HasPrefix(got, "query:") || got == "query:0" {
		t.Fatalf("title__c=%q", got)
	}
}
