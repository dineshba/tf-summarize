package state

import "encoding/json"

type stateV4 struct {
	Version     stateVersionV4           `json:"version"`
	RootOutputs map[string]outputStateV4 `json:"outputs"`
	Resources   []resourceStateV4        `json:"resources"`
}

type outputStateV4 struct {
	ValueRaw     json.RawMessage `json:"value"`
	ValueTypeRaw json.RawMessage `json:"type"`
	Sensitive    bool            `json:"sensitive,omitempty"`
}

type resourceStateV4 struct {
	Module string `json:"module,omitempty"`
	Type   string `json:"type"`
	Name   string `json:"name"`
}

// stateVersionV4 is a weird special type we use to produce our hard-coded
// "version": 4 in the JSON serialization.
type stateVersionV4 struct{}

func (sv stateVersionV4) MarshalJSON() ([]byte, error) {
	return []byte{'4'}, nil
}

func (sv stateVersionV4) UnmarshalJSON([]byte) error {
	// Nothing to do: we already know we're version 4
	return nil
}
