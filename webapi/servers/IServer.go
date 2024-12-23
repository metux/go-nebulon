package servers

// FIXME: move that to base package ?
// FIXME: add CheckRunning() ?
type IServer interface {
	Shutdown(force bool) error
	Serve() error
}
