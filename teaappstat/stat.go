package teaappstat

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
	"teaapp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx/v3"
	"gopkg.in/ini.v1"
	"gopkg.in/mgo.v2/bson"
)

var (
	wg     *sync.WaitGroup
	cfg    *ini.File
	LogRus *log.Logger
)

type userRestaurantModule struct {
	ID   int                    `bson:"_id"`
	Data map[string]interface{} `bson:"data"`
}

type userSkinModule struct {
	ID   int                    `bson:"_id"`
	Data map[string]interface{} `bson:"data"`
}
type userModule struct {
	ID   int                    `bson:"_id"`
	Data map[string]interface{} `bson:"data"`
}
type user struct {
	ID         int `bson:"_id"`
	CreateTime int `bson:"create_time"`
}

func Init() {
	teaapp.InitEnv()
	teaapp.InitConfig()
	teaapp.InitMongo()
	LogRus = teaapp.LogRus
}

func friendGiftStat() {

}

type starStat struct {
	Uid     int
	Lid     int
	StarNum int
}

type StarMap struct {
	Mid         int
	StartTotal  int
	PlayerTotal int
	AverageNum  int
}

type storybook struct {
	Uid  int `bson:"_id"`
	Data map[string]interface{}
}

type storybookHelp struct {
	MapId        int
	Level        int
	RwardTotal   int
	UnLoackTotal int
	Sort         int
	Rate         string
}
type storyConfig struct {
	StoryTitle []storyConfigTitle `json:"story_title"`
}

type storyConfigTitle struct {
	Id         int    `json:"id"`
	Story_name string `json:"story_name"`
	Level      int    `json:"level"`
	Stage      int    `json:"stage"`
}

func PraseStoryTitleConfig() storyConfig {
	content, err := os.ReadFile("./teaappconfig/story.json")
	if err != nil {
		LogRus.Errorf("Read config file err %v", err)
	}
	var contentStrct storyConfig
	err = json.Unmarshal(content, &contentStrct)
	if err != nil {
		LogRus.Errorf("Prase config file err %v", err)
	}
	return contentStrct
}

var (
	storyStat        sync.Map
	storyTitleConfig storyConfig
	wgstat           sync.WaitGroup
)

type UserHome struct {
	ID   int                    `bson:"_id"`
	Data map[string]interface{} `bson:"data"`
}
type UserInfo struct {
	Id                int    `bson:"_id"`
	Status            int    `bson:"status"`
	CreateTime        int    `bson:"create_time"`
	LastLoginTime     int    `bson:"last_login_time"`
	TimeZone          int    `bson:"timezone"`
	SaveMoudleTime    int    `bson:"sava_module_time"`
	DeviceToken       string `bson:"devicetoken"`
	DeviceTokenUpTime int    `bson:"devicetoken_uptime"`
}

/* 统计好友数量 */
func FrindsStat() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	userTable := teaapp.Cfg.Section("mongo").Key("user_m_friendgift").String()

	friendQuery := teaapp.Session.DB(database_1).C(userTable).Find(bson.M{}).Iter()
	var friendGift FriendGift
	//统计每个用户拥有的好友数量，平均好友数量，好友数量中位数
	var friendNumSlice = make([]int, 0)

	for friendQuery.Next(&friendGift) {
		friendNum := len(friendGift.FriendIdList)
		friendNumSlice = append(friendNumSlice, friendNum)

	}
	//平均数
	var friendNumTotalCount int
	for _, count := range friendNumSlice {
		friendNumTotalCount += count
	}
	var totalPlayers = len(friendNumSlice)

	averageNum := float64(friendNumTotalCount) / float64(totalPlayers)
	averageNumStr := strconv.FormatFloat(averageNum, 'f', 2, 64)
	//中位数
	sort.Ints(friendNumSlice)
	var midNum int
	if len(friendNumSlice)%2 == 0 {
		midNum = (friendNumSlice[len(friendNumSlice)/2] + friendNumSlice[len(friendNumSlice)/2-1]) / 2
	} else {
		midNum = friendNumSlice[len(friendNumSlice)/2]
	}
	midNumStr := strconv.Itoa(midNum)
	//输出
	LogRus.Println("averageNum:", averageNumStr, "midNum:", midNumStr)
}

