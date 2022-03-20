package main

import (
	"MagesTools/script"
	"MagesTools/script/format"
	"MagesTools/script/utils"
	"flag"
	"fmt"
	"github.com/go-restruct/restruct"
	"strings"
)

func main() {

	var pType, pSource, pInput, pOutput, pScriptFormat, pCharset, pTbl string
	var pImport, pExport, pSkipChar bool
	flag.StringVar(&pType, "type", "", `[required] Source file type.
    MES(msb) Script: "script"
        Now only MES format scripts are supported
    Diff Binary File: "diff"
        Diff input and output file
`)
	flag.BoolVar(&pExport, "export", false, "[optional] Export mode")
	flag.BoolVar(&pImport, "import", false, "[optional] Import mode")

	flag.StringVar(&pSource, "source", "", `[required] Source file`)

	flag.StringVar(&pInput, "input", "", `[optional] Usually the import mode requires`)
	flag.StringVar(&pOutput, "output", "", `[required] Output file`)

	flag.StringVar(&pScriptFormat, "format", "Npcs", `[script.required] Format of script export and import. Case insensitive
    NPCSManager format: "Npcs"
    NPCSManager Plus format: "NpcsP"`)
	flag.StringVar(&pCharset, "charset", "", `[script.optional] Character set containing only text. Must be utf8 encoding. Choose between "charset" and "tbl"`)
	flag.StringVar(&pTbl, "tbl", "", `[script.optional] Text in TBL format. Must be utf8 encoding. Choose between "charset" and "tbl"`)

	flag.BoolVar(&pSkipChar, "skip", true, "[script.optional] Skip duplicate code table characters.")

	flag.Parse()
	restruct.EnableExprBeta()

	switch pType {
	case "diff":
		if len(pInput) == 0 && len(pOutput) == 0 {
			panic("Diff需要input和output")
		}
		res := utils.FileSameCheck(pInput, pOutput)
		if res {
			fmt.Println("两文件完全相同")
		}
	case "script":
		if !pExport && !pImport {
			panic("必须指定export模式或import模式")
		}
		if len(pSource) == 0 {
			panic("必须指定source源文件")
		}

		var _format format.Format
		switch strings.ToUpper(pScriptFormat) {
		case "NPCS":
			_format = &format.Npcs{}
		case "NPCSP":
			_format = &format.NpcsP{}
		default:
			panic("未知脚本导出格式")
		}
		scr := script.OpenScript(pSource, _format)

		if len(pCharset) > 0 {
			scr.LoadCharset(pCharset, false, pSkipChar)
		} else if len(pTbl) > 0 {
			scr.LoadCharset(pTbl, true, pSkipChar)
		} else {
			panic("必须指定charset文件或tbl文件")
		}
		scr.Read()
		if pExport {
			if len(pOutput) > 0 {
				scr.SaveStrings(pOutput)
			} else {
				panic("必须指定output文件")
			}
		} else if pImport {
			if len(pInput) > 0 {
				scr.LoadStrings(pInput)
			} else {
				panic("必须指定input文件")
			}

			if len(pOutput) > 0 {
				scr.Write(pOutput)
			} else {
				panic("必须指定output文件")
			}
		}
	}

}
