// Package workflowaction is the Vivarcus Record Workflow Action API for wasm guests
// (tinygo / wasip1 safe).
//
// Conceptually aligned with Veeva Vault Java SDK RecordWorkflowAction.
package workflowaction

import (
	"fmt"
	"strings"

	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

// StepType is where a Record Workflow Action may be configured (Veeva WorkflowStepType).
type StepType string

const (
	StepStart StepType = "START"
	StepTask  StepType = "TASK"
)

// Event is the workflow lifecycle hook that invokes Execute (Veeva WorkflowEvent).
type Event string

const (
	EventDisplayParticipants      Event = "DISPLAY_PARTICIPANTS"
	EventGetParticipants          Event = "GET_PARTICIPANTS"
	EventAfterCreate              Event = "AFTER_CREATE"
	EventTaskAfterCreate          Event = "TASK_AFTER_CREATE"
	EventTaskBeforeCompleteDialog Event = "TASK_BEFORE_COMPLETE_DIALOG"
	EventTaskAfterComplete        Event = "TASK_AFTER_COMPLETE"
	EventTaskAfterCancel          Event = "TASK_AFTER_CANCEL"
	EventTaskAfterAssign          Event = "TASK_AFTER_ASSIGN"
)

// StepTypeForEvent returns the workflow step type for an event.
func StepTypeForEvent(ev Event) StepType {
	switch ev {
	case EventDisplayParticipants, EventGetParticipants, EventAfterCreate:
		return StepStart
	default:
		return StepTask
	}
}

// Meta declares component-level static metadata (Java @RecordWorkflowActionInfo).
type Meta struct {
	Name      string
	Label     string
	Object    string
	StepTypes []StepType
	RunAs     action.RunAs
}

// Validate checks required Meta fields.
func (m Meta) Validate() error {
	if strings.TrimSpace(m.Label) == "" {
		return fmt.Errorf("meta label required")
	}
	if len(m.StepTypes) == 0 {
		return fmt.Errorf("meta step_types required")
	}
	for _, st := range m.StepTypes {
		switch st {
		case StepStart, StepTask:
		default:
			return fmt.Errorf("invalid step_type %q", st)
		}
	}
	if m.RunAs != "" && m.RunAs != action.RunAsSystemUser && m.RunAs != action.RunAsRequestOwner {
		return fmt.Errorf("invalid run_as %q", m.RunAs)
	}
	return nil
}

// SupportsStepType reports whether the action may run on the given step type.
func (m Meta) SupportsStepType(st StepType) bool {
	for _, allowed := range m.StepTypes {
		if allowed == st {
			return true
		}
	}
	return false
}

// NormalizedRunAs returns SystemUser when empty.
func (m Meta) NormalizedRunAs() action.RunAs {
	switch m.RunAs {
	case action.RunAsRequestOwner:
		return action.RunAsRequestOwner
	default:
		return action.RunAsSystemUser
	}
}

// ParticipantGroup is one start-step participant control group during execution.
type ParticipantGroup struct {
	Name  string
	Label string
	Users []string
}

// TaskInstance is a snapshot of one workflow task for task-step events.
type TaskInstance struct {
	TaskID         string
	StepAPIName    string
	AssigneeUserID string
	Status         string
}

// TaskChange pairs old/new task rows for assign events.
type TaskChange struct {
	Old TaskInstance
	New TaskInstance
}

// RecordWorkflowActionContext is the runtime context for Record Workflow Action execution.
type RecordWorkflowActionContext struct {
	Event   Event
	VaultID string

	WorkflowInstanceID string
	WorkflowAPIName    string
	WorkflowLabel      string
	ObjectAPIName      string
	RecordIDs          []string
	InitiatingUserID   string
	CurrentUserID      string

	ParticipantGroup *ParticipantGroup
	// Participants is the mutable start-step selection map (group name → user record ids).
	Participants map[string][]string

	TaskStepAPIName string
	TaskChanges     []TaskChange
}

// SetParticipantUsers mutates one participant group (GET_PARTICIPANTS / DISPLAY_PARTICIPANTS).
func (c *RecordWorkflowActionContext) SetParticipantUsers(groupName string, userIDs []string) {
	if c == nil {
		return
	}
	if c.Participants == nil {
		c.Participants = map[string][]string{}
	}
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return
	}
	c.Participants[groupName] = append([]string(nil), userIDs...)
	if c.ParticipantGroup != nil && c.ParticipantGroup.Name == groupName {
		c.ParticipantGroup.Users = append([]string(nil), userIDs...)
	}
}

// RecordWorkflowAction is the required guest interface (Java RecordWorkflowAction).
type RecordWorkflowAction interface {
	Meta() Meta
	Execute(ctx RecordWorkflowActionContext) error
}

// ToWireMeta converts Meta to the describe payload.
func ToWireMeta(m Meta) wire.WorkflowActionMeta {
	steps := make([]string, 0, len(m.StepTypes))
	for _, st := range m.StepTypes {
		steps = append(steps, string(st))
	}
	return wire.WorkflowActionMeta{
		ComponentName: strings.TrimSpace(m.Name),
		Label:         m.Label,
		Object:        m.Object,
		StepTypes:     steps,
		RunAs:         string(m.NormalizedRunAs()),
	}
}

// ContextFromWire converts a host invocation into the guest API shape.
func ContextFromWire(w wire.WorkflowActionContext) RecordWorkflowActionContext {
	participants := w.Participants
	if participants == nil {
		participants = map[string][]string{}
	}
	var pg *ParticipantGroup
	if name := strings.TrimSpace(w.ParticipantGroupName); name != "" {
		pg = &ParticipantGroup{
			Name:  name,
			Label: w.ParticipantGroupLabel,
			Users: append([]string(nil), participants[name]...),
		}
	}
	changes := make([]TaskChange, 0, len(w.TaskChanges))
	for _, c := range w.TaskChanges {
		changes = append(changes, TaskChange{
			Old: TaskInstance{
				TaskID:         c.OldTaskID,
				StepAPIName:    c.StepAPIName,
				AssigneeUserID: c.OldAssigneeUserID,
				Status:         c.OldStatus,
			},
			New: TaskInstance{
				TaskID:         firstNonEmpty(c.NewTaskID, c.OldTaskID),
				StepAPIName:    c.StepAPIName,
				AssigneeUserID: firstNonEmpty(c.NewAssigneeUserID, c.OldAssigneeUserID),
				Status:         firstNonEmpty(c.NewStatus, c.OldStatus),
			},
		})
	}
	return RecordWorkflowActionContext{
		Event:              Event(w.Event),
		VaultID:            w.VaultID,
		WorkflowInstanceID: w.WorkflowInstanceID,
		WorkflowAPIName:    w.WorkflowAPIName,
		WorkflowLabel:      w.WorkflowLabel,
		ObjectAPIName:      w.ObjectAPIName,
		RecordIDs:          append([]string(nil), w.RecordIDs...),
		InitiatingUserID:   w.InitiatingUserID,
		CurrentUserID:      w.CurrentUserID,
		ParticipantGroup:   pg,
		Participants:       participants,
		TaskStepAPIName:    w.TaskStepAPIName,
		TaskChanges:        changes,
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
