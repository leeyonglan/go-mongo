package mergeiostoandroid

import (
	"strconv"
	"teaapp"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Init() {
	teaapp.InitEnv()
	teaapp.InitConfig()
	teaapp.InitMongo()
	importIosUserToAndrod()
	transforTableData()
}

type IdmapNew struct {
	Status      int8   `json:"status"`
	Uid         int    `json:"uid"`
	Is_Official int8   `json:"is_officical"`
	Snsid       string `json:"snsid"`
}
type IdSquence struct {
	Init  int8  `json:"init"`
	Value int32 `json:"value"`
}

func importIosUserToAndrod() {
	database_0_ios := teaapp.Cfg.Section("online_0_ios").Key("database").String()
	database_0_android := teaapp.Cfg.Section("online_0_glo").Key("database").String()

	iosQuery := teaapp.Session.DB(database_0_ios).C("idmapnew").Find(bson.M{"uid": bson.M{"$gt": 100}}).Iter()

	var idItem IdmapNew
	for iosQuery.Next(&idItem) {
		userId := genNextId()
		if userId == 0 {
			break
		}
		hasQuery := teaapp.Session.DB(database_0_android).C("idmapnew").Find(bson.M{"snsid": idItem.Snsid, "status": 1, "ios": 1})
		hasNum, err := hasQuery.Count()
		if err != nil {
			teaapp.LogRus.Errorln(err)
			break
		}
		if hasNum == 0 {
			err := teaapp.Session.DB(database_0_android).C("idmapnew").Insert(
				bson.M{"snsid": idItem.Snsid, "is_officical": 0, "uid": userId, "status": idItem.Status, "iosuid": idItem.Uid, "ios": 1})
			if err != nil {
				teaapp.LogRus.Errorln(err)
				break
			} else {
				teaapp.LogRus.Printf("transfor ios uid %d succ ,new uid:%d", idItem.Uid, userId)
			}
		} else {
			teaapp.LogRus.Printf("transfor ios uid %d has exists snsid:%s", idItem.Uid, idItem.Snsid)
		}
	}
}

func genNextId() int32 {
	change := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"value": 1}},
		ReturnNew: true,
	}
	var doc IdSquence
	database_0_android := teaapp.Cfg.Section("online_0_glo").Key("database").String()
	_, err := teaapp.Session.DB(database_0_android).C("idsquence").Find(bson.M{
		"init": 1,
	}).Apply(change, &doc)
	if err != nil {
		teaapp.LogRus.Errorln(err)
		return 0
	}
	return doc.Value
}

var (
	uidMap map[int]int = make(map[int]int, 100000)
)

