// Package trigger is the Vivarcus Record Trigger API for wasm guests (tinygo / wasip1 safe).
//
// Conceptually aligned with Veeva Vault Java SDK RecordTrigger.
package trigger

import (
	"fmt"
	"strings"

	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

// Event is a Record Trigger event (Java RecordEvent).
type Event string

const (
	BeforeInsert Event = "BEFORE_INSERT"
	AfterInsert  Event = "AFTER_INSERT"
	BeforeUpdate Event = "BEFORE_UPDATE"
	AfterUpdate  Event = "AFTER_UPDATE"
	BeforeDelete Event = "BEFORE_DELETE"
	AfterDelete  Event = "AFTER_DELETE"
)

// EventSegment is PRE_CUSTOM or POST_CUSTOM.
type EventSegment string

const (
	PreCustom      EventSegment = "PRE_CUSTOM"
	PostCustom     EventSegment = "POST_CUSTOM"
	UnspecifiedSeg EventSegment = "UNSPECIFIED"
)

// Meta declares component-level static metadata (Java @RecordTriggerInfo).
type Meta struct {
	Name         string
	Label        string
	Object       string
	Events       []Event
	EventSegment EventSegment
	Order        int
	RunAs        action.RunAs
}

// Validate checks required Meta fields.
func (m Meta) Validate() error {
	if strings.TrimSpace(m.Label) == "" {
		return fmt.Errorf("meta label required")
	}
	if len(m.Events) == 0 {
		return fmt.Errorf("meta events required")
	}
	for _, ev := range m.Events {
		switch ev {
		case BeforeInsert, AfterInsert, BeforeUpdate, AfterUpdate, BeforeDelete, AfterDelete:
		default:
			return fmt.Errorf("invalid event %q", ev)
		}
	}
	switch m.EventSegment {
	case "", PreCustom, PostCustom, UnspecifiedSeg:
	default:
		return fmt.Errorf("invalid event_segment %q", m.EventSegment)
	}
	if m.Order < 0 || m.Order > 10 {
		return fmt.Errorf("invalid order %d", m.Order)
	}
	if m.RunAs != "" && m.RunAs != action.RunAsSystemUser && m.RunAs != action.RunAsRequestOwner {
		return fmt.Errorf("invalid run_as %q", m.RunAs)
	}
	return nil
}

// RecordTriggerContext is the runtime context for Record Trigger execution.
type RecordTriggerContext struct {
	RecordChanges    []RecordChange
	EventSegment     EventSegment
	TriggerEvent     Event
	VaultID          string
	CurrentUserID    string
	InitiatingUserID string
}

// RecordChange is one record mutation in this DML (Phase 1 length is always 1).
type RecordChange struct {
	Old *action.Record
	New *action.Record
	err *Error
}

// SetError marks this change as failed (whole transaction rolls back).
func (c *RecordChange) SetError(subtype, message string) {
	if c == nil {
		return
	}
	c.err = &Error{Subtype: subtype, Message: message}
}

// Err returns the error recorded by SetError, if any.
func (c *RecordChange) Err() error {
	if c == nil || c.err == nil {
		return nil
	}
	return c.err
}

// Error is a user-visible Record Trigger failure.
type Error struct {
	Subtype string
	Message string
}

func (e *Error) Error() string {
	if e == nil {
		return "record trigger failed"
	}
	msg := strings.TrimSpace(e.Message)
	if msg != "" {
		return msg
	}
	sub := strings.TrimSpace(e.Subtype)
	if sub != "" {
		return sub
	}
	return "record trigger failed"
}

// RecordTrigger is the guest entry interface (Java RecordTrigger).
type RecordTrigger interface {
	Meta() Meta
	Execute(ctx RecordTriggerContext) error
}

// ToWireMeta converts Meta to the describe payload.
func ToWireMeta(m Meta) wire.TriggerMeta {
	events := make([]string, 0, len(m.Events))
	for _, ev := range m.Events {
		events = append(events, string(ev))
	}
	runAs := string(m.RunAs)
	if runAs == "" {
		runAs = string(action.RunAsSystemUser)
	}
	seg := string(m.EventSegment)
	if seg == "" || seg == string(UnspecifiedSeg) {
		seg = string(PostCustom)
	}
	return wire.TriggerMeta{
		ComponentName: strings.TrimSpace(m.Name),
		Label:         m.Label,
		Object:        m.Object,
		Events:        events,
		EventSegment:  seg,
		Order:         m.Order,
		RunAs:         runAs,
	}
}

// ContextFromWire converts a host trigger invocation into the guest API shape.
func ContextFromWire(w wire.TriggerContext, host action.FieldWriter) RecordTriggerContext {
	change := RecordChange{}
	if w.Old != nil {
		rec := action.NewRecord(w.Old.ID, w.Old.Object, w.Old.Fields)
		action.AttachFieldWriter(&rec, host)
		change.Old = &rec
	}
	if w.New != nil {
		rec := action.NewRecord(w.New.ID, w.New.Object, w.New.Fields)
		action.AttachFieldWriter(&rec, host)
		change.New = &rec
	}
	return RecordTriggerContext{
		RecordChanges:    []RecordChange{change},
		EventSegment:     EventSegment(w.EventSegment),
		TriggerEvent:     Event(w.Event),
		VaultID:          w.VaultID,
		CurrentUserID:    w.CurrentUserID,
		InitiatingUserID: w.InitiatingUserID,
	}
}
