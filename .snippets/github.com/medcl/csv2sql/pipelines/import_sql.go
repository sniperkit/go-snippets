/*
Copyright 2016 Medcl (m AT medcl.net)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pipelines

import (
	"database/sql"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/infinitbyte/framework/core/pipeline"
	"github.com/infinitbyte/framework/core/util"
	"strings"
)

type ImportSQLJoint struct {
	pipeline.Parameters
}

const mysqlConn = "mysql_conn"
const enableRollback = "rollback_enabled"

func (joint ImportSQLJoint) Name() string {
	return "import_sql"
}

func (joint ImportSQLJoint) Process(c *pipeline.Context) error {

	sqlTextMap := c.MustGetMap(sqlKey)

	for k, sqlText := range sqlTextMap {

		log.Debug("start execute: ", k, ",sql: ", sqlText)
		db, err := sql.Open("mysql", joint.MustGetString(mysqlConn))
		if err != nil {
			log.Error(err)
			panic(err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Error(err)
			panic(err)
		}

		if c.GetBool(enableRollback, true) {
			defer func() {
				if r := recover(); r != nil {
					log.Info("the database is rolled back.")
					err = tx.Rollback()
					if err != nil {
						log.Error(err)
						panic(err)
					}
				}
			}()
		}

		//插入数据
		sql := sqlText.(string)

		array := strings.Split(sql, ";")
		for _, v := range array {
			if v == "" {
				continue
			}
			result, err := tx.Exec(v)

			if err != nil {
				log.Error(err, result, v)
				util.FileAppendNewLine("log/executed_sql_error.log", sql)
				util.FileAppendNewLine("log/executed_sql_error.log", err.Error())
				panic(err)
			}

			rc, _ := result.RowsAffected()
			l, _ := result.RowsAffected()
			log.Infof("sql execute success, %v rows affected, lastInsertID: %v", rc, l)
		}

		err = tx.Commit()
		if err != nil {
			log.Error(err)
			panic(err)
		}

	}

	return nil
}
