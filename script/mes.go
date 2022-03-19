package script

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
)

type Mes struct {
	Magic      string   `struct:"[4]byte"`
	Version    int      `struct:"int32"`
	Count      int      `struct:"int32"`
	TextOffset int      `struct:"int32"`
	Offsets    []Entry  `struct:"sizefrom=Count"`
	Strings    []string `struct:"-"`
	Raw        []byte   `struct:"-"`
}

func LoadMes(data []byte) *Mes {
	mes := &Mes{
		Raw: data,
	}
	return mes
}

func (m *Mes) ReadStrings(readString func([]byte) string) {
	// 读取offset

	err := restruct.Unpack(m.Raw, binary.LittleEndian, m)
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
		m.Offsets[i].Length = nextOffset - offset.Offset
		m.Offsets[i].Offset += m.TextOffset
	}

	// 读取文本
	m.Strings = make([]string, m.Count)
	for i, offset := range m.Offsets {
		m.Strings[i] = readString(m.Raw[offset.Offset : offset.Offset+offset.Length])
	}
}
func (m *Mes) GetStrings() []string {
	return m.Strings
}

func (m *Mes) SetStrings(strings []string) {
	if m.Count != len(strings) {
		panic(fmt.Sprintf("导入文本行数不匹配。原脚本：%d，导入：%d", m.Count, len(strings)))
		return
	}
	m.Strings = strings

}

func (m *Mes) WriteStrings(writeString func(string) []byte) {
	data := bytes.NewBuffer(nil)

	offset := 0
	for i, str := range m.Strings {
		line := writeString(str)

		m.Offsets[i].Offset = offset
		offset += len(line)

		data.Write(line)
	}

	offsetData, err := restruct.Pack(binary.LittleEndian, m)
	if err != nil {
		panic(err)
	}
	m.Raw = make([]byte, 0, len(offsetData)+data.Len())
	m.Raw = append(m.Raw, offsetData...)
	m.Raw = append(m.Raw, data.Bytes()...)
}
func (m *Mes) GetRaw() []byte {
	return m.Raw
}
