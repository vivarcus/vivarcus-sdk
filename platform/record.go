package platform

import (
	"fmt"

	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

// GetRecordFunc is wired by the reactor to host vivarcus:platform/record.get.
// encoded is length-prefixed JSON wire.Record; errCode is 0 on success.
var GetRecordFunc func(object, recordID string) (encoded []byte, errCode int32)

// UpdateRecordFunc is wired by the reactor to host vivarcus:platform/record.update.
var UpdateRecordFunc func(object, recordID string, fieldsJSON []byte) int32

// CreateRecordFunc is wired by the reactor when the host implements create.
var CreateRecordFunc func(object string, fieldsJSON []byte) (id string, errCode int32)

// DeleteRecordFunc is wired by the reactor when the host implements delete.
var DeleteRecordFunc func(object, recordID string) int32

// Get loads one record via the host (Java RecordService / QueryService subset).
func Get(object, recordID string) (action.Record, error) {
	if GetRecordFunc == nil {
		return action.Record{}, fmt.Errorf("record_action_host_call_failed: get unavailable")
	}
	raw, code := GetRecordFunc(object, recordID)
	if code != 0 {
		return action.Record{}, fmt.Errorf("record_action_host_call_failed: get status %d", code)
	}
	var rec wire.Record
	if len(raw) > 0 {
		if err := wire.Decode(raw, &rec); err != nil {
			return action.Record{}, fmt.Errorf("record_action_host_call_failed: %w", err)
		}
	}
	if rec.ID == "" {
		rec.ID = recordID
	}
	if rec.Object == "" {
		rec.Object = object
	}
	return action.NewRecord(rec.ID, rec.Object, rec.Fields), nil
}

// Update writes fields on a record via the host (Java RecordService.batchSaveRecords).
func Update(object, recordID string, fields map[string]any) error {
	if UpdateRecordFunc == nil {
		return fmt.Errorf("record_action_host_call_failed: update unavailable")
	}
	typed := make([]wire.Field, 0, len(fields))
	for name, v := range fields {
		typed = append(typed, wire.FieldFromAny(name, v))
	}
	raw, err := wire.Encode(typed)
	if err != nil {
		return fmt.Errorf("record_action_host_call_failed: %w", err)
	}
	if code := UpdateRecordFunc(object, recordID, raw); code != 0 {
		return fmt.Errorf("record_action_host_call_failed: update status %d", code)
	}
	return nil
}

// Create inserts a record via the host. Phase 1 host may return NOT_IMPLEMENTED.
func Create(object string, fields map[string]any) (string, error) {
	if CreateRecordFunc == nil {
		return "", fmt.Errorf("record_action_host_call_failed: create unavailable")
	}
	typed := make([]wire.Field, 0, len(fields))
	for name, v := range fields {
		typed = append(typed, wire.FieldFromAny(name, v))
	}
	raw, err := wire.Encode(typed)
	if err != nil {
		return "", fmt.Errorf("record_action_host_call_failed: %w", err)
	}
	id, code := CreateRecordFunc(object, raw)
	if code != 0 {
		return "", fmt.Errorf("record_action_host_call_failed: create status %d", code)
	}
	return id, nil
}

// Delete removes a record via the host. Phase 1 host may return NOT_IMPLEMENTED.
func Delete(object, recordID string) error {
	if DeleteRecordFunc == nil {
		return fmt.Errorf("record_action_host_call_failed: delete unavailable")
	}
	if code := DeleteRecordFunc(object, recordID); code != 0 {
		return fmt.Errorf("record_action_host_call_failed: delete status %d", code)
	}
	return nil
}
