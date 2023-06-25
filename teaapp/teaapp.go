package teaapp

import (
	"flag"
	"fmt"
	"math"
	"mongo"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const rootPEM = `-----BEGIN CERTIFICATE-----
MIIECTCCAvGgAwIBAgICEAAwDQYJKoZIhvcNAQELBQAwgZUxCzAJBgNVBAYTAlVT
MRAwDgYDVQQHDAdTZWF0dGxlMRMwEQYDVQQIDApXYXNoaW5ndG9uMSIwIAYDVQQK
DBlBbWF6b24gV2ViIFNlcnZpY2VzLCBJbmMuMRMwEQYDVQQLDApBbWF6b24gUkRT
MSYwJAYDVQQDDB1BbWF6b24gUkRTIGV1LXNvdXRoLTEgUm9vdCBDQTAeFw0xOTEw
MzAyMDIxMzBaFw0yNDEwMzAyMDIxMzBaMIGQMQswCQYDVQQGEwJVUzETMBEGA1UE
CAwKV2FzaGluZ3RvbjEQMA4GA1UEBwwHU2VhdHRsZTEiMCAGA1UECgwZQW1hem9u
IFdlYiBTZXJ2aWNlcywgSW5jLjETMBEGA1UECwwKQW1hem9uIFJEUzEhMB8GA1UE
AwwYQW1hem9uIFJEUyBldS1zb3V0aC0xIENBMIIBIjANBgkqhkiG9w0BAQEFAAOC
AQ8AMIIBCgKCAQEAtEyjYcajx6xImJn8Vz1zjdmL4ANPgQXwF7+tF7xccmNAZETb
bzb3I9i5fZlmrRaVznX+9biXVaGxYzIUIR3huQ3Q283KsDYnVuGa3mk690vhvJbB
QIPgKa5mVwJppnuJm78KqaSpi0vxyCPe3h8h6LLFawVyWrYNZ4okli1/U582eef8
RzJp/Ear3KgHOLIiCdPDF0rjOdCG1MOlDLixVnPn9IYOciqO+VivXBg+jtfc5J+L
AaPm0/Yx4uELt1tkbWkm4BvTU/gBOODnYziITZM0l6Fgwvbwgq5duAtKW+h031lC
37rEvrclqcp4wrsUYcLAWX79ZyKIlRxcAdvEhQIDAQABo2YwZDAOBgNVHQ8BAf8E
BAMCAQYwEgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQU7zPyc0azQxnBCe7D
b9KAadH1QSEwHwYDVR0jBBgwFoAUFBAFcgJe/BBuZiGeZ8STfpkgRYQwDQYJKoZI
hvcNAQELBQADggEBAFGaNiYxg7yC/xauXPlaqLCtwbm2dKyK9nIFbF/7be8mk7Q3
MOA0of1vGHPLVQLr6bJJpD9MAbUcm4cPAwWaxwcNpxOjYOFDaq10PCK4eRAxZWwF
NJRIRmGsl8NEsMNTMCy8X+Kyw5EzH4vWFl5Uf2bGKOeFg0zt43jWQVOX6C+aL3Cd
pRS5MhmYpxMG8irrNOxf4NVFE2zpJOCm3bn0STLhkDcV/ww4zMzObTJhiIb5wSWn
EXKKWhUXuRt7A2y1KJtXpTbSRHQxE++69Go1tWhXtRiULCJtf7wF2Ksm0RR/AdXT
1uR1vKyH5KBJPX3ppYkQDukoHTFR0CpB+G84NLo=
-----END CERTIFICATE-----`

type Resturant struct {
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
	Id             int            `bson:"_id"`
	FriendGift     map[string]int `bson:"friendgift"` //是否需要发送好友送礼物通知
	LastSendTime   int            `bson:"lastsendtime"`
	Version        string         `bson:"newversion"`
	RankReward     int            `bson:"rankreward"`
	RankRestart    string         `bson:"rankrestart"`
	PowerDay       string         `bson:"powerday"`
	PowerTimes     int            `bson:"powertimes"`
	CallBackTimes  int            `bson:"callback"`
	LastAllMailId  string         `bson:"lastallmailid"`
	LastPartMailId string         `bson:"lastpartmailid"`
}

type mailListInfo struct {
	Id       bson.ObjectId `bson:"_id"`
	Type     int           `bson:"type"`
	Isuse    int           `bson:"is_use"`
	UidArray []int         `bson:"uid_array"`
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

var (
	wg              *sync.WaitGroup
	Cfg             *ini.File
	Sys             string
	Env             string
	FixtimezoneFlag string
	LogRus          *log.Logger
	Session         *mgo.Session
)

func InitConfig() {
	cfgv, err := ini.Load("./config.ini")
	Cfg = cfgv
	if err != nil {
		fmt.Printf("Fail to read file:%v", err)
		os.Exit(1)
	}
	LogRus = log.New()
	LogRus.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	LogRus.Formatter.(*logrus.TextFormatter).DisableTimestamp = false
	LogRus.Level = logrus.InfoLevel
	LogRus.Out = os.Stdout

}
func InitMongo() {
	var ca string = ""
	if Env == "pro" {
		ca = rootPEM
	}
	host := Cfg.Section("mongo").Key("host").String()
	port := Cfg.Section("mongo").Key("port").String()
	user := Cfg.Section("mongo").Key("user").String()
	password := Cfg.Section("mongo").Key("password").String()
	mongo := mongo.GetMongo(host, port, user, password, ca)
	Session = mongo.ConSession
}

func ReleaseMongo() {
	Session.Close()
}

func Init() {
	flag.StringVar(&FixtimezoneFlag, "fixTimeZone", "", "fixTimeZone")
	flag.StringVar(&Sys, "sys", "", "android")
	flag.StringVar(&Env, "env", "pro", "devOrPro")
	flag.Parse()

	if Sys == "" {
		fmt.Println("please input system")
		os.Exit(0)
	}
	defer ReleaseMongo()
	InitConfig()
	InitMongo()
	NotiUser()
}

func NotiUser() {
	if Sys == "" {
		fmt.Println("please input system")
		os.Exit(0)
	}

	wg = new(sync.WaitGroup)
	now := int32(time.Now().Unix())

	database_1 := Cfg.Section("mongo").Key("database_1").String()
	database_0 := Cfg.Section("mongo").Key("database_0").String()

	userTable := Cfg.Section("mongo").Key("user_table").String()
	// userDataTable := Cfg.Section("mongo").Key("user_m_data").String()
	versionTable := Cfg.Section("mongo").Key("version_table").String()
	notificationTable := Cfg.Section("mongo").Key("notifaction_table").String()
	activityTable := Cfg.Section("mongo").Key("activity_table").String()
	propsTabel := Cfg.Section("mongo").Key("props_table").String()
	initversion := Cfg.Section("init").Key("version").String()
	mailTable := Cfg.Section("mongo").Key("mail_table").String()

	if FixtimezoneFlag != "" {
		FixTimezoneToInt(Session, database_1, userTable)
		os.Exit(0)
	}

	query := Session.DB(database_1).C(userTable).Find(bson.M{"devicetoken": bson.M{"$exists": true, "$ne": ""}})
	//版本信息
	verQuery := Session.DB(database_0).C(versionTable).Find(bson.M{})
	var versioninfo versionInfo
	err := verQuery.One(&versioninfo)
	if err != nil {
		LogRus.Debugf("query version table err %v", err)
	}

	//系统邮件信息  从2023-5-19 开始，之前的不处理
	year := 2023
	month := time.May
	day := 19
	hour := 12
	minute := 30
	second := 0
	startTime := time.Date(year, month, day, hour, minute, second, 0, time.UTC).Unix()
	//只检查7天内的
	endTime := now - 7*24*3600
	if int32(startTime) > endTime {
		endTime = int32(startTime)
	}
	mailQuery := Session.DB(database_0).C(mailTable).Find(bson.M{"create_time": bson.M{"$gt": endTime}, "is_use": 1, "type": bson.M{"$in": []int{1, 2}}}).Sort("-create_time")
	var mailList []mailListInfo
	err = mailQuery.All(&mailList)
	if err != nil {
		LogRus.Debugf("query mail err %v", err)
	}

	//排行榜信息
	acitivityQuery := Session.DB(database_0).C(activityTable).Find(bson.M{"name": "rank"})
	var activityinfo acitivityInfo
	err = acitivityQuery.One(&activityinfo)
	if err != nil {
		LogRus.Debugf("query activity table err %v", err)
	}

	var users []userInfo
	err = query.All(&users)
	if err != nil {
		fmt.Printf("err:%v", err)
	}
	LogRus.Info("total:", len(users))
	for _, userinfo := range users {
		if userinfo.DeviceToken == "null" {
			continue
		}
		// 判断所在时区现在是否可以发 中午10-14pm, 晚上6-8pm, 10-12am
		userTimeZone := userinfo.TimeZone
		timeTime := GetLocalTime(userTimeZone)
		if isForbiddenTime(timeTime) {
			LogRus.Info(userinfo.Id, " in forbiddenTime")
			continue
		}
		if !isInTime(timeTime) {
			LogRus.Info(userinfo.Id, " not in notification time")
			continue
		}
		notiQuery := Session.DB(database_1).C(notificationTable).Find(bson.M{"_id": userinfo.Id})
		var notiItem notificationInfo
		err = notiQuery.One(&notiItem)
		if err != nil {
			LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("query notification err:%v", err)
			notiItem = notificationInfo{
				Id:           userinfo.Id,
				Version:      initversion,
				LastSendTime: 0,
			}
		}
		if notiItem.Version == "" {
			notiItem.Version = initversion
		}
		// 判断是否间隔2个小时
		lastSendTime := notiItem.LastSendTime
		if lastSendTime != 0 && lastSendTime+2*3600 >= int(now) {
			LogRus.Info(userinfo.Id, " in 2 hours cd")
			continue
		}
		// 优先级:版本更新的时候 无限关卡排行榜奖励 无限关卡排行榜重新开启 提示新内容可玩  邮箱收到其他信息的时候 体力值恢复满的时候  好友送礼物的时候  召回功能

		//判断是否需要发新版本通知
		if isTrue := isNeedNewVersionNoti(versioninfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_NEWVERSION, userinfo.DeviceToken)
			if err == nil {
				_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{NOTI_NEWVERSION: versioninfo.Name, "lastsendtime": now}})
				if err != nil {
					LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:", NOTI_NEWVERSION, err)
				}
			}
			LogRus.Info(userinfo.Id, NOTI_NEWVERSION)
			continue
		}
		//无限关卡排行榜奖励通知
		if isTrue := isNeedRankRewardNoti(activityinfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_RANK_REWARD, userinfo.DeviceToken)
			if err == nil {
				_, err = Session.DB(database_1).C(notificationTable).Upsert(
					bson.M{"_id": userinfo.Id},
					bson.M{"$unset": bson.M{NOTI_RANK_REWARD: 1}, "$set": bson.M{"lastsendtime": now}},
				)
				if err != nil {
					LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_RANK_REWARD, err)
				}
			}
			LogRus.Info(userinfo.Id, NOTI_RANK_REWARD)
			continue
		}
		//无限关卡排行榜开始通知
		// userDataQuery := session.DB(database_1).C(userDataTable).Find(bson.M{"_id": userinfo.Id})

		if isTrue := isNeedRankStartNoti(activityinfo, notiItem); isTrue {
			err = sendNotification(userinfo.Id, NOTI_RANK_RESTART, userinfo.DeviceToken)
			if err == nil {
				_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{NOTI_RANK_RESTART: activityinfo.AcitivityId, "lastsendtime": now}})
				if err != nil {
					LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_RANK_RESTART, err)
				}
			}
			LogRus.Info(userinfo.Id, NOTI_RANK_RESTART)
			continue
		}

		// 系统邮件通知
		if needSendAll, needSendPart, lastAllMailId, lastPartMailId := isNeedMailNoti(userinfo.Id, mailList, notiItem); needSendAll || needSendPart {
			if needSendAll {
				err = sendNotification(userinfo.Id, NOTI_SYSTEMINFO, userinfo.DeviceToken)
				if err == nil {
					_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{"lastallmailid": lastAllMailId, "lastsendtime": now}})
					if err != nil {
						LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_SYSTEMINFO, err)
					}
				}
				continue
			}
			if needSendPart {
				err = sendNotification(userinfo.Id, NOTI_SYSTEMINFO, userinfo.DeviceToken)
				if err == nil {
					_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{"lastpartmailid": lastPartMailId, "lastsendtime": now}})
					if err != nil {
						LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_SYSTEMINFO, err)
					}
				}
				continue
			}
		}

		//体力恢复满通知
		propsQuery := Session.DB(database_1).C(propsTabel).Find(bson.M{"_id": userinfo.Id})
		propsQuery.Select(bson.M{"_id": 1, "data.userTokenData.power": 1, "data.userTokenData.powerStartTime": 1})
		var propsinfo propsInfo
		err = propsQuery.One(&propsinfo)
		if err != nil {
			LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("query prorps table err :%v", err)
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
				_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$set": bson.M{NOTI_POWER: day, "powertimes": times, "lastsendtime": now}})
				if err != nil {
					LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_POWER, err)
				}
			}
			LogRus.Info(userinfo.Id, NOTI_POWER)
			continue
		}

		//好友送礼物
		if _, ok := notiItem.FriendGift["time"]; ok {
			err = sendNotification(userinfo.Id, NOTI_FRIENDGIFT, userinfo.DeviceToken)
			if err == nil {
				_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$unset": bson.M{NOTI_FRIENDGIFT: 1}, "$set": bson.M{"lastsendtime": now}})
				if err != nil {
					LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_FRIENDGIFT, err)
				}
			}
			LogRus.Info(userinfo.Id, NOTI_FRIENDGIFT)
			continue
		}

		//召回
		if isTrue := isNeedCallNoti(userinfo, notiItem); isTrue {
			callNotiType := NOTI_CALLBACK + "_" + strconv.Itoa(notiItem.CallBackTimes+1)
			err = sendNotification(userinfo.Id, callNotiType, userinfo.DeviceToken)
			if err == nil {
				_, err = Session.DB(database_1).C(notificationTable).Upsert(bson.M{"_id": userinfo.Id}, bson.M{"$inc": bson.M{NOTI_CALLBACK: 1}, "$set": bson.M{"lastsendtime": now}})
				if err != nil {
					LogRus.WithFields(logrus.Fields{"uid": userinfo.Id}).Infof("update notification %s err:%v", NOTI_CALLBACK, err)
				}
			}
			LogRus.Info(userinfo.Id, callNotiType)
			continue
		}
	}
	wg.Wait()
	LogRus.Info("succ!")
}

