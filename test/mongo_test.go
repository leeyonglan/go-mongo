package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/leeyonglan/go-mongo/mongo"
)

func TestConn(t *testing.T) {
	var dbconf []*mongo.DbConf
	dbconf = append(dbconf, &mongo.DbConf{Host: "127.0.0.1", Port: "27017"})

	var dbconfs = &mongo.DbConnConf{
		Confs: dbconf,
	}
	fmt.Println(dbconfs)
	dbconfs.Init()

	var conn = dbconfs.GetConn(context.WithValue(context.Background(), "uid", 10))
	var mongo = &mongo.Mongo{ConSession: conn}
	var where = make(map[string]interface{})
	where["_id"] = 4
	var _, err = mongo.Find("chaingame", "userinfo", where)
	if err != nil {
		t.Error("find from db not passed", err)
	} else {
		t.Log("find from db passed")
	}

}
