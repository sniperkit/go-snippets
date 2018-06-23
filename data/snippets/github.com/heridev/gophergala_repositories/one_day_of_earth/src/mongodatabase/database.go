package mongodatabase

//For Database I've choose MongoDB !!
import (
	//	"APIs/youtube"
	"config"
	//	"fmt"
	"gopkg.in/mgo.v2"
)

type Mongo struct {
	Session *mgo.Session
	DB      *mgo.Database
}

//Connect using config constants
func (m *Mongo) Connect() (err error) {
	m.Session, err = mgo.Dial(config.MONGODB_CONNECTION)
	if err != nil {
		return
	}
	m.Session.SetSafe(&mgo.Safe{})
	m.DB = m.Session.DB(config.MONGO_DATABASE)
	return
}

func (m *Mongo) Insert(collection string, doc interface{}) (err error) {
	c := m.DB.C(collection)
	err = c.Insert(doc)
	return
}

func (m *Mongo) Remove(collection string, conditions map[string]interface{}) (err error) {
	c := m.DB.C(collection)
	err = c.Remove(conditions)
	return
}

func (m *Mongo) Update(collection string, conditions map[string]interface{}, doc interface{}) (err error) {
	c := m.DB.C(collection)
	err = c.Update(conditions, doc)
	return
}

func (m *Mongo) Find(collection string, conditions map[string]interface{}) *mgo.Query {
	c := m.DB.C(collection)
	return c.Find(conditions)
}

func (m *Mongo) CloseConnection() {
	m.Session.Close()
}

func (m *Mongo) FindOne(collection string, conditions map[string]interface{}, obj interface{}) (bool, error) {
	query := m.Find(collection, conditions)
	err := query.One(obj)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return true, err
}

func (m *Mongo) FindAll(collection string, conditions map[string]interface{}, obj interface{}) (bool, error) {
	query := m.Find(collection, conditions)
	err := query.All(obj)
	if err == mgo.ErrNotFound {
		return false, nil
	}
	return true, err
}
