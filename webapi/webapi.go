// Package webapi is the Vivarcus Custom Web API guest interface (tinygo / wasip1 safe).
package webapi

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/vivarcus/vivarcus-sdk/action"
	"github.com/vivarcus/vivarcus-sdk/wire"
)

const (
	StatusSuccess = "SUCCESS"
	StatusFailure = "FAILURE"
	StatusWarning = "WARNING"
)

// Meta declares component-level static metadata for one Custom Web API.
type Meta struct {
	Name           string
	Label          string
	EndpointName   string
	MinimumVersion string
	APIGroup       string
	RunAs          action.RunAs
}

// Validate checks required Meta fields.
func (m Meta) Validate() error {
	if strings.TrimSpace(m.Label) == "" {
		return fmt.Errorf("meta label required")
	}
	ep := strings.TrimSpace(m.EndpointName)
	if ep == "" {
		return fmt.Errorf("meta endpoint_name required")
	}
	if !validEndpointName(ep) {
		return fmt.Errorf("meta endpoint_name %q must be snake_case", ep)
	}
	if strings.TrimSpace(m.MinimumVersion) == "" {
		return fmt.Errorf("meta minimum_version required")
	}
	if strings.TrimSpace(m.APIGroup) == "" {
		return fmt.Errorf("meta api_group required")
	}
	return nil
}

func validEndpointName(s string) bool {
	if s == "" || s[0] < 'a' || s[0] > 'z' {
		return false
	}
	for _, r := range s {
		if r == '_' || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			continue
		}
		if unicode.IsUpper(r) {
			return false
		}
		return false
	}
	return true
}

// NormalizedRunAs returns RequestOwner when empty (Custom Web API default).
func (m Meta) NormalizedRunAs() action.RunAs {
	switch m.RunAs {
	case action.RunAsSystemUser:
		return action.RunAsSystemUser
	default:
		return action.RunAsRequestOwner
	}
}

// Context is the runtime context for WebApi.Execute.
type Context struct {
	APIVersion       string
	Endpoint         string
	JSON             map[string]any
	VaultID          string
	CurrentUserID    string
	InitiatingUserID string
}

// Error is one FAILURE envelope error.
type Error struct {
	Type    string
	Message string
}

// Response is returned from Execute.
type Response struct {
	Status string
	Data   map[string]any
	Errors []Error
}

// WebApi is the required guest interface (Java WebApi.execute).
type WebApi interface {
	Meta() Meta
	Execute(ctx Context) (Response, error)
}

// ToWireMeta converts Meta to the describe payload.
func ToWireMeta(m Meta) wire.WebApiMeta {
	return wire.WebApiMeta{
		ComponentName:  strings.TrimSpace(m.Name),
		Label:          m.Label,
		EndpointName:   strings.TrimSpace(m.EndpointName),
		MinimumVersion: strings.TrimSpace(m.MinimumVersion),
		APIGroup:       strings.TrimSpace(m.APIGroup),
		RunAs:          string(m.NormalizedRunAs()),
	}
}

// ContextFromWire converts a host invocation into the guest API shape.
func ContextFromWire(w wire.WebApiContext) Context {
	return Context{
		APIVersion:       w.APIVersion,
		Endpoint:         w.Endpoint,
		JSON:             w.JSON,
		VaultID:          w.VaultID,
		CurrentUserID:    w.CurrentUserID,
		InitiatingUserID: w.InitiatingUserID,
	}
}
