package cronTask

import (
	"APIs"
	"APIs/instagram"
	"config"
	"lib"
	"mongodatabase"
	"time"
)

func Instagram_Cron(lat, lng, date, distance string) (cerr *lib.CError) {
	index := lib.MD5strings(lat, lng, date, distance)
	loc_hash := lib.MD5strings(lat, lng)
	instagram_collection := mongodatabase.InstagramCollection{}
	m := mongodatabase.Mongo{}
	m.Connect()
	elem_condition := map[string]interface{}{
		"index": index,
	}
	found, err := m.FindOne(config.INSTAGRAM_DB_COLLECTION, elem_condition, &instagram_collection)

	if err != nil {
		cerr = &lib.CError{}
		cerr.SetMessage(err.Error())
		return
	}

	images, rerr := instagram.SearchImages(lat, lng, date, distance)
	if rerr != nil {
		cerr = rerr
		return
	}
	if !found {
		//CleanUP last Data
		clean_data := make([]mongodatabase.FlickrCollection, 0)
		q := m.Find(config.INSTAGRAM_DB_COLLECTION, map[string]interface{}{
			"locationhash": lib.MD5strings(lat, lng),
		})
		q.Sort("createdate")
		err := q.All(&clean_data)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
		if len(clean_data) > 2 {
			m.Remove(config.INSTAGRAM_DB_COLLECTION, map[string]interface{}{
				"index": clean_data[0].Index,
			})
		}

		instagram_collection.CreateDate = time.Now().UTC()
		instagram_collection.Index = index
		instagram_collection.DateStr = date
		instagram_collection.LocationHash = loc_hash
		instagram_collection.Images = append(instagram_collection.Images, images...)
		err = m.Insert(config.INSTAGRAM_DB_COLLECTION, instagram_collection)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
	} else {
		var (
			temp_images []APIs.ApiImage
			contains    bool
		)
		for _, t := range images {
			contains = false
			for i, v := range instagram_collection.Images {
				if t.Id == v.Id {
					instagram_collection.Images[i] = t
					contains = true
					break
				}
			}
			if !contains {
				temp_images = append(temp_images, t)
			}
		}
		instagram_collection.Images = append(instagram_collection.Images, temp_images...)
		err = m.Update(config.INSTAGRAM_DB_COLLECTION, elem_condition, instagram_collection)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
	}

	return
}
