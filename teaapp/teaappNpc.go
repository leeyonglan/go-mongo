package teaapp

import (
	"context"
	"fmt"
	"math"
	"notification"
	"os"
	"strconv"
	"time"

	"github.com/leeyonglan/go-mongo/mongo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gopkg.in/mgo.v2/bson"
)

type resturant struct {
	Id   int                                          `bson:"_id"`
	Data map[string]map[string]map[string]interface{} `bson:"data"`
}

// notification 类型
const (
	NOTI_FRIENDGIFT   = "friendgift"
	NOTI_POWER        = "powerday"
	NOTI_RANK_REWARD  = "rankreward"
	NOTI_RANK_RESTART = "rankrestart"
	NOTI_NEWVERSION   = "newversion"
	NOTI_NEWCONTENT   = "newcontent"
	NOTI_SYSTEMINFO   = "sysinfo"
	NOTI_CALLBACK     = "callback"
)

//user 表
type userInfo struct {
	Id                int    `bson:"_id"`
	Status            int    `bson:"status"`
	CreateTime        int    `bson:"create_time"`
	LastLoginTime     int    `bson:"last_login_time"`
	TimeZone          int    `bson:"timezone"`
	SaveMoudleTime    int    `bson:"sava_module_time"`
	DeviceToken       string `bson:"devicetoken"`
	DeviceTokenUpTime int    `bson:"devicetoken_uptime"`
}

//notification 表
type notificationInfo struct {
	Id            int            `bson:"_id"`
	FriendGift    map[string]int `bson:"friendgift"` //是否需要发送好友送礼物通知
	LastSendTime  int            `bson:"lastsendtime"`
	Version       string         `bson:"newversion"`
	RankReward    int            `bson:"rankreward"`
	RankRestart   string         `bson:"rankrestart"`
	PowerDay      string         `bson:"powerday"`
	PowerTimes    int            `bson:"powertimes"`
	CallBackTimes int            `bson:"callback"`
}

// version 表
type versionInfo struct {
	Name       string `bson:"name"`        //版本号 “1.3.56”
	UpdateTime int    `bson:"update_time"` //更新时间
}

type acitivityInfo struct {
	AcitivityId string `bson:"hd_id"`
	Name        string `bson:"name"`
	StartTime   int    `bson:"b_time"`
	HasReward   int8   `bson:"is_send_reward"`
	RewardTime  int    `bson:"is_send_reward_time"`
}
type usertoken struct {
	Power     int `bson:"power"`
	PowerTime int `bson:"powerStartTime"`
}
type propsdata struct {
	UserTokenData usertoken `bson:"userTokenData"`
}
type propsInfo struct {
	Id   string    `bson:"_id"`
	Data propsdata `bson:"data"`
}

