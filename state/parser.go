package state

import (
	"encoding/json"
	"fmt"
)

type StateParser interface {
	Parse() (stateV4, error)
}

func CreateParser(data []byte, fileName string) (StateParser, error) {
	return DefaultStateParser{data: data}, nil
}

type DefaultStateParser struct {
	data []byte
}

func (d DefaultStateParser) Parse() (stateV4, error) {
	state := stateV4{}
	err := json.Unmarshal(d.data, &state)
	if err != nil {
		return stateV4{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return state, nil
}
