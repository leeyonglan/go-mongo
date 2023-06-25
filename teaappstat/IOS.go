package teaappstat

import (
	"fmt"
	"teaapp"
	"time"

	"github.com/tealeg/xlsx"
	"gopkg.in/mgo.v2/bson"
)

type IOS struct {
	Name string
}

func (I IOS) getPlatName() string {
	return I.Name
}

func (I IOS) getAdList() (error, [][]string) {
	reportad, err := xlsx.FileToSlice(adreportPath)
	if err != nil {
		err = fmt.Errorf("read report xlsx err:%w", err)
	}
	list := make([][]string, 0)
	reportAdContent := reportad[0][1:]
	for _, row := range reportAdContent {
		if row[2] == I.Name {
			row[0] = formatDate(row[0])
			orderList := []string{
				row[0],
				row[1],
				row[2],
				row[4],
				row[7],
				"",
			}
			list = append(list, orderList)
		}
	}
	return err, list
}

func (I IOS) getPayList() map[string]*userOrderPay {
	database_0 := teaapp.Cfg.Section("mongo").Key("database_0_ios").String()
	orderpay := teaapp.Cfg.Section("mongo").Key("order_pay_table").String()
	userquery := teaapp.Session.DB(database_0).C(orderpay).Find(bson.M{"is_online_pay": 1}).Iter()

	var userOrderPayInfo *userOrderPay

	iapDict := make(map[string]*userOrderPay)
	for userquery.Next(&userOrderPayInfo) {
		timeTime := time.Unix(int64(userOrderPayInfo.CreateTime), 0)
		if userOrderPayInfo.CreateTime == 0 {
			ctime := userOrderPayInfo.GoodInfo["create_time"].(int)
			timeTime = time.Unix(int64(ctime), 0)
		}
		datestr := timeTime.Format("2006/01/02")

		if orderInfo, ok := iapDict[datestr]; ok {
			orderInfo.DayTotalPrice += getpriceById(userOrderPayInfo.GoodsId)
		} else {
			price := getpriceById(userOrderPayInfo.GoodsId)
			userOrderInfo := &userOrderPay{
				DayTotalPrice: price,
			}
			iapDict[datestr] = userOrderInfo
		}
		fmt.Println("IOS GOODS ID:", datestr, userOrderPayInfo.GoodsId)
	}
	return iapDict
}
