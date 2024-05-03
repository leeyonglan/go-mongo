package deviceidfix

import (
	"strings"
	"teaapp"

	"gopkg.in/mgo.v2/bson"
)

func Init() {
	teaapp.InitEnv()
	teaapp.InitConfig()
	teaapp.InitMongo()
	queryNum()
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

func queryNum() {
	database_0_android := teaapp.Cfg.Section("online_0_glo").Key("database").String()

	iosQuery := teaapp.Session.DB(database_0_android).C("idmapnew").Find(bson.M{"uid": bson.M{"$gt": 100}, "status": 1}).Iter()
	var totalCount int
	var errorCount int
	var idItem IdmapNew
	for iosQuery.Next(&idItem) {
		snsId := idItem.Snsid
		if strings.HasPrefix(snsId, "gpc_") {
			totalCount++
			teaapp.LogRus.Infof("gpc snsid:%s", snsId)
			originSnsId := strings.TrimPrefix(snsId, "gpc_")

			hasQuery := teaapp.Session.DB(database_0_android).C("idmapnew").Find(bson.M{"snsid": originSnsId})
			hasNum, err := hasQuery.Count()
			if err != nil {
				teaapp.LogRus.Errorln(err)
				break
			}
			if hasNum > 0 {
				teaapp.LogRus.Infof("pgs origin snsid:%s has exists", originSnsId)
				errorCount++
			}
		}
	}
	teaapp.LogRus.Infof("total count:%d,error count:%d", totalCount, errorCount)
}
