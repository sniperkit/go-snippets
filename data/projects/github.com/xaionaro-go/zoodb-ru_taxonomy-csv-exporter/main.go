package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

type category struct {
	Id       int
	ParentId int
	Name     string
}

type item struct {
	Id           int
	Name         string
	CategoryId   int
	Synonyms     []string
	CustomFields map[string]string
}

func getCategoryMap(db *sql.DB) map[int]category {
	results, err := db.Query("SELECT id, parentid, name, alt_name FROM dle_category")
	if err != nil {
		panic(err.Error())
	}

	categoryMap := map[int]category{}
	for results.Next() {
		var cat category
		var altName string
		err = results.Scan(&cat.Id, &cat.ParentId, &cat.Name, &altName)
		if err != nil {
			panic(err.Error())
		}
		if strings.Index(cat.Name, "(") != -1 {
			cat.Name = strings.Split(strings.Split(cat.Name, "(")[1], ")")[0]
		} else {
			cat.Name = altName
		}
		categoryMap[cat.Id] = cat
	}

	return categoryMap
}

func getItems(db *sql.DB) []item {
	results, err := db.Query("SELECT id, title, xfields, category FROM dle_post")
	if err != nil {
		panic(err.Error())
	}

	items := []item{}
	for results.Next() {
		var item item
		var xFields string
		err = results.Scan(&item.Id, &item.Name, &xFields, &item.CategoryId)
		if err != nil {
			panic(err.Error())
		}
		item.Name = strings.Split(item.Name, " - ")[0]

		item.CustomFields = map[string]string{}

		for _, field := range strings.Split(xFields, "||") {
			parts := strings.Split(field, "|")
			if len(parts) < 2 {
				continue
			}
			k := parts[0]
			v := parts[1]
			switch k {
			case "english", "latin", "russian", "inotherlanguages":
			default:
				continue
			}
			item.CustomFields[k] = v
		}

		if item.CustomFields["latin"] != "" && strings.Index(item.CustomFields["latin"], ",") == -1 {
			item.Name = item.CustomFields["latin"]
		}

		for k, v := range item.CustomFields {
			switch k {
			case "english", "latin", "russian", "inotherlanguages":
				newSynonyms := strings.Split(v, ",")
				for _, newSynonym := range newSynonyms {
					newSynonym = strings.Trim(newSynonym, " ")
					if newSynonym == item.Name {
						continue
					}
					item.Synonyms = append(item.Synonyms, newSynonym)
				}
			}
		}

		items = append(items, item)
	}

	return items
}

func jsonOut(fileName string, obj interface{}) {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "    ")
	err = enc.Encode(obj)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func main() {
	//mysqlLogin    := os.Getenv("MYSQL_LOGIN")
	//mysqlPassword := os.Getenv("MYSQL_PASSWORD")

	db, err := sql.Open("mysql", "root:@unix(/var/run/mysqld/mysqld.sock)/zoodb")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	categoryMap := getCategoryMap(db)
	items := getItems(db)

	jsonOut("categoryMap.json", categoryMap)
	jsonOut("items.json", items)
}
