package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/zhangjunfang/softwareScaffolding/codeGenerator/vender/lib"
)

func Write() {
	for _, v := range lib.GetAllTables() {
		f, err := os.OpenFile(lib.Src()+lib.HeadToUpper(lib.GetDBName())+"/controllers/"+lib.HeadToUpper(v.TableName)+".go", os.O_RDWR, 0766)
		if err != nil {
			logs.Error(err.Error())
			os.Exit(-1)
		}
		f.WriteString(Package() + Import() + Struct(v.TableName) + I(v.TableName) + C(v.TableName) + U(v.TableName) + R(v.TableName) + D(v.TableName))
		f.Close()
		lib.GoFmt(lib.Src() + lib.HeadToUpper(lib.GetDBName()) + "/controllers/" + lib.HeadToUpper(v.TableName) + ".go")
	}
}

func Package() string {
	return "package %s/controllers\n"
}

func Import() string {
	return fmt.Sprintf(`
    import(
	    "%s/models"
        "%s/utils"
	    "strconv"
	    "github.com/astaxie/beego"
	)`, lib.URL+lib.HeadToUpper(lib.GetDBName()), lib.URL+lib.HeadToUpper(lib.GetDBName()))
}

func Struct(tablename string) string {
	return fmt.Sprintf(`
	type %sController struct {
	beego.Controller
	}
	`, lib.HeadToUpper(tablename))
}

