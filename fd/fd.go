package fd

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
)

const ()

// Fds is container of file descriptors
type FdList struct {
	List        []*Fd
	dir         string
	excludedfds []int
}

type Fd struct {
	Id   int
	Path string
}

// GetFdList
func GetFdList() *FdList {
	fd := new(FdList)
	fd.excludedfds = make([]int, 4)
	fd.excludedfds[0] = syscall.Stdin
	fd.excludedfds[1] = syscall.Stdout
	fd.excludedfds[2] = syscall.Stderr

	// OS dependent settings
	switch runtime.GOOS {
	case "linux":
		fd.dir = "/proc/self/fd"
	case "darwin":
		fd.dir = "/dev/fd"
	default:
		log.Fatalf("GOOS %s not supported.", runtime.GOOS)
	}
	fd.fetchFdList()
	return fd
}

func (fds *FdList) isExcludedFd(fdpath string) bool {
	for _, f := range fds.excludedfds {
		if fdpath == filepath.Join(fds.dir, fmt.Sprint(f)) {
			return true
		}
	}
	return false
}

func (fds *FdList) fetchFdList() {
	files, err := filepath.Glob(filepath.Join(fds.dir, "*"))
	if err != nil {
		panic(err)
	}
	for _, fdpath := range files {
		if !fds.isExcludedFd(fdpath) {
			fd := new(Fd)
			id, err := strconv.Atoi(filepath.Base(fdpath))
			if err != nil {
				panic(err)
			}
			fd.Id = id
			fd.Path = fdpath

			fds.List = append(fds.List, fd)
		}
	}
}
