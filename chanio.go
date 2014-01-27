// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The chanio package allows for using the io.Reader,
// io.Writer, and io.Closer interfaces to interact with
// byte channels.
package chanio

import (
	"io"
)

// Reader implements the io.Reader interface.
type Reader <-chan byte

// Read reads bytes from r into p, blocking until at least
// 1 byte is available. After that, it reads as many bytes
// into p as it can without blocking, or until p is full.
//
// If r is closed, Read will return io.EOF. If there are any
// bytes available to be read from r, it will return those
// bytes and a nil error, and return io.EOF on the subsequent
// call.
//
// If len(p) == 0, Read returns 0 and a nil error immediately.
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

// WriteCloser implements the io.WriteCloser interface.
type WriteCloser chan<- byte

// Write writes bytes from p into w, blocking until all
// bytes have been written. Write will never return a
// non-nil error.
func (w WriteCloser) Write(p []byte) (int, error) {
	if len(p) < 1 {
		return 0, nil
	}
	for _, b := range p {
		w <- b
	}
	return len(p), nil
}

// Close closes w. It will never return a non-nil error.
func (w WriteCloser) Close() error {
	close(w)
	return nil
}

// ReadWriteCloser implements the io.ReadWriteCloser interface.
type ReadWriteCloser chan byte

// Read reads bytes from r into p, blocking until at least
// 1 byte is available. After that, it reads as many bytes
// into p as it can without blocking, or until p is full.
//
// If r is closed, Read will return io.EOF. If there are any
// bytes available to be read from r, it will return those
// bytes and a nil error, and return io.EOF on the subsequent
// call.
//
// If len(p) == 0, Read returns 0 and a nil error immediately.
func (r ReadWriteCloser) Read(p []byte) (int, error) {
	return Reader(chan byte(r)).Read(p)
}

// Write writes bytes from p into w, blocking until all
// bytes have been written. Write will never return a
// non-nil error.
func (r ReadWriteCloser) Write(p []byte) (int, error) {
	return WriteCloser(chan byte(r)).Write(p)
}

// Close closes w. It will never return a non-nil error.
func (r ReadWriteCloser) Close() error {
	return WriteCloser(chan byte(r)).Close()
}
