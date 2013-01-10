// Copyright (c) 2013 jnml. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fileutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)

func n(s []byte) byte {
	for i, c := range s {
		if c < '0' || c > '9' {
			s = s[:i]
			break
		}
	}
	v, _ := strconv.Atoi(string(s))
	return byte(v)
}

func init() {
	b, err := ioutil.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		panic(err)
	}

	tokens := bytes.Split(b, []byte("."))
	if n := len(tokens); n < 3 {
		panic(n)
	}

	// Supported since kernel 2.6.38
	if bytes.Compare([]byte{n(tokens[0]), n(tokens[1]), n(tokens[2])}, []byte{2, 6, 38}) < 0 {
		puncher = func(*os.File, int64, int64) error { return nil }
	}
}

var puncher = func(f *os.File, off, len int64) error {
	const (
		/*
			/usr/include/linux$ grep FL_ falloc.h
		*/
		_FALLOC_FL_KEEP_SIZE  = 0x01 // default is extend size
		_FALLOC_FL_PUNCH_HOLE = 0x02 // de-allocates range
	)

	_, _, errno := syscall.Syscall6(
		syscall.SYS_FALLOCATE,
		uintptr(f.Fd()),
		uintptr(_FALLOC_FL_KEEP_SIZE|_FALLOC_FL_PUNCH_HOLE),
		uintptr(off),
		uintptr(len),
		0, 0)
	return os.NewSyscallError("SYS_FALLOCATE", errno)
}

// PunchHole deallocates space inside a file in the byte range starting at
// offset and continuing for len bytes. No-op for kernels < 2.6.38.
func PunchHole(f *os.File, off, len int64) error {
	return puncher(f, off, len)
}

// Fadvise predeclares an access pattern for file data.  See also 'man 2
// posix_fadvise'.
func Fadvise(f *os.File, off, len int64, advice FadviseAdvice) error {
	_, _, errno := syscall.Syscall6(
		syscall.SYS_FADVISE64,
		uintptr(f.Fd()),
		uintptr(off),
		uintptr(len),
		uintptr(advice),
		0, 0)
	return os.NewSyscallError("SYS_FADVISE64", errno)
}
