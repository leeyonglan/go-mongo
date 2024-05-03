package teaappstat

import (
	"fmt"
	"sort"
	"teaapp"

	"gopkg.in/mgo.v2/bson"
)

type UserProLog struct {
	UserId int `bson:"uid"`
	EndNum int `bson:"end_num"`
}

type userProData struct {
	ID   int                    `bson:"_id"`
	Data map[string]interface{} `bson:"data"`
}

func DiamondStat() {
	Init()
	var totalUser int = 0
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_table").String()
	usermProTable := teaapp.Cfg.Section("mongo").Key("props_table").String()
	userModuleTable := teaapp.Cfg.Section("mongo").Key("user_m_data").String()

	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"create_time": bson.M{"$gt": 1692892800}})
	var user user
	userqueryIter := userquery.Iter()
	var diamondRetMap map[string]*map[int]int = make(map[string]*map[int]int)
	// var diamondRange []int = []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100} // 钻石范围
	var diamondRange []int = []int{0, 1000, 2000, 3000, 4000, 5000, 10000} // 钻石范围

	for userqueryIter.Next(&user) {
		fmt.Println("U:", user.ID)

		//在哪个店铺
		userquery := teaapp.Session.DB(database_1).C(userModuleTable).Find(bson.M{"_id": user.ID, "data.gameVersion": bson.M{"$in": []string{"1.4.76", "1.4.78"}}})
		var usermodule userModule
		err := userquery.One(&usermodule)
		if err != nil {
			LogRus.Errorf("Query user:%d  err:%v", user.ID, err)
			continue
		}
		userData, ok := usermodule.Data["userData"].(map[string]interface{})
		if !ok {
			LogRus.Errorf("Assert userData err")
			continue
		}
		var maxStageId, sok = userData["maxStageId"]
		if !sok {
			continue
		}
		var maxRestaurantId, mok = userData["maxRestaurantId"]
		if !mok {
			continue
		}
		LogRus.Infof("User:%d maxStageId:%d maxRestaurantId:%d", user.ID, maxStageId, maxRestaurantId)
		if maxRestaurantId.(int) == 1 {
			if maxStageId.(int) <= 28 {
				continue
			} else {
				maxRestaurantId = 2
			}
		} else if maxRestaurantId.(int) >= 2 {
			if maxStageId.(int) >= 55 {
				maxRestaurantId = maxRestaurantId.(int) + 1
			} else if maxStageId.(int) > 30 {
				continue
			}
		}
		//钻石，金币
		query := teaapp.Session.DB(database_1).C(usermProTable).Find(bson.M{"_id": user.ID})
		var userProData userProData
		err = query.One(&userProData)
		if err != nil {
			continue
		}
		diamond, ok := userProData.Data["userTokenData"].(map[string]interface{})["token_coin"]
		if !ok {
			continue
		}
		LogRus.Infof("User:%d maxStageId:%d maxRestaurantId:%d,diamond:%d", user.ID, maxStageId, maxRestaurantId, diamond.(int))
		totalUser++

		rid := fmt.Sprintf("%d", maxRestaurantId)
		var diamondMap map[int]int
		if _, ok := diamondRetMap[rid]; ok {
			diamondMap = *diamondRetMap[rid]
		} else {
			diamondMap = make(map[int]int)
			diamondRetMap[rid] = &diamondMap
		}
		for j := len(diamondRange) - 1; j >= 0; j-- {
			if diamond.(int) > diamondRange[j] {
				diamondMap[j]++
				break
			}
		}

	}
	fmt.Println("totalUser:", totalUser)
	//output
	var title []string = []string{"ID", "NUM"}
	var mapkeys []string
	for key, _ := range diamondRetMap {
		mapkeys = append(mapkeys, key)
	}
	sort.Strings(mapkeys)
	for _, key := range mapkeys {
		var mapvalue map[int]int = *diamondRetMap[key]
		var values [][]string = make([][]string, 0)
		var mapvaluekeys []int
		for key, _ := range mapvalue {
			mapvaluekeys = append(mapvaluekeys, key)
		}
		sort.Ints(mapvaluekeys)
		for _, vkey := range mapvaluekeys {

			var rangeValue string = fmt.Sprintf("%d", diamondRange[vkey])
			var rangeValueNext string = ""
			if vkey == len(diamondRange)-1 {
				rangeValueNext = "∞"
			} else {
				rangeValueNext = fmt.Sprintf("%d", diamondRange[vkey+1])
			}

			var vkeyStr = fmt.Sprintf("%s-%s", rangeValue, rangeValueNext)
			var vkeyValStr = fmt.Sprintf("%d", mapvalue[vkey])

			value := [2]string{vkeyStr, vkeyValStr}
			values = append(values, value[:])
		}
		outputXlsx("coinStat.xlsx", "r_"+key, title, values)
	}

}
