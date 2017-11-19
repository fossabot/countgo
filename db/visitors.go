package db

import (
	"log"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"github.com/tomasen/realip"
	"time"
)

var mgoSession *mgo.Session

const (
	c_visitors   = "visitors"
	d_ip         = "ip"
	d_date       = "date"
	d_user_agent = "user_agent"
)

type Database struct {
	dbconfig Conf
}

type Conf struct {
	Host     string `yaml: "host"`
	Database string `yaml: "database"`
	Username string `yaml: "username"`
	Password string `yaml: "password"`
}

func NewDb(c Conf) *Database {
	mgoSession = initMgoSession(c)
	return &Database{c}
}

func initMgoSession(c Conf) *mgo.Session {
	if mgoSession == nil {
		var err error
		info := &mgo.DialInfo{
			Addrs:    []string{c.Host},
			Database: c.Database,
			Username: c.Username,
			Password: c.Password,
		}
		mgoSession, err = mgo.DialWithInfo(info)
		if err != nil {
			log.Fatal("Failed to start the Mongo session")
		}
	}
	return mgoSession.Clone()
}

func (db Database) InsertVisitor(r *http.Request) error {

	data := bson.M{}
	data[d_ip] = realip.RealIP(r)
	data[d_date] = time.Now()
	for k, v := range r.Header {
		data[k] = v
	}

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	err := c.Insert(data)

	return err
}

func (db Database) GetNumberOfVisitors() (int, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	totalNum, err := c.Count()

	return totalNum, err
}

func (db Database) GetDistinctPublicIPs() ([]string, error) {

	c := mgoSession.DB(db.dbconfig.Database).C(c_visitors)
	var result []string
	err := c.Find(nil).Distinct("ip", &result)

	return result, err
}
