package basic

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// ReadUint8 read an uint8
func ReadUint8(r io.Reader) (uint8, error) {
	buf := []byte{0}
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	} else if bytes != 1 {
		return 0, fmt.Errorf("failed to read uint8 (%d instead of 1)", bytes)
	}
	return uint8(buf[0]), nil
}

// WriteUint8 an uint8
func WriteUint8(i uint8, w io.Writer) error {
	buf := []byte{i}
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 1 {
		return fmt.Errorf("failed to write uint16 (%d instead of 1)", bytes)
	}
	return nil
}

// ReadInt8 reads a little endian int8
func ReadInt8(r io.Reader) (int8, error) {
	i, err := ReadUint8(r)
	return int8(i), err
}

// WriteInt8 writes a little endian int8
func WriteInt8(i int8, w io.Writer) error {
	return WriteUint8(uint8(i), w)
}

// ReadUint16 reads a little endian uint16
func ReadUint16(r io.Reader) (uint16, error) {
	buf := []byte{0, 0}
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	} else if bytes != 2 {
		return 0, fmt.Errorf("failed to read uint16 (%d instead of 2)", bytes)
	}
	return binary.LittleEndian.Uint16(buf), nil
}

// WriteUint16 writes a little endian uint16
func WriteUint16(i uint16, w io.Writer) error {
	buf := []byte{0, 0}
	binary.LittleEndian.PutUint16(buf, i)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 2 {
		return fmt.Errorf("failed to write uint16 (%d instead of 2)", bytes)
	}
	return nil
}

// ReadInt16 reads a little endian int16
func ReadInt16(r io.Reader) (int16, error) {
	i, err := ReadUint16(r)
	return int16(i), err
}

// WriteInt16 writes a little endian int16
func WriteInt16(i int16, w io.Writer) error {
	return WriteUint16(uint16(i), w)
}

// ReadUint32 reads a little endian uint32
func ReadUint32(r io.Reader) (uint32, error) {
	buf := []byte{0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	} else if bytes != 4 {
		return 0, fmt.Errorf("failed to read uint32 (%d instead of 4)", bytes)
	}
	return binary.LittleEndian.Uint32(buf), nil
}

// WriteUint32 writes a little endian uint32
func WriteUint32(i uint32, w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(buf, i)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 4 {
		return fmt.Errorf("failed to write uint32 (%d instead of 4)", bytes)
	}
	return nil
}

// ReadUint64 read a little endian uint64
func ReadUint64(r io.Reader) (uint64, error) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	} else if bytes != 8 {
		return 0, fmt.Errorf("failed to read uint32 (%d instead of 8)", bytes)
	}
	return binary.LittleEndian.Uint64(buf), nil
}

// WriteUint64 writes a little endian uint64
func WriteUint64(i uint64, w io.Writer) error {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(buf, i)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 8 {
		return fmt.Errorf("failed to write uint32 (%d instead of 8)", bytes)
	}
	return nil
}

// ReadInt32 reads a little endian int32
func ReadInt32(r io.Reader) (int32, error) {
	i, err := ReadUint32(r)
	return int32(i), err
}

// WriteInt32 writes a little endian int32
func WriteInt32(i int32, w io.Writer) error {
	return WriteUint32(uint32(i), w)
}

// ReadInt64 reads a little endian int64
func ReadInt64(r io.Reader) (int64, error) {
	i, err := ReadUint64(r)
	return int64(i), err
}

// WriteInt64 writes a little endian int64
func WriteInt64(i int64, w io.Writer) error {
	return WriteUint64(uint64(i), w)
}

// ReadFloat32 read a little endian float32
func ReadFloat32(r io.Reader) (float32, error) {
	buf := []byte{0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	} else if bytes != 4 {
		return 0, fmt.Errorf("failed to read float32 (%d instead of 4)", bytes)
	}
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}

// WriteFloat32 writes a little endian float32
func WriteFloat32(f float32, w io.Writer) error {
	buf := []byte{0, 0, 0, 0}
	bits := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf, bits)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 4 {
		return fmt.Errorf("failed to write float32 (%d instead of 4)", bytes)
	}
	return nil
}

// ReadFloat64 read a little endian float64
func ReadFloat64(r io.Reader) (float64, error) {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0, err
	} else if bytes != 8 {
		return 0, fmt.Errorf("failed to read float64 (%d instead of 4)", bytes)
	}
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}

// WriteFloat64 writes a little endian float64
func WriteFloat64(f float64, w io.Writer) error {
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	bits := math.Float64bits(f)
	binary.LittleEndian.PutUint64(buf, bits)
	bytes, err := w.Write(buf)
	if err != nil {
		return err
	} else if bytes != 8 {
		return fmt.Errorf("failed to write float64 (%d instead of 4)", bytes)
	}
	return nil
}

// ReadBool read a one byte size binary value
func ReadBool(r io.Reader) (bool, error) {
	u, err := ReadUint8(r)
	if u == 0 {
		return false, err
	}
	return true, err
}

// WriteBool writes a one byte size binary value
func WriteBool(b bool, w io.Writer) error {
	if b {
		return WriteUint8(1, w)
	}
	return WriteUint8(0, w)
}

// ReadString reads a string: first the size of the string is read
// using ReadUint32, then the bytes of the string.
func ReadString(r io.Reader) (string, error) {
	size, err := ReadUint32(r)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read string size: %s", err)
	}
	if size == 0 {
		return "", nil
	}
	// FIXME: do not allocate everything at one, read by block of
	// 4094 until either the reader fail or size is reached.
	buf := make([]byte, size)
	bytes, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	} else if uint32(bytes) != size {
		return "", fmt.Errorf("failed to read string: %d instead of %d",
			bytes, size)
	}
	return string(buf), nil
}

// WriteString writes a string: first the size of the string is
// written using WriteUint32, then the bytes of the string.
func WriteString(s string, w io.Writer) error {
	if err := WriteUint32(uint32(len(s)), w); err != nil {
		return fmt.Errorf("failed to write string size: %s", err)
	}
	bytes, err := w.Write([]byte(s))
	if err != nil {
		return err
	} else if bytes != len(s) {
		return fmt.Errorf("failed to write string data (%d instead of %d)", bytes, len(s))
	}
	return nil
}
