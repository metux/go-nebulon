package httpd

import (
	"errors"
	"fmt"
	"log"
	"net/textproto"
	"strconv"
	"strings"
)

// errNoOverlap is returned by serveContent's parseRange if first-byte-pos of
// all of the byte-range-spec values is greater than the content size.
var errNoOverlap = errors.New("invalid range: failed to overlap")

var errInvalidRange = errors.New("invalid range")

// httpRange specifies the byte range to be sent to the client.
type httpRange struct {
	start, length int64
}

func (r httpRange) contentRange(size int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", r.start, r.start+r.length-1, size)
}

//func (r httpRange) mimeHeader(contentType string, size int64) textproto.MIMEHeader {
//	return textproto.MIMEHeader{
//		"Content-Range": {r.contentRange(size)},
//		"Content-Type":  {contentType},
//	}
//}

// parseRange parses a Range header string as per RFC 7233.
// errNoOverlap is returned if none of the ranges overlap.
func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil // header not present
	}
	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errInvalidRange
	}
	var ranges []httpRange
	noOverlap := false
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = textproto.TrimString(ra)
		log.Printf("trying ra=\"%s\"\n", ra)
		if ra == "" {
			continue
		}
		start, end, ok := strings.Cut(ra, "-")
		if !ok {
			log.Printf("missing '-'\n")
			return nil, errInvalidRange
		}
		start, end = textproto.TrimString(start), textproto.TrimString(end)
		var r httpRange
		if start == "" {
			log.Printf("from the end of the stream ?\n")
			// If no start is specified, end specifies the
			// range start relative to the end of the file,
			// and we are dealing with <suffix-length>
			// which has to be a non-negative integer as per
			// RFC 7233 Section 2.1 "Byte-Ranges".
			if end == "" || end[0] == '-' {
				return nil, errInvalidRange
			}
			i, err := strconv.ParseInt(end, 10, 64)
			if i < 0 || err != nil {
				return nil, errInvalidRange
			}
			if i > size {
				i = size
			}
			r.start = size - i
			r.length = size - r.start
		} else {
			log.Printf("normal range start %s\n", start)
			i, err := strconv.ParseInt(start, 10, 64)
			if err != nil || i < 0 {
				log.Printf("invalid int range %s\n", start)
				return nil, errInvalidRange
			}
			if i >= size {
				// If the range begins after the size of the content,
				// then it does not overlap.
				log.Printf("i >= size  -- i=%d size=%d\n", i, size)
				noOverlap = true
				continue
			}
			r.start = i
			if end == "" {
				// If no end is specified, range extends to end of the file.
				r.length = size - r.start
			} else {
				i, err := strconv.ParseInt(end, 10, 64)
				if err != nil || r.start > i {
					log.Printf("failed parsing chunk end: %s\n", start)
					return nil, errInvalidRange
				}
				if i >= size {
					i = size - 1
				}
				r.length = i - r.start + 1
			}
		}
		log.Printf("appending range part %+v\n", r)
		ranges = append(ranges, r)
	}

	//	log.Printf("noOverlap: %s\n", noOverlap)
	log.Printf("ranges=%+v\n", ranges)

	if noOverlap && len(ranges) == 0 {
		// The specified ranges did not overlap with the content.
		return nil, errNoOverlap
	}
	return ranges, nil
}
