package openaitrans

import (
	"main/teaappstat"
	"os/exec"
	"strings"
	"teaapp"

	"gopkg.in/mgo.v2/bson"
)

type TransItem struct {
	Id           string `bson:"_id"`
	UpdateTime   int    `bson:"update_time"`
	ZhTxt        string `bson:"zh_txt"`
	ZhType       int8   `bson:"zh_type"`
	EnTxt        string `bson:"en_txt"`
	EnType       int8   `bson:"en_type"`
	Remark       string `bson:"remark"`
	WordLimit    int    `bson:"word_limit"`
	From         string `bson:"from"`
	CategoryName string `bson:"category_name"`
	CreateTime   int    `bson:"create_time"`
}

var (
	pythonExec       string
	pythonScriptPath string
)

func Do() {

	teaappstat.Init()

	// dstLang := "SP"
	langs := []string{"VI", "SP", "FR", "MA", "TH", "PO", "AR"}
	for _, dstLang := range langs {
		database_1 := teaapp.Cfg.Section("mongo").Key("database_trans").String()
		transTable := teaapp.Cfg.Section("mongo").Key("trans_table").String()
		keyPrefix := strings.ToLower(dstLang)
		langKey := keyPrefix + "_txt"
		userquery := teaapp.Session.DB(database_1).C(transTable).Find(bson.M{"en_txt": bson.M{"$ne": ""}, langKey: bson.M{"$eq": ""}}).Iter()

		pythonExec = teaapp.Cfg.Section("trans").Key("pythonexec").String()
		pythonScriptPath = teaapp.Cfg.Section("trans").Key("pythonscript").String()

		var transItem TransItem
		var transucc int
		for userquery.Next(&transItem) {
			teaapp.LogRus.Printf("start trans en to %s : %s", dstLang, transItem.EnTxt)
			tranTxt := ExecOpenAI(transItem.EnTxt, dstLang)
			if len(tranTxt) > 0 {
				teaapp.LogRus.Printf("success trans: %s", tranTxt)
				transucc++
				teaapp.LogRus.Printf("success trans total:%d", transucc)

				langType := keyPrefix + "_type"
				err := teaapp.Session.DB(database_1).C(transTable).UpdateId(transItem.Id, bson.M{"$set": bson.M{langKey: tranTxt, langType: 0}})
				if err != nil {
					teaapp.LogRus.Errorf("update trans %s err %v", transItem.Id, err)
					continue
				}
			}
		}
	}
}

func ExecOpenAI(srcTxt string, lang string) string {
	cmd := exec.Command(pythonExec, pythonScriptPath, srcTxt, lang)
	output, err := cmd.Output()
	if err != nil {
		teaapp.LogRus.Errorf("trans error srcTxt %s,err %v", srcTxt, err)
		return ""
	}
	return string(output)
}
