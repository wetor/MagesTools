package script

import (
	"bufio"
	"bytes"
	"fmt"
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
	ReadOffset()
	ReadStrings(readString func([]byte) *string)
	GetStrings() []*string
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
	defer f.Close()
	if err != nil {
		panic(err)
	}
	script := &Script{}

	magic := make([]byte, 3)
	f.Read(magic)
	f.Seek(0, io.SeekStart)
	switch string(magic) {
	case "MES":
		script.Strings = LoadMes(f)
	}
	script.Strings.ReadOffset()

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
			k := BytesToUint16(HexToBytes(line[0:4]))
			s.DecodeCharset[k] = line[5:]
			s.EncodeCharset[line[5:]] = k
		}
	} else {
		data, _ := io.ReadAll(f)
		for i, char := range string(data) {
			k := uint16(0x8000 + i)
			v := string(char)
			s.DecodeCharset[k] = v
			s.EncodeCharset[v] = k
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

// Save
//  Description 保存文本，需要先执行script.Read
//  Receiver s *Script
//  Param filename string
//
func (s *Script) Save(filename string) {
	f, _ := os.Create(filename)
	defer f.Close()

	strings := s.Strings.GetStrings()
	for _, str := range strings {
		f.WriteString(*str + "\n")
	}
}

func (s *Script) readString(data []byte) *string {
	inName := false
	haveName := false
	// inColor := false
	name := ""
	text := ""

	for i := 0; i < len(data); {
		switch data[i] {
		case LineBreak:
			//if inColor {
			//	text += "#>"
			//} else {
			//text += FormatByte(data[i])
			//}
			//inColor = false
			text += FormatByte(data[i])
			i++
		case NameStart:
			haveName = true
			inName = true
			//text += FormatByte(data[i])
			i++
		case LineStart:
			inName = false
			//text += FormatByte(data[i])
			i++
		case Present:
			text += FormatByte(data[i])
			i++
		case SetColor:
			//inColor = true
			text += FormatBytes(data[i : i+4])
			i += 4
		case PresentUnknown05:
			text += FormatByte(data[i])
			i++
		case PresentResetAlignment:
			text += FormatByte(data[i])
			i++
		case RubyBaseStart:
			text += FormatByte(data[i])
			i++
		case RubyTextStart:
			text += FormatByte(data[i])
			i++
		case RubyTextEnd:
			text += FormatByte(data[i])
			i++
		case SetFontSize:
			text += FormatBytes(data[i : i+3])
			i += 3
		case PrintInParallel:
			text += FormatByte(data[i])
			i++
		case PrintInCenter:
			text += FormatByte(data[i])
			i++
		case SetMarginTop:
			text += FormatBytes(data[i : i+3])
			i += 3
		case SetMarginLeft:
			text += FormatBytes(data[i : i+3])
			i += 3
		case GetHardcodedValue:
			text += FormatBytes(data[i : i+3])
			i += 3
		case EvaluateExpression:
			// 15 29 0A A4 B5 14 14 00 81 00 00 08 FF
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
			i++
			text += FormatBytes(tmp.Bytes())

		case PresentUnknown18:
			text += FormatByte(data[i])
			i++
		case AutoForward:
			text += FormatByte(data[i])
			i++
		case AutoForward1A:
			text += FormatByte(data[i])
			i++
		case RubyCenterPerChar:
			text += FormatByte(data[i])
			i++
		case AltLineBreak:
			text += FormatByte(data[i])
			i++
		case Terminator:
			inName = false
			text += FormatByte(data[i])
			i++
		default:
			index := BytesToUint16(data[i : i+2])
			if char, has := s.DecodeCharset[index]; has {
				if inName {
					name += char
				} else {
					text += char
				}
			} else {
				text += FormatBytes(data[i : i+2])
			}
			i += 2
		}

	}

	if haveName {
		text = fmt.Sprintf(":[%s]:%s", name, text)
	}
	return &text
}
