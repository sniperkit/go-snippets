package gorm_recursive_fetcher

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"fmt"
	"flag"
	"testing"
	"os"
)

func InitDB() *gorm.DB {
  db, err := gorm.Open("sqlite3", ":memory:")
  if err != nil {
    panic("failed to connect database")
  }
  db.Exec("DROP TABLE IF EXISTS rs;")
  
	db.SingularTable(true)
	
  // Migrate the schema
  db.AutoMigrate(&rs{})
  db.AutoMigrate(&rs_sub{})
  
  populateDb(db)
  //db.LogMode(true)
  return db

}

type rs struct {
	//Parent	*rs	`gorm:"ForeignKey:ID;AssociationForeignKey:ParentID"`
	Childs	[]*rs	`gorm:"ForeignKey:ID;AssociationForeignKey:ParentID" walkrec:"true"`
	ID	int64
	ParentID	int64	`gorm:"column:ParentID"`
	Value	string
	Sub	[]rs_sub	`gorm:"ForeignKey:ID;AssociationForeignKey:Rs_ID" walkrec:"true"`
}

type rs_sub struct {
	ID	int64
	Rs_ID	int64	`gorm:"column:rs_id"`
	Value	string
}

var db *gorm.DB

func populateDb(db *gorm.DB) {
	
	rs1 := rs{ID: 1, ParentID: 0, Value: "root"}
	db.Save(&rs1)
	rs2 := rs{ID: 2, ParentID: 1, Value: "Child1"}
	db.Save(&rs2)
	rs3 := rs{ID: 3, ParentID: 2, Value: "Child2"}
	db.Save(&rs3)
	rs4 := rs{ID: 4, ParentID: 3, Value: "Child3"}
	db.Save(&rs4)
	rs5 := rs{ID: 5, ParentID: 3, Value: "Child4"}
	db.Save(&rs5)
	rs_s := rs_sub{ID:1, Rs_ID:4, Value: "SubChild1"}
	db.Save(&rs_s)
}

func getParams() (id int) {
	flag.IntVar(&id, "id", 1, "the id to fetch")
	flag.Parse()
	return
}

func fetch(db *gorm.DB, id interface{}) (d rs, found bool) {
	
	//db.First(&d, id)
	found = false
	found = !db.Find(&d, id).RecordNotFound()
	if found {
		fetchRec(db, &d)
	}
	return
}

func TestMain(t *testing.T) {
	db = InitDB()
	defer db.Close()
	id := getParams()
	PStdErr("Loading data with ID %d\n", id)
	rs, found := fetch(db, id)
	if found {
		fmt.Print(XmlMarshal(rs) + "\n")
		os.Exit(0)
	}
	os.Exit(1)
}
