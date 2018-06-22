package cronTask

import (
	"APIs/youtube"
	"config"
	"lib"
	"mongodatabase"
	"time"
)

func Youtube_Cron(lat, lng, date, distance string) (cerr *lib.CError) {
	index := lib.MD5strings(lat, lng, date, distance)
	loc_hash := lib.MD5strings(lat, lng)
	youtube_collection := mongodatabase.YoutubeCollection{}
	m := mongodatabase.Mongo{}
	m.Connect()
	elem_condition := map[string]interface{}{
		"index": index,
	}
	found, err := m.FindOne(config.YOUTUBE_DB_COLLECTION, elem_condition, &youtube_collection)

	if err != nil {
		cerr = &lib.CError{}
		cerr.SetMessage(err.Error())
		return
	}

	videos, rerr := youtube.SearchVideos(lat, lng, date, distance, true, "")
	if rerr != nil {
		cerr = rerr
		return
	}
	if !found {
		//CleanUP last Data
		clean_data := make([]mongodatabase.FlickrCollection, 0)
		q := m.Find(config.YOUTUBE_DB_COLLECTION, map[string]interface{}{
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
			m.Remove(config.YOUTUBE_DB_COLLECTION, map[string]interface{}{
				"index": clean_data[0].Index,
			})
		}

		youtube_collection.CreateDate = time.Now().UTC()
		youtube_collection.Index = index
		youtube_collection.DateStr = date
		youtube_collection.LocationHash = loc_hash
		youtube_collection.Videos = append(youtube_collection.Videos, videos...)
		err = m.Insert(config.YOUTUBE_DB_COLLECTION, youtube_collection)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
	} else {
		var (
			temp_videos []youtube.Video
			contains    bool
		)
		for _, t := range videos {
			contains = false
			for i, v := range youtube_collection.Videos {
				if t.ID == v.ID {
					youtube_collection.Videos[i] = t
					contains = true
					break
				}
			}
			if !contains {
				temp_videos = append(temp_videos, t)
			}
		}
		youtube_collection.Videos = append(youtube_collection.Videos, temp_videos...)
		err = m.Update(config.YOUTUBE_DB_COLLECTION, elem_condition, youtube_collection)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
	}

	return
}
