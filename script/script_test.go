package script

import "testing"

func TestOpenScript(t *testing.T) {
	scr := OpenScript("../data/CC/script/mes00/cc_01_01_00.msb")
	scr.LoadCharset("../data/CC/MJPN.txt", true)
	scr.Read()
	scr.Save("../data/CC/txt/cc_01_01_00.msb.txt")
}
