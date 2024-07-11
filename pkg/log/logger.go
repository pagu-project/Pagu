package log

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/pagu-project/Pagu/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalInst *logger
	logLevel   zerolog.Level
)

type logger struct {
	subs   map[string]*SubLogger
	writer io.Writer
}

type SubLogger struct {
	logger zerolog.Logger
	name   string
}

func InitGlobalLogger(cfg *config.Logger) {
	if globalInst == nil {
		writers := []io.Writer{}

		if slices.Contains(cfg.Targets, "file") {
			// File writer.
			fw := &lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				Compress:   cfg.Compress,
			}
			writers = append(writers, fw)
		}

		if slices.Contains(cfg.Targets, "console") {
			// Console writer.
			writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})
		}

		globalInst = &logger{
			subs:   make(map[string]*SubLogger),
			writer: io.MultiWriter(writers...),
		}

		// Set the global log level from the configuration.
		level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
		if err != nil {
			level = zerolog.InfoLevel // Default to info level if parsing fails.
		}
		zerolog.SetGlobalLevel(level)

		log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()
	}
}

// NewLoggerLevel initializes the logger level.
func NewLoggerLevel(level zerolog.Level) {
	logLevel = level
}

func getLoggersInst() *logger {
	if globalInst == nil {
		globalInst = &logger{
			subs:   make(map[string]*SubLogger),
			writer: zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"},
		}
		log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()
	}

	return globalInst
}

// SetLoggerLevel sets logger level based on env.
func SetLoggerLevel(level string) {
	parsedLevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		parsedLevel = zerolog.InfoLevel // Default to info level if parsing fails
	}
	logLevel = parsedLevel
}

func GetCurrentLogLevel() zerolog.Level {
	return logLevel
}

func addFields(event *zerolog.Event, keyvals ...any) *zerolog.Event {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}

		value := keyvals[i+1]
		switch v := value.(type) {
		case fmt.Stringer:
			if isNil(v) {
				event.Any(key, v)
			} else {
				event.Stringer(key, v)
			}
		case error:
			event.AnErr(key, v)
		case []byte:
			event.Str(key, hex.EncodeToString(v))
		default:
			event.Any(key, v)
		}
	}

	return event
}

func NewSubLogger(name string) *SubLogger {
	inst := getLoggersInst()
	sl := &SubLogger{
		logger: zerolog.New(inst.writer).With().Timestamp().Logger(),
		name:   name,
	}

	inst.subs[name] = sl

	return sl
}

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...any) {
	addFields(event, keyvals...).Msg(msg)
}

func (sl *SubLogger) Trace(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Error(), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Fatal(), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...any) {
	sl.logObj(sl.logger.Panic(), msg, keyvals...)
}

func Trace(msg string, keyvals ...any) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...any) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...any) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...any) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...any) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...any) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...any) {
	addFields(log.Panic(), keyvals...).Msg(msg)
}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
