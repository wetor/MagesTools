package script

import (
	"encoding/hex"
	"fmt"
	"github.com/go-restruct/restruct"
	"os"
	"strings"
	"testing"
)

func TestLoadMes(t *testing.T) {
	restruct.EnableExprBeta()
	f, _ := os.Open("../data/CC/script/mes00/cc_01_01_00.msb")
	defer f.Close()
	mes := LoadMes(f)
	mes.ReadOffset()
	mes.ReadStrings(func(data []byte) *string {
		fmt.Println(data)
		return nil
	})
}

func Test001(t *testing.T) {
	src := []byte{1, 0, 123, 44}
	encodedStr := hex.EncodeToString(src)
	encodedStr = strings.ToUpper(encodedStr)
	fmt.Println(encodedStr)

	test, _ := hex.DecodeString(encodedStr)
	fmt.Println(test)
}
