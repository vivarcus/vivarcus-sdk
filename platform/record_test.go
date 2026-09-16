package platform_test

import (
	"testing"

	"github.com/vivarcus/vivarcus-sdk/platform"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

func TestGetUpdateUnavailable(t *testing.T) {
	t.Cleanup(func() {
		platform.GetRecordFunc = nil
		platform.UpdateRecordFunc = nil
		platform.CreateRecordFunc = nil
		platform.DeleteRecordFunc = nil
	})
	platform.GetRecordFunc = nil
	platform.UpdateRecordFunc = nil
	if _, err := platform.Get("o", "r"); err == nil {
		t.Fatal("expected get error")
	}
	if err := platform.Update("o", "r", nil); err == nil {
		t.Fatal("expected update error")
	}
	if _, err := platform.Create("o", nil); err == nil {
		t.Fatal("expected create error")
	}
	if err := platform.Delete("o", "r"); err == nil {
		t.Fatal("expected delete error")
	}
}

func TestGetUpdateRoundTrip(t *testing.T) {
	t.Cleanup(func() {
		platform.GetRecordFunc = nil
		platform.UpdateRecordFunc = nil
	})
	rec := wire.Record{ID: "r1", Object: "demo__c", Fields: []wire.Field{wire.FieldFromAny("title__c", "n")}}
	raw, err := wire.Encode(rec)
	if err != nil {
		t.Fatal(err)
	}
	platform.GetRecordFunc = func(object, recordID string) ([]byte, int32) {
		if object != "demo__c" || recordID != "r1" {
			t.Fatalf("%s %s", object, recordID)
		}
		return raw, 0
	}
	got, err := platform.Get("demo__c", "r1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != "r1" || got.GetString("title__c") != "n" {
		t.Fatalf("%+v", got)
	}
	var seen []byte
	platform.UpdateRecordFunc = func(object, recordID string, fieldsJSON []byte) int32 {
		seen = append([]byte(nil), fieldsJSON...)
		return 0
	}
	if err := platform.Update("demo__c", "r1", map[string]any{"title__c": "n"}); err != nil {
		t.Fatal(err)
	}
	if len(seen) == 0 {
		t.Fatal("expected update payload")
	}
}

func TestGetUpdateWhenFuncsWired(t *testing.T) {
	platform.GetRecordFunc = func(object, recordID string) ([]byte, int32) {
		raw, err := wire.Encode(wire.Record{
			ID:     recordID,
			Object: object,
			Fields: []wire.Field{{Name: "name__v", Kind: "string", S: "N"}},
		})
		if err != nil {
			t.Fatal(err)
		}
		return raw, 0
	}
	t.Cleanup(func() { platform.GetRecordFunc = nil })
	rec, err := platform.Get("study__v", "1")
	if err != nil {
		t.Fatal(err)
	}
	if rec.GetString("name__v") != "N" {
		t.Fatalf("%+v", rec)
	}
	platform.UpdateRecordFunc = func(object, recordID string, fieldsJSON []byte) int32 { return 0 }
	t.Cleanup(func() { platform.UpdateRecordFunc = nil })
	if err := platform.Update("study__v", "1", map[string]any{"name__v": "X"}); err != nil {
		t.Fatal(err)
	}
}
