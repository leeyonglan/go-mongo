package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"teaapp"
	"time"
	"unsafe"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/leeyonglan/go-mongo"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

func TestConn() {
	var dbconf []*mongo.DbConf
	// dbconf = append(dbconf, &mongo.DbConf{Host: "49.232.8.38", Port: "27017", User: "admin", Pass: "Rh8qie8KftgkWYND"})
	dbconf = append(dbconf, &mongo.DbConf{Host: "127.0.0.1", Port: "27017", User: "transuser", Pass: "LdtduhfUbZ21sGUS"})

	var dbconfs = &mongo.DbConnConf{
		Confs: dbconf,
	}
	dbconfs.Init()

	var conn = dbconfs.GetConn(context.WithValue(context.Background(), `uid`, 10))
	var mongo = &mongo.Mongo{ConSession: conn}
	var where = make(map[string]interface{})
	var d, err = mongo.Find("miaocha_trans", "cl_user_trans", where)
	var totalLen int
	pattern := regexp.MustCompile("掌柜")

	if err != nil {
		fmt.Println("find from db not passed", err)
	} else {
		for _, value := range d {
			a := pattern.FindSubmatch([]byte(value.Zh_txt))
			if len(a) > 0 {
				updateWhere := make(map[string]string)
				updateWhere["_id"] = value.Id
				updateValue := strings.ReplaceAll(value.Zh_txt, `掌柜`, `老板`)
				fmt.Println("relace string:", updateValue)
				err := mongo.Update("miaocha_trans", "cl_user_trans", bson.M{"_id": value.Id}, bson.M{"$set": bson.M{"zh_txt": updateValue}})
				if err != nil {
					fmt.Println(err, value.Id)
					break
				} else {
					totalLen++
					fmt.Println(value.Zh_txt, value.Id+` update succ`)
				}
				// totalLen++
				fmt.Println(`update total `, totalLen)
			}

		}
		fmt.Println(`total:`, totalLen)
	}
}

type Employee struct {
	Position string
	Salary   uint
}

func EmployeeByID(id int) Employee {
	return Employee{Position: "coo", Salary: 10000}
}

func TestStruct() {
	s := EmployeeByID(1)
	s.Salary = 1000
	fmt.Println(s.Position, s.Salary)
}

type Animal struct {
}
type Dog struct {
	Animal
}

func (a *Dog) Name() {
	fmt.Println("dog")
}
func (a *Animal) Name() {
	fmt.Println("animal")
}

var wg sync.WaitGroup

func wgRountine(i int) {
	defer wg.Done()
	fmt.Printf("test %d \n", i)
}
func testWg() {
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go wgRountine(i)
	}
	wg.Wait()
	println("main finish")
}

func testSyncMutex() {
	a := 0
	var lock sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			lock.Lock()
			defer lock.Unlock()
			a += 1
			fmt.Printf("goroutine %d,a=%d \n", idx, a)
			wg.Done()
		}(a)
	}
	wg.Wait()
	println("test syncMutex finisth")
}
func testSignal() {
	route := gin.Default()
	route.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "hello Gin Server")
	})
	// srv := &http.Server{
	// 	Addr:    ":8080",
	// 	Handler: route,
	// }
	// go func() {
	// if err := endless.ListenAndServe(":8080", route); err != nil {
	// 	log.Fatalf("listen: %s\n", err)
	// }

	// // 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	// quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// // kill 默认会发送 syscall.SIGTERM 信号
	// // kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// // kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// // signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit
	// log.Println("Shutdown Server ...")
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatal("Server Shutdown:", err)
	// }
	log.Println("Server exiting")
}

func getSepicalCharNum(s string) int32 {
	// var sepicalChars [5]string = [5]string{"a", "e", "i", "o", "u"}
	var sepicalChars map[rune]bool = map[rune]bool{'a': true, 'e': true, 'i': true, 'o': true, 'u': true}
	sslice := []rune(s)
	var count int32 = 0
	for _, v := range sslice {
		if _, ok := sepicalChars[v]; ok {
			count++
		}
	}
	return count
}

type Foo struct {
	a int32
	b int32
}

func testPointer() {
	foo := &Foo{}
	bp := uintptr(unsafe.Pointer(foo)) + 4
	*(*int32)(unsafe.Pointer(bp)) = 1
	fmt.Println(foo.b)
}

// func main() {

// 	/**
// 	src := os.Args[1]
// 	dest := os.Args[2]
// 	fileio.CP(src, dest)
// 	**/
// 	// TestConn()
// 	// Main2()
// 	// testWg()
// 	// testSyncMutex()
// 	// teaapp.UpdateNpcStarTotal()
// 	// testSignal()
// 	// var count = getSepicalCharNum("adfadadsafsdfa")
// 	// fmt.Printf("total:%d", count)
// 	// a := 'A'
// 	// fmt.Println(a)
// 	// goexcel.ReplaceContent()
// 	// testPointer()
// 	// ecom.Do()
// 	ecom.TestRestful()
// }

// func main() {
// 	// fmt.Println(f3())
// 	// versiondiff.Diff()
// 	// notification.DoPush()
// 	// teaapp.NotiUser()
// 	var deviceToken string
// 	flag.StringVar(&deviceToken, "token", "", "deviceToken")
// 	flag.Parse()
// 	if deviceToken == "" {
// 		fmt.Println("please input device token")
// 		os.Exit(0)
// 	}
// 	notification.InitFcm(deviceToken)
// 	// notification.DoAndroidPush("noti_newversion", "d2qQcdtUQcWwKR2NCwhKuM:APA91bF3nKmr_zDkJgM8Wu8DyDhTPdrbt8v1hLRig2_W2wa6rp1sEiBmwGdzqUIrNsHg6Myi8v1z030Pi_yy6AU4y5KgC0dSRE-KPg86riCK564yKJ6TZ58qxHAOx2x4FDDTgcosWPhq")
// }

var Log = logrus.New()
var Sys string

func main() {
	teaapp.NotiUser()
}
