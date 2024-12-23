package util

import (
	"fmt"
	"testing"
)

func testRange_start(t *testing.T, startPos int64) {
	str := fmt.Sprintf("bytes=%d-", startPos)
	t.Logf("parsing http range string: \"%s\"\n", str)

	r, err := ParseHttpRange(fmt.Sprintf("bytes=%d-", startPos))
	if err != nil {
		t.Fatalf("parse error %s\n", err)
	}
	if r.Unit != "bytes" {
		t.Fatalf("Unit is not \"bytes\"\n")
	}
	if !r.HasStart {
		t.Fatalf("HasStart not set\n")
	}
	if r.StartPos != startPos {
		t.Fatalf("StartPos != %d\n", startPos)
	}
	if r.HasEnd {
		t.Fatalf("HasEnd must not be set\n")
	}
}

func testRange_end(t *testing.T, endPos int64) {
	str := fmt.Sprintf("bytes=-%d", endPos)
	t.Logf("parsing http range string: \"%s\"\n", str)

	r, err := ParseHttpRange(fmt.Sprintf("bytes=-%d", endPos))
	if err != nil {
		t.Fatalf("parse error %s\n", err)
	}
	if r.Unit != "bytes" {
		t.Fatalf("Unit is not \"bytes\"\n")
	}
	if !r.HasEnd {
		t.Fatalf("HasEnd not set\n")
	}
	if r.EndPos != endPos {
		t.Fatalf("EndPos != %d\n", endPos)
	}
	if r.HasStart {
		t.Fatalf("HasStart must not be set\n")
	}
}

func Test_Range_1(t *testing.T) {
	testRange_start(t, 0)
}

func Test_Range_2(t *testing.T) {
	testRange_start(t, 3182)
}

func Test_Range_3(t *testing.T) {
	testRange_end(t, 0)
}

func Test_Range_4(t *testing.T) {
	testRange_end(t, 3182)
}
