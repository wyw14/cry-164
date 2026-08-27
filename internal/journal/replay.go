package journal

import (
	"encoding/json"
	"github.com/wyw14/cry-164/internal/model"
	"os"
)

func Replay(path string) ([]model.Incident, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var entries []model.Incident
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