// 拥有至少2张桌子
func HaveTwoDesk() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermhome := teaapp.Cfg.Section("mongo").Key("user_m_home").String()
	userTable := teaapp.Cfg.Section("mongo").Key("user_table").String()

	query := teaapp.Session.DB(database_1).C(userTable).Find(bson.M{"last_login_time": bson.M{"$gt": 1700648567}})

	var users []UserInfo
	err := query.All(&users)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	var totalDeskUser int
	for _, user := range users {
		var userId = user.Id

		userquery := teaapp.Session.DB(database_1).C(usermhome).Find(bson.M{"_id": userId}).Iter()

		var userHome UserHome

		for userquery.Next(&userHome) {
			if userHome.ID <= 100 {
				continue
			}
			decorMap, isType := userHome.Data["userHomeDecorMap"].(map[string]interface{})
			if !isType {
				continue
			}
			var deskCount int
			for k, decor := range decorMap {
				intK, _ := strconv.Atoi(k)
				if intK <= 5 {
					if decor.(map[string]interface{})["hasBuy"] == true {
						deskCount++
					}
				}
			}
			if deskCount >= 2 {
				totalDeskUser++
				println("uid:", userHome.ID)
			}
		}
	}
	fmt.Println("totalDeskUser:", totalDeskUser)
}

// 故事书领奖解锁统计
func StoryBookStat() {
	Init()
	storyTitleConfig = PraseStoryTitleConfig()

	for i := 0; i < 11; i++ {
		wgstat.Add(1)
		start := i * 10000
		end := (i + 1) * 10000
		if i == 0 {
			start = 101
		}
		go StoryQuery(start, end)
	}
	wgstat.Wait()

	excelTitel := []string{"id", "mapid", "level", "rewarded", "unlock", "rate"}

	//处理排序
	storyList := []*storybookHelp{}
	storyStat.Range(func(key, value any) bool {
		item := value.(*storybookHelp)
		item.Sort = item.MapId*1000 + item.Level*10
		storyList = append(storyList, item)
		return true
	})
	sort.Slice(storyList, func(i int, j int) bool {
		return storyList[i].Sort < storyList[j].Sort
	})
	excelContent := [][]string{}
	for i, storyItem := range storyList {
		rate := float64(storyItem.RwardTotal) / float64(storyItem.UnLoackTotal)
		ratestr := getFloatStr(rate)
		excelContent = append(excelContent, []string{
			strconv.Itoa(i),
			strconv.Itoa(storyItem.MapId),
			strconv.Itoa(storyItem.Level),
			strconv.Itoa(storyItem.RwardTotal),
			strconv.Itoa(storyItem.UnLoackTotal),
			ratestr,
		})
	}
	outputXlsx("storystat.xlsx", "story", excelTitel, excelContent)
}

func StoryQuery(idStart int, idEnd int) {
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_m_data").String()
	storybook_table := teaapp.Cfg.Section("mongo").Key("storybook_table").String()
	userquery := teaapp.Session.DB(database_1).C(storybook_table).Find(bson.M{"_id": bson.M{"$gte": idStart, "$lt": idEnd}}).Iter()

	var storyBook storybook

	for userquery.Next(&storyBook) {
		LogRus.Infof("uid:%d", storyBook.Uid)
		userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"_id": storyBook.Uid})
		var user userModule
		err := userquery.One(&user)
		if err != nil {
			LogRus.Errorf("Query user:%d  err:%v", storyBook.Uid, err)
			continue
		}
		userData, ok := user.Data["userData"].(map[string]interface{})
		if !ok {
			LogRus.Errorf("Assert userData err")
			continue
		}
		_, ok = userData["maxStageId"]
		if !ok {
			continue
		}
		_, ok = userData["maxRestaurantId"]
		if !ok {
			continue
		}
		MaxStageId, ok := userData["maxStageId"].(int)
		MaxRestaurantId := userData["maxRestaurantId"].(int)
		LogRus.Infof("maxstageId %d,maxRestaurantId %d", MaxStageId, MaxRestaurantId)

		//处理已经领奖的数量
		if rewardStat, ok := storyBook.Data["readRewardState"]; ok {
			rewardStatArr, ok := rewardStat.([]interface{})
			if !ok {
				continue
			}
			for i, v := range rewardStatArr {
				rewardItem, ok := v.(map[string]interface{})
				if !ok {
					continue
				}
				//如果需要进一步统计到餐厅第几关
				for level, _ := range rewardItem {
					levelInt, _ := strconv.Atoi(level)
					key := strconv.Itoa(i) + "_" + level
					// istr, _ := strconv.Atoi(i)
					storyItem, _ := storyStat.LoadOrStore(key, &storybookHelp{
						MapId:      i,
						Level:      levelInt,
						RwardTotal: 0,
					})
					st := storyItem.(*storybookHelp)
					st.RwardTotal++

				}
			}
		}
		//处理已经解锁的数量
		for _, configItem := range storyTitleConfig.StoryTitle {
			key := strconv.Itoa(configItem.Stage) + "_" + strconv.Itoa(configItem.Level)
			storyItem, _ := storyStat.LoadOrStore(key, &storybookHelp{
				MapId:      configItem.Stage,
				Level:      configItem.Level,
				RwardTotal: 0,
			})
			st := storyItem.(*storybookHelp)
			if MaxStageId > configItem.Stage {
				st.UnLoackTotal++
				continue
			}
			if MaxStageId == configItem.Stage {
				if MaxRestaurantId >= configItem.Level {
					st.UnLoackTotal++
				}
			}
		}
	}
	wgstat.Done()
	return
}

