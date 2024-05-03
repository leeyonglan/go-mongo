package utils

import "time"

func ConvertDateFormat(input string) (string, error) {
	// 定义输入和输出的时间格式
	inputFormat := "2-Jan-06"
	outputFormat := "2006/01/02"

	// 将输入字符串解析为时间对象
	t, err := time.Parse(inputFormat, input)
	if err != nil {
		return "", err
	}

	// 格式化时间为指定输出格式的字符串
	output := t.Format(outputFormat)

	return output, nil
}

func ConvertTimeStampToDateFormat(input int) (string, error) {
	// 定义输入和输出的时间格式
	outputFormat := "2006-01-02 15:04:05"

	// 将输入字符串解析为时间对象
	t := time.Unix(int64(input), 0)

	// 格式化时间为指定输出格式的字符串
	output := t.Format(outputFormat)

	return output, nil
}
