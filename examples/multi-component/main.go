// Module github.com/acme.corp.sdkdemo — multi-entry customer SDK layout (ADR-21).
//
// One local go.mod, one deploy module: Record Action types in actions/, Record Trigger
// types in triggers/, shared helpers in shared/ (each .go → Sdkcode). See README.md.
package main

import (
	_ "github.com/acme.corp.sdkdemo/actions"
	_ "github.com/acme.corp.sdkdemo/triggers"
)

func main() {}
