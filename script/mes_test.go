package script

import (
	"fmt"
	"github.com/go-restruct/restruct"
	"io"
	"os"
	"testing"
)

func TestLoadMes(t *testing.T) {
	restruct.EnableExprBeta()
	f, _ := os.Open("../data/CC/script/mes00/cc_01_01_00.msb")
	defer f.Close()
	data, _ := io.ReadAll(f)
	mes := LoadMes(data)
	mes.ReadStrings(func(data []byte) string {
		fmt.Println(data)
		return ""
	})
}
