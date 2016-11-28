package splitter

// Universal interface to
type Splitter interface {
	Flush() error
	Close()
	HandleLine(line string) error
	Error() error
}