// 是否需要版本通知
func isNeedNewVersionNoti(version versionInfo, notification notificationInfo) bool {
	now := time.Now().Unix()
	notiVersion := notification.Version
	if now > int64(version.UpdateTime)+24*3600 {
		return false
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

func isNeedMailNoti(uid int, mailList []mailListInfo, notification notificationInfo) (needSendAll bool, needSendPart bool, lastAllMailId string, lastPartMailId string) {
	//遍历所有邮件，取出最新的 发给所有用户的邮件和发给部分用户的邮件
	for _, v := range mailList {
		//发给所有用户的邮件
		if v.Type == 2 && lastAllMailId == "" {
			lastAllMailId = v.Id.Hex()
			continue
		}
		if v.Type == 2 {
			continue
		}
		//发给部分用户的邮件
		for _, userid := range v.UidArray {
			if uid == userid {
				lastPartMailId = v.Id.Hex()
				break
			}
		}
		if lastPartMailId != "" && lastAllMailId != "" {
			break
		}
	}
	if notification.LastAllMailId != lastAllMailId {
		needSendAll = true
	}
	if notification.LastPartMailId != lastPartMailId {
		needSendPart = true
	}
	return
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
	//如果已经超过一天了，不再发通知
	if now-int64(powerRestoreTime) > 24*3600 {
		return false
	}
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
	LogRus.WithFields(log.Fields{
		"uid":         uid,
		"deviceTOken": deviceToken,
		"notiType":    notiType,
	}).Info("start send")

	if Sys == "ios" {
		wg.Add(1)
		go DoPush(notiType, deviceToken)
	} else if Sys == "android" {
		wg.Add(1)
		go DoAndroidPush(notiType, deviceToken)
	}
	return
}

type userTableInfo struct {
	Id                int    `bson:"_id"`
	Status            int    `bson:"status"`
	CreateTime        int    `bson:"create_time"`
	LastLoginTime     int    `bson:"last_login_time"`
	TimeZone          string `bson:"timezone"`
	SaveMoudleTime    int    `bson:"sava_module_time"`
	DeviceToken       string `bson:"devicetoken"`
	DeviceTokenUpTime int    `bson:"devicetoken_uptime"`
}

func FixTimezoneToInt(mongo *mgo.Session, database string, userTable string) {
	query := mongo.DB(database).C(userTable).Find(bson.M{}).Iter()
	var user userTableInfo
	for query.Next(&user) {
		if user.TimeZone != "" {
			LogRus.Printf("FixTimezoneToInt id:%d, timezoneType %T \n", user.Id, user.TimeZone)
			timezoneInt, err := strconv.Atoi(user.TimeZone)
			if err != nil {
				continue
			}
			mongo.DB(database).C(userTable).UpdateId(user.Id, bson.M{"$set": bson.M{"timezone": timezoneInt}})
		}
	}
	if err := query.Close(); err != nil {
		LogRus.Println("close iter err:", err)
	}
}

func UpdateProp() {
	defer ReleaseMongo()
	InitConfig()
	InitMongo()

	database := Cfg.Section("mongo").Key("database").String()
	usermdata := Cfg.Section("mongo").Key("user_m_data").String()
	usermresturant := Cfg.Section("mongo").Key("user_m_game_props").String()
	query := Session.DB(database).C(usermresturant).Find(bson.M{"data.infiniteEndTimeMap": bson.M{"$exists": true}})
	var resturantRes []Resturant
	err := query.All(&resturantRes)
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

		err = Session.DB(database).C(usermdata).Update(bson.M{"_id": value.Id}, bson.M{"$set": bson.M{"data.userData.starcount": totalLevel}})
		if err != nil {
			fmt.Printf("update db error %v \n", err)
		} else {
			fmt.Printf("%d update succ starcount:%d \n", value.Id, totalLevel)
		}
	}
}

func UpdateNpcStarTotal() {
	defer ReleaseMongo()
	InitConfig()
	InitMongo()

	database := Cfg.Section("mongo").Key("database").String()
	usermdata := Cfg.Section("mongo").Key("user_m_data").String()
	usermresturant := Cfg.Section("mongo").Key("user_m_resturant").String()

	query := Session.DB(database).C(usermresturant).Find(bson.M{"_id": bson.M{"$lt": 100}})
	var resturantRes []Resturant
	err := query.All(&resturantRes)
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

		err = Session.DB(database).C(usermdata).Update(bson.M{"_id": value.Id}, bson.M{"$set": bson.M{"data.userData.starcount": totalLevel}})
		if err != nil {
			fmt.Printf("update db error %v \n", err)
		} else {
			fmt.Printf("%d update succ starcount:%d \n", value.Id, totalLevel)
		}
	}
}
