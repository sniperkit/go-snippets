package models

import (
	_ "errors"
	"fmt"
	_ "reflect"
	_ "strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type User struct {
	Id       int
	Name     string `orm:"size(128)"`
	UserName string `orm:"size(128)"`
	Email    string `orm:"size(128)"`
	Token    string `orm:"size(128)"`
	Avatar   string `orm:"size(128)"`
	Admin    string `orm:"size(128)"`
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

func CreateUser(m *User) bool {
	o := orm.NewOrm()

	var rows int
	var result int64 = -1
	beego.Info(m.Admin)
	o.Raw("select count(*) as Count from user").QueryRow(&rows)
	beego.Info(rows)
	if rows > 0 {
		m.Admin = "no"
	}
	beego.Info(m.Admin)

	user := User{Email: m.Email}

	err := o.Read(&user, "Email")
	if err == orm.ErrNoRows {
		beego.Info("no result found")

		id, _ := o.Insert(m)
		result = id

	} else if err == orm.ErrMissPK {
		beego.Info("no primary key found")
	} else {
		beego.Info(user.Name)
	}

	beego.Info(result)
	if result == -1 {
		var users int
		o.Raw("select count(*) as Count from user where admin = ? and name = ?", "yes", m.Name).QueryRow(&users)
		if users > 0 {
			beego.Info("read admin")
			return true
		} else {
			beego.Info("fake admin")
			return false
		}
	} else if result == 1 {
		beego.Info("admin")
		return true
	} else {
		beego.Info("user")
		return false
	}
}

func IsAdmin(m *User) {
	//o := orm.NewOrm()
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetUserByEmail(email string) (v *User, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM user where email = ?", email).QueryRow(&v)
	if err == nil {
	}else{
		beego.Error(err)
		return nil, err
	}
	return v, nil
}

func GetUserByUsername(username string) (v *User, err error) {
	o := orm.NewOrm()

	err = o.Raw("SELECT * FROM user where user_name = ?", username).QueryRow(&v)
	if err == nil {
	}else{
		beego.Error(err)
		return nil, err
	}
	return v, nil
}

func GetAdminToken() (v *User, err error){
	o := orm.NewOrm()
	admin := "yes"
	err = o.Raw("SELECT * FROM user where admin = ?", admin).QueryRow(&v)
	if err == nil {
	}else{
		beego.Error(err)
		return nil, err
	}
	return v, nil
}

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
