package manager

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DefaultIdColumn = "id"

type EntityDbManager struct {
	Db        *sql.DB
	EntityMap map[string]string
}

func NewEntityDbManager(db *sql.DB) *EntityDbManager {
	return NewEntityDbManagerWithEntityMap(db, map[string]string{})
}

func NewEntityDbManagerWithEntityMap(db *sql.DB, entityMap map[string]string) *EntityDbManager {
	return &EntityDbManager{
		db,
		entityMap,
	}
}

func (em *EntityDbManager) GetIdColumn(entity string) string {
	if v, ok := em.EntityMap[entity]; ok {
		return v
	}
	return DefaultIdColumn
}

func (em *EntityDbManager) GetEntities(entity string, filterParams map[string]string, limit string, offset string, orderBy string, orderDir string) ([]map[string]interface{}, int, error) {
	var whereClause string = ""
	if len(filterParams) > 0 {
		var whereConditions []string
		r := strings.NewReplacer("*", "%")

		for filterParamKey, filterParamVal := range filterParams {
			wherePiece := fmt.Sprintf("`%s` LIKE '%s'", filterParamKey, r.Replace(filterParamVal))
			whereConditions = append(whereConditions, wherePiece)
		}

		whereClause = fmt.Sprintf("WHERE %s", strings.Join(whereConditions, " AND "))
	}

	query := fmt.Sprintf(
		"SELECT * FROM `%s` %s ORDER BY %s %s LIMIT %s, %s",
		entity,
		whereClause,
		orderBy,
		orderDir,
		offset,
		limit,
	)

	allResults, err := em.retrieveAllResultsByQuery(query)
	if err != nil {
		return make([]map[string]interface{}, 0), 0, err
	}

	var countResult string
	countQuery := fmt.Sprintf(
		"SELECT count(%s) FROM `%s`%s",
		em.GetIdColumn(entity),
		entity,
		whereClause,
	)

	countErr := em.Db.QueryRow(countQuery).Scan(&countResult)
	if countErr != nil {
		return make([]map[string]interface{}, 0), 0, countErr
	}

	count, _ := strconv.Atoi(countResult)

	return allResults, count, nil
}

func (em *EntityDbManager) GetEntity(entity string, id string) (map[string]interface{}, error) {
	result, err := em.retrieveSingleResultById(entity, id)
	if err != nil {
		return make(map[string]interface{}), err
	}

	return result, nil
}

func (em *EntityDbManager) PostEntity(entity string, postData map[string]interface{}) (int64, error) {
	columnsQuery := fmt.Sprintf(
		"SHOW COLUMNS FROM `%s`",
		entity,
	)

	var columns []string
	var values []string

	columnsResult, err := em.retrieveAllResultsByQuery(columnsQuery)

	// if there is an error in SHOW COLUMNS all fields are required
	if err != nil {
		for postDataKey, postDataVal := range postData {
			columns = append(columns, fmt.Sprintf("`%s`", postDataKey))
			values = append(values, em.convertJsonValue(postDataVal))
		}
	} else {
		for _, columnsRow := range columnsResult {
			column := columnsRow["Field"].(string)
			if column == em.GetIdColumn(entity) {
				continue
			}

			_, ok := postData[column]
			if ok {
				columns = append(columns, fmt.Sprintf("`%s`", column))
				values = append(values, em.convertJsonValue(postData[column]))
			}
		}
	}

	insertQuery := fmt.Sprintf(
		"INSERT INTO `%s` (%s) VALUES(%s)",
		entity,
		strings.Join(columns, ", "),
		strings.Join(values, ", "),
	)

	res, err := em.Db.Exec(insertQuery)
	if err != nil {
		return 0, err
	}

	newId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return newId, nil
}

