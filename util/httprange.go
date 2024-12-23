package util

import (
	"fmt"
	"net/http"
	"strings"
)

type HttpRange struct {
	Unit     string
	HasStart bool
	HasEnd   bool
	StartPos int64
	EndPos   int64
}

func ParseRangeHeader(header http.Header) ([]HttpRange, error) {
	h, ok := header["Range"]
	if !ok {
		return []HttpRange{}, nil
	}
	return ParseHttpRangeList(h)
}

func ParseHttpRangeList(values []string) ([]HttpRange, error) {
	ret := []HttpRange{}

	for _, str := range values {
		r, err := ParseHttpRange(str)
		if err != nil {
			return ret, err
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func ParseHttpRange(r string) (HttpRange, error) {
	res := HttpRange{}

	spl_main := strings.Split(r, "=")
	if len(spl_main) != 2 {
		return res, fmt.Errorf("syntax error: must have form: <unit>=<range>")
	}

	res.Unit = spl_main[0]

	spl_interval := strings.Split(spl_main[1], "-")
	if len(spl_interval) != 2 {
		return res, fmt.Errorf("syntax error: interval must have form [<startpos>]-[<endpos>]")
	}

	// parse startpos
	if spl_interval[0] != "" {
		res.HasStart = true
		if _, err := fmt.Sscan(spl_interval[0], &res.StartPos); err != nil {
			return res, fmt.Errorf("syntax error: startpos must be empty or int [%w]", err)
		}
	}

	// parse endpos
	if spl_interval[1] != "" {
		res.HasEnd = true
		if _, err := fmt.Sscan(spl_interval[1], &res.EndPos); err != nil {
			return res, fmt.Errorf("syntax error: endpos must be empty or int [%w]", err)
		}
	}

	return res, nil
}
