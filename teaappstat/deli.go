package teaappstat

import (
	"fmt"
	"teaapp"

	"gopkg.in/mgo.v2/bson"
)

type RecipeTable struct {
	// Implement RecipeTable methods if needed
}

type HomeDecorShopModel struct {
	// Implement HomeDecorShopModel methods if needed
}

type UserHomeModule struct {
	// Implement UserHomeModule methods if needed
}

func (m *HomeDecorShopModel) getTypeListMap() map[string][]struct {
	function int
	func_var int
	// Add other fields based on your actual data structure
} {
	// Implement getTypeListMap logic based on your requirements
	// This is a placeholder; you may need to fetch data from a database or elsewhere
	return make(map[string][]struct {
		function int
		func_var int
	})
}

func (m *UserHomeModule) getDecorSelectIndexByType(typeVar int) int {
	// Implement getDecorSelectIndexByType logic based on your requirements
	// This is a placeholder; you may need to fetch data from a database or elsewhere
	return 0
}

func getDeliciousValue(allHandBook map[string]int) int {
	delicious := 0
	configDB := &ConfigDB{}
	for key, value := range allHandBook {
		recipeData := configDB.getRecipeItemByKeyID("recipe_id", key)

		lv := getCardLvAndExp(recipeData.Rarity, value, configDB)["lv"].(int)
		rare := recipeData.Rarity
		dValue := rare * lv

		if key == "110001" && lv >= 1 {
			dValue += 1
		}

		delicious += dValue
	}

	// _decorShopModel := &HomeDecorShopModel{}
	// _decormap := _decorShopModel.getTypeListMap()
	// userHomeModule := &UserHomeModule{}

	// for key, values := range _decormap {
	// 	typeVar := 0
	// 	fmt.Sscanf(key, "%d", &typeVar)

	// 	index := userHomeModule.getDecorSelectIndexByType(typeVar)
	// 	if index >= 0 {
	// 		data := values[index]
	// 		if data.function == 4 {
	// 			delicious += data.func_var // 增加美味值
	// 		}
	// 	}
	// }

	return delicious
}

type UserHomeData struct {
	ID         int   `bson:"_id"`
	Data       Data  `bson:"data"`
	UpdateTime int64 `bson:"update_time"`
	CreateTime int64 `bson:"create_time"`
}

type Data struct {
	FoodList         []Food                   `bson:"foodList"`
	HomeDecoShow     []bool                   `bson:"homeDecoShow"`
	UserHomeDecorMap map[string]UserHomeDecor `bson:"userHomeDecorMap"`
}

type Food struct {
	ID      int   `bson:"id"`
	EndTime int64 `bson:"endTime"`
	Timeout bool  `bson:"timeout"`
}

type UserHomeDecor struct {
	ID       int  `bson:"id"`
	HasBuy   bool `bson:"hasBuy"`
	Selected bool `bson:"selected"`
}

type Lottery struct {
	ID   int `bson:"_id"`
	Data struct {
		AllHandBook map[string]int `bson:"allHandBook"`
		DataVersion int            `bson:"dataVersion"`
		GameVersion string         `bson:"gameVersion"`
		UserID      string         `bson:"userId"`
	} `bson:"data"`
}

var (
	starLevel = [][]int{{0, 200}, {201, 600}, {601, 1200}, {1201, 2000}, {2001, 3000}}
)

func Caculate() {
	Init()
	LoadAndParse()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermhome := teaapp.Cfg.Section("mongo").Key("user_m_home").String()
	userTable := teaapp.Cfg.Section("mongo").Key("user_table").String()
	lotteryTable := teaapp.Cfg.Section("mongo").Key("user_m_lottery").String()

	query := teaapp.Session.DB(database_1).C(userTable).Find(bson.M{"last_login_time": bson.M{"$gt": 1699511860}})

	var users []UserInfo
	err := query.All(&users)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	var totalDeskUser int
	for _, user := range users {
		var userId = user.Id

		userquery := teaapp.Session.DB(database_1).C(lotteryTable).Find(bson.M{"_id": userId})
		var userLottery Lottery
		var levelstart int
		err := userquery.One(&userLottery)
		if err != nil {
			fmt.Printf("err:%v,userId %d ", err, userId)
			fmt.Println()
			continue
		}
		allHandBook := userLottery.Data.AllHandBook
		delicious := getDeliciousValue(allHandBook)

		homeQuery := teaapp.Session.DB(database_1).C(usermhome).Find(bson.M{"_id": userId})
		var userHome UserHomeData
		err = homeQuery.One(&userHome)
		decorDelicious := 0
		if err != nil {
			fmt.Printf("err:%v,userId %d ", err, userId)
			fmt.Println()
		}
		for _, v := range userHome.Data.UserHomeDecorMap {
			if v.HasBuy {
				decorConfig := DecorMaps[v.ID]
				if decorConfig.Function == 4 {
					decorDelicious += decorConfig.FuncVar
				}

			}
		}
		delicious += decorDelicious
		for level, v := range starLevel {
			if delicious >= v[0] && delicious <= v[1] {
				levelstart = level
				break
			}
		}

		if levelstart == 3 {
			totalDeskUser++
		}
		// var userHome UserHome
	}
	println("totoal user:", totalDeskUser)
}
