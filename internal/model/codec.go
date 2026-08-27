package model

import "encoding/json"

func EncodeCycle(cycle Cycle) ([]byte, error) { return json.Marshal(cycle) }
func DecodeCycle(data []byte) (Cycle, error) {
	var cycle Cycle
	err := json.Unmarshal(data, &cycle)
	return cycle, err
}
