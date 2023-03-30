package fileio

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var (
// prefabFileList =
)

func CP(src string, dest string) {
	fmt.Printf(`src:%s, dest:%s`, src, dest)
	ret := recursiveCheck(src)

	if ret != nil {
		fmt.Println(ret)
	}
}

func recursiveCheck(dirpath string) error {
	dir, _ := os.ReadDir(dirpath)
	for _, file := range dir {
		file_path := filepath.Join(dirpath, file.Name())
		if file.IsDir() {
			recursiveCheck(file_path)
		}
		if path.Ext(file_path) == ".prefab" {
			fmt.Printf("file:%s \n", file_path)
			// fbyte, err := os.ReadFile(file_path)
			// if err != nil {
			// 	fmt.Printf("read file err,file:%s,err:%s", file_path, err.Error())
			// }
			// os.WriteFile("aa", fbyte)
		}

	}
	return nil
	// return filepath.Walk(dir,
	// 	func(fpath string, f os.FileInfo, err error) error {
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if f.IsDir() {
	// 			return recursiveCheck(fpath)
	// 		}
	// 		fmt.Printf("%s", path.Ext(f.Name()))
	// 		if path.Ext(f.Name()) == "prefab" {
	// 			fmt.Printf("prefab file:%s \r\n", f.Name())
	// 		}
	// 		return nil
	// 	})
}
