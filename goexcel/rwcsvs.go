package goexcel

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// readCSV 读取 CSV 文件并返回关卡数据的映射
func readCSV(filename string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// 读取 CSV 文件
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	// 创建一个地图关卡数据的映射
	levelData := make(map[string]int)

	// 初始化地图关卡数据
	for i := 1; i <= 16; i++ {
		maxLevel := 60
		if i == 1 {
			maxLevel = 40
		}
		for j := 1; j <= maxLevel; j++ {
			key := fmt.Sprintf("%d_%d", i, j)
			levelData[key] = 0 // 初始化为 0
		}
	}

	// 补全现有的地图关卡数据
	level := 9000
	maxLevel := 15
	for i := 1; i <= maxLevel; i++ {
		key := fmt.Sprintf("%d_%d", level, i)
		levelData[key] = 0 // 你可以根据实际情况修改这里的逻辑，比如设为特殊值

	}
	level = 10000
	maxLevel = 10
	for i := 1; i <= maxLevel; i++ {
		key := fmt.Sprintf("%d_%d", level, i)
		levelData[key] = 0 // 你可以根据实际情况修改这里的逻辑，比如设为特殊值
	}

	// 补全现有的地图关卡数据
	for _, record := range records {
		level := record[0]
		value, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("Error converting value:", err)
			continue
		}

		key := strings.Trim(level, " ")
		levelData[key] = value // 你可以根据实际情况修改这里的逻辑，比如设为特殊值
	}

	return levelData, nil
}

// mergeCSVs 读取多个 CSV 文件，合并数据并生成一个新的 CSV 文件
func mergeCSVs(filenames []string, outputFilename string) error {
	levelData := make(map[string][]int)

	for idx, filename := range filenames {
		data, err := readCSV(filename)
		if err != nil {
			return err
		}

		for key, value := range data {
			if _, ok := levelData[key]; !ok {
				levelData[key] = make([]int, len(filenames))
			}
			levelData[key][idx] = value
		}
	}

	// 获取所有地图关卡数据的键并排序
	keys := make([]string, 0, len(levelData))
	for key := range levelData {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		mapI, _ := strconv.Atoi(keys[i][:strings.Index(keys[i], "_")])
		mapJ, _ := strconv.Atoi(keys[j][:strings.Index(keys[j], "_")])
		if mapI == mapJ {
			levelI, _ := strconv.Atoi(keys[i][strings.Index(keys[i], "_")+1:])
			levelJ, _ := strconv.Atoi(keys[j][strings.Index(keys[j], "_")+1:])
			return levelI < levelJ
		}
		return mapI < mapJ
	})

	file, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("error creating new file %s: %v", outputFilename, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	header := make([]string, len(filenames)+1)
	header[0] = "关卡"
	for idx := range filenames {
		header[idx+1] = fmt.Sprintf("文件%d", idx+1)
	}
	writer.Write(header)

	// 写入数据
	for _, key := range keys {
		row := make([]string, len(filenames)+1)
		row[0] = key
		for idx, value := range levelData[key] {
			row[idx+1] = strconv.Itoa(value)
		}
		writer.Write(row)
	}

	return nil
}

func DoCsv() {
	// 获取所有符合条件的 CSV 文件
	files, err := filepath.Glob("goexcel/csv/*.csv")
	if err != nil {
		fmt.Println("Error finding CSV files:", err)
		return
	}

	err = mergeCSVs(files, "merged_levels.csv")
	if err != nil {
		fmt.Println("Error merging CSVs:", err)
		return
	}

	fmt.Println("Merged CSV file created successfully.")
}
