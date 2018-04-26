package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

// Firebasekey Enviroment variable - 12Factor III. Config
var Firebasekey = os.Getenv("Getenv")

type payload struct {
	FcmToken   string `json:"fcm_token"`
	Content    string `json:"content"`
	DeviceType string `json:"device_type"`
	Title      string `json:"title"`
	Link       string `json:"link"`
}

type push struct {
	Status     int
	FcmToken   string
	Title      string
	CreatedAt  time.Time `bson:"createdAt" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `bson:"updatedAt" json:"updatedAt,omitempty"`
	Content    string
	Link       string
	DeviceType string
}

// NotificationIOS struct of payload to IOS
type NotificationIOS struct {
	Notification struct {
		Title string
		Text  string
		Tag   string
		Badge string
	}
	Data struct {
		Link string
	}
	RegistrationIds struct {
		FcmToken string
	}
}

// NotificationAndroid struct of payload to Android
type NotificationAndroid struct {
	Tag   string
	Badge string
	Data  struct {
		Title string
		Body  string
		Link  string
	}
	RegistrationIds struct {
		FcmToken string
	}
}

func main() {

	log.Println("\nServer listing in 8080 port")
	http.HandleFunc("/push", Push)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//Push Receive datas to push in servers
func Push(w http.ResponseWriter, req *http.Request) {
	var payload payload

	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	err = json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
	}

	if payload.DeviceType == "ios" {
		var fireBasePayload NotificationIOS

		fireBasePayload.Notification.Title = payload.Title
		fireBasePayload.Notification.Text = payload.Content
		fireBasePayload.Notification.Tag = "Tag"
		fireBasePayload.Notification.Badge = "1"
		fireBasePayload.Data.Link = payload.Link
		fireBasePayload.RegistrationIds.FcmToken = payload.FcmToken
		b, _ := json.Marshal(fireBasePayload)

		go sendPostToFirebase(b)

	} else {
		var fireBasePayload NotificationAndroid

		fireBasePayload.Tag = "Tag"
		fireBasePayload.Badge = "1"
		fireBasePayload.Data.Title = payload.Title
		fireBasePayload.Data.Body = payload.Content
		fireBasePayload.Data.Link = payload.Link
		fireBasePayload.RegistrationIds.FcmToken = payload.FcmToken
		b, _ := json.Marshal(fireBasePayload)

		go sendPostToFirebase(b)

	}
	persistPush(payload)

}

func sendPostToFirebase(a []byte) error {
	d := bytes.NewReader(a)

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", d)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("key=%s", Firebasekey))
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}
	fmt.Printf(".")
	return nil
}

func persistPush(p payload) error {
	session, err := mgo.Dial("localhost:27017")
	c := session.DB("notification_center").C("push")

	if err != nil {
		panic(err)
	}
	defer session.Close()

	var payPush push

	payPush.Status = 1
	payPush.FcmToken = p.FcmToken
	payPush.Title = p.Title
	payPush.CreatedAt = time.Now()
	payPush.UpdatedAt = time.Now()
	payPush.Content = p.Content
	payPush.Link = p.Link
	payPush.DeviceType = p.DeviceType

	err = c.Insert(payPush)

	return err
}