func NotiUser() {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file:%v", err)
		os.Exit(1)
	}

	now := int32(time.Now().Unix())
	host := cfg.Section("mongo").Key("host").String()
	port := cfg.Section("mongo").Key("port").String()
	user := cfg.Section("mongo").Key("user").String()
	password := cfg.Section("mongo").Key("password").String()
	mongo := getMongo(host, port, user, password)

	database_1 := cfg.Section("mongo").Key("database_1").String()
	database_0 := cfg.Section("mongo").Key("database_0").String()

	userTable := cfg.Section("mongo").Key("user_table").String()
	// userDataTable := cfg.Section("mongo").Key("user_m_data").String()
	versionTable := cfg.Section("mongo").Key("version_table").String()
	notificationTable := cfg.Section("mongo").Key("notifaction_table").String()
	activityTable := cfg.Section("mongo").Key("activity_table").String()
	propsTabel := cfg.Section("mongo").Key("props_table").String()

	session := mongo.ConSession
	defer session.Close()

	query := session.DB(database_1).C(userTable).Find(bson.M{"devicetoken": bson.M{"$exists": true}})
	//版本信息
	verQuery := session.DB(database_0).C(versionTable).Find(bson.M{})
	var versioninfo versionInfo
	err = verQuery.One(&versioninfo)
	if err != nil {
		log.Debugf("query version table err %v", err)
	}
	//排行榜信息
	acitivityQuery := session.DB(database_0).C(activityTable).Find(bson.M{"name": "rank"})
	var activityinfo acitivityInfo
	err = acitivityQuery.One(&activityinfo)
	if err != nil {
		log.Debugf("query activity table err %v", err)
	}

	var users []userInfo
	err = query.All(&users)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	log.Info("total:", len(users))
	for _, userinfo := range users {
		// 判断所在时区现在是否可以发 中午12-2pm, 晚上6-8pm, 10-12am
		log.Info(userinfo.Id, " start send")
		userTimeZone := userinfo.TimeZone
		timeTime := GetLocalTime(userTimeZone)
		if isForbiddenTime(timeTime) {
			log.Info(userinfo.Id, " in forbiddenTime")
			continue
		}
		if !isInTime(timeTime) {
			log.Info(userinfo.Id, " not in notification time")
			// continue
		}
		notiQuery := session.DB(database_1).C(notificationTable).Find(bson.M{"_id": userinfo.Id})
		var notiItem notificationInfo
		err = notiQuery.One(&notiItem)
		if err != nil {
			log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("query notification err:%v", err)
			notiItem = notificationInfo{
				Id:           userinfo.Id,
				Version:      "",
				LastSendTime: 0,
			}
		}
		// 判断是否间隔2个小时
		lastSendTime := notiItem.LastSendTime
		if lastSendTime != 0 && lastSendTime+2*3600 <= int(now) {
			log.Info(userinfo.Id, " in 2 hours cd")
			// continue
		}

		// 优先级:版本更新的时候 无限关卡排行榜奖励 无限关卡排行榜重新开启 提示新内容可玩  邮箱收到其他信息的时候 体力值恢复满的时候  好友送礼物的时候  召回功能

		//判断是否需要发新版本通知
		if isTrue := isNeedNewVersionNoti(versioninfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_NEWVERSION, userinfo.DeviceToken)
			if err == nil {
				_, err = session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{NOTI_NEWVERSION: versioninfo.Name, "lastsendtime": now}})
				if err != nil {
					log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("update notification %s err:", NOTI_NEWVERSION, err)
				}
			}
			continue
		}
		//无限关卡排行榜奖励通知
		if isTrue := isNeedRankRewardNoti(activityinfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_RANK_REWARD, userinfo.DeviceToken)
			if err == nil {
				_, err = session.DB(database_1).C(notificationTable).Upsert(
					bson.M{"_id": userinfo.Id},
					bson.M{"$unset": bson.M{NOTI_RANK_REWARD: 1}, "$set": bson.M{"lastsendtime": now}},
				)
				if err != nil {
					log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_RANK_REWARD, err)
				}
			}
			continue
		}
		//无限关卡排行榜开始通知
		// userDataQuery := session.DB(database_1).C(userDataTable).Find(bson.M{"_id": userinfo.Id})

		if isTrue := isNeedRankStartNoti(activityinfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_RANK_RESTART, userinfo.DeviceToken)
			if err == nil {
				_, err = session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{NOTI_RANK_RESTART: activityinfo.AcitivityId, "lastsendtime": now}})
				if err != nil {
					log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_RANK_RESTART, err)
				}
			}
			continue
		}
		//体力恢复满通知
		propsQuery := session.DB(database_1).C(propsTabel).Find(bson.M{"_id": userinfo.Id})
		propsQuery.Select(bson.M{"_id": 1, "data.userTokenData.power": 1})
		var propsinfo propsInfo
		err = propsQuery.One(&propsinfo)
		if err != nil {
			log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("query prorps table err :%v", err)
		}
		if isTrue := isNeedPowerNoti(propsinfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_POWER, userinfo.DeviceToken)
			if err == nil {
				day := time.Now().Format("20060102")
				times := notiItem.PowerTimes
				if notiItem.PowerDay != day {
					times = 1
				} else {
					times++
				}
				_, err = session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{NOTI_POWER: day, "powertimes": times, "lastsendtime": now}})
				if err != nil {
					log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_POWER, err)
				}
			}
			continue
		}

		//好友送礼物
		if _, ok := notiItem.FriendGift["time"]; ok {
			err = sendNotification(userinfo.Id, NOTI_FRIENDGIFT, userinfo.DeviceToken)
			if err == nil {
				_, err = session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$unset": bson.M{NOTI_FRIENDGIFT: 1, "lastsendtime": now}})
				if err != nil {
					log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_FRIENDGIFT, err)
				}
			}
			continue
		}

		//召回
		if isTrue := isNeedCallNoti(userinfo, notiItem); isTrue {
			callNotiType := NOTI_CALLBACK + "_" + strconv.Itoa(notiItem.CallBackTimes+1)
			err = sendNotification(userinfo.Id, callNotiType, userinfo.DeviceToken)
			if err == nil {
				_, err = session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$inc": bson.M{NOTI_CALLBACK: 1}, "$set": bson.M{"lastsendtime": now}})
				if err != nil {
					log.WithFields(log.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_CALLBACK, err)
				}
			}
			continue
		}
	}
	log.Info("succ!")
}

