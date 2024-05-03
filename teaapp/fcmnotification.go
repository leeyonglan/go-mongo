package teaapp

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

var MessagingInstance *messaging.Client
var ctx context.Context

func InitFcmInstance() {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		LogRus.Fatalf("error initializing app: %v\n", err)
	}
	// Obtain a messaging.Client from the App.
	ctx = context.Background()
	MessagingInstance, err = app.Messaging(ctx)
	if err != nil {
		LogRus.Fatalf("error getting Messaging client: %v\n", err)
	}

}

// TODO and to env
//
//export GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/service-account-file.json"
func InitFcm(registrationToken string) {

	// This registration token comes from the client FCM SDKs.
	//registrationToken := "d2qQcdtUQcWwKR2NCwhKuM:APA91bF3nKmr_zDkJgM8Wu8DyDhTPdrbt8v1hLRig2_W2wa6rp1sEiBmwGdzqUIrNsHg6Myi8v1z030Pi_yy6AU4y5KgC0dSRE-KPg86riCK564yKJ6TZ58qxHAOx2x4FDDTgcosWPhq"

	notification := messaging.Notification{
		Title: "test",
		Body:  "testataa",
	}
	androidConfig := messaging.AndroidConfig{
		Priority: "high",
	}

	// See documentation on defining a message payload.
	message := &messaging.Message{
		// Data: map[string]string{
		// 	"notitype": "noti_newversion",
		// },
		Token:        registrationToken,
		Android:      &androidConfig,
		Notification: &notification,
	}
	client := MessagingInstance
	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		LogRus.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)
}

func DoAndroidPush(notiType string, deviceToken string) (err error) {
	defer wg.Done()
	client := MessagingInstance
	// This registration token comes from the client FCM SDKs.
	//registrationToken := "d2qQcdtUQcWwKR2NCwhKuM:APA91bF3nKmr_zDkJgM8Wu8DyDhTPdrbt8v1hLRig2_W2wa6rp1sEiBmwGdzqUIrNsHg6Myi8v1z030Pi_yy6AU4y5KgC0dSRE-KPg86riCK564yKJ6TZ58qxHAOx2x4FDDTgcosWPhq"

	// notification := messaging.Notification{
	// 	Title: "test",
	// 	Body:  "testataa",
	// }
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: map[string]string{
			"notitype": notiType,
		},
		Token: deviceToken,
		// Notification: &notification,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)

	if err != nil {
		LogRus.Error(err)
		ExpireDeviceToken(deviceToken)
		return
	}
	// Response is a message ID string.
	LogRus.Info("Successfully sent message:", response)
	return
}
