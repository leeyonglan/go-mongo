package teaappstat

import (
	"main/utils"
	"sort"
	"teaapp"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func LostUserCoinStat() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_table").String()
	usermProTable := teaapp.Cfg.Section("mongo").Key("props_table").String()
	var coinsArr []int
	var now = time.Now().Unix()
	var weekAgo = now - 86400*14
	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"create_time": bson.M{"$gt": 1701360000, "$lt": weekAgo}}).Iter()
	var user UserInfo
	for userquery.Next(&user) {
		uid := user.Id
		registTime := user.CreateTime
		lastLoginTime := user.LastLoginTime
		// 离注册时间超过7天，流失用户
		if (lastLoginTime - registTime) < 86400*14 {
			// if lastLoginTime < int(weekAgo) {
			registTimeStr, _ := utils.ConvertTimeStampToDateFormat(registTime)
			lastLoginStr, _ := utils.ConvertTimeStampToDateFormat(lastLoginTime)

			LogRus.Infof("uid:%d, registTime:%s, lastLoginTime:%s", uid, registTimeStr, lastLoginStr)

			query := teaapp.Session.DB(database_1).C(usermProTable).Find(bson.M{"_id": uid})
			var userProData userProData
			err := query.One(&userProData)
			if err != nil {
				continue
			}
			token_coin, ok := userProData.Data["userTokenData"].(map[string]interface{})["token_coin"]
			if !ok {
				continue
			}
			coin, ok := token_coin.(int)

			coinsArr = append(coinsArr, coin)
			LogRus.Infof("uid:%d, coins:%d", uid, coin)
		}
	}
	//金币范围
	var coinRange []int = []int{0, 50, 100, 300, 500, 1000} // 钻石范围
	var coinRangeMap = make(map[int]int)
	//总人数
	var totalUser = len(coinsArr)
	// 排序
	sort.Ints(coinsArr)
	//中位数
	median := coinsArr[len(coinsArr)/2]
	// 平均数
	var sum int
	for _, v := range coinsArr {
		sum += v
		if v <= coinRange[1] {
			coinRangeMap[coinRange[1]]++
		}
		if v > coinRange[1] && v <= coinRange[2] {
			coinRangeMap[coinRange[2]]++
		}
		if v > coinRange[2] && v <= coinRange[3] {
			coinRangeMap[coinRange[3]]++
		}
		if v > coinRange[3] && v <= coinRange[4] {
			coinRangeMap[coinRange[4]]++
		}
		if v > coinRange[4] && v <= coinRange[5] {
			coinRangeMap[coinRange[5]]++
		}
		if v > coinRange[5] {
			coinRangeMap[1001]++
		}
	}
	avg := sum / len(coinsArr)

	LogRus.Infof("median:%d, avg:%d", median, avg)
	LogRus.Infof("totalNum:%d,coinRangeMap:%v", totalUser, coinRangeMap)

	LogRus.Infof("SUCCESS")
}
