package utils

import (
	"reflect"

	"github.com/rs/zerolog/log"
)

func GetType(v interface{}) string {
	if t := reflect.TypeOf(v); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func LogInterface(v any) {
	log.Info().Interface(GetType(v), v).Str("status", "logging.interface").Send()
}
