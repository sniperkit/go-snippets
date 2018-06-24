package model

import (
	"fmt"
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Write() {
	for _, v := range lib.GetAllTables() {
		f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/models/"+lib.HeadToUpper(v.TableName)+".go", os.O_RDWR, 0766)
		if err != nil {
			logs.Error(err.Error())
			os.Exit(-1)
		}
		f.WriteString(Package() + Import() + Struct(v.TableName) + TableName(v.TableName) + Init(v.TableName) + Attributes(v.TableName) + All(v.TableName) + One(v.TableName) + C(v.TableName) + U(v.TableName) + D(v.TableName) + R(v.TableName))
		f.Close()
		lib.GoFmt(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/models/" + lib.HeadToUpper(v.TableName) + ".go")
	}
}

func Package() string {
	return "package models\n"
}

func Import() string {
	return fmt.Sprintf(`
	import (
		"github.com/astaxie/beego/orm"
	)
	`)
}

func Struct(tablename string) string {
	var s string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			s += fmt.Sprintf("ID int `orm:\"column(%s);pk;auto\"`", v.Name) + "\n"
		} else {
			if strings.Index(v.Type, "int") > -1 {
				s += fmt.Sprintf("%s int `orm:\"column(%s)\"`", lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else if strings.Index(v.Type, "float") > -1 {
				s += fmt.Sprintf("%s float32 `orm:\"column(%s)\"`", lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else if strings.Index(v.Type, "double") > -1 {
				s += fmt.Sprintf("%s float64 `orm:\"column(%s)\"`", lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else {
				s += fmt.Sprintf("%s string `orm:\"column(%s)\"`", lib.HeadToUpper(v.Name), v.Name) + "\n"
			}
		}
	}
	return fmt.Sprintf(`
	type %sModel struct {
		%s
	}	
	`, lib.HeadToUpper(tablename), s)
}

func TableName(tablename string) string {
	return fmt.Sprintf(`
	func (t *%sModel) TableName() string {
		return "%s"
	}
	`, lib.HeadToUpper(tablename), tablename)
}

func Init(tablename string) string {
	return fmt.Sprintf(`
	func init() {
		orm.RegisterModel(new(%sModel))
	}
	`, lib.HeadToUpper(tablename))
}

func Attributes(tablename string) string {
	var s string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			s += fmt.Sprintf(`"ID":"ID",`) + "\n"
		} else {
			s += fmt.Sprintf(`"%s":"%s",`, lib.HeadToUpper(v.Name), lib.HeadToUpper(v.Name)) + "\n"
		}
	}
	return fmt.Sprintf(`
	func %sAttributes() map[string]string {
		attributes := map[string]string{
			%s
		}
		return attributes
	}
	`, lib.HeadToUpper(tablename), s)
}

func All(tablename string) string {
	return fmt.Sprintf(`
	func GetAll%s() (ms []%sModel, err error) {

		qb, err := orm.NewQueryBuilder("mysql")
		if err != nil {
			return ms, err
		}

		qb.Select("*").From("%s")
		o := orm.NewOrm()
		o.Raw(qb.String()).QueryRows(&ms)

		return ms, nil
	}
	`, lib.HeadToUpper(tablename), lib.HeadToUpper(tablename), tablename)
}

func One(tablename string) string {
	return fmt.Sprintf(`
	func Get%sById(id int) (%sModel, error) {
		m := %sModel{ID: id}
		o := orm.NewOrm()
		err := o.Read(&m)
		return m, err
	}
	`, lib.HeadToUpper(tablename), lib.HeadToUpper(tablename), lib.HeadToUpper(tablename))
}

func C(tablename string) string {
	return fmt.Sprintf(`
	func Create%s(m *%sModel) (int, error) {
		o := orm.NewOrm()
		n, err := o.Insert(m)
		return int(n), err
	}
	`, lib.HeadToUpper(tablename), lib.HeadToUpper(tablename))
}

func U(tablename string) string {
	return fmt.Sprintf(`
	func Update%s(m *%sModel) error {
		o := orm.NewOrm()
		v := %sModel{ID: m.ID}
		if err := o.Read(&v); err != nil {
			return err
		}
		if _, err := o.Update(m); err != nil {
			return err
		}
		return nil
	}
	`, lib.HeadToUpper(tablename), lib.HeadToUpper(tablename), lib.HeadToUpper(tablename))
}

func D(tablename string) string {
	return fmt.Sprintf(`
	func Delete%sById(id int) (int, error) {
		o := orm.NewOrm()
		v := %sModel{ID: id}
		n, err := o.Delete(&v)
		return int(n), err
	}
	`, lib.HeadToUpper(tablename), lib.HeadToUpper(tablename))
}

func R(tablename string) string {
	return fmt.Sprintf(`
	func Search%s(fields []string, limit int, offset int, orderby []string, sort string) []%sModel {
		var f string
		var r string
		var ms []%sModel
		qb, _ := orm.NewQueryBuilder("mysql")

		if len(fields) > 0 {
			for _, v := range fields {
				f = f + v
			}
			qb.Select(f).From("%s")
		} else {
			qb.Select("*").From("%s")
		}

		if len(orderby) > 0 {
			for _, v := range orderby {
				r = r + v
			}
			if sort == "ASC" || sort == "DESC" {
				qb.OrderBy(r + " " + sort)
			} else {
				qb.OrderBy(r)
			}
		} else {
			if sort == "ASC" || sort == "DESC" {
				qb.OrderBy(sort)
			}
			qb.OrderBy(r)
		}

		qb.Limit(limit).Offset(offset)

		sql := qb.String()

		o := orm.NewOrm()
		o.Raw(sql).QueryRows(&ms)
		return ms
	}
	`, lib.HeadToUpper(tablename), lib.HeadToUpper(tablename), lib.HeadToUpper(tablename), tablename, tablename)
}
