package util

import (
	"net/http"
	"time"
)

var (
	gmtTimeLoc = time.FixedZone("GMT", 0)
)

func HttpTime(t time.Time) string {
	return t.In(gmtTimeLoc).Format(http.TimeFormat)
}
