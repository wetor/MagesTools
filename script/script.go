package script

import (
	"MagesTools/script/format"
	"bufio"
	"io"
	"os"
)

const (
	LineBreak             = 0x00
	NameStart             = 0x01
	LineStart             = 0x02
	Present               = 0x03
	SetColor              = 0x04
	PresentUnknown05      = 0x05
	PresentResetAlignment = 0x08
	RubyBaseStart         = 0x09
	RubyTextStart         = 0x0A
	RubyTextEnd           = 0x0B
	SetFontSize           = 0x0C
	PrintInParallel       = 0x0E
	PrintInCenter         = 0x0F
	SetMarginTop          = 0x11
	SetMarginLeft         = 0x12
	GetHardcodedValue     = 0x13
	EvaluateExpression    = 0x15
	PresentUnknown18      = 0x18
	AutoForward           = 0x19
	AutoForward1A         = 0x1A
	RubyCenterPerChar     = 0x1E
	AltLineBreak          = 0x1F
	Terminator            = 0xFF
)

type Entry struct {
	Index  int `struct:"int32"`
	Offset int `struct:"int32"`
	Length int `struct:"-"`
}

type Strings interface {
	ReadStrings(readString func([]byte) string)
	GetStrings() []string
	WriteStrings(writeString func(string) []byte)
	SetStrings(strings []string)
	GetRaw() []byte
}

type Script struct {
	Strings       Strings
	Format        format.Format
	DecodeCharset map[uint16]string
	EncodeCharset map[string]uint16
}

// OpenScript
//  Description 打开脚本文件
//  Param filename string
//  Return *Script
//
func OpenScript(filename string, format format.Format) *Script {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	script := &Script{
		Format: format,
	}
	switch string(data[0:3]) {
	case "MES":
		script.Strings = LoadMes(data)
	}

	return script
}

// LoadCharset
//  Description 载入码表/字符集
//  Receiver s *Script
//  Param filename string 文件名
//  Param isTBL bool 是否为码表。否则为字符集，字符集从0x8000开始
//  Param skipExist bool 是否检查并跳过重复出现的字符，仅以第一次出现为准
//
func (s *Script) LoadCharset(filename string, isTBL, skipExist bool) {

	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	decodeCharset := make(map[uint16]string, 65535)
	encodeCharset := make(map[string]uint16, 65535)
	if isTBL {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 5 {
				k := BytesToUint16Big(HexToBytes(line[0:4]))
				v := line[5:]
				if skipExist {
					if _, has := encodeCharset[v]; !has {
						decodeCharset[k] = v
						encodeCharset[v] = k
					}
				} else {
					decodeCharset[k] = v
					encodeCharset[v] = k
				}

			}
		}
	} else {
		data, _ := io.ReadAll(f)
		runes := []rune(string(data))
		for i, char := range runes {
			k := uint16(0x8000 + i)
			v := string(char)
			if skipExist {
				if _, has := encodeCharset[v]; !has {
					decodeCharset[k] = v
					encodeCharset[v] = k
				}
			} else {
				decodeCharset[k] = v
				encodeCharset[v] = k
			}
		}
	}
	s.Format.SetCharset(decodeCharset, encodeCharset)

}

// Read
//  Description 导出文本，需要先执行script.LoadCharset载入码表
//  Receiver s *Script
//
func (s *Script) Read() {

	s.Strings.ReadStrings(s.Format.DecodeLine)
}

// SaveStrings
//  Description 保存文本，需要先执行script.Read
//  Receiver s *Script
//  Param filename string
//
func (s *Script) SaveStrings(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	strings := s.Strings.GetStrings()
	for _, str := range strings {
		f.WriteString(str + "\n")
	}
}

// LoadStrings
//  Description 载入文本并导入
//  Receiver s *Script
//  Param filename string
//
func (s *Script) LoadStrings(filename string) {
	f, _ := os.Open(filename)
	defer f.Close()

	strings := make([]string, 0, len(s.Strings.GetStrings()))
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		strings = append(strings, line)
	}
	s.Strings.SetStrings(strings)
}

// Write
//  Description 保存为脚本
//  Receiver s *Script
//  Param filename string
//
func (s *Script) Write(filename string) {
	s.Strings.WriteStrings(s.Format.EncodeLine)

	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(s.Strings.GetRaw())
}
