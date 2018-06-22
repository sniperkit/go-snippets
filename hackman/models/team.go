package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	"reflect"
	"strings"
	"time"
)

type Team struct {
	Id          int64  `orm:"auto"`
	Name        string `orm:"size(128)"`
	RepoName    string `orm:"size(128)"`
	UserId1     int
	UserId2     int
	UserId3     int
	UserId4     int
	HackathonId int
	CreatorId   int
	AccByU1     bool
	AccByU2     bool
	AccByU3     bool
	AccByU4     bool
	CreatedAt   time.Time `orm:"type(datetime)"`
}

func init() {
	orm.RegisterModel(new(Team))
}

// AddTeam insert a new Team into database and returns
// last inserted Id on success.
func AddTeam(m *Team) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetTeamById retrieves Team by Id. Returns error if
// Id doesn't exist
func GetTeamById(id int64) (v *Team, err error) {
	o := orm.NewOrm()
	v = &Team{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetTeamByName(name string) (v *Team, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM team where name = ?", name).QueryRow(&v)
	if err == nil {
	}else{
		return nil, err
	}
	return v, nil
}

func GetTeamByUIdHId(UId int, HId int) (v *Team, err error) {
	o := orm.NewOrm()
	v = &Team{UserId2: UId, HackathonId: HId}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}


/* To check the team od a user */
func GetTeamByUserId1(userId int, hackathonId int) (v *Team, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM team where user_id1 = ? and hackathon_id = ?", userId, hackathonId).QueryRow(&v)
	if err == nil {
	}else{
		return nil, err
	}
	return v, nil
}

func GetTeamByUserId2(userId int, hackathonId int) (v *Team, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM team where user_id2 = ? and hackathon_id = ?", userId, hackathonId).QueryRow(&v)
	if err == nil {
	}else{
		return nil, err
	}
	return v, nil
}

func GetTeamByUserId3(userId int, hackathonId int) (v *Team, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM team where user_id3 = ? and hackathon_id = ?", userId, hackathonId).QueryRow(&v)
	if err == nil {
	}else{
		return nil, err
	}
	return v, nil
}

func GetTeamByUserId4(userId int, hackathonId int) (v *Team, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM team where user_id4 = ? and hackathon_id = ?", userId, hackathonId).QueryRow(&v)
	if err == nil {
	}else{
		return nil, err
	}
	return v, nil
}

func GetAllTeamByHackathonId(hackathonId int) (v []Team, err error) {
	o := orm.NewOrm()
	var id int64

	id,err = o.Raw("SELECT * FROM team where hackathon_id = ?", hackathonId).QueryRows(&v)
	if err == nil {
		beego.Info(id)
		beego.Info(v)
	}else{
		beego.Error(err)
		return nil, err
	}
	return v, nil
}

/*Till here*/

// GetAllTeam retrieves all Team matches certain condition. Returns empty list if
// no records exist
func GetAllTeam(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Team))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Team
	qs = qs.OrderBy(sortFields...)
	if _, err := qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateTeam updates Team by Id and returns error if
// the record to be updated doesn't exist
func UpdateTeamById(m *Team) (err error) {
	o := orm.NewOrm()
	v := Team{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteTeam deletes Team by Id and returns error if
// the record to be deleted doesn't exist
func DeleteTeam(id int64) (err error) {
	o := orm.NewOrm()
	v := Team{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Team{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
