package common

import (
	"log"
)

type StoreBase struct {
	Name   string
	Logger *log.Logger
}

func (sb StoreBase) Logf(format string, v ...any) {
	format = "[" + sb.Name + "] " + format

	if sb.Logger == nil {
		log.Printf(format, v...)
	} else {
		sb.Logger.Printf(format, v...)
	}
}
