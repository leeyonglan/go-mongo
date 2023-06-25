package teaapp

import (
	"errors"
	"strings"
	"sync"

	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
)

type ApnsClientSingleton struct {
	Client *apns2.Client
}

var instance *ApnsClientSingleton
var once sync.Once

func GetInstance() *ApnsClientSingleton {
	once.Do(func() {
		apnkey := Cfg.Section("ios").Key("apnkey").String()
		keyid := Cfg.Section("ios").Key("keyid").String()
		teamid := Cfg.Section("ios").Key("teamid").String()
		authKey, err := token.AuthKeyFromFile(apnkey)
		if err != nil {
			LogRus.Fatal("token error:", err)
		}
		token := &token.Token{
			AuthKey: authKey,
			// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
			KeyID: keyid,
			// TeamID from developer account (View Account -> Membership)
			TeamID: teamid,
		}
		instance = &ApnsClientSingleton{}
		// instance.Client = apns2.NewTokenClient(token).Development()
		instance.Client = apns2.NewTokenClient(token).Production()
	})
	return instance
}
func DoPushBody() {
	deviceToken := "cbe4a70dedafea04b63e01a730b508a972635aa7bcdb15d479bd67b04781924b"
	notification := &apns2.Notification{}
	notification.DeviceToken = deviceToken
	notification.Topic = Cfg.Section("ios").Key("bundleid").String()
	payLoad := `{
		"aps" : {
		   "alert" : {
			  "title" : "test",
			  "body" : "testat"
		   },
		   "sound":"default"
		}
	 }`
	notification.Payload = []byte(payLoad)
	client := GetInstance().Client
	res, err := client.Push(notification)
	if err != nil {
		LogRus.Info("There was an error", err)
		return
	}
	if res.Sent() {
		LogRus.Info("Sent:", res.ApnsID)
	} else {
		LogRus.Infof("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		err = errors.New(res.Reason)
	}
}

func DoPush(notiType string, deviceToken string) (err error) {
	defer wg.Done()
	notification := &apns2.Notification{}
	notification.DeviceToken = deviceToken //"5bebadcbb5cab711af22822f97438438b0d1c190405d32965a68f0ac0a7c9061"
	notification.Topic = Cfg.Section("ios").Key("bundleid").String()
	// notification.Payload = []byte(`{"aps":{"alert":"Hello!"}}`) // See Payload section below
	notiType = "noti_" + notiType
	payLoad := `{
		"aps" : {
		   "alert" : {
			  "loc-key" : "@type",
			  "loc-args" : []
		   },
		   "sound":"default"
		}
	 }`
	payLoad = strings.Replace(payLoad, "@type", notiType, 1)
	notification.Payload = []byte(payLoad)
	client := GetInstance().Client
	res, err := client.Push(notification)
	if err != nil {
		LogRus.Info("There was an error", err)
		return
	}
	if res.Sent() {
		LogRus.Info("Sent:", res.ApnsID)
	} else {
		LogRus.Infof("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		err = errors.New(res.Reason)
	}
	return
}
