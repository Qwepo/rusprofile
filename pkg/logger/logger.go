package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

func NewLogger(lvl, filename string) zerolog.Logger {

	level := getZerologLevel(lvl)
	writer := zeroLoggWriter(filename)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(writer).
		Level(zerolog.Level(level)).
		With().
		Timestamp().
		Caller().
		Logger()

	return logger

}
func getZerologLevel(lvl string) int8 {
	switch lvl {
	case "debug":
		return 0
	case "info":
		return 1
	case "warn":
		return 2
	case "error":
		return 3
	case "fatal":
		return 4
	case "panic":
		return 5
	case "nolevel":
		return 6
	case "disable":
		return 6
	default:
		return 3
	}

}

func zeroLoggWriter(filename string) zerolog.LevelWriter {
	l := lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	writer := zerolog.MultiLevelWriter(os.Stdout, &l)
	return writer

}
