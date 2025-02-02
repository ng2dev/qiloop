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
		Magic:   0x42dead42,
		ID:      3,   // id
		Size:    110, // size
		Version: 0,   // version
		Type:    1,   // type call
		Flags:   0,   // flags
		Service: 0,   // service
		Object:  0,   // object
		Action:  8,   // authenticate action
	}
	helpParseHeader(t, filename, expected)
}

func TestParseReplyHeader(t *testing.T) {
	filename := "header-reply-authenticate.bin"
	expected := net.Header{
		Magic:   0x42dead42,
		ID:      3,   // id
		Size:    138, // size
		Version: 0,   // version
		Type:    2,   // type reply
		Flags:   0,   // flags
		Service: 0,   // service
		Object:  0,   // object
		Action:  8,   // authenticate action
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
	var buf bytes.Buffer
	if err := input.Write(&buf); err != nil {
		t.Errorf("write net. %s", err)
	}
	var output net.Message
	if err := output.Read(&buf); err != nil {
		t.Errorf("read net. %s", err)
	}
	if input.Header != output.Header {
		t.Errorf("expected %#v, got %#v", input.Header, output.Header)
	}
}

func LimitedReader(msg net.Message, size int) io.Reader {
	var buf bytes.Buffer
	msg.Write(&buf)
	return &io.LimitedReader{
		R: &buf,
		N: int64(size),
	}
}

type LimitedWriter struct {
	size int
}

func (b *LimitedWriter) Write(buf []byte) (int, error) {
	if len(buf) <= b.size {
		b.size -= len(buf)
		return len(buf), nil
	}
	oldSize := b.size
	b.size = 0
	return oldSize, io.EOF
}

func NewLimitedWriter(size int) io.Writer {
	return &LimitedWriter{
		size: size,
	}
}

func TestWriterHeaderError(t *testing.T) {
	hdr := net.NewHeader(net.Call, 1, 1, 1, 1)
	for i := 1; i < int(net.HeaderSize); i++ {
		w := NewLimitedWriter(int(net.HeaderSize) - i)
		err := hdr.Write(w)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	w := NewLimitedWriter(int(net.HeaderSize))
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
		r := LimitedReader(msg, i)
		err := hdr.Read(r)
		if err == nil {
			panic(fmt.Errorf("not expecting a success at %d", i))
		}
	}
	r := LimitedReader(msg, max)
	err := hdr.Read(r)
	if err != nil {
		panic(err)
	}
}
