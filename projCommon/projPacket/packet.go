package projPacket

import (
	"bytes"
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"math"
	"strconv"
	"sync"
)

var PacketPool = newPool()

type Pool struct {
	Pool *sync.Pool
}

func (pp *Pool) Get() *Packet {
	return pp.Pool.Get().(*Packet)
}

func (pp *Pool) Put(p *Packet) {
	p.Reset()
	pp.Pool.Put(p)
}

type Packet struct {
	buf *bytes.Buffer
}

func (p *Packet) Reset() {
	p.buf.Reset()
}

func (p *Packet) Len() uint16 {
	return uint16(p.buf.Len())
}

func (p *Packet) Buffer() *bytes.Buffer {
	return p.buf
}

func (p *Packet) CopyBytes() []byte {
	return append([]byte{}, p.buf.Bytes()...)
}

func (p *Packet) Init(new *Packet) {
	p.Reset()
	_, err := p.WriteBytes(new.CopyBytes())
	if err != nil {
		return
	}
}

func (p *Packet) Next(n int) []byte {
	return p.buf.Next(n)
}

func (p *Packet) Bytes() (b []byte) {
	return p.buf.Bytes()
}

func (p *Packet) ReadString(delim byte) (s string, err error) {
	return p.buf.ReadString(delim)
}

func (p *Packet) ReadWChar() (s string, err error) {
	return p.ReadWCharByLen(int(p.Len() / 2))
}

func (p *Packet) ReadRune() (v rune, size int, err error) {
	return p.buf.ReadRune()
}

func (p *Packet) ReadWCharAndLen() (s string, err error) {
	l, _ := p.ReadUint16()
	s, err = p.ReadWCharByLen(int(l))
	return s, err
}

func (p *Packet) ReadWCharByLen(l int) (s string, err error) {
	for i := 0; i < l; i++ {
		var v uint16
		var err error
		if v, err = p.ReadUint16(); err != nil {
			return "", err
		}

		s += strconv.Itoa(int(v))
	}

	return s, nil
}

func (p *Packet) ReadBool() (v bool, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadUint8() (v uint8, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadInt8() (v int8, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadUint16() (v uint16, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadInt16() (v int16, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadUint32() (v uint32, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadInt32() (v int32, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadUint64() (v uint64, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadInt64() (v int64, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadFloat32() (v float32, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadFloat64() (v float64, err error) {
	err = binary.Read(p.buf, binary.LittleEndian, &v)
	return
}

func (p *Packet) ReadAny(v any) error {
	return binary.Read(p.buf, binary.LittleEndian, &v)
}

func (p *Packet) ReadReStruct(v any) error {
	buf := p.buf.Next(binary.Size(v))
	return restruct.Unpack(buf, binary.LittleEndian, &v)
}

func (p *Packet) WriteBytes(b []byte) (n int, err error) {
	return p.buf.Write(b)
}

func (p *Packet) WriteString(s string) (n int, err error) {
	return p.buf.WriteString(s)
}

func (p *Packet) WriteRune(r rune) (n int, err error) {
	return p.buf.WriteRune(r)
}

func (p *Packet) WriteWChar(s string) (n int, err error) {
	return p.WriteWCharByLen(s, len(s))
}

func (p *Packet) WriteWCharAndLen(s string) (n int, err error) {
	n = len(s)
	err = p.WriteUint16(uint16(len(s)))
	if err != nil {
		return -1, err
	}

	_, err = p.WriteWCharByLen(s, len(s))
	if err != nil {
		return -1, err
	}

	return len(s), nil
}

func (p *Packet) WriteWCharByLen(s string, l int) (n int, err error) {
	for i := 0; i < l; i++ {
		if err := p.WriteUint16(uint16(s[i])); err != nil {
			return -1, err
		}

		n++
	}

	return n, nil
}

func (p *Packet) WriteBool(v bool) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint8(v uint8) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint16(v uint16) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint32(v uint32) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint64(v uint64) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteFloat32(v float32) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, math.Float32bits(v))
}

func (p *Packet) WriteFloat64(v float64) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, math.Float64bits(v))
}

func (p *Packet) WriteAny(v any) (err error) {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteReStruct(v any) error {
	buf, err := restruct.Pack(binary.LittleEndian, v)
	err = p.WriteAny(buf)
	return err
}

func newPool() *Pool {
	return &Pool{Pool: &sync.Pool{New: func() any { return newPacket() }}}
}

func newPacket() *Packet {
	p := &bytes.Buffer{}
	p.Grow(1024 * 8)
	return &Packet{
		buf: p,
	}
}
