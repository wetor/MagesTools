package script

import (
	"MagesTools/script/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
)

type Sc3 struct {
	Magic       string   `struct:"[4]byte"`
	OffsetStart int      `struct:"int32"`
	OffsetEnd   int      `struct:"int32"`
	Count       int      `struct:"-"`
	Offsets     []Entry  `struct:"-"`
	Strings     []string `struct:"-"`
	Raw         []byte   `struct:"-"`
}

func LoadSc3(data []byte) *Sc3 {
	sc3 := &Sc3{
		Raw: data,
	}
	return sc3
}

func (s *Sc3) ReadStrings(readString func([]byte) string) {
	// 读取offset

	err := restruct.Unpack(s.Raw, binary.LittleEndian, s)
	if err != nil {
		panic(err)
	}
	s.Count = (s.OffsetEnd - s.OffsetStart) / 4
	s.Offsets = make([]Entry, s.Count)
	offset := s.OffsetStart
	for i := 0; i < s.Count; i++ {
		s.Offsets[i].Offset = int(utils.BytesToUint32(s.Raw[offset : offset+4]))
		offset += 4
		if i >= 1 {
			s.Offsets[i-1].Length = s.Offsets[i].Offset - s.Offsets[i-1].Offset
		}
	}
	s.Offsets[s.Count-1].Length = len(s.Raw) - s.Offsets[s.Count-1].Offset

	// 读取文本
	s.Strings = make([]string, s.Count)
	for i, offset := range s.Offsets {
		s.Strings[i] = readString(s.Raw[offset.Offset : offset.Offset+offset.Length])
	}

}
func (s *Sc3) GetStrings() []string {
	return s.Strings
}

func (s *Sc3) SetStrings(strings []string) {
	if s.Count != len(strings) {
		panic(fmt.Sprintf("导入文本行数不匹配。原脚本：%d，导入：%d", s.Count, len(strings)))
		return
	}
	s.Strings = strings

}

func (s *Sc3) WriteStrings(writeString func(string) []byte) {
	if s.Count == 0 {
		return
	}
	data := bytes.NewBuffer(nil) // 新文本段数据

	offset := s.Offsets[0].Offset // 文本段开始
	for i, str := range s.Strings {
		line := writeString(str)
		s.Offsets[i].Offset = offset
		offset += len(line)

		data.Write(line)
	}
	offset = s.Offsets[0].Offset           // 文本段开始
	s.Raw = s.Raw[:offset]                 // 清除原文本段
	s.Raw = append(s.Raw, data.Bytes()...) // 写入新文本段

	offset = s.OffsetStart         // 偏移段开始
	for i := 0; i < s.Count; i++ { // 写入新偏移段
		copy(s.Raw[offset:offset+4], utils.Uint32ToBytes(uint32(s.Offsets[i].Offset)))
		offset += 4
	}
}
func (s *Sc3) GetRaw() []byte {
	return s.Raw
}
