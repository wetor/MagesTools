package script

import (
	"fmt"
	"github.com/go-restruct/restruct"
	"io"
	"os"
	"testing"
)

func TestLoadSc3(t *testing.T) {
	restruct.EnableExprBeta()
	f, _ := os.Open("../data/CCLCC/script/claa01.scx")
	defer f.Close()
	data, _ := io.ReadAll(f)
	mes := LoadSc3(data)
	mes.ReadStrings(func(data []byte) string {
		fmt.Println(data)
		return ""
	})
}
