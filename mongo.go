package main

import (
	"context"
	"fmt"
	openaitrans "main/openaiTrans"
	"mongo"
	"strings"

	"regexp"

	"gopkg.in/mgo.v2/bson"
)

func TestConn() {
	var dbconf []*mongo.DbConf
	// dbconf = append(dbconf, &mongo.DbConf{Host: "49.232.8.38", Port: "27017", User: "admin", Pass: "Rh8qie8KftgkWYND"})
	dbconf = append(dbconf, &mongo.DbConf{Host: "127.0.0.1", Port: "27017", User: "transuser", Pass: "LdtduhfUbZ21sGUS"})

	var dbconfs = &mongo.DbConnConf{
		Confs: dbconf,
	}
	dbconfs.Init()

	var conn = dbconfs.GetConn(context.WithValue(context.Background(), `uid`, 10))
	var mongo = &mongo.Mongo{ConSession: conn}
	var where = make(map[string]interface{})
	var d, err = mongo.Find("miaocha_trans", "cl_user_trans", where)
	var totalLen int
	pattern := regexp.MustCompile("掌柜")

	if err != nil {
		fmt.Println("find from db not passed", err)
	} else {
		for _, value := range d {
			a := pattern.FindSubmatch([]byte(value.Zh_txt))
			if len(a) > 0 {
				updateWhere := make(map[string]string)
				updateWhere["_id"] = value.Id
				updateValue := strings.ReplaceAll(value.Zh_txt, `掌柜`, `老板`)
				fmt.Println("relace string:", updateValue)
				err := mongo.Update("miaocha_trans", "cl_user_trans", bson.M{"_id": value.Id}, bson.M{"$set": bson.M{"zh_txt": updateValue}})
				if err != nil {
					fmt.Println(err, value.Id)
					break
				} else {
					totalLen++
					fmt.Println(value.Zh_txt, value.Id+` update succ`)
				}
				// totalLen++
				fmt.Println(`update total `, totalLen)
			}

		}
		fmt.Println(`total:`, totalLen)
	}
}

type Foo struct {
	a int32
	b int32
}

func main() {
	// testmap()
	// teaappstat.StoryBookStat()
	// teaappstat.StarStat()
	// teaappstat.DoMachineStat()
	// teaappstat.ClothStat()
	// teaapp.Init()
	// teaappstat.Do()
	// teaappstat.StageStat()
	// teaapp.InitConfig()
	// teaapp.DoPushBody()
	openaitrans.Do()
	// testmapslice()
}

func testmapslice() {
	starmidstat := make(map[string][]int)

	// 创建一个切片并将其分配给 map 的键
	starmidstat["a"] = []int{1, 2, 3}

	// 获取切片并进行修改
	midstar := starmidstat["a"]
	midstar = append(midstar, 4)

	// 输出 map 中的值
	fmt.Println(starmidstat["a"])
}
