package script

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"io"
)

type Mes struct {
	Magic      string    `struct:"[4]byte"`
	Version    int       `struct:"int32"`
	Count      int       `struct:"int32"`
	TextOffset int       `struct:"int32"`
	Offsets    []*Entry  `struct:"sizefrom=Count"`
	Strings    []*string `struct:"-"`
	Reader     io.Reader `struct:"-"`
	Raw        []byte    `struct:"-"`
}

func LoadMes(r io.Reader) *Mes {
	mes := &Mes{
		Reader: r,
	}
	return mes
}

func (m *Mes) ReadOffset() {
	data, err := io.ReadAll(m.Reader)
	m.Raw = data
	if err != nil {
		panic(err)
	}
	err = restruct.Unpack(m.Raw, binary.LittleEndian, m)
	if err != nil {
		panic(err)
	}

	nextOffset := 0
	for i, offset := range m.Offsets {
		if i+1 < m.Count {
			nextOffset = m.Offsets[i+1].Offset
		} else {
			nextOffset = len(m.Raw) - m.TextOffset
		}
		offset.Length = nextOffset - offset.Offset
		offset.Offset += m.TextOffset
	}
}

func (m *Mes) ReadStrings(readString func([]byte) *string) {
	m.Strings = make([]*string, m.Count)
	for i, offset := range m.Offsets {
		m.Strings[i] = readString(m.Raw[offset.Offset : offset.Offset+offset.Length])
	}
}

func (m *Mes) GetStrings() []*string {
	return m.Strings
}
