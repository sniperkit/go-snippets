package models

import (
	_ "errors"
	"fmt"
	_ "reflect"
	_ "strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type Hackathon struct {
	Id           int
	Name         string `orm:"size(128)"`
	Description  string `orm:"size(128)"`
	Organization string `orm:"size(128)"`
	CreatedAt   time.Time `orm:"type(datetime)"`
	StartedTime   time.Time `orm:"type(datetime)"`
	EndTime   time.Time `orm:"type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Hackathon))
}

// AddHackathon insert a new Hackathon into database and returns
// last inserted Id on success.
func AddHackathon(m *Hackathon) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// CreateHackathon checks for an existing hackathon and create one
// if it's not already there by same name.
func CreateHackathon(m *Hackathon) {
	o := orm.NewOrm()
	hackathon := Hackathon{Name: m.Name}

	err := o.Read(&hackathon, "Name")
	if err == orm.ErrNoRows {
		beego.Info("no result found")
		o.Insert(m)
	} else if err == orm.ErrMissPK {
		beego.Info("no primary key found")
	} else {
		beego.Info(hackathon.Name)
	}
}

// GetAllHackathon retrieves all Hackathon events as a slice object
func GetAllHackathon() []Hackathon {
	o := orm.NewOrm()
	var hackathons []Hackathon

	_, err := o.Raw("SELECT * FROM hackathon").QueryRows(&hackathons)
	if err == nil {
	}
	return hackathons
}

// GetHackathonById retrieves Hackathon by Id. Returns error if
// Id doesn't exist
func GetHackathonById(id int) (v *Hackathon, err error) {
	o := orm.NewOrm()
	v = &Hackathon{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetHackathonByName(name string) (v *User, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM hackathon where name = ?", name).QueryRow(&v)
	if err == nil {
	}else{
		beego.Error(err)
		return nil, err
	}
	return v, nil
}

// UpdateHackathon updates Hackathon by Id and returns error if
// the record to be updated doesn't exist
func UpdateHackathonById(m *Hackathon) (err error) {
	o := orm.NewOrm()
	v := Hackathon{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHackathon deletes Hackathon by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHackathon(id int) (err error) {
	o := orm.NewOrm()
	v := Hackathon{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Hackathon{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
