package tilnewtab

import (
	"encoding/json"
	"errors"
	"flickr"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"imgur"
	"math/rand"
	"net/http"
	"reddit"
	"regexp"
	"strconv"
)

const tilKey string = "todayilearned"
const epKey string = "earthporn"
const minImgHeight int = 1000
const minImgWidth int = 1000

var titleSizeRegexp *regexp.Regexp = regexp.MustCompile(`\[(\d{4})\s*x\s*(\d{4})\]`)
var imgurRegexp *regexp.Regexp = regexp.MustCompile(`imgur.com/([a-zA-Z0-9]{7})\.?[a-z]*`)
var flickrRegexp *regexp.Regexp = regexp.MustCompile(`www\.flickr\.com.*([0-9]{11})`)
var instaRegexp *regexp.Regexp = regexp.MustCompile(`instagram\.com`)
var deviantRegexp *regexp.Regexp = regexp.MustCompile(`deviantart\.com`)

type filter func(context.Context, []reddit.Post) []reddit.Post

func init() {
	http.HandleFunc("/get_random_img", randomImageHandler)
	http.HandleFunc("/get_all_img", allImageHandler)
	http.HandleFunc("/get_random_til", randomTilHandler)
	http.HandleFunc("/update_saved_posts", updateSavedPostsHandler)
}

/*
=========================================
=============== HANDLERS ================
=========================================
*/

func updateSavedPostsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := deleteSavedPosts(ctx, tilKey)
	if err != nil && err != memcache.ErrCacheMiss {
		log.Errorf(ctx, "Error deleting saved posts %s: %s", tilKey, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = deleteSavedPosts(ctx, epKey)
	if err != nil && err != memcache.ErrCacheMiss {
		log.Errorf(ctx, "Error deleting saved posts %s: %s", epKey, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// repopulate db/cache
	_, err = getRandomPost(ctx, epKey, filterImgPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = getRandomPost(ctx, tilKey, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func randomImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	img, err := getRandomPost(ctx, epKey, filterImgPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func allImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	imgs, err := getAllPosts(ctx, epKey, filterImgPosts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&imgs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func randomTilHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	til, err := getRandomPost(ctx, tilKey, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(&til)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
=========================================
=============== HELPERS =================
=========================================
*/

func getRandomPost(c context.Context, key string, postFilter filter) (*reddit.Post, error) {
	posts, err := getAllPosts(c, key, postFilter)
	if err != nil {
		return nil, err
	}
	if len(posts) > 0 {
		// pick a rand element from posts
		return &posts[rand.Intn(len(posts))], nil
	}

	errorStr := fmt.Sprintf("Could not get any %s from anywhere :(", key)
	return nil, errors.New(errorStr)
}

func getAllPosts(c context.Context, key string, postFilter filter) ([]reddit.Post, error) {
	var err error
	var posts []reddit.Post
	needCacheSave := false
	needDBSave := false
	// first hit memcache
	// if the value is in memcache
	if cached_posts, err := memcache.Get(c, key); err == nil {
		log.Infof(c, "Cache hit for %s", key)
		if err = json.Unmarshal(cached_posts.Value, &posts); err != nil {
			log.Errorf(c, "Unmarshalling json from cache for %s failed with error: %s", key, err)
		}
		// if it was a cache miss
	} else if err == memcache.ErrCacheMiss {
		needCacheSave = true
		// if there was a bad error thrown
	} else {
		log.Errorf(c, "Memcache threw an error that was not cache miss when fetching %s: %s", key, err)
	}

	// if there are still no TILs, need to go to DB
	if len(posts) == 0 {
		log.Infof(c, "Checking DB for %s", key)
		q := datastore.NewQuery(key).Ancestor(getParentKey(c, key))
		_, err = q.GetAll(c, &posts)
		if err != nil {
			// if the db call came back with an error, log the error but keep going
			log.Errorf(c, "Fetching %s from db failed with error: %s", key, err)
		} else if len(posts) == 0 {
			needDBSave = true
		}
	}

	// if there are still no TILs, hit reddit directly
	if len(posts) == 0 {
		log.Infof(c, "Hitting reddit for %s", key)
		posts, err = reddit.QueryRedditTop(c, key, 40, "week")
		if err != nil {
			// if there is an error here, game over
			return nil, err
		}
		if postFilter != nil {
			posts = postFilter(c, posts)
			if err != nil {
				// if there is an error here, game over
				return nil, err
			}
		}
	}

	if len(posts) != 0 {
		// save the data where ever it is missing from
		if needCacheSave {
			log.Infof(c, "Attempting to save %s to cache", key)
			if err = savePostsToCache(c, posts, key); err != nil {
				log.Errorf(c, "Failed to save %s posts to cache with error: %s", key, err)
			}
		}
		if needDBSave {
			log.Infof(c, "Attempting to save %s to DB", key)
			if err = savePostsToDB(c, posts, key); err != nil {
				log.Errorf(c, "Failed to save %s posts to db with error: %s", key, err)
			}
		}
	}
	return posts, nil
}

func checkImgDimensions(width int, height int) bool {
	ratio := float64(width) / float64(height)
	return width >= minImgWidth && height >= minImgHeight && ratio >= 1 && ratio <= 2
}

func filterImgPosts(c context.Context, posts []reddit.Post) []reddit.Post {
	var filtered []reddit.Post
	for _, post := range posts {
		// first check if it is imgur
		match := imgurRegexp.FindStringSubmatch(post.URL)
		if len(match) > 1 {
			log.Infof(c, "Checking URL with imgur: %s", post.URL)
			id := match[1]
			info, err := imgur.GetImageInfo(c, id)
			if err != nil {
				log.Errorf(c, "Error while checking imgur image info for url %s: %s", post.URL, err)
				continue
			}
			if checkImgDimensions(info.Width, info.Height) {
				log.Infof(c, "Keeping image based on imgur")
				post.URL = info.Link
				filtered = append(filtered, post)
			}
		} else {
			// check flickr
			match = flickrRegexp.FindStringSubmatch(post.URL)
			if len(match) > 1 {
				log.Infof(c, "Checking URL with flickr: %s", post.URL)
				id := match[1]
				info, err := flickr.GetImageInfo(c, id)
				if err != nil {
					log.Errorf(c, "Error while checking flickr image info for url %s: %s", post.URL, err)
					continue
				}
				if checkImgDimensions(info.Width, info.Height) {
					log.Infof(c, "Keeping image based on flickr")
					post.URL = info.Link
					filtered = append(filtered, post)
				}
			} else if !instaRegexp.MatchString(post.URL) && !deviantRegexp.MatchString(post.URL) {
				// then try to scrape the size from the title
				match := titleSizeRegexp.FindStringSubmatch(post.Title)
				if len(match) > 2 {
					log.Infof(c, "Checking URL with title: %s", post.URL)
					width, err := strconv.Atoi(match[1])
					if err != nil {
						log.Errorf(c, "Error while parsing image size: %s, %s", match[1], err)
						continue
					}
					height, err := strconv.Atoi(match[2])
					if err != nil {
						log.Errorf(c, "Error while parsing image size: %s, %s", match[2], err)
						continue
					}
					if checkImgDimensions(width, height) {
						log.Infof(c, "Keeping image based on title: %s", post.Title)
						filtered = append(filtered, post)
					}
				}
			}
		}
	}
	return filtered
}

func savePostsToDB(c context.Context, posts []reddit.Post, key string) error {
	var keys []*datastore.Key
	for _, element := range posts {
		keys = append(keys, datastore.NewKey(c, key, element.Title, 0, getParentKey(c, key)))
	}
	_, err := datastore.PutMulti(c, keys, posts)
	return err
}

func savePostsToCache(c context.Context, posts []reddit.Post, key string) error {
	encoded, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	item := &memcache.Item{
		Key:   key,
		Value: encoded,
	}
	err = memcache.Set(c, item)
	return err
}

func deleteSavedPosts(c context.Context, key string) error {
	// flush db
	q := datastore.NewQuery(key).Ancestor(getParentKey(c, key)).KeysOnly()
	keys, err := q.GetAll(c, nil)
	if err != nil {
		return err
	}
	if err = datastore.DeleteMulti(c, keys); err != nil {
		return err
	}
	// flush cache
	return memcache.Delete(c, key)
}

func getParentKey(c context.Context, key string) *datastore.Key {
	return datastore.NewKey(c, key+"parent", key+"parent", 0, nil)
}
