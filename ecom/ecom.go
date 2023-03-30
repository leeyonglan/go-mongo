package ecom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var (
	sellAccount       string = "apidevnew2"
	interfaceAcctount string = "wlijun2-test"
	secret            string = "d8a42f9ef"
)

type Payload struct {
	Sid        string `json:"sid"`
	Appkey     string `json:"appkey"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
	Start_time string `json:"start_time"`
	End_time   string `json:"end_time"`
}

func Do() {
	requrl := "https://sandbox.wangdian.cn/openapi2/stock_query.php"
	timenow := time.Now().Unix()
	starttime := timenow - 24*3600
	startTimeFmt := time.Unix(starttime, 0).Format("2006-01-02 15:04:05")
	endTimeFmt := time.Unix(timenow, 0).Format("2006-01-02 15:04:05")
	timestr := strconv.Itoa(int(timenow))
	params := url.Values{}
	params.Add("sid", sellAccount)
	params.Add("appkey", interfaceAcctount)
	params.Add("timestamp", timestr)
	params.Add("sign", secret)
	params.Add("start_time", startTimeFmt)
	params.Add("end_time", endTimeFmt)
	payload := Payload{Sid: sellAccount, Appkey: interfaceAcctount, Timestamp: timestr, Sign: secret, Start_time: startTimeFmt, End_time: endTimeFmt}

	payloadBytes, err := json.Marshal(payload)

	fmt.Println("requrl:", requrl)
	fmt.Println("params:", payloadBytes)
	req, err := http.NewRequest("POST", requrl, bytes.NewBuffer(payloadBytes))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	fmt.Println(string(body))
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = []User{
	{ID: "1", Name: "Alice", Email: "alice@example.com"},
	{ID: "2", Name: "Bob", Email: "bob@example.com"},
	{ID: "3", Name: "Charlie", Email: "charlie@example.com"},
}

func TestRestful() {
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, user := range users {
		if user.ID == params["id"] {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "User not found")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Error parsing request body")
		return
	}
	user.ID = fmt.Sprintf("%d", len(users)+1)
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, user := range users {
		if user.ID == params["id"] {
			var updatedUser User
			err := json.NewDecoder(r.Body).Decode(&updatedUser)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Error parsing request body")
				return
			}
			updatedUser.ID = user.ID
			users[i] = updatedUser
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "User not found")
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for i, user := range users {
		if user.ID == params["id"] {
			users = append(users[:i], users[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "User not found")
}
