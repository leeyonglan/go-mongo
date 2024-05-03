package goexcel

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func Do() {
	// 打开原始 CSV 文件
	file, err := os.Open("levels.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 读取 CSV 文件
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
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
		level, err := strconv.Atoi(record[0])
		if err != nil {
			fmt.Println("Error converting level:", err)
			continue
		}
		value, err := strconv.Atoi(record[1])
		if err != nil {
			fmt.Println("Error converting value:", err)
			continue
		}

		key := fmt.Sprintf("%d", level)
		levelData[key] = value // 你可以根据实际情况修改这里的逻辑，比如设为特殊值
	}

	// 将地图关卡数据写入新的 CSV 文件
	newFile, err := os.Create("new_levels.csv")
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	for key, value := range levelData {
		writer.Write([]string{key, strconv.Itoa(value)})
	}
}
