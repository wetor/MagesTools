package script

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
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
	DecodeCharset map[uint16]string
	EncodeCharset map[string]uint16
}

// OpenScript
//  Description 打开脚本文件
//  Param filename string
//  Return *Script
//
func OpenScript(filename string) *Script {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	script := &Script{}
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
//
func (s *Script) LoadCharset(filename string, isTBL bool) {
	s.DecodeCharset = make(map[uint16]string, 65535)
	s.EncodeCharset = make(map[string]uint16, 65535)

	f, _ := os.Open(filename)
	defer f.Close()
	if isTBL {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if len(line) > 5 {
				k := BytesToUint16(HexToBytes(line[0:4]))
				v := line[5:]
				if _, has := s.EncodeCharset[v]; !has {
					s.DecodeCharset[k] = v
					s.EncodeCharset[v] = k
				}
			}
		}
	} else {
		data, _ := io.ReadAll(f)
		for i, char := range string(data) {
			k := uint16(0x8000 + i)
			v := string(char)
			if _, has := s.EncodeCharset[v]; !has {
				s.DecodeCharset[k] = v
				s.EncodeCharset[v] = k
			}
		}
	}

}

// Read
//  Description 导出文本，需要先执行script.LoadCharset载入码表
//  Receiver s *Script
//
func (s *Script) Read() {
	s.Strings.ReadStrings(s.readString)
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
func (s *Script) readString(data []byte) string {

	text := bytes.NewBuffer(nil)

	for i := 0; i < len(data); {
		switch data[i] {
		case LineBreak:
			text.WriteString(FormatByte(data[i]))
			i++
		case NameStart:
			text.WriteString(":[")
			i++
		case LineStart:
			text.WriteString("]:")
			i++
		case Present:
			text.WriteString(FormatByte(data[i]))
			i++
		case SetColor:
			text.WriteString(FormatBytes(data[i : i+4]))
			i += 4
		case PresentUnknown05:
			text.WriteString(FormatByte(data[i]))
			i++
		case PresentResetAlignment:
			text.WriteString(FormatByte(data[i]))
			i++
		case RubyBaseStart:
			text.WriteString(FormatByte(data[i]))
			i++
		case RubyTextStart:
			text.WriteString(FormatByte(data[i]))
			i++
		case RubyTextEnd:
			text.WriteString(FormatByte(data[i]))
			i++
		case SetFontSize:
			text.WriteString(FormatBytes(data[i : i+3]))
			i += 3
		case PrintInParallel:
			text.WriteString(FormatByte(data[i]))
			i++
		case PrintInCenter:
			text.WriteString(FormatByte(data[i]))
			i++
		case SetMarginTop:
			text.WriteString(FormatBytes(data[i : i+3]))
			i += 3
		case SetMarginLeft:
			text.WriteString(FormatBytes(data[i : i+3]))
			i += 3
		case GetHardcodedValue:
			text.WriteString(FormatBytes(data[i : i+3]))
			i += 3
		case EvaluateExpression:
			tmp := bytes.NewBuffer(nil)
			tmp.WriteByte(data[i])
			i++
			for !(data[i] == 0 && data[i+1] == 0) {
				switch data[i] & 0x60 {
				case 0: //1 byte
					tmp.WriteByte(data[i])
					i++
				case 0x20: //2 byte
					tmp.Write(data[i : i+2])
					i += 2
				case 0x40: //3 byte
					tmp.Write(data[i : i+3])
					i += 3
				case 0x60: // le int32 4 byte
					tmp.Write(data[i : i+4])
					i += 4
				}
			}
			tmp.WriteByte(data[i])
			text.WriteString(FormatBytes(tmp.Bytes()))
			i++
		case PresentUnknown18:
			text.WriteString(FormatByte(data[i]))
			i++
		case AutoForward:
			text.WriteString(FormatByte(data[i]))
			i++
		case AutoForward1A:
			text.WriteString(FormatByte(data[i]))
			i++
		case RubyCenterPerChar:
			text.WriteString(FormatByte(data[i]))
			i++
		case AltLineBreak:
			text.WriteString(FormatByte(data[i]))
			i++
		case Terminator:
			text.WriteString(FormatByte(data[i]))
			i++
		default:
			index := BytesToUint16(data[i : i+2])
			if char, has := s.DecodeCharset[index]; has {
				text.WriteString(char)
			} else {
				text.WriteString(FormatBytes(data[i : i+2]))
			}
			i += 2
		}

	}
	return text.String()
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
	s.Strings.WriteStrings(s.writeString)

	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(s.Strings.GetRaw())
}

func (s *Script) stringToBytes(str string) []byte {
	data := bytes.NewBuffer(nil)
	for _, char := range str {
		if index, has := s.EncodeCharset[string(char)]; has {
			data.Write(Uint16ToBytes(index))
		} else {
			panic("不存在的字符：" + string(char))
		}
	}
	return data.Bytes()
}
func (s *Script) writeString(str string) []byte {

	data := bytes.NewBuffer(nil)

	line := []rune(strings.TrimSpace(str))
	i := 0
	inBytes := false
	tempStr := ""
	for i < len(line) {
		switch line[i] {
		case ':':
			if i == 0 {
				i += 2
				for line[i] != ']' {
					tempStr += string(line[i])
					i++
				}
				i += 2
				data.WriteByte(NameStart)
				data.Write(s.stringToBytes(tempStr))
				data.WriteByte(LineStart)
				tempStr = ""
			}
		case '[':
			if len(tempStr) > 0 {
				data.Write(s.stringToBytes(tempStr))
				tempStr = ""
			}
			inBytes = true
			i += 3 //跳过[0x
		case ']':
			if inBytes {
				data.Write(HexToBytes(tempStr))
				inBytes = false
				tempStr = ""
			} else {
				panic("错误的 ] 结束符号，在：" + str)
			}
			i++
		default:
			tempStr += string(line[i])
			i++
		}
	}
	if len(tempStr) > 0 {
		data.Write(s.stringToBytes(tempStr))
	}
	return data.Bytes()
}
