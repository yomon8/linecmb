package readworker

import (
	"bufio"
	"io"

	"github.com/yomon8/linecmb/printer"
)

const (
	readBufferSize = 64 * 1024
	separator      = '\n'
)

// ReadWorker have own fd and read data from it
type ReadWorker struct {
	fd int
}

// NewReadWorker return ReadWorker object
func NewReadWorker(fd int) *ReadWorker {
	rw := new(ReadWorker)
	rw.fd = fd
	return rw
}

// Run read data from file descriptor and send it to printer object
func (rw *ReadWorker) Run(reader *bufio.Reader) {
	isLastBufferSeparated := false
	lastWriteBytes := 0
	for {
		buf := make([]byte, readBufferSize)
		n, err := reader.Read(buf)
		if err == io.EOF {
			if n == 0 {
				if !isLastBufferSeparated && lastWriteBytes != 0 {
					rw.sendSpool(1, []byte{separator})
				}
				break
			} else {
				rw.sendSpool(n, buf)
			}
		} else if err != nil {
			panic(err)
		}
		if n != 0 {
			isLastBufferSeparated = buf[n-1] == separator
			lastWriteBytes = n
			rw.sendSpool(n, buf)
		}
	}
}

func (rw *ReadWorker) sendSpool(pos int, buf []byte) {
	s := new(printer.Spool)
	s.Data = make([]byte, len(buf[:pos]))
	copy(s.Data, buf[:pos])
	s.Id = rw.fd
	printer.Get().Queue <- s
}
