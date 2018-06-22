package models

import (
	_ "errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "reflect"
	_ "strings"
	"time"
)

type Announcement struct {
	Id           int
	Category     string    `orm:"size(128)"`
	Announcement string    `orm:"size(128)"`
	Time         time.Time `orm:"auto_now_add;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Announcement))
}

// AddAnnouncement insert a new Announcement into database and returns
// last inserted Id on success.
func AddAnnouncement(m *Announcement) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAnnouncementById retrieves Announcement by Id. Returns error if
// Id doesn't exist
func GetAnnouncementById(id int) (v *Announcement, err error) {
	o := orm.NewOrm()
	v = &Announcement{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHackathon retrieves all Hackathon events as a slice object
func GetAllAnnouncement() []Announcement {
	o := orm.NewOrm()
	var announcements []Announcement

	_, err := o.Raw("SELECT * FROM announcement").QueryRows(&announcements)
	if err == nil {
	}
	return announcements
}

// UpdateAnnouncement updates Announcement by Id and returns error if
// the record to be updated doesn't exist
func UpdateAnnouncementById(m *Announcement) (err error) {
	o := orm.NewOrm()
	v := Announcement{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAnnouncement deletes Announcement by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAnnouncement(id int) (err error) {
	o := orm.NewOrm()
	v := Announcement{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Announcement{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
