package versiondiff

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/**

{
    "assets": {
        "assets/internal/config.json": {
            "md5": "9d9cc56e0930d4e3e7415aff3ce5c139",
            "size": 1312
        },
        "assets/internal/import/02/0275e94c-56a7-410f-bd1a-fc7483f7d14a.json": {
            "md5": "7871fc157af8a1f7ebef099323476b67",
            "size": 78
        },
	}
	"packageUrl": "https://acdn.gamesupernova.com/hotupdate_android/",
    "remoteManifestUrl": "https://acdn.gamesupernova.com/ts_web/manifest/android/project.manifest",
    "remoteVersionUrl": "https://acdn.gamesupernova.com/ts_web/manifest/android/version.manifest",
    "searchPaths": [],
    "version": "1.0.2303172042"
}
**/

type ResourceVal struct {
	Md5  string `json:"md5"`
	Size int32  `json:"size"`
}

type Versionjson struct {
	Assets            map[string]ResourceVal `json:"assets"`
	PackageUrl        string                 `json:"packageUrl"`
	RemoteManifestUrl string                 `json:"remoteVersionUrl"`
	SearchPaths       []string               `json:"searchPaths"`
	Version           string                 `json:"version"`
}

func Diff() {
	verfile := readFile("project_format.manifest")
	verfile2 := readFile("project_format.manifest_1.4.61")
	count := 0
	var notexistFile *[]string = new([]string)
	for fileName2, val2 := range verfile2.Assets {
		val, ok := verfile.Assets[fileName2]
		if ok {
			if val.Md5 != val2.Md5 {
				fmt.Printf("file md5 different:%s \n", fileName2)
				count++
			}
		} else {
			count++
			*notexistFile = append(*notexistFile, fileName2)
			// fmt.Printf("file not exists:%s \n", fileName2)
		}
	}

	for _, filename := range *notexistFile {
		fmt.Printf("new file:%s \n", filename)
	}

	fmt.Println("total diff file:", count)
}

func readFile(fileName string) (versionjson *Versionjson) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("error open json file")
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("error reading json file")
	}
	versionjson = &Versionjson{}
	err = json.Unmarshal(jsonData, versionjson)
	if err != nil {
		fmt.Println(err)
	}
	return
}