func StageStat() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_m_data").String()
	usermresturant := teaapp.Cfg.Section("mongo").Key("user_m_resturant").String()
	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"_id": bson.M{"$gt": 100}})
	var user userModule
	var count int
	userqueryIter := userquery.Iter()
	for userqueryIter.Next(&user) {
		userData, ok := user.Data["userData"].(map[string]interface{})
		if !ok {
			LogRus.Errorf("Assert userData err")
			continue
		}
		if maxid, ok := userData["maxRestaurantId"]; ok {
			MaxRestaurantId := maxid.(int)
			if MaxRestaurantId > 1 {
				fmt.Println("U:", user.ID)
				query := teaapp.Session.DB(database_1).C(usermresturant).Find(bson.M{"_id": user.ID})
				var resturantRes teaapp.Resturant
				err := query.One(&resturantRes)
				if err != nil {
					LogRus.Printf("err: %v", err)
				}

				if stage1, ok := resturantRes.Data["userRestaurantMap"]["1"]; ok {
					if stageLevel, levelok := stage1["stageLevel"]; levelok {
						level := stageLevel.(int)
						if level == 1 {
							count++
							fmt.Println("BUG USER ID:", user.ID)
							continue
						}
					}

				}
			}
		}
	}
	fmt.Println("count:", count)
}

func StarStat() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_m_data").String()
	usermresturant := teaapp.Cfg.Section("mongo").Key("user_m_resturant").String()

	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"_id": bson.M{"$gt": 100}})
	// userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{"_id": 18307})
	var user userModule
	userqueryIter := userquery.Iter()

	//平均数
	var starstat = make(map[string]*StarMap)
	//中位数
	var starmidstat = make(map[string]*[]int)
	for userqueryIter.Next(&user) {

		LogRus.Infof("uid:%d", user.ID)
		query := teaapp.Session.DB(database_1).C(usermresturant).Find(bson.M{"_id": user.ID})
		var resturantRes teaapp.Resturant
		err := query.One(&resturantRes)
		if err != nil {
			LogRus.Printf("err: %v", err)
		}

		for mid, mapdata := range resturantRes.Data["userRestaurantMap"] {
			midInt, _ := strconv.Atoi(mid)
			if midInt > 200 {
				continue
			}
			midstar, ok := starmidstat[mid]
			if !ok {
				midstar = &[]int{}
				starmidstat[mid] = midstar
			}
			starmap, ok := starstat[mid]
			if !ok {

				starmap = &StarMap{
					Mid:         midInt,
					StartTotal:  0,
					PlayerTotal: 0,
				}
				starstat[mid] = starmap
			}

			if stageLevel, ok := mapdata["stageLevel"]; ok {
				if intStageLevel, w := stageLevel.(int); w {
					if intStageLevel > 60 {
						total := 0
						if starmapData, ok := mapdata["levelStarMap"].(map[string]interface{}); ok {
							for _, num := range starmapData {
								total += num.(int)
							}
						}
						if midInt == 1 {
							if total < 40 {
								total = 40
							}
						} else {
							if total < 60 {
								total = 60
							}
						}
						*midstar = append(*midstar, total)
						starmap.StartTotal += total
						starmap.PlayerTotal += 1
					}
				}
			}

		}
	}

	for i := range starmidstat {
		sort.Ints(*starmidstat[i])
	}

	starSlice := []*StarMap{}
	for _, starmapItem := range starstat {
		starSlice = append(starSlice, starmapItem)
	}
	sort.Slice(starSlice, func(i int, j int) bool {
		return starSlice[i].Mid < starSlice[j].Mid
	})
	outdataList := [][]string{}
	for _, item := range starSlice {
		if item.PlayerTotal > 0 {
			item.AverageNum = int(math.Floor(float64(item.StartTotal) / float64(item.PlayerTotal)))
		} else {
			item.AverageNum = 0
		}
		midstr := strconv.Itoa(item.Mid)
		//中位数
		midnum := 0
		midstat, ok := starmidstat[midstr]
		if ok && len(*midstat) > 0 {
			len := math.Floor(float64(len(*midstat) / 2))
			midnum = (*midstat)[int(len)]
		}

		outItem := []string{
			midstr,
			strconv.Itoa(item.AverageNum),
			strconv.Itoa(midnum),
		}
		outdataList = append(outdataList, outItem)
	}
	title := []string{"mapId", "average", "midnum"}
	outputXlsx("starstat.xlsx", "star", title, outdataList)
	LogRus.Infoln("succ")
}