func (em *EntityDbManager) UpdateEntity(entity string, id string, updateData map[string]interface{}) (int64, map[string]interface{}, error) {
	entityToUpdate, err := em.retrieveSingleResultById(entity, id)
	if err != nil {
		return 0, make(map[string]interface{}), err
	} else if len(entityToUpdate) <= 0 {
		return 0, entityToUpdate, nil
	}

	var updateSet []string
	for updKey, _ := range entityToUpdate {
		_, ok := updateData[updKey]

		if ok {
			entityToUpdate[updKey] = em.convertJsonValue(updateData[updKey])
			updateSet = append(updateSet, fmt.Sprintf("`%s` = %s", updKey, entityToUpdate[updKey]))
		}
	}

	if len(updateSet) <= 0 {
		return 0, entityToUpdate, nil
	}

	updQuery := fmt.Sprintf(
		"UPDATE `%s` SET %s WHERE %s = %s",
		entity,
		strings.Join(updateSet, ", "),
		em.GetIdColumn(entity),
		id,
	)

	res, err := em.Db.Exec(updQuery)
	if err != nil {
		return 0, make(map[string]interface{}), err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, make(map[string]interface{}), err
	}

	return rowsAffected, entityToUpdate, nil
}

func (em *EntityDbManager) DeleteEntity(entity string, id string) (int64, error) {
	query := fmt.Sprintf(
		"DELETE FROM `%s` WHERE %s = %s",
		entity,
		em.GetIdColumn(entity),
		id,
	)

	res, err := em.Db.Exec(query)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (em *EntityDbManager) retrieveAllResultsByQuery(query string) ([]map[string]interface{}, error) {
	allResults := make([]map[string]interface{}, 0)
	rows, err := em.Db.Query(query)
	if err != nil {
		return allResults, err
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = cerr
		}
	}()

	cols, err := rows.Columns()
	if err != nil {
		return allResults, err
	}

	rawResult := make([]interface{}, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make(map[string]interface{})
		err = rows.Scan(dest...)
		if err != nil {
			return allResults, err
		}

		for i, raw := range rawResult {
			result[cols[i]] = em.convertDbValue(raw)
		}

		allResults = append(allResults, result)
	}

	return allResults, nil
}

func (em *EntityDbManager) retrieveSingleResultById(entity string, id string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	query := fmt.Sprintf(
		"SELECT * FROM `%s` WHERE %s = %s",
		entity,
		em.GetIdColumn(entity),
		id,
	)

	rows, err := em.Db.Query(query)
	if err != nil {
		return result, err
	}

	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = cerr
		}
	}()

	cols, err := rows.Columns()
	if err != nil {
		return result, err
	}

	rawResult := make([]interface{}, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	resultCount := 0
	for rows.Next() {
		resultCount++
		err = rows.Scan(dest...)
		if err != nil {
			return result, err
		}

		for i, raw := range rawResult {
			result[cols[i]] = em.convertDbValue(raw)
		}

		if resultCount > 1 {
			continue
		}
	}

	return result, nil
}

func (em *EntityDbManager) convertDbValue(dbValue interface{}) interface{} {
	switch t := dbValue.(type) {
	default:
		fmt.Printf("[EntityDbManager] Unexpected db type %T: %#v\n", t, dbValue)
		return ""
	case bool:
		return dbValue.(bool)
	case int:
		return dbValue.(int)
	case int64:
		return dbValue.(int64)
	case []byte:
		return string(dbValue.([]byte))
	case string:
		return dbValue.(string)
	case time.Time:
		return dbValue.(time.Time).String()
	case nil:
		return nil
	}
}

func (em *EntityDbManager) convertJsonValue(jsonValue interface{}) string {
	switch t := jsonValue.(type) {
	default:
		fmt.Printf("[EntityDbManager] Unexpected json type %T: %#v\n", t, jsonValue)
		return ""
	case bool:
		return strconv.Itoa(em.Btoi(jsonValue.(bool)))
	case int:
		return strconv.Itoa(jsonValue.(int))
	case int64:
		return strconv.FormatInt(jsonValue.(int64), 10)
	case float64:
		return strconv.Itoa(int(jsonValue.(float64)))
	case string:
		return fmt.Sprintf("'%s'", strings.Replace(jsonValue.(string), "'", "\\'", -1))
	case nil:
		return "NULL"
	}
}

func (em *EntityDbManager) Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
