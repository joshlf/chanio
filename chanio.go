// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chanio

import (
	"io"
)

type Reader <-chan byte

func (r Reader) Read(p []byte) (int, error) {
	if len(p) < 1 {
		return 0, nil
	}

	// Read at least one byte
	var ok bool
	p[0], ok = <-r
	if !ok {
		return 0, io.EOF
	}

	// Keep reading until the read would block
	n := 1
	for _ = range p[1:] {
		select {
		case p[n] = <-r:
			n++
		default:
			break
		}
	}
	return n, nil
}

type WriteCloser chan<- byte

func (w WriteCloser) Write(p []byte) (int, error) {
	if len(p) < 1 {
		return 0, nil
	}
	for _, b := range p {
		w <- b
	}
	return len(p), nil
}

func (w WriteCloser) Close() error {
	close(w)
	return nil
}

type ReadWriteCloser chan byte

func (r ReadWriteCloser) Read(p []byte) (int, error) {
	return Reader(chan byte(r)).Read(p)
}

func (r ReadWriteCloser) Write(p []byte) (int, error) {
	return WriteCloser(chan byte(r)).Write(p)
}

func (r ReadWriteCloser) Close() error {
	return WriteCloser(chan byte(r)).Close()
}
