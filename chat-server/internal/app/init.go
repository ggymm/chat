package app

import (
	"io"
	"log/slog"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"chat-server/config"
)

func Init() {
	config.Init()

	// 初始化日志
	log := slog.New(slog.NewTextHandler(io.MultiWriter(
		&lumberjack.Logger{
			Filename:   config.AppLog(),
			MaxAge:     30,  // days
			MaxSize:    256, // megabytes
			MaxBackups: 128, // files
		},
		io.MultiWriter(os.Stdout),
	), &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.TimeKey {
				if t, ok := attr.Value.Any().(time.Time); ok {
					return slog.Attr{
						Key:   attr.Key,
						Value: slog.StringValue(t.Format("2006-01-02 15:04:05")),
					}
				}
			}
			return attr
		},
	}))
	slog.SetDefault(log)
}
