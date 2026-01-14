package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Setup() {
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "15:04:05",
		NoColor:    false,
	}

	// Format Level (DEBUG, INFO, etc.)
	output.FormatLevel = func(i interface{}) string {
		var l string
		if ll, ok := i.(string); ok {
			switch ll {
			case "debug":
				l = "\033[36mDEBUG\033[0m"
			case "info":
				l = "\033[32mINFO \033[0m"
			case "warn":
				l = "\033[33mWARN \033[0m"
			case "error":
				l = "\033[31mERROR\033[0m"
			case "fatal":
				l = "\033[35mFATAL\033[0m"
			default:
				l = strings.ToUpper(ll)
			}
		} else {
			if i == nil {
				l = "???"
			} else {
				l = strings.ToUpper(fmt.Sprintf("%s", i))
			}
		}
		return l
	}

	// Format Message
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("\033[1m%s\033[0m", i)
	}

	// Format Field Name (the Key)
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\033[34m%s:\033[0m", i)
	}

	// Format Field Value (The Fix: Use Sprintf to handle structs/ints/etc)
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	log.Logger = zerolog.New(output).With().Timestamp().Caller().Logger()
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(key string, value interface{}, msg string) {
	log.Info().Interface(key, value).Msg(msg)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Error(err error, msg string) {
	log.Error().Err(err).Msg(msg)
}

func Fatal(err error, msg string) {
	log.Fatal().Err(err).Msg(msg)
}
