package teaappstat

import (
	"main/utils"
	"teaapp"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func SkinStat() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_table").String()
	usermSkinTable := teaapp.Cfg.Section("mongo").Key("skin_table").String()

	var now = time.Now().Unix()
	var weekAgo = now - 86400*14
	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"last_login_time": bson.M{"$gt": weekAgo}}).Iter()
	var user UserInfo
	var total = 0
	for userquery.Next(&user) {
		uid := user.Id
		if uid < 100 {
			continue
		}
		registTime := user.CreateTime
		lastLoginTime := user.LastLoginTime
		// 离注册时间超过7天，流失用户
		// if lastLoginTime < int(weekAgo) {
		registTimeStr, _ := utils.ConvertTimeStampToDateFormat(registTime)
		lastLoginStr, _ := utils.ConvertTimeStampToDateFormat(lastLoginTime)

		LogRus.Infof("uid:%d, registTime:%s, lastLoginTime:%s", uid, registTimeStr, lastLoginStr)

		userQuery := teaapp.Session.DB(database_1).C(usermSkinTable).Find(bson.M{"_id": uid})
		userQuery.Select(bson.M{"_id": 1, "data.myPartList": 1})
		queryIter := userQuery.Iter()
		var userSkin userSkinModule
		for queryIter.Next(&userSkin) {
			userSkindata, ok := userSkin.Data["myPartList"].([]interface{})
			if !ok {
				continue
			}
			for _, val := range userSkindata {
				if valint, ok := val.(int); ok {
					if valint == 22 {
						total++
						break
					}
				}

			}
		}
	}
	LogRus.Infof("total:%d", total)
	LogRus.Infof("SUCCESS")
}
