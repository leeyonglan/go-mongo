package mongo

import (
	"context"
	"fmt"

	"gopkg.in/mgo.v2"
)

var (
	connPool map[string]*mgo.Session
)

type DbConf struct {
	Host string
	Port string
}

type DbConnConf struct {
	confs []*DbConf
}

func (conf *DbConnConf) Init() {
	connPool = make(map[string]*mgo.Session, 20)
}

func (conf *DbConnConf) GetConn(ctx context.Context) *mgo.Session {
	dbConf := conf.getConnConf(ctx)
	confIdent := conf.getConnIdent(dbConf)
	var conn *mgo.Session
	if _, ok := connPool[confIdent]; !ok {
		var err error
		var url = "mongodb://" + dbConf.Host + ":" + dbConf.Port
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

func (conf *DbConnConf) getConnIdent(dbconf *DbConf) string {
	return "xxx"
}

func (conf *DbConnConf) getConnConf(ctx context.Context) *DbConf {
	//TODO get session key by ctx,defalt return the first element
	return conf.confs[0]
}
