package timemachinelog

import (
	"testing"
)

func TestLogger(t *testing.T) {
	log := Log{
		NormalLevel:  InfoLevel,
		MinLevel:     DebugLevel,
		TriggerLevel: ErrorLevel,
		ContextKey:   "contextID",
	}

	log.Debug().Str("contextID", "1234567890").Msg("debug log-1")
	log.Info().Str("contextID", "1234567890").Msg("info log-1")
	log.Warn().Str("contextID", "1234567890").Msg("warn log-1")
}
