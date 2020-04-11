package main

import (
	timemachinelog "github.com/laqiiz/timemachine-logger"
	"time"
)

func main() {
	log := timemachinelog.Log{
		MinLevel:     timemachinelog.DebugLevel,
		NormalLevel:  timemachinelog.InfoLevel,
		TriggerLevel: timemachinelog.ErrorLevel,
		ContextKey:   "contextID",
	}
	defer log.Close()

	log.Trace().Str("contextID", "1234567890").Msg("🔧trace log-1")
	log.Trace().Str("contextID", "1234567890").Msg("🔧trace log-2")
	time.Sleep(1 * time.Second)

	log.Debug().Str("contextID", "1234567890").Msg("🟢debug log-1")
	time.Sleep(1 * time.Second)

	log.Debug().Str("contextID", "1234567890").Msg("🟢debug log-2")
	time.Sleep(1 * time.Second)

	log.Info().Str("contextID", "1234567890").Msg("🔵info log-1")
	log.Info().Str("contextID", "1234567890").Msg("🔵info log-2")
	time.Sleep(1 * time.Second)

	log.Warn().Str("contextID", "1234567890").Msg("🚧warn log-1")
	time.Sleep(1 * time.Second)

	log.Error().Str("contextID", "1234567890").Msg("🚨error log-1")
	time.Sleep(1 * time.Second)
}
