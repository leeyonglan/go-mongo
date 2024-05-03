package teaappstat

import (
	"fmt"
	"teaapp"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UserData struct {
	RTime   int64 `bson:"rtime"`
	IsClaim int64 `bson:"isclaim"`
}

type FriendGift struct {
	ID           int                 `bson:"_id"`
	ClaimNum     int                 `bson:"claimNum"`
	ClaimTime    int64               `bson:"claimTime"`
	Request      map[string]UserData `bson:"request"`
	Invite       map[string]UserData `bson:"invite"`
	FRequest     map[string]UserData `bson:"frequest"`
	FInvite      map[string]UserData `bson:"finvite"`
	FriendIdList []int               `bson:"friendIdList"`
	HelpList     map[string]UserData `bson:"helpList"`
}

func DoFor_20240115() {
	DoUnlockFriendUpdate()
	AutoClaim()
}

func DoUnlockFriendUpdate() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	usermdata := teaapp.Cfg.Section("mongo").Key("user_m_data").String()
	userTable := teaapp.Cfg.Section("mongo").Key("user_table").String()

	var storyBook storybook

	LogRus.Infof("uid:%d", storyBook.Uid)
	userquery := teaapp.Session.DB(database_1).C(usermdata).Find(bson.M{}).Iter()
	var user userModule
	for userquery.Next(&user) {
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
		MaxStageId, _ := userData["maxStageId"].(int)
		MaxRestaurantId, _ := userData["maxRestaurantId"].(int)
		if MaxRestaurantId >= 1 && MaxStageId >= 14 {
			teaapp.Session.DB(database_1).C(userTable).Update(bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"unlock.friend": 1}})
			LogRus.Infof("maxstageId %d,maxRestaurantId %d Friend Unlock Id %d", MaxStageId, MaxRestaurantId, user.ID)
		}

	}
}

func DelExpireRquest() {
	Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	userFriendGfit := teaapp.Cfg.Section("mongo").Key("user_m_friendgift").String()
	userquery := teaapp.Session.DB(database_1).C(userFriendGfit).Find(bson.M{}).Iter()
	var friendGift FriendGift

	//循环遍历 判断rtime 超过24小时的删除
	currentTime := time.Now().Unix()
	for userquery.Next(&friendGift) {
		var requestChange bool
		var update = bson.M{}
		for key, value := range friendGift.Request {
			if value.RTime+24*3600 < currentTime {
				requestChange = true
				delete(friendGift.Request, key)
			}
		}
		if requestChange {
			update["request"] = friendGift.Request
		}
		var inviteChange bool
		for key, value := range friendGift.Invite {
			if value.RTime+24*3600 < currentTime {
				inviteChange = true
				delete(friendGift.Invite, key)
			}
		}
		if inviteChange {
			update["invite"] = friendGift.Invite
		}
		var frequestChange bool
		for key, value := range friendGift.FRequest {
			if value.RTime+24*3600 < currentTime {
				frequestChange = true
				delete(friendGift.FRequest, key)
			}
		}
		if frequestChange {
			update["frequest"] = friendGift.FRequest
		}
		var finviteChange bool
		for key, value := range friendGift.FInvite {
			if value.RTime+24*3600 < currentTime {
				finviteChange = true
				delete(friendGift.FInvite, key)
			}
		}
		if finviteChange {
			update["finvite"] = friendGift.FInvite
		}

		if len(update) > 0 {
			LogRus.Infof("DEL UID %d EXPIRE REQUEST", friendGift.ID)
			teaapp.Session.DB(database_1).C(userFriendGfit).Update(bson.M{"_id": friendGift.ID}, bson.M{"$set": update})
		}
	}
	LogRus.Infof("DEL EXPIRE REQUEST DONE")
}

func AutoClaim() {
	// Init()
	database_1 := teaapp.Cfg.Section("mongo").Key("database_1").String()
	userFriendGfit := teaapp.Cfg.Section("mongo").Key("user_m_friendgift").String()
	userquery := teaapp.Session.DB(database_1).C(userFriendGfit).Find(bson.M{}).Iter()
	var friendGift FriendGift
	for userquery.Next(&friendGift) {
		var totalClaim int

		var update = bson.M{}
		for _, value := range friendGift.FriendIdList {
			var friendUpdate = bson.M{}
			var friendId = value
			fuserQuery := teaapp.Session.DB(database_1).C(userFriendGfit).Find(bson.M{"_id": friendId})
			var ffriendGift FriendGift
			err := fuserQuery.One(&ffriendGift)
			if err != nil {
				continue
			}
			for key, value := range ffriendGift.HelpList {
				isClaim := value.IsClaim
				if key == fmt.Sprint(friendGift.ID) && isClaim == 0 {
					totalClaim++
					friendUpdate["helpList."+key+".isclaim"] = 1
				}
			}
			if len(friendUpdate) > 0 {
				LogRus.Infof("UID %d SET CLAIMED %d", friendId, 1)
				teaapp.Session.DB(database_1).C(userFriendGfit).Update(bson.M{"_id": friendId}, bson.M{"$set": friendUpdate})

			}
		}

		update["autoClaim"] = totalClaim
		if len(update) > 0 {
			LogRus.Infof("UID %d AUTO CLAIM %d", friendGift.ID, totalClaim)
			teaapp.Session.DB(database_1).C(userFriendGfit).Update(bson.M{"_id": friendGift.ID}, bson.M{"$set": update})
		}
	}
}
