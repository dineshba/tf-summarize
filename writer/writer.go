package writer

import "io"

type Writer interface {
	Write(writer io.Writer) error
}
