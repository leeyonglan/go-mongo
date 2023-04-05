package teaapp

import (
	"errors"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

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
		apnkey := cfg.Section("ios").Key("apnkey").String()
		keyid := cfg.Section("ios").Key("keyid").String()
		teamid := cfg.Section("ios").Key("teamid").String()
		authKey, err := token.AuthKeyFromFile(apnkey)
		if err != nil {
			log.Fatal("token error:", err)
		}
		token := &token.Token{
			AuthKey: authKey,
			// KeyID from developer account (Certificates, Identifiers & Profiles -> Keys)
			KeyID: keyid,
			// TeamID from developer account (View Account -> Membership)
			TeamID: teamid,
		}
		instance = &ApnsClientSingleton{}
		instance.Client = apns2.NewTokenClient(token).Development()
	})
	return instance
}

func DoPush(notiType string, deviceToken string) (err error) {
	defer wg.Done()
	notification := &apns2.Notification{}
	notification.DeviceToken = deviceToken //"5bebadcbb5cab711af22822f97438438b0d1c190405d32965a68f0ac0a7c9061"
	notification.Topic = cfg.Section("ios").Key("bundleid").String()
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
		log.Info("There was an error", err)
		return
	}
	if res.Sent() {
		log.Info("Sent:", res.ApnsID)
	} else {
		log.Infof("Not Sent: %v %v %v\n", res.StatusCode, res.ApnsID, res.Reason)
		err = errors.New(res.Reason)
	}
	return
}
