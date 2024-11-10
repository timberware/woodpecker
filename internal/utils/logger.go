package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	Log = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatLevel: func(i interface{}) string {
			level, ok := i.(string)
			if ok && level == "update" {
				return "| UPDATE |"
			}
			if ok && level == "error" {
				return "| ERROR |"
			}
			return "| INFO |"
		},
		FormatMessage: func(i interface{}) string {
			return i.(string)
		},
	}).With().Timestamp().Logger()
}
