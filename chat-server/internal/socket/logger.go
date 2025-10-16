package socket

import (
	"chat-server/config"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/ggymm/gopkg/rolling"
)

type SocketLogger struct {
	log *slog.Logger
}

func (l *SocketLogger) Debugf(format string, args ...any) {
	l.log.Debug(format, args...)
}
func (l *SocketLogger) Infof(format string, args ...any) {
	l.log.Info(format, args...)
}
func (l *SocketLogger) Warnf(format string, args ...any) {
	l.log.Warn(format, args...)
}
func (l *SocketLogger) Errorf(format string, args ...any) {
	l.log.Error(format, args...)
}
func (l *SocketLogger) Fatalf(format string, args ...any) {
	l.log.Error(format, args...)
}

func newLog() *slog.Logger {
	writer := io.MultiWriter(
		&rolling.Logger{
			Filename:   config.AppLog("socket"),
			MaxAge:     30,  // days
			MaxSize:    256, // megabytes
			MaxBackups: 128, // files
		},
		io.MultiWriter(os.Stdout),
	)
	opt := &slog.HandlerOptions{
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
	}
	return slog.New(slog.NewTextHandler(writer, opt))
}
