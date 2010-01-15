package goldorak

import (
	"log"
	"redis"
	"strconv"
	"strings"
	"os"
)

type Connection struct {
	app     string
	db      int
	client  redis.Client
}

type Model struct {
	conn    *Connection
	name    string
}

type Instance struct {
	model   *Model
	id      int64
	Param   string
}

const keySeparator = ":"
const modelNamespace = "_models"
const keyNextId = "_id"


/**************/
/* Connection */
/**************/

func Connect() *Connection {
	fullname := GetConfig("appname")
	database := GetConfig("database")
	db, err  := strconv.Atoi(database)
	if err != nil {
		log.Stderrf("Can't read the database", err)
		return nil
	}
	spec     := redis.DefaultSpec().Db(db)
	cli, err := redis.NewSynchClientWithSpec(spec)
	if err != nil {
		log.Stderrf("Can't connect to the database", err)
		return nil
	}
	err       = cli.Ping()
	if err != nil {
		log.Stderrf("Can't ping the database", err)
		return nil
	}
	return &Connection{fullname, db, cli}
}

func (this *Connection) NewModel(name string) *Model {
	return &Model{this, name}
}


/*********/
/* Model */
/*********/

func (this *Model) FullKey(key string) string {
	return this.conn.app + keySeparator + this.name + keySeparator + key
}

func (this *Model) NextId() (id int64, err os.Error) {
	key := this.conn.app + keySeparator + modelNamespace + keySeparator + this.name
	return this.conn.client.Incr(key)
}

func (this *Model) Get(key string) (string, os.Error) {
	value, err := this.conn.client.Get(this.FullKey(key))
	return string(value), err
}

func (this *Model) Set(key string, value string) os.Error {
	return this.conn.client.Set(this.FullKey(key), strings.Bytes(value))
}

func (this *Model) Create(param string) *Instance {
	p := Parameterize(param)
	id, err := this.NextId()
	if err != nil {
		log.Stderrf("Impossible to create an instance of %s", this.name, err)
		return nil
	}
	err = this.Set(p, string(id))
	if err != nil {
		log.Stderrf("Impossible to create an instance of %s", this.name, err)
		return nil
	}
	return &Instance{this, id, p}
}

func (this *Model) Find(param string) *Instance {
	value, err := this.Get(param)
	if err != nil {
		log.Stderrf("Impossible to find %s (%s)", param, this.name, err)
		return nil
	}
	id, err := strconv.Atoi64(value)
	if err != nil {
		log.Stderrf("Impossible to convert %s to int", value, err)
		return nil
	}
	return &Instance{this, id, param};
}


/************/
/* Instance */
/************/

func (this *Instance) FieldKey(field string) string {
	return this.Param + keySeparator + field
}

func (this *Instance) Get(field string) string {
	value, err := this.model.Get(this.FieldKey(field))
	if err != nil {
		log.Stderrf("Impossible to get %s", field, err)
		return ""
	}
	return value
}

func (this *Instance) Set(field string, value string) os.Error {
	return this.model.Set(this.FieldKey(field), value)
}

