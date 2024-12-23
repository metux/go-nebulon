package wire

// FIXME: move this to base package ?

type Header map[string]string

const (
	Header_ContentType = "Content-Type"
)

const (
	ContentType_MP4 = "video/mp4"
)

func (h Header) SetContentType(c string) {
	h[Header_ContentType] = c
}

func (h Header) ContentType() string {
	val, _ := h[Header_ContentType]
	return val
}
