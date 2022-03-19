package script

import (
	"github.com/go-restruct/restruct"
	"testing"
)

func TestOpenScript(t *testing.T) {

	restruct.EnableExprBeta()
	scr := OpenScript("../data/CC/script/mes00/cc_01_01_00.msb")
	scr.LoadCharset("../data/CC/MJPN.txt", true)
	scr.Read()

	// 导出
	scr.SaveStrings("../data/CC/txt/cc_01_01_00.msb.txt")
}
func TestScript_LoadStrings(t *testing.T) {
	restruct.EnableExprBeta()
	scr := OpenScript("../data/CC/script/mes00/cc_01_01_00.msb")
	scr.LoadCharset("../data/CC/MJPN.txt", true)
	scr.Read()

	// 导入
	scr.LoadStrings("../data/CC/txt/cc_01_01_00.msb.txt")
	scr.Write("../data/CC/cc_01_01_00.msb")

}