// 是否需要版本通知
func isNeedNewVersionNoti(version versionInfo, notification notificationInfo) bool {
	now := time.Now().Unix()
	notiVersion := notification.Version
	if now > int64(version.UpdateTime)+24*3600 {
		// return false
	}
	if notiVersion == version.Name {
		return false
	}
	return true
}

//是否需要排行榜奖励通知
func isNeedRankRewardNoti(activity acitivityInfo, notification notificationInfo) bool {
	if notification.RankReward != 1 {
		return false
	}
	now := time.Now().Unix()
	if now > int64(activity.RewardTime)+24*3600 {
		return false
	}
	return true
}

// 是否需要排行榜开始通知
func isNeedRankStartNoti(activity acitivityInfo, notification notificationInfo) bool {
	if notification.RankRestart == activity.AcitivityId {
		return false
	}
	now := time.Now().Unix()
	if now > int64(activity.StartTime)+24*3600 {
		return false
	}
	return true
}

func isNeedPowerNoti(propsinfo propsInfo, notification notificationInfo) bool {
	if propsinfo.Data.UserTokenData.Power >= 5 {
		return false
	}
	powerday := notification.PowerDay
	now := time.Now().Unix()
	day := time.Now().Format("20060102")
	if day == powerday && notification.PowerTimes >= 3 {
		return false
	}
	powerRestoreTime := propsinfo.Data.UserTokenData.PowerTime
	powerRestoreCount := math.Floor(float64((now-int64(powerRestoreTime))/900)) + 1
	if propsinfo.Data.UserTokenData.Power+int(powerRestoreCount) < 5 {
		return false
	}
	return true
}

var callTimesAarray [4]int = [4]int{1, 3, 5, 7}

func isNeedCallNoti(userinfo userInfo, notification notificationInfo) bool {
	now := time.Now().Unix()
	lastLoginTime := userinfo.LastLoginTime
	hasCallTimes := notification.CallBackTimes
	if hasCallTimes == 4 {
		return false
	}
	intervalDay := callTimesAarray[hasCallTimes]
	intervalTime := intervalDay * 24 * 3600
	if now <= int64(lastLoginTime)+int64(intervalTime) {
		return false
	}
	return true
}

var notiList map[string][]string = make(map[string][]string)

func sendNotification(uid int, notiType string, deviceToken string) (err error) {
	log.Info(deviceToken, " send ", notiType)
	err = notification.DoPush(notiType, deviceToken)
	return
	// notiTypeList, ok := notiList[notiType]
	// if ok {
	// 	if notiTypeList == nil {
	// 		notiTypeList = make([]string, 100)
	// 	}
	// 	notiTypeList = append(notiTypeList, deviceToken)
	// 	notiList[notiType] = notiTypeList
	// }
}

