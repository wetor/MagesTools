package format

import (
	"bytes"
	"fmt"
	"strings"

	"MagesTools/script/utils"
)

type NpcsP struct {
	decodeCharset map[uint16]string
	encodeCharset map[string]uint16
}

func (f *NpcsP) SetCharset(decode map[uint16]string, encode map[string]uint16) {
	f.decodeCharset = decode
	f.encodeCharset = encode
}
func (f *NpcsP) stringToBytes(str string) []byte {
	data := bytes.NewBuffer(nil)
	for _, char := range str {
		if index, has := f.encodeCharset[string(char)]; has {
			data.Write(utils.Uint16ToBytesBig(index))
		} else {
			panic("不存在的字符：" + string(char))
		}
	}
	return data.Bytes()
}

func (f *NpcsP) DecodeLine(data []byte) string {

	text := bytes.NewBuffer(nil)

	for i := 0; i < len(data); {
		switch data[i] {
		case LineBreak:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case NameStart:
			text.WriteString(":[")
			i++
		case LineStart:
			text.WriteString("]:")
			i++
		case Present:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case SetColor:
			text.WriteString(utils.FormatBytes(data[i : i+4]))
			i += 4
		case PresentUnknown05:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case PresentResetAlignment:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case RubyBaseStart:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case RubyTextStart:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case RubyTextEnd:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case SetFontSize:
			text.WriteString(utils.FormatBytes(data[i : i+3]))
			i += 3
		case PrintInParallel:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case PrintInCenter:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case SetMarginTop:
			text.WriteString(utils.FormatBytes(data[i : i+3]))
			i += 3
		case SetMarginLeft:
			text.WriteString(utils.FormatBytes(data[i : i+3]))
			i += 3
		case GetHardcodedValue:
			text.WriteString(utils.FormatBytes(data[i : i+3]))
			i += 3
		case EvaluateExpression:
			tmp := bytes.NewBuffer(nil)
			tmp.WriteByte(data[i])
			i++
			for {
				//for !(data[i] == 0 && data[i+1] == 0) {
				if (data[i] & 0x80) == 0x80 {
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
				} else {
					if data[i] == 0 {
						break
					} else {
						tmp.WriteByte(data[i])
						i++
					}
				}
			}
			tmp.WriteByte(data[i])
			text.WriteString(utils.FormatBytes(tmp.Bytes()))
			i++
		case PresentUnknown18:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case AutoForward:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case AutoForward1A:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case RubyCenterPerChar:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case AltLineBreak:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		case Terminator:
			text.WriteString(utils.FormatByte(data[i]))
			i++
		default:
			index := utils.BytesToUint16Big(data[i : i+2])
			if char, has := f.decodeCharset[index]; has {
				text.WriteString(char)
			} else {
				if utils.ShowWarning && data[i] >= 0x80 {
					fmt.Printf("Warning: 字库可能缺少 [%02X %02X] 对应的字符！\n", data[i], data[i+1])
				}
				text.WriteString(utils.FormatBytes(data[i : i+2]))
			}
			i += 2
		}

	}
	return text.String()
}
func (f *NpcsP) EncodeLine(str string) []byte {

	data := bytes.NewBuffer(nil)

	line := []rune(strings.TrimSpace(str))
	i := 0
	inBytes := false
	inName := false
	tempStr := ""
	for i < len(line) {
		if line[i] == ':' && i+1 < len(line) && line[i+1] == '[' && !(i+3 < len(line) && line[i+3] == 'x') {
			inName = true
			i += 2
		} else if line[i] == ']' && i+1 < len(line) && line[i+1] == ':' {
			if inName {
				data.WriteByte(NameStart)
				data.Write(f.stringToBytes(tempStr))
				data.WriteByte(LineStart)
				tempStr = ""
				inName = false
			} else {
				panic("错误的 ]: 结束符号，在：" + str)
			}
			i += 2
		} else if line[i] == '[' && !inName && line[i+1] == '0' && line[i+2] == 'x' {
			if len(tempStr) > 0 {
				data.Write(f.stringToBytes(tempStr))
				tempStr = ""
			}
			inBytes = true
			i += 3 //跳过[0x
		} else if line[i] == ']' {
			if inBytes {
				data.Write(utils.HexToBytes(tempStr))
				inBytes = false
				tempStr = ""
			} else {
				panic("错误的 ] 结束符号，在：" + str)
			}
			i++
		} else {
			tempStr += string(line[i])
			i++
		}
	}
	if len(tempStr) > 0 {
		data.Write(f.stringToBytes(tempStr))
	}
	return data.Bytes()
}
