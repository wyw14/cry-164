package runtime

import (
	"github.com/wyw14/cry-164/internal/model"
	"os"
	"strconv"
)

func LoadConfig() model.Config {
	config := model.DefaultConfig()
	if value, err := strconv.Atoi(os.Getenv("AMMONIALOOP_PORT")); err == nil && value > 0 {
		config.Port = value
	}
	return config
}
