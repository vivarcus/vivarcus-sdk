package actions

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/platform"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

func TestSetTitleExecute(t *testing.T) {
	t.Cleanup(func() {
		platform.UpdateRecordFunc = nil
		platform.LogFunc = nil
	})
	var gotObject, gotID string
	var gotFields []wire.Field
	platform.UpdateRecordFunc = func(object, recordID string, fieldsJSON []byte) int32 {
		gotObject, gotID = object, recordID
		if err := wire.Decode(fieldsJSON, &gotFields); err != nil {
			t.Fatal(err)
		}
		return 0
	}

	rec := action.NewRecord("r1", "sdk_demo__c", nil)
	if _, err := (SetTitle{}).Execute(action.RecordActionContext{
		Records: []action.Record{rec},
	}); err != nil {
		t.Fatal(err)
	}
	if gotObject != "sdk_demo__c" || gotID != "r1" {
		t.Fatalf("update %s %s", gotObject, gotID)
	}
	found := false
	for _, f := range gotFields {
		if f.Name == "title__c" && f.Any() == "from-sdk" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("title__c not updated: %+v", gotFields)
	}
}

func TestSetTitleIsExecutable(t *testing.T) {
	if (SetTitle{}).IsExecutable(action.RecordActionContext{}) {
		t.Fatal("expected false with no records")
	}
	rec := action.NewRecord("r1", "sdk_demo__c", nil)
	if !(SetTitle{}).IsExecutable(action.RecordActionContext{Records: []action.Record{rec}}) {
		t.Fatal("expected true")
	}
}
