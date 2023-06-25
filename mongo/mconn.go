package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"

	"gopkg.in/mgo.v2"
)

var (
	connPool map[string]*mgo.Session
)

type DbConf struct {
	Host string
	Port string
	User string
	Pass string
	Ca   string
}
type MongoItem struct {
	Id            string `bson:"_id"`
	Category_name string `bson:"category_name"`
	Create_time   int32  `bson:"create_time"`
	En_txt        string `bson:"en_txt"`
	En_type       int8   `bson:"en_type"`
	French_txt    string
	From          string `bson:"from"`
	German_txt    string
	Img_path      string
	Remark        string
	Word_limit    int32
	Zh_txt        string `bson:"zh_txt"`
	Zh_type       int8   `bson:"zh_type"`
}

type DbConnConf struct {
	Confs []*DbConf
}

func (conf *DbConnConf) Init() {
	connPool = make(map[string]*mgo.Session, 20)
}

func (conf *DbConnConf) GetConn(ctx context.Context) *mgo.Session {
	dbConf := conf.getConnConf(ctx)
	confIdent := conf.getConnIdent(dbConf)
	var conn *mgo.Session
	fmt.Println("GetConn")
	if _, ok := connPool[confIdent]; !ok {
		var err error
		var url = "mongodb://" + dbConf.User + ":" + dbConf.Pass + "@" + dbConf.Host + ":" + dbConf.Port
		fmt.Println(`mongo url:`, url)
		conn, err = mgo.Dial(url)
		if err != nil {
			fmt.Println("err:", err)
		}
		connPool[confIdent] = conn
	} else {
		conn = connPool[confIdent]
	}
	return conn.Clone()
}

func (conf *DbConnConf) GetSSLCon(ctx context.Context) *mgo.Session {
	dbConf := conf.getConnConf(ctx)
	rootCerts := x509.NewCertPool()
	ok := rootCerts.AppendCertsFromPEM([]byte(dbConf.Ca))
	if !ok {
		panic("failed to parse root certificat")
	}
	tlsConfig := &tls.Config{
		RootCAs:            rootCerts,
		InsecureSkipVerify: true,
	}
	var url = "mongodb://" + dbConf.User + ":" + dbConf.Pass + "@" + dbConf.Host + ":" + dbConf.Port
	fmt.Println(`mongo url ssl:`, url)
	dialInfo, err := mgo.ParseURL(url)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("dialInfo:", dialInfo)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		if err != nil {
			log.Println(err)
		}
		return conn, err
	}
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	return session
}

func (conf *DbConnConf) getConnIdent(dbconf *DbConf) string {
	return "xxx"
}

func (conf *DbConnConf) getConnConf(ctx context.Context) *DbConf {
	//TODO get session key by ctx,defalt return the first element
	return conf.Confs[0]
}

// ca 证书
func GetMongo(host string, port string, user string, password string, ca string) *Mongo {
	var dbconf []*DbConf
	dbconf = append(
		dbconf,
		&DbConf{
			Host: host,
			Port: port,
			User: user,
			Pass: password,
			Ca:   ca,
		},
	)

	var dbconfs = &DbConnConf{
		Confs: dbconf,
	}
	dbconfs.Init()

	var conn *mgo.Session
	if ca == "" {
		conn = dbconfs.GetConn(context.TODO())
	} else {
		conn = dbconfs.GetSSLCon(context.TODO())
	}
	return &Mongo{ConSession: conn}
}
