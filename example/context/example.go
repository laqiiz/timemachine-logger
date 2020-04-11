package main

import (
	"context"
	timemachinelog "github.com/laqiiz/timemachine-logger"
	"time"
)

func main() {
	log := timemachinelog.Log{
		MinLevel:     timemachinelog.TraceLevel,
		NormalLevel:  timemachinelog.InfoLevel,
		TriggerLevel: timemachinelog.ErrorLevel,
		ContextKey:   "contextID",
	}

	start, ctx := log.Start(context.Background())
	defer start.End()

	log.Trace().Context(ctx).Msg("🔧trace log-1")
	log.Trace().Context(ctx).Msg("🔧trace log-2")
	time.Sleep(1 * time.Second)

	log.Debug().Context(ctx).Msg("🟢debug log-1")
	time.Sleep(1 * time.Second)

	log.Debug().Context(ctx).Msg("🟢debug log-2")
	time.Sleep(1 * time.Second)

	log.Info().Context(ctx).Msg("🔵info log-1")
	log.Info().Context(ctx).Msg("🔵info log-2")
	time.Sleep(1 * time.Second)

	log.Warn().Context(ctx).Msg("🚧warn log-1")
	time.Sleep(1 * time.Second)

	log.Error().Context(ctx).Msg("🚨error log-1")
	time.Sleep(1 * time.Second)

}
