package teaappstat

import (
	"fmt"
	"main/utils"

	"github.com/tealeg/xlsx"
)

type Android struct {
	IOS
}

/*
* 获取充值信息
 */
func (A Android) getPayList() map[string]*userOrderPay {

	var userOrderPayInfo *userOrderPay

	iapDict := make(map[string]*userOrderPay)
	//从文件里读
	reportad, err := xlsx.FileToSlice(androidPayPath)
	if err != nil {
		err = fmt.Errorf("read android pay file err:%w", err)
	}
	payContent := reportad[0][1:]
	if err != nil {
		err = fmt.Errorf("read report xlsx err:%w", err)
	}
	for _, v := range payContent {
		datestr := v[1]
		formatDateStr, err := utils.ConvertDateFormat(datestr)
		if err != nil {
			err = fmt.Errorf("convert date err:%w", err)
		}
		if orderInfo, ok := iapDict[formatDateStr]; ok {
			orderInfo.DayTotalPrice += getpriceById(v[9])
		} else {
			price := getpriceById(v[9])
			userOrderPayInfo = &userOrderPay{
				DayTotalPrice: price,
			}
			iapDict[formatDateStr] = userOrderPayInfo
		}
	}
	return iapDict
}
