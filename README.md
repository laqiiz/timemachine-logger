# timemachinelog

timemachine-loggerはzerologのHookでログレベルを時系列にさかのぼって調整します。
通常時は指定されたログレベル（INFOなど）で出力しますが、ある一定のログレベル（ERROR）が出力された場合に、TRACE, INFOなどの低レベルのログも時をさかのぼって出力します。

必要な場合に遡って出力する、ドライブレコーダのような働きを目指しています。

## Usage

```go
log.SetPutput(timemachinelog.Logger {
    MinLevel:     "DEBUG",
    NormalLevel:  "INFO",
    TriggerLevel: "ERROR",
    ContextKey:   "contextID"
})
```

```go
zerolog.SetGlobalLevel(zerolog.InfoLevel)
defer timemachinelog.End("transactionID", "1234567890")

log.Debug().Str("contextID", "1234567890").Msg("debug1 log")
log.Info().Str("contextID", "1234567890").Msg("info log")
log.Debug().Str("contextID", "1234567890").Msg("debug2 log")
log.Warn().Str("contextID", "1234567890").Msg("warn log")
log.Debug().Str("contextID", "1234567890").Msg("debug3 log")
log.Error().Str("contextID", "1234567890").Msg("error log")
```

```log
=> Before Output
info log
warn log
error log


=> After Output
==== Timemachine ====
📛 debug1 log
📛 info log
📛 debug2 log
📛 warn log
📛 debug3 log
📛 error log
```

## Options

```go
log.SetPutput(timemachinelog.Logger {
    MinLevel: "DEBUG",              // 緊急時はDEBUGレベル以上を再生する
    NromalLevel: "INFO",            // 通常時のレベル
    TriggerLevel: "ERROR",          // 出力トリガー
    ContextKey: "requestID"         // ログコンテキストを一意に定めるキー
    MaxLogLength: 10 * 1024 * 1024, // 10MB 
    Expires: 10 * time.Minute,      // 10 minute
})
```

## Context Usage

ContextKeyは context.Contextを引き回すことで省略可能。

```go
ctx, tml := timemachinelog.Start(context.Background())
defer tml.End()

// 省略可能
log.Context(ctx).Msg("debug1 log")
```


## License

MIT
