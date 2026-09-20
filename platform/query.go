package platform

import (
	"fmt"

	"github.com/vivarcus/vivarcus-sdk/wire"
)

// QueryFunc is wired by the reactor to host vivarcus:platform/query.execute.
var QueryFunc func(vql string) (encoded []byte, errCode int32)

// Query runs read-only VQL via the host (DAC-aware, row-capped).
func Query(vql string) (wire.QueryResult, error) {
	if QueryFunc == nil {
		return wire.QueryResult{}, fmt.Errorf("record_action_host_call_failed: query unavailable")
	}
	raw, code := QueryFunc(vql)
	if code != 0 {
		return wire.QueryResult{}, fmt.Errorf("record_action_host_call_failed: query status %d", code)
	}
	var out wire.QueryResult
	if len(raw) > 0 {
		if err := wire.Decode(raw, &out); err != nil {
			return wire.QueryResult{}, fmt.Errorf("record_action_host_call_failed: %w", err)
		}
	}
	return out, nil
}
