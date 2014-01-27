// Copyright 2014 The Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chanio

import (
	"io"
	"testing"
)

var (
	_ io.Reader          = Reader(nil)
	_ io.WriteCloser     = WriteCloser(nil)
	_ io.ReadWriteCloser = ReadWriteCloser(nil)
)

func TestRead(t *testing.T) {
	str := "hello, world!"
	out := []byte(str)
	in := make([]byte, len(out))
	c := make(chan byte)
	go write(out, c)
	r := Reader(c)
	tmp := in
	for {
		n, err := r.Read(tmp)
		if err != nil {
			if err != io.EOF {
				t.Errorf("Error reading: %v", err)
			}
			break
		}
		tmp = tmp[n:]
		if len(tmp) == 0 {
			break
		}
	}
	if string(in) != str {
		t.Errorf("Expected \"%s\"; got \"%s\"", str, string(in))
	}
}

func write(p []byte, c chan byte) {
	for _, b := range p {
		c <- b
	}
}

func TestWrite(t *testing.T) {
	str := "hello, world!"
	out := []byte(str)
	in := make([]byte, len(out))
	c := make(chan byte)
	go writeWithWriter(out, c)

	for i := range in {
		in[i] = <-c
	}
}

func writeWithWriter(p []byte, c chan byte) {
	for n, err := WriteCloser(c).Write(p); len(p) > 0 && err == nil; n, err = WriteCloser(c).Write(p) {
		p = p[n:]
	}
}