func CPost(tablename string) string {
	var s string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			continue
		} else {
			if strings.Index(v.Type, "int") > -1 {
				s += fmt.Sprintf(`m.%s, _ = c.GetInt("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else if strings.Index(v.Type, "float") > -1 {
				s += fmt.Sprintf(`m.%s, _ = c.GetFloat("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else if strings.Index(v.Type, "double") > -1 {
				s += fmt.Sprintf(`m.%s, _ = c.GetFloat("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else {
				s += fmt.Sprintf(`m.%s = c.GetString("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			}
		}
	}
	return s
}

func UPost(tablename string) string {
	var s string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			s += fmt.Sprintf(`m.ID, _ = c.GetInt("%s")`, v.Name) + "\n"
		} else {
			if strings.Index(v.Type, "int") > -1 {
				s += fmt.Sprintf(`m.%s, _ = c.GetInt("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else if strings.Index(v.Type, "float") > -1 {
				s += fmt.Sprintf(`m.%s, _ = c.GetFloat("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else if strings.Index(v.Type, "double") > -1 {
				s += fmt.Sprintf(`m.%s, _ = c.GetFloat("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			} else {
				s += fmt.Sprintf(`m.%s = c.GetString("%s")`, lib.HeadToUpper(v.Name), v.Name) + "\n"
			}
		}
	}
	return s
}

func I(tablename string) string {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	return fmt.Sprintf(`
	func (c *%sController) Get() {
		// Data
		var tpl utils.Tpl
		tpl.Title = "%s"
		tpl.Attributes = models.%sAttributes()
		tpl.Model, _ = models.GetAll%s()
		tpl.Success = c.GetSession("success")
		tpl.Danger = c.GetSession("danger")
		c.DelSession("success")
		c.DelSession("danger")
		c.Data["Tpl"] = tpl
		// Page
		c.Layout = "public/tpl.html"
		c.TplName = "%s/index.html"
	}	
	`, u, u, u, u, l)
}

func C(tablename string) string {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	return fmt.Sprintf(`
	// Create ...
	func (c *%sController) Create() {
		var m models.%sModel

		// Post form submit
		if c.Ctx.Input.IsPost() {
			%s
			num, err := models.Create%s(&m)
			if err != nil {
				utils.ToLog("Error", err.Error(), 3)
				c.SetSession("danger", err.Error())
			} else {
				utils.ToLog("Success", "Create ID:"+strconv.Itoa(num), 5)
				c.SetSession("success", "Create ID:"+strconv.Itoa(num))
				c.Redirect("/%s/view?%s_id="+strconv.Itoa(m.ID), 302)
			}
		}

		// Data
		var tpl utils.Tpl
		tpl.Title = "%s Create"
		tpl.Attributes = models.%sAttributes()
		tpl.Action = c.URLFor("%sController.Create")
		tpl.Model = m
		c.Data["Tpl"] = tpl
		// Page
		c.Layout = "public/tpl.html"
		c.TplName = "%s/create.html"
	}	
	`, u, u, CPost(tablename), u, l, l, u, u, u, l)
}

func U(tablename string) string {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	var id string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			id = v.Name
			break
		} else {
			logs.Error(tablename + ":The primary key lost")
			os.Exit(-1)
		}
	}
	return fmt.Sprintf(`
	// Update ...
	func (c *%sController) Update() {

		id, _ := c.GetInt("%s")

		// Post form submit
		if c.Ctx.Input.IsPost() {
			var m models.%sModel
			%s
			err := models.Update%s(&m)
			if err != nil {
				utils.ToLog("Error", err.Error(), 3)
				c.SetSession("danger", err.Error())
				c.Redirect("/%s/view?%s="+strconv.Itoa(m.ID), 302)
			} else {
				utils.ToLog("Success", "Update ID:"+strconv.Itoa(m.ID), 5)
				c.SetSession("success", "Update ID:"+strconv.Itoa(m.ID))
				c.Redirect("/%s/view?%s="+strconv.Itoa(m.ID), 302)
			}
		}

		// Data
		var tpl utils.Tpl
		tpl.Title = "%s Update"
		tpl.Attributes = models.%sAttributes()
		tpl.Model, _ = models.Get%sById(id)
		tpl.Action = c.URLFor("%sController.Update")
		c.Data["Tpl"] = tpl
		// Page
		c.Layout = "public/tpl.html"
		c.TplName = "%s/update.html"
	}	
	`, u, id, u, UPost(tablename), u, l, id, l, id, u, u, u, u, l)
}

func R(tablename string) string {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	var id string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			id = v.Name
			break
		} else {
			logs.Error(tablename + ":The primary key lost")
			os.Exit(-1)
		}
	}
	return fmt.Sprintf(`
	// View ...
	func (c *%sController) View() {
		id, _ := c.GetInt("%s")
		// Data
		var tpl utils.Tpl
		tpl.Title = "%s View"
		tpl.Attributes = models.%sAttributes()
		tpl.Model, _ = models.Get%sById(id)
		tpl.Action = c.URLFor("%sController.View")
		tpl.Success = c.GetSession("success")
		tpl.Danger = c.GetSession("danger")
		c.DelSession("success")
		c.DelSession("danger")
		c.Data["Tpl"] = tpl
		// Page
		c.Layout = "public/tpl.html"
		c.TplName = "%s/view.html"
	}	
	`, u, id, u, u, u, u, l)
}

func D(tablename string) string {
	u := lib.HeadToUpper(tablename)
	l := lib.OutPerfix(tablename)
	var id string
	for _, v := range lib.GetAllFields(tablename) {
		if v.Key == "PRI" {
			id = v.Name
			break
		} else {
			logs.Error(tablename + ":The primary key lost")
			os.Exit(-1)
		}
	}
	return fmt.Sprintf(`
	// Del ...
	func (c *%sController) Del() {
		id, _ := c.GetInt("%s")
		_, err := models.Delete%sById(id)
		if err != nil {
			utils.ToLog("Error", err.Error(), 3)
			c.SetSession("danger", err.Error())
		} else {
			utils.ToLog("Success", "Delete ID:"+strconv.Itoa(id), 5)
			c.SetSession("success", "Delete ID:"+strconv.Itoa(id))
		}
		c.Redirect("/%s", 302)
	}
	`, u, id, u, l)
}
