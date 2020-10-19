package main

import (
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	ConSession *mgo.Session
}

func (mongo *Mongo) Find(dbName string, collectionName string, where map[string]interface{}) (result []interface{}, err error) {
	session := mongo.ConSession
	defer session.Close()

	var query *mgo.Query
	query = session.DB(dbName).C(collectionName).Find(where)
	err = query.All(&result)
	return
}

func (mongo *Mongo) FindId(dbName string, collectionName string, id interface{}) (result interface{}, err error) {
	session := mongo.ConSession
	defer session.Close()

	var query *mgo.Query
	query = session.DB(dbName).C(collectionName).FindId(id)
	err = query.One(&result)
	return
}

func (mongo *Mongo) Insert(dbName string, collectionName string, doc ...interface{}) (err error) {
	session := mongo.ConSession
	defer session.Close()
	err = session.DB(dbName).C(collectionName).Insert(doc)
	return
}

func (mongo *Mongo) Update(dbName string, collectionName string, where interface{}, update interface{}) (err error) {
	session := mongo.ConSession
	defer session.Close()
	err = session.DB(dbName).C(collectionName).Update(where, update)
	return
}
