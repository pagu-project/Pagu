package log

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalInst *logger

type logger struct {
	subs   map[string]*SubLogger
	writer io.Writer
}

type SubLogger struct {
	logger zerolog.Logger
	name   string
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

func InitGlobalLogger() {
	if globalInst == nil {
		writers := []io.Writer{}
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05"})

		fw := &lumberjack.Logger{
			Filename: "RoboPac.log",
			MaxSize:  15,
		}
		writers = append(writers, fw)

		globalInst = &logger{
			subs:   make(map[string]*SubLogger),
			writer: io.MultiWriter(writers...),
		}
		log.Logger = zerolog.New(globalInst.writer).With().Timestamp().Logger()
	}
}

func addFields(event *zerolog.Event, keyvals ...interface{}) *zerolog.Event {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "!MISSING-VALUE!")
	}

	for i := 0; i < len(keyvals); i += 2 {
		key, ok := keyvals[i].(string)
		if !ok {
			key = "!INVALID-KEY!"
		}
		///
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

func (sl *SubLogger) logObj(event *zerolog.Event, msg string, keyvals ...interface{}) {
	addFields(event, keyvals...).Msg(msg)
}

func (sl *SubLogger) Trace(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Trace(), msg, keyvals...)
}

func (sl *SubLogger) Debug(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Debug(), msg, keyvals...)
}

func (sl *SubLogger) Info(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Info(), msg, keyvals...)
}

func (sl *SubLogger) Warn(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Warn(), msg, keyvals...)
}

func (sl *SubLogger) Error(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Error(), msg, keyvals...)
}

func (sl *SubLogger) Fatal(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Fatal(), msg, keyvals...)
}

func (sl *SubLogger) Panic(msg string, keyvals ...interface{}) {
	sl.logObj(sl.logger.Panic(), msg, keyvals...)
}

func Trace(msg string, keyvals ...interface{}) {
	addFields(log.Trace(), keyvals...).Msg(msg)
}

func Debug(msg string, keyvals ...interface{}) {
	addFields(log.Debug(), keyvals...).Msg(msg)
}

func Info(msg string, keyvals ...interface{}) {
	addFields(log.Info(), keyvals...).Msg(msg)
}

func Warn(msg string, keyvals ...interface{}) {
	addFields(log.Warn(), keyvals...).Msg(msg)
}

func Error(msg string, keyvals ...interface{}) {
	addFields(log.Error(), keyvals...).Msg(msg)
}

func Fatal(msg string, keyvals ...interface{}) {
	addFields(log.Fatal(), keyvals...).Msg(msg)
}

func Panic(msg string, keyvals ...interface{}) {
	addFields(log.Panic(), keyvals...).Msg(msg)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