func transforTableData() {
	var idTables = []string{
		"user",
		"user_config",
		"user_m_activity",
		"user_m_battle_pass",
		"user_m_data",
		"user_m_game_props",
		"user_m_home",
		"user_m_lottery",
		"user_m_notification",
		"user_m_pawprime",
		"user_m_piggybank",
		"user_m_rank",
		"user_m_restaurant",
		"user_m_skin_map",
		"user_m_story_book",
		"user_m_troubadour",
		"user_m_world_map",
		"user_recipe",
		"user_m_friendgift", //friendIdList fuids
	}
	var uid_F_uidTables = []string{
		"user_friend_content", //uid,f_uid
		"user_msg",            // fromUid toUid
		"user_login_reward",   //uid
		"user_rank_data",      //uid
	}
	database_ios := teaapp.Cfg.Section("online_ios").Key("database").String()
	database_android := teaapp.Cfg.Section("online_glo").Key("database").String()
	for _, tableName := range idTables {
		tableQuery := teaapp.Session.DB(database_ios).C(tableName).Find(bson.M{"_id": bson.M{"$gt": 100}}).Iter()

		var item bson.M
		for tableQuery.Next(&item) {
			_id := item["_id"].(int)
			oldUid := int(_id)

			newUid := getNewUidId(oldUid)
			if newUid == 0 {
				continue
			}
			item["_id"] = newUid

			// 好友列表数据
			if tableName == "user_m_friendgift" {
				var fuids []int = make([]int, 0)
				_, ok := item["friendIdList"]
				friends, tok := item["friendIdList"].([]interface{})
				if ok && tok {
					teaapp.LogRus.Infof("table:%s insert olduid:%d data err %v", tableName, oldUid, item)
					for _, fuid := range friends {
						fuidInt, Ok := fuid.(int)
						if !Ok {
							teaapp.LogRus.Errorf("table:%s insert olduid:%d interface to int fail,data err %v", tableName, oldUid, fuid)
							continue
						}
						if fuidInt <= 100 {
							fuids = append(fuids, int(fuidInt))
						} else {
							newFuid := getNewUidId(int(fuidInt))
							fuids = append(fuids, int(newFuid))
						}
					}
					item["friendIdList"] = fuids
				}
				_, hok := item["helpList"]
				helps, hhok := item["helpList"].(bson.M)
				newHelps := bson.M{}
				if hok && hhok && helps != nil {
					for fuid, v := range helps {
						fuidInt, err := strconv.Atoi(fuid)
						if err != nil {
							teaapp.LogRus.Errorf("table:%s convert fuid to int  olduid:%d string to int fail,data err %v", tableName, oldUid, fuid)
							continue
						}
						newFuid := getNewUidId(int(fuidInt))
						newHelps[strconv.Itoa(newFuid)] = v
					}
					item["helpList"] = newHelps
				} else {
					delete(item, "helpList")
					// teaapp.LogRus.Errorf("table:%s convert helpList olduid:%d data err %v", tableName, oldUid, item)
				}
			}

			err := teaapp.Session.DB(database_android).C(tableName).Insert(item)
			if err != nil {
				teaapp.LogRus.Errorf("table:%s insert olduid:%d data err %s", tableName, oldUid, err.Error())
				continue
			} else {
				teaapp.LogRus.Infof("table:%s insert olduid:%d data newuid:%d,succ", tableName, oldUid, newUid)
			}
		}
	}

	for _, tableName := range uid_F_uidTables {
		tableQuery := teaapp.Session.DB(database_ios).C(tableName).Find(bson.M{}).Iter()
		var item bson.M
		for tableQuery.Next(&item) {
			var oldUid int
			var newUid int
			if tableName == "user_friend_content" || tableName == "user_login_reward" || tableName == "user_rank_data" {
				uid, ok := item["uid"].(int)
				if !ok {
					teaapp.LogRus.Errorf("table:%s interface to int uid olduid:%d data err %v", tableName, uid, item)
					continue
				}
				if uid <= 100 {
					continue
				}
				teaapp.LogRus.Infof("table:%s start insert olduid:%d", tableName, uid)
				newUid = getNewUidId(int(uid))
				item["uid"] = newUid
				oldUid = uid

			}
			if tableName == "user_friend_content" {
				f_uid, ok := item["f_uid"].(int)
				if !ok {
					teaapp.LogRus.Errorf("table:%s interface to int f_uid :%d data err %v", tableName, f_uid, item)
					continue
				}
				newF_uid := getNewUidId(int(f_uid))
				item["f_uid"] = newF_uid

			}

			if tableName == "user_msg" {
				fromUid, ok := item["fromUid"].(int)
				if !ok {
					teaapp.LogRus.Errorf("table:%s interface to int fromUid :%d data err %v", tableName, fromUid, item)
					continue
				}
				toUid, ok := item["toUid"].(int)
				if !ok {
					teaapp.LogRus.Errorf("table:%s interface to int toUid :%d data err %v", tableName, toUid, item)
					continue
				}
				newFromUid := getNewUidId(int(fromUid))
				newToUid := getNewUidId(int(toUid))
				item["fromUid"] = newFromUid
				item["toUid"] = newToUid
				oldUid = fromUid
				newUid = int(newFromUid)
			}
			delete(item, "_id")
			err := teaapp.Session.DB(database_android).C(tableName).Insert(item)
			if err != nil {
				teaapp.LogRus.Errorf("table:%s insert olduid:%d data err %v", tableName, oldUid, err)
				continue
			} else {
				teaapp.LogRus.Infof("table:%s insert olduid:%d data newuid:%d,succ", tableName, oldUid, newUid)
			}
		}
	}

}

func getNewUidId(oldUid int) int {
	newUid, ok := uidMap[oldUid]
	if ok {
		return newUid
	}
	database_0_android := teaapp.Cfg.Section("online_0_glo").Key("database").String()

	query := teaapp.Session.DB(database_0_android).C("idmapnew").Find(bson.M{"status": 1, "iosuid": oldUid})
	var idmapInfo IdmapNew
	err := query.One(&idmapInfo)
	if err != nil {
		teaapp.LogRus.Errorln(err, oldUid)
		return 0
	}
	uidMap[oldUid] = idmapInfo.Uid
	return idmapInfo.Uid
}
