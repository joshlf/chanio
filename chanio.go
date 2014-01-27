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

type Writer chan<- byte

func (w Writer) Write(p []byte) (int, error) {
	if len(p) < 1 {
		return 0, nil
	}
	for _, b := range p {
		w <- b
	}
	return len(p), nil
}

func (w Writer) Close() error {
	close(w)
	return nil
}

type ReadWriter chan byte

func (r ReadWriter) Read(p []byte) (int, error) {
	return Reader(chan byte(r)).Read(p)
}

func (r ReadWriter) Write(p []byte) (int, error) {
	return Writer(chan byte(r)).Write(p)
}

func (r ReadWriter) Close() error {
	return Writer(chan byte(r)).Close()
}
