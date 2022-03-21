package script

import (
	"MagesTools/script/format"
	"MagesTools/script/utils"
	"bufio"
	"io"
	"os"
)

type Entry struct {
	Index  int `struct:"int32"`
	Offset int `struct:"int32"`
	Length int `struct:"-"`
}

type Strings interface {
	// ReadStrings
	//  Description 读取解析脚本全部字符串
	//  Param readString
	ReadStrings(readString func([]byte) string)
	// GetStrings
	//  Description 取出全部字符串
	//  Return []string
	GetStrings() []string
	// SetStrings
	//  Description 替换全部字符串
	//  Param strings
	SetStrings(strings []string)
	// WriteStrings
	//  Description 写到导入字符串
	//  Param writeString
	WriteStrings(writeString func(string) []byte)
	// GetRaw
	//  Description 获取脚本数据
	//  Return []byte
	GetRaw() []byte
}

type Script struct {
	Strings       Strings
	Format        format.Format
	DecodeCharset map[uint16]string
	EncodeCharset map[string]uint16
}

// NewScript
//  Description 打开脚本文件
//  Param filename string
//  Return *Script
//
func NewScript(filename string, format format.Format) *Script {
	script := &Script{}
	script.Open(filename, format)
	return script
}

// Open
//  Description 打开脚本文件，如果已经使用LoadCharset载入码表，则不需要重新调用LoadCharset
//  Receiver s *Script
//  Param filename string
//  Param format format.Format
//
func (s *Script) Open(filename string, format format.Format) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	s.Format = format
	switch string(data[0:3]) {
	case "MES":
		s.Strings = LoadMes(data)
	case "SC3":
		s.Strings = LoadSc3(data)
	default:
		panic("不支持的文件类型！ " + filename)
	}
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
				k := utils.BytesToUint16Big(utils.HexToBytes(line[0:4]))
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
	s.DecodeCharset = decodeCharset
	s.EncodeCharset = encodeCharset
}

// Read
//  Description 解析文本，需要至少执行一次script.LoadCharset载入码表
//  Receiver s *Script
//
func (s *Script) Read() {
	if s.DecodeCharset != nil && s.EncodeCharset != nil {
		s.Format.SetCharset(s.DecodeCharset, s.EncodeCharset)
	}
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
