package script

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

func BytesToUint16(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
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
		return fmt.Sprintf("[%dx%s]", len(data), BytesToHex(data))
	}
}
func FormatByte(data byte) string {
	return fmt.Sprintf("[1x%02X]", data)

}
