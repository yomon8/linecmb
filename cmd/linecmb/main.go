package main

import (
	"bufio"
	"os"
	"sync"

	//------------------------------
	//---- Memory Profiler
	//------------------------------
	//"log"
	//"net/http"
	//"net/http/pprof"

	"github.com/yomon8/linecmb/fd"
	"github.com/yomon8/linecmb/printer"
	"github.com/yomon8/linecmb/readworker"
)

func main() {
	//------------------------------
	//---- Memory Profiler
	//------------------------------
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	wg := new(sync.WaitGroup)
	go printer.Get().Run()
	fds := fd.GetFdList()
	for _, f := range fds.List {
		wg.Add(1)
		go func(fd *fd.Fd) {
			defer func() {
				wg.Done()
			}()
			file, err := os.Open(fd.Path)
			if err != nil || file == nil {
				return
			}
			rw := readworker.NewReadWorker(fd.Id)
			rw.Run(bufio.NewReader(file))
		}(f)
	}
	wg.Wait()
	printer.Get().Close()
}
