package purge

import "time"

type AuditEntry struct {
	At      time.Time
	Applied int
	Failed  bool
}

func NewAudit(result Result) AuditEntry {
	return AuditEntry{At: time.Now().UTC(), Applied: len(result.Applied), Failed: result.Err != nil}
}
