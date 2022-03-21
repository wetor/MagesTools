package script

import (
	"MagesTools/script/format"
	"MagesTools/script/utils"
	"fmt"
	"github.com/go-restruct/restruct"
	"path"
	"testing"
)

func TestScript_CCLCC_Export(t *testing.T) {

	restruct.EnableExprBeta()
	scr := NewScript("../data/CCLCC/script/claa01.scx", &format.NpcsP{})
	scr.LoadCharset("../data/CCLCC/MJPN.txt", true, false)
	scr.Read()

	// 导出
	scr.SaveStrings("../data/CCLCC/txt/claa01.scx.txt")
}
func TestScript_CCLCC_Import(t *testing.T) {
	restruct.EnableExprBeta()
	scr := NewScript("../data/CCLCC/script/claa01.scx", &format.NpcsP{})
	scr.LoadCharset("../data/CCLCC/MJPN.txt", true, false)
	scr.Read()

	// 导入
	scr.LoadStrings("../data/CCLCC/txt/claa01.scx.txt")
	scr.Write("../data/CCLCC/claa01.scx")

}

// TestScript_CCLCC_Check 检查CC导出导入后变化
func TestScript_CCLCC_Check(t *testing.T) {
	restruct.EnableExprBeta()

	file := "../data/CCLCC/script/claa01.scx"
	charset := "../data/CCLCC/MJPN.txt"

	scr := NewScript(file, &format.NpcsP{})
	scr.LoadCharset(charset, true, false)
	scr.Read()

	// 导出
	scr.SaveStrings("../data/temp/1.txt")
	// 导入
	scr.LoadStrings("../data/temp/1.txt")
	scr.Write("../data/temp/1.scx")

	utils.FileSameCheck(file, "../data/temp/1.scx")

}

func TestScript_CC_Export(t *testing.T) {
	utils.ShowWarning = true
	restruct.EnableExprBeta()
	scr := NewScript("../data/CC/script/mes00/cc_01_01_00.msb", &format.NpcsP{})
	scr.LoadCharset("../data/CC/MJPN.txt", true, true)
	scr.Read()

	// 导出
	scr.SaveStrings("../data/CC/txt/cc_01_01_00.msb.txt")
}

func TestScript_CC_Export_All(t *testing.T) {
	utils.ShowWarning = true
	dir := "../data/CC/script/mes00"
	tbl := "../data/CC/MJPN.txt"
	saveDir := "../data/CC/txt/"

	restruct.EnableExprBeta()

	scr := &Script{}
	scr.LoadCharset(tbl, true, false)
	files, _ := utils.GetDirFileList(dir)
	for _, file := range files {
		fmt.Println(file)
		scr.Open(file, &format.NpcsP{})
		scr.Read()
		// 导出
		scr.SaveStrings(saveDir + path.Base(file) + ".txt")
	}

}

func TestScript_CC_Import(t *testing.T) {
	restruct.EnableExprBeta()
	scr := NewScript("../data/CC/script/mes00/cc_01_01_00.msb", &format.NpcsP{})
	scr.LoadCharset("../data/CC/MJPN.txt", true, true)
	scr.Read()

	// 导入
	scr.LoadStrings("../data/CC/txt/cc_01_01_00.msb.txt")
	scr.Write("../data/CC/cc_01_01_00.msb")

}

// TestScript_CC_Check 检查CC导出导入后变化
func TestScript_CC_Check(t *testing.T) {
	restruct.EnableExprBeta()

	file := "../data/CC/script/mes00/cc_01_01_00.msb"
	charset := "../data/CC/MJPN.txt"

	scr := NewScript(file, &format.NpcsP{})
	scr.LoadCharset(charset, true, true)
	scr.Read()

	// 导出
	scr.SaveStrings("../data/temp/1.txt")
	// 导入
	scr.LoadStrings("../data/temp/1.txt")
	scr.Write("../data/temp/1.msb")

	utils.FileSameCheck(file, "../data/temp/1.msb")

}

func TestScript_RNE_Export(t *testing.T) {

	restruct.EnableExprBeta()
	scr := NewScript("../data/RNE/script/mes00/RN05_20A_00.msb", &format.Npcs{})
	scr.LoadCharset("../data/RNE/Charset_PSV_JP.utf8", false, false)
	scr.Read()
	// 导出
	scr.SaveStrings("../data/RNE/txt/RN05_20A_00.msb.txt")
}
func TestScript_RNE_Import(t *testing.T) {
	restruct.EnableExprBeta()
	scr := NewScript("../data/RNE/script/mes00/RN05_20A_00.msb", &format.Npcs{})
	scr.LoadCharset("../data/RNE/Charset_PSV_JP.utf8", false, false)
	scr.Read()
	// 导入
	scr.LoadStrings("../data/RNE/txt/RN05_20A_00.msb_tool.txt")
	scr.Write("../data/RNE/RN05_20A_00.msb")

}

// TestScript_CC_Check 检查RNE导出导入后变化
func TestScript_RNE_Check(t *testing.T) {
	restruct.EnableExprBeta()

	file := "../data/RNE/script/mes00/RN05_20A_00.msb"
	charset := "../data/RNE/Charset_PSV_JP.utf8"

	scr := NewScript(file, &format.NpcsP{})
	scr.LoadCharset(charset, false, true)
	scr.Read()

	// 导出
	scr.SaveStrings("../data/temp/1.txt")
	// 导入
	scr.LoadStrings("../data/temp/1.txt")
	scr.Write("../data/temp/1.msb")

	utils.FileSameCheck(file, "../data/temp/1.msb")

}
