package journal

import (
	"github.com/wyw14/cry-164/internal/model"
	"time"
)

func IncidentFor(message, severity string) model.Incident {
	return model.Incident{ID: model.NewOperation("incident").ID, Severity: severity, Message: message, At: time.Now().UTC()}
}
