package wire

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

// Length-prefixed JSON ABI (Phase 1 core-module). WIT is the IDL source of truth;
// this codec is the temporary wire until wit-bindgen codegen lands.

// Context is the Record Action context crossing the host/guest boundary.
type Context struct {
	Action             string            `json:"action,omitempty"`
	Records            []Record          `json:"records"`
	UserInputRecord    *Record           `json:"user_input_record,omitempty"`
	Configuration      map[string]string `json:"configuration,omitempty"`
	PreExecutionInputs map[string]string `json:"pre_execution_inputs,omitempty"`
	Usage              string            `json:"usage,omitempty"`
	UsageEntryPoint    string            `json:"usage_entry_point,omitempty"`
	UsageLabel         string            `json:"usage_label,omitempty"`
	VaultID            string            `json:"vault_id"`
	CurrentUserID      string            `json:"current_user_id"`
	InitiatingUserID   string            `json:"initiating_user_id"`
}

// Record is one object record with typed fields (no map[string]any).
type Record struct {
	ID     string  `json:"id"`
	Object string  `json:"object"`
	Fields []Field `json:"fields"`
}

// Field is a typed field value.
type Field struct {
	Name string  `json:"name"`
	Kind string  `json:"kind"` // string|int|bool|float|null|time
	S    string  `json:"s,omitempty"`
	I    int64   `json:"i,omitempty"`
	B    bool    `json:"b,omitempty"`
	F    float64 `json:"f,omitempty"`
	T    string  `json:"t,omitempty"` // RFC3339
}

// Meta is one Record Action's metadata inside a Describe payload.
type Meta struct {
	ComponentName       string   `json:"component_name,omitempty"`
	Label               string   `json:"label"`
	Object              string   `json:"object,omitempty"`
	ObjectAction        string   `json:"object_action,omitempty"`
	Usages              []string `json:"usages,omitempty"`
	Icon                string   `json:"icon,omitempty"`
	UserInputObject     string   `json:"user_input_object,omitempty"`
	UserInputObjectType string   `json:"user_input_object_type,omitempty"`
	RunAs               string   `json:"run_as,omitempty"`
}

// TriggerMeta is one Record Trigger's metadata inside a Describe payload.
type TriggerMeta struct {
	ComponentName string   `json:"component_name,omitempty"`
	Label         string   `json:"label"`
	Object        string   `json:"object,omitempty"`
	Events        []string `json:"events"`
	EventSegment  string   `json:"event_segment,omitempty"`
	Order         int      `json:"order,omitempty"`
	RunAs         string   `json:"run_as,omitempty"`
}

// Describe is the payload returned by __sdk_describe (one wasm, many entries).
type Describe struct {
	APIVersion string        `json:"api_version"`
	Actions    []Meta        `json:"actions,omitempty"`
	Triggers   []TriggerMeta `json:"triggers,omitempty"`
}

// TriggerContext is the Record Trigger context crossing the host/guest boundary.
type TriggerContext struct {
	Trigger          string  `json:"trigger,omitempty"`
	Event            string  `json:"event"`
	EventSegment     string  `json:"event_segment,omitempty"`
	Old              *Record `json:"old,omitempty"`
	New              *Record `json:"new,omitempty"`
	VaultID          string  `json:"vault_id"`
	CurrentUserID    string  `json:"current_user_id"`
	InitiatingUserID string  `json:"initiating_user_id"`
}

// TriggerResult is returned from guest execute_trigger.
type TriggerResult struct {
	Error           string     `json:"error,omitempty"`
	Mutations       []Mutation `json:"mutations,omitempty"`
	SetErrorSubtype string     `json:"set_error_subtype,omitempty"`
	SetErrorMessage string     `json:"set_error_message,omitempty"`
}

// ExecuteResult is returned from guest execute.
type ExecuteResult struct {
	Message   string     `json:"message,omitempty"`
	Error     string     `json:"error,omitempty"`
	Mutations []Mutation `json:"mutations,omitempty"`
}

// UIResult is returned from optional on_pre_execute / on_post_execute.
type UIResult struct {
	Title          string `json:"title,omitempty"`
	ConfirmMessage string `json:"confirm_message,omitempty"`
	Message        string `json:"message,omitempty"`
	Error          string `json:"error,omitempty"`
}

// Mutation is a staged field write from guest Execute/SetValue.
type Mutation struct {
	RecordID string `json:"record_id"`
	Field    string `json:"field"`
	Value    Field  `json:"value"`
}

// Encode returns u32-le length + JSON payload.
func Encode(v any) ([]byte, error) {
	body, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	out := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(out[:4], uint32(len(body)))
	copy(out[4:], body)
	return out, nil
}

// Decode strips optional length prefix and unmarshals JSON into dest.
func Decode(payload []byte, dest any) error {
	body := payload
	if len(payload) >= 4 {
		n := binary.LittleEndian.Uint32(payload[:4])
		if int(n) == len(payload)-4 {
			body = payload[4:]
		}
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return fmt.Errorf("wire decode: %w", err)
	}
	return nil
}

// FieldFromAny converts a host any value into a typed Field.
func FieldFromAny(name string, v any) Field {
	if v == nil {
		return Field{Name: name, Kind: "null"}
	}
	switch t := v.(type) {
	case string:
		return Field{Name: name, Kind: "string", S: t}
	case bool:
		return Field{Name: name, Kind: "bool", B: t}
	case int:
		return Field{Name: name, Kind: "int", I: int64(t)}
	case int32:
		return Field{Name: name, Kind: "int", I: int64(t)}
	case int64:
		return Field{Name: name, Kind: "int", I: t}
	case float32:
		return Field{Name: name, Kind: "float", F: float64(t)}
	case float64:
		return Field{Name: name, Kind: "float", F: t}
	case json.Number:
		if i, err := t.Int64(); err == nil {
			return Field{Name: name, Kind: "int", I: i}
		}
		if f, err := t.Float64(); err == nil {
			return Field{Name: name, Kind: "float", F: f}
		}
		return Field{Name: name, Kind: "string", S: t.String()}
	default:
		return Field{Name: name, Kind: "string", S: fmt.Sprint(t)}
	}
}

// Any returns the Go value for a typed Field.
func (f Field) Any() any {
	switch f.Kind {
	case "null", "":
		return nil
	case "string":
		return f.S
	case "int":
		return f.I
	case "bool":
		return f.B
	case "float":
		return f.F
	case "time":
		return f.T
	default:
		return f.S
	}
}
