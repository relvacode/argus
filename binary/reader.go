package binary

import (
	"math"
	"unicode/utf16"
)

const (
	Uint8   = 1
	Uint16  = 2
	Uint32  = 4
	Uint64  = 8
	Float64 = Uint64
)

func New(buf []byte) *Reader {
	return &Reader{
		buf: buf,
	}
}

type Reader struct {
	buf []byte
	off int
}

func (r *Reader) Pos() int { return r.off }

func (r *Reader) Seek(to int) {
	r.off = to
}

func (r *Reader) Uint8() (v uint8) {
	v = r.buf[r.off]
	r.off += 1
	return
}

func (r *Reader) Uint16() (v uint16) {
	v = ByteOrder.Uint16(r.buf[r.off:])
	r.off += 2
	return
}

func (r *Reader) Uint32() (v uint32) {
	v = ByteOrder.Uint32(r.buf[r.off:])
	r.off += 4
	return
}

func (r *Reader) Uint64() (v uint64) {
	v = ByteOrder.Uint64(r.buf[r.off:])
	r.off += 8
	return
}

func (r *Reader) Float64() float64 {
	return math.Float64frombits(r.Uint64())
}

func (r *Reader) Utf16String() string {
	var chars []uint16
	for {
		chr := r.Uint16()
		if chr == 0 {
			break
		}

		chars = append(chars, chr)
	}

	return string(utf16.Decode(chars))
}
