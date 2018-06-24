package sqlite

import (
	"fmt"
)

func (c *Config) AddProfile(name string) error {
	c.DB.Create(&Profile{Name: name})

	return nil
}

func (c *Config) RemoveProfile(name string) error {
	if c.profileUsedInHost(name) {
		return fmt.Errorf("profile \"%s\" cannot be removed. It is being used by an host", name)
	}

	c.DB.Where("name = ?", name).Delete(Profile{})

	return nil
}

func (c *Config) ListProfile(filter map[string]string) map[string][]string {
	result := make(map[string][]string)

	sql := c.DB.Table("profiles").Select("profiles.name, data.name").Joins("LEFT JOIN profile_datas ON profile_datas.profile_id = profiles.id").Joins("LEFT JOIN data ON data.id = profile_datas.data_id")

	if v, ok := filter["name"]; ok {
		sql = sql.Where("profiles.name = ?", v)
	}

	if v, ok := filter["elem"]; ok {
		sql = sql.Where("data.name = ?", v)
	}

	rows, _ := sql.Rows()
	defer rows.Close()

	for rows.Next() {
		var profile string
		var data string

		rows.Scan(&profile, &data)

		result[profile] = append(result[profile], data)
	}

	return result
}

func (c *Config) CountProfile() int {
	var count int64

	c.DB.Model(&Profile{}).Count(&count)

	return int(count)
}

func (c *Config) profileUsedInHost(name string) bool {
	var count int64

	c.DB.Table("hosts").Joins("JOIN profiles ON profiles.id = hosts.profile_id").Where("profiles.name = ?", name).Count(&count)

	if count > 0 {
		return true
	}

	return false
}
