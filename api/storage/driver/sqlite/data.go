package sqlite

import (
	"fmt"

	"github.com/kassisol/metadata/api/types"
)

func (c *Config) AddData(name, dtype, value, description string) error {
	c.DB.Create(&Data{
		Name:        name,
		Type:        dtype,
		Value:       value,
		Description: description,
	})

	return nil
}

func (c *Config) UpdateData(name, value, description string) error {
	data := Data{}
	c.DB.Where("name = ?", name).First(&data)

	if len(value) > 0 {
		data.Value = value
	}

	if len(description) > 0 {
		data.Description = description
	}

	c.DB.Save(&data)

	return nil
}

func (c *Config) RemoveData(name string) error {
	if c.memberOfProfile(name) {
		return fmt.Errorf("data \"%s\" cannot be removed. It is being used by a profile", name)
	}

	c.DB.Where("name = ?", name).Delete(Data{})

	return nil
}

func (c *Config) ListData(filter map[string]string) []types.Data {
	var result []types.Data

	sql := c.DB.Table("data").Select("id, name, type, value, description")

	if v, ok := filter["name"]; ok {
		sql = sql.Where("name = ?", v)
	}

	if v, ok := filter["type"]; ok {
		sql = sql.Where("type = ?", v)
	}

	if v, ok := filter["value"]; ok {
		sql = sql.Where("value = ?", v)
	}

	rows, _ := sql.Rows()
	defer rows.Close()

	for rows.Next() {
		var d_id uint
		var d_name string
		var d_type string
		var d_value string
		var d_description string

		rows.Scan(&d_id, &d_name, &d_type, &d_value, &d_description)

		dr := types.Data{
			ID:          d_id,
			Name:        d_name,
			Type:        d_type,
			Value:       d_value,
			Description: d_description,
		}

		result = append(result, dr)
	}

	return result
}

func (c *Config) AddDataToProfile(profile, data string) {
	p := Profile{}
	d := Data{}

	c.DB.Where("name = ?", profile).Find(&p)
	c.DB.Where("name = ?", data).Find(&d)

	c.DB.Model(&p).Association("Datas").Append(&d)
}

func (c *Config) RemoveDataFromProfile(profile, data string) {
	p := Profile{}
	d := Data{}

	c.DB.Where("name = ?", profile).Find(&p)
	c.DB.Where("name = ?", data).Find(&d)

	c.DB.Model(&p).Association("Datas").Delete(&d)
}

func (c *Config) CountData() int {
	var count int64

	c.DB.Model(&Data{}).Count(&count)

	return int(count)
}

func (c *Config) memberOfProfile(name string) bool {
	var count int64

	c.DB.Table("profiles").Joins("JOIN profile_datas ON profile_datas.profile_id = profiles.id").Joins("JOIN data ON data.id = profile_datas.data_id").Where("data.name = ?", name).Count(&count)

	if count > 0 {
		return true
	}

	return false
}
