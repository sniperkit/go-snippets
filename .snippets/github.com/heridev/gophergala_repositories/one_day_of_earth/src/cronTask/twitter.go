package cronTask

import (
	"APIs/twitter"
	"config"
	"lib"
	"mongodatabase"
	"time"
)

func Twitter_Cron(lat, lng, date, distance string) (cerr *lib.CError) {
	index := lib.MD5strings(lat, lng, date, distance)
	loc_hash := lib.MD5strings(lat, lng)
	twitter_collection := mongodatabase.TwitterCollection{}
	m := mongodatabase.Mongo{}
	m.Connect()
	elem_condition := map[string]interface{}{
		"index": index,
	}
	found, err := m.FindOne(config.TWITTER_DB_COLLECTION, elem_condition, &twitter_collection)

	if err != nil {
		cerr = &lib.CError{}
		cerr.SetMessage(err.Error())
		return
	}

	tweets, rerr := twitter.SearchTweets(lat, lng, date, distance, true)
	if rerr != nil {
		cerr = rerr
		return
	}
	if !found {
		//CleanUP last Data
		clean_data := make([]mongodatabase.FlickrCollection, 0)
		q := m.Find(config.TWITTER_DB_COLLECTION, map[string]interface{}{
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
			m.Remove(config.TWITTER_DB_COLLECTION, map[string]interface{}{
				"index": clean_data[0].Index,
			})
		}

		twitter_collection.CreateDate = time.Now().UTC()
		twitter_collection.Index = index
		twitter_collection.DateStr = date
		twitter_collection.LocationHash = loc_hash
		twitter_collection.Tweets = append(twitter_collection.Tweets, tweets...)
		err = m.Insert(config.TWITTER_DB_COLLECTION, twitter_collection)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
	} else {
		var (
			temp_tweets []twitter.Tweet
			contains    bool
		)
		for _, t := range tweets {
			contains = false
			for i, v := range twitter_collection.Tweets {
				if t.ID == v.ID {
					twitter_collection.Tweets[i] = t
					contains = true
					break
				}
			}
			if !contains {
				temp_tweets = append(temp_tweets, t)
			}
		}
		twitter_collection.Tweets = append(twitter_collection.Tweets, temp_tweets...)
		err = m.Update(config.TWITTER_DB_COLLECTION, elem_condition, &twitter_collection)
		if err != nil {
			cerr = &lib.CError{}
			cerr.SetMessage(err.Error())
			return
		}
	}

	return
}
