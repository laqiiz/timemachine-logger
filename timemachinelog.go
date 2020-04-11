package timemachinelog

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"io"
	"os"
	"sync"
)

const (
	defaultBuffer                 = 256
	defaultConcurrentLogUsageSize = 10
)

type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
)

var defaultLogger = zerolog.New(os.Stderr).With().Logger()

type Log struct {
	MinLevel     Level
	NormalLevel  Level
	TriggerLevel Level
	ContextKey   string
	Output       io.Writer

	zlog    zerolog.Logger
	mu      sync.Mutex
	history map[string][]LogHistory
}

type LogHistory struct {
	e   *Event
	msg string
}

type Event struct {
	*zerolog.Event
	TransactionID string
	Logger        *Log
	Level         Level
}

type LogTransaction struct {
	log           *Log
	transactionID string
}

func (l *Log) Start(ctx context.Context) (LogTransaction, context.Context) {
	transactionID := uuid.New().String()
	return LogTransaction{
		log:           l,
		transactionID: transactionID,
	}, context.WithValue(ctx, l.ContextKey, transactionID)
}

func (t *LogTransaction) End() {
	histories := t.log.history[t.transactionID]
	for _, v := range histories {
		// Trigger となるLogLevelのEventが無かったということなので、Normalレベルで絞る
		if v.e.Level < v.e.Logger.NormalLevel {
			continue
		}
		v.e.Event.Msg(v.msg)
	}
	t.log.history[t.transactionID] = histories[:0] // clear
}

func (l *Log) Trace() *Event {
	l.setup()
	e := l.zlog.Trace()
	return &Event{Level: TraceLevel, Event: e, Logger: l}
}

func (l *Log) Debug() *Event {
	l.setup()
	e := l.zlog.Debug()
	return &Event{Level: DebugLevel, Event: e, Logger: l}
}

func (l *Log) Info() *Event {
	l.setup()
	e := l.zlog.Info()
	return &Event{Level: InfoLevel, Event: e, Logger: l}
}

func (l *Log) Warn() *Event {
	l.setup()
	e := l.zlog.Warn()
	return &Event{Level: WarnLevel, Event: e, Logger: l}
}

func (l *Log) Error() *Event {
	l.setup()
	e := l.zlog.Error()
	return &Event{Level: ErrorLevel, Event: e, Logger: l}
}

func (l *Log) Fatal() *Event {
	l.setup()
	e := l.zlog.Fatal()
	return &Event{Level: FatalLevel, Event: e, Logger: l}
}

func (l *Log) Panic() *Event {
	l.setup()
	e := l.zlog.Panic()
	return &Event{Level: PanicLevel, Event: e, Logger: l}
}

func (l *Log) Close() {
	for _, items := range l.history {
		for _, v := range items {
			if v.e.Level >= l.NormalLevel {
				v.e.Event.Msg(v.msg)
			}
		}
	}
}

func (l *Log) setup() {
	if l.history == nil {
		l.history = make(map[string][]LogHistory, defaultConcurrentLogUsageSize)
	}

	if l.Output == nil {
		l.zlog = defaultLogger
	} else {
		l.zlog = zerolog.New(l.Output).With().Logger()
	}
}

func (e *Event) Str(key, val string) *Event {
	if e.Logger.ContextKey == key {
		e.TransactionID = val
	} else {
		// いったんtransactionIDの場合はログ出力させない
		e.Event = e.Event.Str(key, val)
	}
	return e
}

func (e *Event) Context(ctx context.Context) *Event {
	v := ctx.Value(e.Logger.ContextKey)
	if v != nil {
		e.TransactionID = fmt.Sprint(v)
	}
	return e
}

func (e *Event) Err(err error) *Event {
	e.Event = e.Event.Err(err)
	return e
}

func (e *Event) Msg(msg string) {
	if e.Level < e.Logger.MinLevel {
		return
	}

	if e.TransactionID == "" { // 未設定の場合はそのまま出力する仕様
		e.Event.Msg(msg)
		return
	}

	if _, ok := e.Logger.history[e.TransactionID]; !ok {
		e.Logger.history[e.TransactionID] = make([]LogHistory, 0, defaultBuffer)
	}

	// 呼び出されたときのTimestampを設定
	e.Timestamp()

	e.Logger.history[e.TransactionID] = append(e.Logger.history[e.TransactionID], LogHistory{
		e:   e,
		msg: msg,
	})

	if e.Level >= e.Logger.TriggerLevel {
		e.Logger.Info().Event.Msg("======== Timemachine ============")
		histories := e.Logger.history[e.TransactionID]
		for _, v := range histories {
			// この場合は全件出力
			v.e.Event.Msg(v.msg)
		}
		e.Logger.history[e.TransactionID] = histories[:0] // clear
	}
}

func (e *Event) Msgf(format string, v ...interface{}) {
	if e.Level < e.Logger.MinLevel {
		return
	}
	e.Msg(fmt.Sprintf(format, v...))
}

func (e *Event) Send() {
	if e.Level < e.Logger.MinLevel {
		return
	}
	e.Msg("")
}
