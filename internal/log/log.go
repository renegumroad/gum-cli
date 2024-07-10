package log

import (
	"os"

	"github.com/fatih/color"
	"github.com/pkg/errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	writer      zerolog.ConsoleWriter
	initialized = false

	blue   = color.New(color.FgBlue, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
	green  = color.New(color.FgGreen, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
)

func IsInitialized() bool {
	return initialized
}

func Initialize(level string) error {
	writer = zerolog.ConsoleWriter{Out: os.Stdout}

	writer.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = prefix(blue, ll)
			case "info":
				l = prefix(green, ll)
			case "warn":
				l = prefix(yellow, ll)
			case "error":
				l = prefix(red, ll)
			case "fatal":
				l = prefix(red, ll)
			}
		}

		return l
	}

	log.Logger = zerolog.New(writer).With().Timestamp().Logger()

	if err := SetLogLevel(level); err != nil {
		return err
	}

	initialized = true
	Debugln("Logging initialized")
	return nil
}

func prefix(c *color.Color, msg string) string {
	return c.SprintfFunc()(msg + ":")
}

func SetLogLevel(level string) error {
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	default:
		return errors.Errorf("Unexpected log level: %s", level)
	}

	return nil
}

func Debugln(a ...interface{}) {
	log.Debug().Msgf("%s", a...)
}

func Debugf(format string, a ...interface{}) {
	log.Debug().Msgf(format, a...)
}

func Infoln(a ...interface{}) {
	log.Info().Msgf("%s", a...)
}

func Infof(format string, a ...interface{}) {
	log.Info().Msgf(format, a...)
}

func Warnln(a ...interface{}) {
	log.Warn().Msgf("%s", a...)
}

func Warnf(format string, a ...interface{}) {
	log.Warn().Msgf(format, a...)
}

func Errorln(a ...interface{}) {
	log.Error().Msgf("%s", a...)
}

func Errorf(format string, a ...interface{}) {
	log.Error().Msgf(format, a...)
}

func Fatalln(a ...interface{}) {
	log.Fatal().Msgf("%s", a...)
}

func Fatalf(format string, a ...interface{}) {
	log.Fatal().Msgf(format, a...)
}
