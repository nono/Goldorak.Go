package goldorak

import (
	"log"
	"redis"
	"strconv"
	"strings"
)

type Model struct {
	db     int
	name   string
	client redis.Client
}

const keySeparator = ":"

func NewModel(name string) *Model {
	fullname := GetConfig("appname") + keySeparator + name
	database := GetConfig("database")
	db, err  := strconv.Atoi(database)
	if err != nil {
		log.Exitf("Can't read the database", err)
	}
	spec  := redis.DefaultSpec().Db(db)
	m     := new(Model)
	m.db   = db
	m.name = fullname
	m.client, err = redis.NewSynchClientWithSpec(spec)
	if err != nil {
		log.Exitf("Can't create the client", err)
	}
	return m
}

func (this *Model) FullKey(key string) string {
	return this.name + keySeparator + key
}

func (this *Model) Get(key string) string {
	value, err := this.client.Get(this.FullKey(key))
	if err != nil {
		log.Stderr("Error on Get", err)
		// TODO return something
	}
	return string(value);
}

func (this *Model) Set(key string, value string) {
	this.client.Set(this.FullKey(key), strings.Bytes(value))
	// TODO error handling
}

