package database

import "../../common"

func Query(s Select) (result map[string]string, err error) {
	rows, err := common.Query(s.String(), s.params...)
	if err != nil {
		return
	}
	result = Fetch(rows)
	return
}

func QueryAll(s Select) (result []map[string]interface{}, err error) {
	rows, err := common.Query(s.String(), s.params...)
	if err != nil {
		return
	}
	result = FetchAll(rows)
	return
}
