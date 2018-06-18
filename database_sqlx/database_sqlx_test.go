package database_sqlx

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/go-zero-boilerplate/databases/testing_utils"
)

func TestDbSqlxExec(t *testing.T) {
	Convey("Testing the sqlx exec insert fail", t, func() {
		db, err := testing_utils.GetTestingDb("dummy_table_dbsqlx1")
		So(err, ShouldBeNil)

		var countBefore int

		err = db.QueryRowx(`SELECT COUNT(*) AS cnt FROM dummy_table_dbsqlx1`).Scan(&countBefore)
		So(err, ShouldBeNil)

		ds := NewSqlxDatabase(db)
		//This should fail
		_, err = ds.Exec(`INSERT INTO dummy_table_dbsqlx1 (name, age)`)
		So(err, ShouldNotBeNil)

		var countAfter int
		err = db.QueryRowx(`SELECT COUNT(*) AS cnt FROM dummy_table_dbsqlx1`).Scan(&countAfter)
		So(err, ShouldBeNil)

		So(countAfter, ShouldEqual, countBefore)
	})

	Convey("Testing the sqlx exec insert success", t, func() {
		db, err := testing_utils.GetTestingDb("dummy_table_dbsqlx2")
		So(err, ShouldBeNil)

		var countBefore, countAfter int
		err = db.QueryRowx(`SELECT COUNT(*) AS cnt FROM dummy_table_dbsqlx2`).Scan(&countBefore)
		So(err, ShouldBeNil)
		ds := NewSqlxDatabase(db)
		_, err = ds.Exec(`INSERT INTO dummy_table_dbsqlx2 (name, age) VALUES (?, ?)`, "Hallo 1", 23)
		So(err, ShouldBeNil)
		err = db.QueryRowx(`SELECT COUNT(*) AS cnt FROM dummy_table_dbsqlx2`).Scan(&countAfter)
		So(err, ShouldBeNil)
		So(countAfter, ShouldEqual, countBefore+1)

		result, err := ds.Exec(`INSERT INTO dummy_table_dbsqlx2 (name, age) VALUES (?, ?)`, "Hallo 2", 23)
		So(err, ShouldBeNil)

		lastInsertedId, err := result.LastInsertId()
		So(err, ShouldBeNil)
		So(lastInsertedId, ShouldEqual, 2)

		rowsAffected, err := result.RowsAffected()
		So(err, ShouldBeNil)
		So(rowsAffected, ShouldEqual, 1)

		result, err = ds.Exec(`INSERT INTO dummy_table_dbsqlx2 (name, age) VALUES (?, ?)`, "Hallo 3", 23)
		So(err, ShouldBeNil)

		lastInsertedId, err = result.LastInsertId()
		So(err, ShouldBeNil)
		So(lastInsertedId, ShouldEqual, 3)

		result, err = ds.Exec(`INSERT INTO dummy_table_dbsqlx2 (name, age) VALUES (?, ?)`, "Hallo 4", 23)
		So(err, ShouldBeNil)

		rowsAffected, err = result.RowsAffected()
		So(err, ShouldBeNil)
		So(rowsAffected, ShouldEqual, 1)
	})
}
