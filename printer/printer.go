package printer

import (
	"bufio"
	"bytes"
	"os"
	"sync"
)

const (
	queueSize       = 8 * 1024
	spoolBufferSize = 8 * 1024
	maxConcurrency  = 1024
	endOfLine       = '\n'
	writeBufferSize = 4096
)

var (
	singleton = newPrinter()
)

// Printer output []byte to stdout
type Printer struct {
	Queue  chan *Spool
	writer *bufio.Writer
	buffer map[int]*spoolData
	quit   bool
	wg     *sync.WaitGroup
}

func newPrinter() *Printer {
	p := Printer{
		Queue:  make(chan *Spool, queueSize),
		buffer: make(map[int]*spoolData, maxConcurrency),
		writer: bufio.NewWriterSize(os.Stdout, writeBufferSize),
		quit:   false,
		wg:     new(sync.WaitGroup),
	}
	for i := 0; i < maxConcurrency; i++ {
		p.buffer[i] = newSpoolData()
	}
	return &p
}

// Get Printer singleton object
func Get() *Printer {
	if singleton == nil {
		return newPrinter()
	}
	return singleton
}

// Close wait all print and close instance
func (p *Printer) Close() {
	close(p.Queue)
	p.quit = true
	p.wg.Wait()
}

func (p *Printer) write(data []byte) {
	_, err := p.writer.Write(data)
	if err != nil {
		panic(err)
	}
}

func (p *Printer) putBuffer(id int, data []byte) {
	pos := p.buffer[id].pos
	p.buffer[id].data[pos] = data
	p.buffer[id].pos++
}

func (p *Printer) deleteBuffer(id int) {
	// when EOF receive without EOL write out remaining data
	for j := 0; j < p.buffer[id].pos; j++ {
		p.write(p.buffer[id].data[j])
	}
	p.write([]byte{endOfLine})
	delete(p.buffer, id)
}

func (p *Printer) printLine(data []byte) []byte {
	if sep := bytes.IndexByte(data, endOfLine); sep == -1 {
		return data
	} else {
		p.write(data[:sep+1])
		if len(data) == 0 {
			return nil
		}
		return p.printLine(data[sep+1:])
	}
}

// Run print worker. it should be goroutine
func (p *Printer) Run() {
	p.wg.Add(1)
	defer func() {
		p.wg.Done()
	}()

	for {
		spo, more := <-p.Queue
		if spo != nil {
			if sep := bytes.IndexByte(spo.Data, endOfLine); sep == -1 {
				// line not found
				p.putBuffer(spo.Id, spo.Data)
				continue
			} else {
				// line found
				if p.buffer[spo.Id].pos > 0 {
					// print buffered data
					for j := 0; j < p.buffer[spo.Id].pos; j++ {
						p.write(p.buffer[spo.Id].data[j])
					}
				}
				p.buffer[spo.Id] = newSpoolData()
				remain := p.printLine(spo.Data)
				p.putBuffer(spo.Id, remain)
			}
			err := p.writer.Flush()
			if err != nil {
				panic(err)
			}
		}
		if !more {
			if p.quit {
				break
			}
			continue
		}
	}
}

type spoolData struct {
	data [][]byte
	pos  int
}

func newSpoolData() *spoolData {
	s := spoolData{
		data: make([][]byte, spoolBufferSize),
		pos:  0,
	}
	return &s
}

type Spool struct {
	Id   int
	Data []byte
}
