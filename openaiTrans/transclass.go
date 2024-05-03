package openaitrans

import (
	"main/teaappstat"
	"strings"
	"teaapp"

	"gopkg.in/mgo.v2/bson"
)

func UpdateTransClass() {

	teaappstat.Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_trans").String()
	transTable := teaapp.Cfg.Section("mongo").Key("trans_table").String()

	userquery := teaapp.Session.DB(database_1).C(transTable).Find(bson.M{}).Iter()
	var transItem TransItem
	for userquery.Next(&transItem) {
		_id := transItem.From
		_idArr := strings.Split(_id, ">")
		if len(_idArr) != 3 {
			continue
		}
		sheetName := _idArr[1]
		propertyId := _idArr[2]
		categoryName := ""
		if sheetName == "story_detail" && propertyId == "text" {
			categoryName = "story_detail"
		}
		if sheetName == "event_detail" && propertyId == "text" {
			categoryName = "event_detail"
		}
		if categoryName == "" {
			continue
		}

		err := teaapp.Session.DB(database_1).C(transTable).UpdateId(transItem.Id, bson.M{"$set": bson.M{"category_name": categoryName}})
		if err != nil {
			teaapp.LogRus.Errorf("update trans %s err %v", transItem.Id, err)
			continue
		}
		teaapp.LogRus.Printf("from:%s, category_name:%s", _id, categoryName)
	}
}
