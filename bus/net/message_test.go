package net_test

import (
	"bytes"
	"fmt"
	"github.com/lugu/qiloop/bus/net"
	"io"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

func helpParseHeader(t *testing.T, filename string, expected net.Header) {
	var m net.Message
	path := filepath.Join("testdata", filename)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = m.Read(file); err != nil {
		t.Error(err)
	}
	if m.Header != expected {
		t.Errorf("expected %#v, got %#v", expected, m.Header)
	}
}

func TestParseCallHeader(t *testing.T) {
	filename := "header-call-authenticate.bin"
	expected := net.Header{
		0x42dead42,
		3,   // id
		110, // size
		0,   // version
		1,   // type call
		0,   // flags
		0,   // service
		0,   // object
		8,   // authenticate action
	}
	helpParseHeader(t, filename, expected)
}

func TestParseReplyHeader(t *testing.T) {
	filename := "header-reply-authenticate.bin"
	expected := net.Header{
		0x42dead42,
		3,   // id
		138, // size
		0,   // version
		2,   // type call
		0,   // flags
		0,   // service
		0,   // object
		8,   // authenticate action
	}
	helpParseHeader(t, filename, expected)
}

func TestConstants(t *testing.T) {

	if uint32(unsafe.Sizeof(net.Header{})) != net.HeaderSize {
		t.Error("invalid header size")
	}
	if 0x42dead42 != net.Magic {
		t.Error("invalid magic definition")
	}
	if 0 != net.Version {
		t.Error("invalid version definition")
	}
	if 1 != net.Call {
		t.Error("invalid call definition")
	}
	if 2 != net.Reply {
		t.Error("invalid call definition")
	}
	if 3 != net.Error {
		t.Error("invalid error definition")
	}
}

func TestMessageConstructor(t *testing.T) {
	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	m := net.NewMessage(h, make([]byte, 99))

	if m.Header.Magic != net.Magic {
		t.Errorf("invalid magic: %d", m.Header.Magic)
	}
	if m.Header.Version != net.Version {
		t.Errorf("invalid version: %d", m.Header.Version)
	}
	if m.Header.Type != net.Call {
		t.Errorf("invalid type: %d", m.Header.Type)
	}
	if m.Header.Flags != 0 {
		t.Errorf("invalid flags: %d", m.Header.Flags)
	}
	if m.Header.Service != 1 {
		t.Errorf("invalid service: %d", m.Header.Service)
	}
	if m.Header.Object != 2 {
		t.Errorf("invalid service: %d", m.Header.Object)
	}
	if m.Header.Action != 3 {
		t.Errorf("invalid service: %d", m.Header.Action)
	}
	if m.Header.ID != 4 {
		t.Errorf("invalid id: %d", m.Header.ID)
	}
	if m.Header.Size != 99 {
		t.Errorf("invalid size: %d", m.Header.Size)
	}
}

func TestWriteReadMessage(t *testing.T) {
	h := net.NewHeader(net.Call, 1, 2, 3, 4)
	input := net.NewMessage(h, make([]byte, 99))
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := input.Write(buf); err != nil {
		t.Errorf("failed to write net. %s", err)
	}
	var output net.Message
	if err := output.Read(buf); err != nil {
		t.Errorf("failed to read net. %s", err)
	}
	if input.Header != output.Header {
		t.Errorf("expected %#v, got %#v", input.Header, output.Header)
	}
}

func NewBrokenReader(msg net.Message, size int) io.Reader {
	var buf bytes.Buffer
	msg.Write(&buf)
	data := buf.Bytes()
	if size > len(data) {
		panic(fmt.Errorf("size is too large: %d", size))
	}
	return bytes.NewBuffer(data[:size])
}

type BrokenWriter struct {
	size int
}

func (b *BrokenWriter) Write(buf []byte) (int, error) {
	if len(buf) <= b.size {
		b.size -= len(buf)
		return len(buf), nil
	}
	old_size := b.size
	b.size = 0
	return old_size, io.EOF
}

func NewBrokenWriter(size int) io.Writer {
	return &BrokenWriter{
		size: size,
	}
}

func TestWriterHeaderError(t *testing.T) {
	hdr := net.NewHeader(net.Call, 1, 1, 1, 1)
	for i := 1; i < int(net.HeaderSize); i++ {
		w := NewBrokenWriter(int(net.HeaderSize) - i)
		err := hdr.Write(w)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	w := NewBrokenWriter(int(net.HeaderSize))
	err := hdr.Write(w)
	if err != nil {
		panic(err)
	}
}

func TestReadHeaderError(t *testing.T) {
	hdr := net.NewHeader(net.Call, 1, 1, 1, 1)
	data := make([]byte, 0)
	max := int(net.HeaderSize)
	msg := net.NewMessage(hdr, data)
	for i := 0; i < max; i++ {
		r := NewBrokenReader(msg, i)
		err := hdr.Read(r)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	r := NewBrokenReader(msg, max)
	err := hdr.Read(r)
	if err != nil {
		panic(err)
	}
}
