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

func DiamondStat() {
	outputXlsx("diamondStat1.xlsx", "r_1", []string{"ID", "NUM"}, [][]string{{"1", "2"}, {"3", "4"}})
	return
	Init()
	var totalUser int = 0
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_table").String()
	diamondLogTable := teaapp.Cfg.Section("mongo").Key("diamond_log_table").String()
	usermdataTable := teaapp.Cfg.Section("mongo").Key("user_m_data").String()

	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"create_time": bson.M{"$gt": 10}})
	var user user
	userqueryIter := userquery.Iter()
	var diamondRetMap map[string]*map[int]int = make(map[string]*map[int]int)
	var diamondRange []int = []int{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100} // 钻石范围
	var restMaxId int = 15
	for userqueryIter.Next(&user) {
		fmt.Println("U:", user.ID)

		query := teaapp.Session.DB(database_1).C(usermdataTable).Find(bson.M{"_id": user.ID, "data.gameVersion": "1.0.6"})
		var userData interface{}
		err := query.One(userData)
		if err != nil {
			continue
		}
		totalUser++
		for i := 2; i < restMaxId; i++ {
			rid := fmt.Sprintf("%d", i)
			var diamondMap map[int]int
			if _, ok := diamondRetMap[rid]; ok {
				diamondMap = *diamondRetMap[rid]
			} else {
				diamondMap = make(map[int]int)
				diamondRetMap[rid] = &diamondMap
			}
			query := teaapp.Session.DB(database_1).C(diamondLogTable).Find(bson.M{"uid": user.ID, "related_id": rid, "level_id": "1", "prop_id": "10002"}).Sort("-create_time").Limit(1)
			var userProLog UserProLog
			err := query.One(&userProLog)
			if err != nil {
				continue
			}
			for j := len(diamondRange) - 1; j >= 0; j-- {
				if userProLog.EndNum > diamondRange[j] {
					diamondMap[j]++
					break
				}
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
		outputXlsx("diamondStat.xlsx", "r_"+key, title, values)
	}

}
