package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

func BytesToHex(data []byte) string {
	return strings.ToUpper(hex.EncodeToString(data))
}

func HexToBytes(data string) []byte {
	bytes, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return bytes
}

func BytesToUint16Big(data []byte) uint16 {
	return binary.BigEndian.Uint16(data)
}

func BytesToUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}
func Uint16ToBytesBig(data uint16) []byte {
	bytes := make([]byte, 2)
	binary.BigEndian.PutUint16(bytes, data)
	return bytes
}
func Uint16ToBytes(data uint16) []byte {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, data)
	return bytes
}

func FormatBytes(data []byte) string {
	if len(data) == 0 {
		return ""
	} else {
		return fmt.Sprintf("[0x%s]", BytesToHex(data))
	}
}
func FormatByte(data byte) string {
	return fmt.Sprintf("[0x%02X]", data)

}

func FileSameCheck(file1, file2 string) bool {
	data1, _ := os.ReadFile(file1)
	data2, _ := os.ReadFile(file2)
	if len(data1) != len(data2) {
		fmt.Printf("文件大小不一致。file1:%d  file2:%d\n", len(data1), len(data2))
		return false
	}
	same := true
	fmt.Printf("%08X: input > output\n", 0)
	for i := 0; i < len(data1); i++ {
		if data1[i] != data2[i] {
			same = false
			fmt.Printf("%08X:   %02X  >  %02X\n", i, data1[i], data2[i])
		}
	}
	return same
}
