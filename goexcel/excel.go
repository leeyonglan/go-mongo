package goexcel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/tealeg/xlsx"
)

func ReplaceContent() {
	src := os.Args[1]
	if src == "" {
		panic("no excel dir path find")
	}
	filepath.Walk(src, func(fliePath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			fmt.Println("dir:", fliePath)
			return nil
		}

		var ext = filepath.Ext(fliePath)
		if ext == ".DS_Store" {
			return nil
		}
		file, ferr := xlsx.OpenFile(fliePath)
		if ferr != nil {
			panic(ferr.Error())
		}
		for _, sheet := range file.Sheets {
			// fmt.Println(sheet.Name + "\n")
			for _, Row := range sheet.Rows {
				var flag bool = false
				for _, Cell := range Row.Cells {
					IndexRune := strings.IndexRune(Cell.Value, '掌')
					if IndexRune >= 0 {
						GuiCharIndex := strings.IndexRune(Cell.Value, '柜')
						if GuiCharIndex >= 0 {
							// fmt.Println("IndexRune:", IndexRune, "GuiCharIndex:", GuiCharIndex)
							if IndexRune+3 == GuiCharIndex {
								flag = true
								rstring := strings.Replace(Cell.Value, "掌柜", "老板", -1)
								fmt.Println("replace string:", rstring)
								Cell.Value = rstring
							}
						}
					}

				}
				if flag {
					fmt.Println("filename: " + fliePath + " sheet: " + sheet.Name + " id: " + Row.Cells[0].Value + "\n")
					err := file.Save(fliePath)
					if err != nil {
						panic(err.Error())
					}

				}
			}

		}
		return nil
	})
}

func IsChinese(str string) bool {
	var count int
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
			break
		}
	}
	return count > 0
}
