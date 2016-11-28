package dumper

/**
 * Created by andriy on 04/12/15.
 */

type Dumper interface {
	Flush() error
	Close()
	HandleLine(line string) error
	Error() error
}
