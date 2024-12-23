package wire

import (
	"fmt"
)

const (
	HttpV1_Header_BlockRef   = "BlockRef"
	HttpV1_Header_FetchDepth = "Fetch-Depth"
)

func HttpV1_Error(code int) error {
	return fmt.Errorf("http error code: %d", code)
}