func ClothStat() {
	Init()
	cloth := make(map[int]int)
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	// database_0 := teaapp.Cfg.Section("mongo").Key("database_0").String()
	userskinTable := teaapp.Cfg.Section("mongo").Key("skin_table").String()

	userQuery := teaapp.Session.DB(database_1).C(userskinTable).Find(bson.M{})
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
				cloth[valint]++
			}

		}
	}
	//format
	list := [][]string{}
	for id, num := range cloth {
		list = append(list, []string{strconv.Itoa(id), strconv.Itoa(num)})
	}
	title := []string{"id", "count"}
	//output
	outputXlsx("clothstat.xlsx", "cloth", title, list)
	LogRus.Info("succ!")
}

var validDays = 30
var totalUser = 0
var activeUser = 0

func DoMachineStat() {
	Init()
	//总用户
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	userTable := teaapp.Cfg.Section("mongo").Key("user_table").String()
	total, err := teaapp.Session.DB(database_1).C(userTable).Count()
	if err != nil {
		LogRus.Errorln("query total user count err:%v", err)
	}
	totalUser = total

	validTime := time.Now().Unix() - int64(validDays)*24*3600
	activeuser, err := teaapp.Session.DB(database_1).C(userTable).Find(bson.M{"last_login_time": bson.M{"$gt": validTime}}).Count()
	if err != nil {
		LogRus.Errorln("query active user err:%v", err)
	}
	activeUser = activeuser

	var i int = 1
	for i < 16 {
		MachineStat(strconv.Itoa(i))
		i++
	}
}

// 用户设备
func MachineStat(mapid string) {
	machineStar := make(map[string]map[string]int)
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	userTable := teaapp.Cfg.Section("mongo").Key("user_table").String()

	// database_0 := teaapp.Cfg.Section("mongo").Key("database_0").String()
	userResturant := teaapp.Cfg.Section("mongo").Key("user_m_resturant").String()
	userQuery := teaapp.Session.DB(database_1).C(userResturant).Find(bson.M{})
	userQuery.Select(bson.M{"_id": 1, "data.userRestaurantMap": 1})
	queryIter := userQuery.Iter()
	var userRestuarant userRestaurantModule
	now := time.Now().Unix()
	for queryIter.Next(&userRestuarant) {
		//判断是否已经流失
		var userInfo map[string]interface{}
		query := teaapp.Session.DB(database_1).C(userTable).Find(bson.M{"_id": userRestuarant.ID})
		query.Select(bson.M{"_id": 1, "last_login_time": 1})
		query.One(&userInfo)
		lastLoginTime := userInfo["last_login_time"].(int)
		if lastLoginTime+validDays*24*3600 < int(now) {
			continue
		}
		if userRestuarant.ID > 0 {
			rdatamap, ok := userRestuarant.Data["userRestaurantMap"].(map[string]interface{})
			if !ok {
				continue
			}
			for rid, rdata := range rdatamap {
				//不是本店铺的
				if rid != mapid {
					continue
				}
				ridInt, err := strconv.Atoi(rid)
				if err != nil {
					LogRus.Debugf("atoi err %v", err)
					continue
				}

				// 判断地图id
				if ridInt > 200 {
					continue
				}

				machineMap, ok := rdata.(map[string]interface{})["userFacilityMap"]
				if ok {
					for mid, mdata := range machineMap.(map[string]interface{}) {
						mtier := mdata.(map[string]interface{})["facilityTier"]

						_, mok := machineStar[mid]
						if !mok {
							machineStar[mid] = make(map[string]int)
						}

						mtierInt, mitok := mtier.(int)
						if !mitok {
							continue
						}
						if mtierInt > 10 {
							LogRus.Errorf("user data err,userId:%d,mapId:%d,machineId:%s,machinelevel:%d", userRestuarant.ID, ridInt, mid, mtierInt)
							continue
						}
						mtierstr := strconv.Itoa(mtierInt)
						if _, mtok := machineStar[mid][mtierstr]; mtok {

							machineStar[mid][mtierstr]++
						} else {
							machineStar[mid][mtierstr] = 1
						}

						// fmt.Printf("uid:%d,rid:%s,mid:%s,mtier:%v \n", userRestuarant.ID, rid, mid, mtier)
					}
				}
			}
		}
	}
	out(mapid, machineStar)
	LogRus.Info("map %s output succ!", mapid)
}