func batchSend() {
	// for notitype, tokens := range notiList {

	// }
}

func UpdateProp() {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file:%v", err)
		os.Exit(1)
	}
	host := cfg.Section("mongo").Key("host").String()
	port := cfg.Section("mongo").Key("port").String()
	user := cfg.Section("mongo").Key("user").String()
	password := cfg.Section("mongo").Key("password").String()
	mongo := getMongo(host, port, user, password)

	database := cfg.Section("mongo").Key("database").String()
	usermdata := cfg.Section("mongo").Key("user_m_data").String()
	usermresturant := cfg.Section("mongo").Key("user_m_game_props").String()

	session := mongo.ConSession
	query := session.DB(database).C(usermresturant).Find(bson.M{"data.infiniteEndTimeMap": bson.M{"$exists": true}})
	var resturantRes []resturant
	err = query.All(&resturantRes)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	for _, value := range resturantRes {
		totalLevel := 0
		for _, mapdata := range value.Data["infiniteEndTimeMap"] {
			if stageLevel, ok := mapdata["stageLevel"]; ok {
				if intStageLevel, w := stageLevel.(int); w {
					totalLevel += intStageLevel
				}

			}
		}
		fmt.Printf("Id: %d, totalLevel:%d \n", value.Id, totalLevel)

		err = session.DB(database).C(usermdata).Update(bson.M{"_id": value.Id}, bson.M{"$set": bson.M{"data.userData.starcount": totalLevel}})
		if err != nil {
			fmt.Printf("update db error %v \n", err)
		} else {
			fmt.Printf("%d update succ starcount:%d \n", value.Id, totalLevel)
		}
	}
}

func UpdateNpcStarTotal() {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Printf("Fail to read file:%v", err)
		os.Exit(1)
	}
	host := cfg.Section("mongo").Key("host").String()
	port := cfg.Section("mongo").Key("port").String()
	user := cfg.Section("mongo").Key("user").String()
	password := cfg.Section("mongo").Key("password").String()
	mongo := getMongo(host, port, user, password)

	database := cfg.Section("mongo").Key("database").String()
	usermdata := cfg.Section("mongo").Key("user_m_data").String()
	usermresturant := cfg.Section("mongo").Key("user_m_resturant").String()

	session := mongo.ConSession
	query := session.DB(database).C(usermresturant).Find(bson.M{"_id": bson.M{"$lt": 100}})
	var resturantRes []resturant
	err = query.All(&resturantRes)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	for _, value := range resturantRes {
		totalLevel := 0
		for _, mapdata := range value.Data["userRestaurantMap"] {
			if stageLevel, ok := mapdata["stageLevel"]; ok {
				if intStageLevel, w := stageLevel.(int); w {
					totalLevel += intStageLevel
				}

			}
		}
		fmt.Printf("Id: %d, totalLevel:%d \n", value.Id, totalLevel)

		err = session.DB(database).C(usermdata).Update(bson.M{"_id": value.Id}, bson.M{"$set": bson.M{"data.userData.starcount": totalLevel}})
		if err != nil {
			fmt.Printf("update db error %v \n", err)
		} else {
			fmt.Printf("%d update succ starcount:%d \n", value.Id, totalLevel)
		}
	}
}

func getMongo(host string, port string, user string, password string) *mongo.Mongo {
	var dbconf []*mongo.DbConf
	dbconf = append(dbconf, &mongo.DbConf{Host: host, Port: port, User: user, Pass: password})

	var dbconfs = &mongo.DbConnConf{
		Confs: dbconf,
	}
	dbconfs.Init()

	// var conn = dbconfs.GetSSLCon(context.TODO())
	var conn = dbconfs.GetConn(context.TODO())
	return &mongo.Mongo{ConSession: conn}
}
