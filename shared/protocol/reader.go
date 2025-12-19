package protocol

import (
	"encoding/binary"
	"math"
	"strings"
)

type Reader struct { // Little endian
	bytes  []byte
	offset int
}

func NewReader(bytes []byte) *Reader {
	return &Reader{bytes, 0}
}

func (r *Reader) GetI8() int8 {
	r.offset++
	return int8(r.bytes[r.offset-1])
}

func (r *Reader) GetI16() int16 {
	r.offset += 2
	return int16(r.bytes[r.offset-2])<<8 | int16(r.bytes[r.offset-1])
}

func (r *Reader) GetI32() int32 {
	r.offset += 4
	return int32(r.bytes[r.offset-4])<<24 | int32(r.bytes[r.offset-3])<<16 | int32(r.bytes[r.offset-2])<<8 | int32(r.bytes[r.offset-1])
}

func (r *Reader) GetI64() int64 {
	r.offset += 8
	return int64(r.bytes[r.offset-8])<<56 | int64(r.bytes[r.offset-7])<<48 | int64(r.bytes[r.offset-6])<<40 | int64(r.bytes[r.offset-5])<<32 | int64(r.bytes[r.offset-4])<<24 | int64(r.bytes[r.offset-3])<<16 | int64(r.bytes[r.offset-2])<<8 | int64(r.bytes[r.offset-1])
}

func (r *Reader) GetU8() uint8 {
	r.offset++
	return r.bytes[r.offset-1]
}

func (r *Reader) GetU16() uint16 {
	r.offset += 2
	return binary.LittleEndian.Uint16(r.bytes[r.offset-2 : r.offset])
}

func (r *Reader) GetU32() uint32 {
	r.offset += 4
	return binary.LittleEndian.Uint32(r.bytes[r.offset-4 : r.offset])
}

func (r *Reader) GetU64() uint64 {
	r.offset += 8
	return binary.LittleEndian.Uint64(r.bytes[r.offset-8 : r.offset])
}

func (r *Reader) GetF32() float32 {
	r.offset += 4
	return math.Float32frombits(binary.LittleEndian.Uint32(r.bytes[r.offset-4 : r.offset]))
}

func (r *Reader) GetF64() float64 {
	r.offset += 8
	return math.Float64frombits(binary.LittleEndian.Uint64(r.bytes[r.offset-8 : r.offset]))
}

func (r *Reader) GetStringUTF8() string {
	var str strings.Builder

	for {
		if r.bytes[r.offset] == 0 {
			r.offset++
			break
		}

		str.WriteString(string(r.bytes[r.offset]))
		r.offset++
	}

	return str.String()
}