type dataSortHelp struct {
	MachineId string
	StarLevel string
	Count     int
	Sort      int
}

func out(mapid string, data map[string]map[string]int) {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	file, err = xlsx.OpenFile("stat.xlsx")
	if err != nil {
		file = xlsx.NewFile()
	}
	sheet, err = file.AddSheet(mapid)
	if err != nil {
		LogRus.Errorf("newsheet err:%v", err)
	}
	var titleRow []string = []string{"id", "star", "count", "per_active", "per_all"}
	row = sheet.AddRow()

	for _, val := range titleRow {
		cell = row.AddCell()
		cell.Value = val
	}
	//sort map
	midslice := make([]dataSortHelp, 0)
	// for mid, _ := range data {
	// 	midint, _ := strconv.Atoi(mid)
	// 	midslice = append(midslice, midint)
	// }
	// sort.Slice(midslice, func(i, j int) bool {
	// 	return i > j
	// })

	for mid, stardata := range data {
		for starlevel, count := range stardata {
			starInt, _ := strconv.Atoi(starlevel)
			if starInt > 10 {
				continue
			}
			midint, _ := strconv.Atoi(mid)
			sort := midint*1000 + starInt*10
			midslice = append(midslice, dataSortHelp{
				MachineId: mid,
				StarLevel: starlevel,
				Count:     count,
				Sort:      sort,
			})
		}
	}
	sort.Slice(midslice, func(i, j int) bool {
		return midslice[i].Sort < midslice[j].Sort
	})

	for _, data := range midslice {
		newrow := sheet.AddRow()
		newrow.AddCell().Value = data.MachineId
		newrow.AddCell().Value = data.StarLevel

		countstr := strconv.Itoa(data.Count)
		newrow.AddCell().Value = countstr

		activePer := float64(data.Count) / float64(activeUser)
		activePerStr := getFloatStr(activePer)
		newrow.AddCell().Value = activePerStr

		totalPer := float64(data.Count) / float64(totalUser)
		totalPerStr := getFloatStr(totalPer)
		newrow.AddCell().Value = totalPerStr
	}

	err = file.Save("stat.xlsx")
	if err != nil {
		LogRus.Errorf("save xlsx err:%v", err)
	}
}

func getFloatStr(value float64) string {
	activePera, err := strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	if err != nil {
		activePera = 0.00
		LogRus.Errorln("parseFloat err:%v", err)
	}
	activePera = float64(activePera) * float64(100)
	return fmt.Sprintf("%.2f%%", activePera)
}

func outputXlsx(fileName string, sheetName string, title []string, data [][]string) {

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error
	if _, err := os.Stat(fileName); err == nil {
		file, err = xlsx.OpenFile(fileName)
		if err != nil {
			LogRus.Errorf("open file err:%v", err)
		}
	} else {
		file = xlsx.NewFile()
	}
	sheet, err = file.AddSheet(sheetName)
	if err != nil {
		LogRus.Errorf("add newsheet err:%v", err)
	}
	row = sheet.AddRow()
	for _, val := range title {
		cell = row.AddCell()
		cell.SetString(val)
	}
	for _, v := range data {
		newrow := sheet.AddRow()
		for _, vv := range v {
			newrow.AddCell().SetString(vv)
		}
	}
	err = file.Save(fileName)
	if err != nil {
		LogRus.Errorf("save %s err:%v", fileName, err)
	}
}
