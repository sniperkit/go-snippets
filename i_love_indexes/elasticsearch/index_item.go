package elasticsearch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/GopherGala/i_love_indexes/conn_throttler"
	"gopkg.in/errgo.v1"
)

type IndexItem struct {
	Id             string    `json:"_id,omitempty"`
	IndexOfId      string    `json:"index_of_id"`
	Name           string    `json:"name"`
	EscapedName    string    `json:"escaped_name"`
	Path           string    `json:"path"`
	UpdatedAt      time.Time `json:"updated_at"`
	Size           int64     `json:"size,omitempty"`
	URL            string    `json:"url",omitempty`
	MIMEType       string    `json:"mime_type,omitempty"`
	Description    string    `json:"description,omitempty"`
	LastSeenAt     time.Time `json:"last_seen_at"`
	LastModifiedAt time.Time `json:"last_modified_at"`
}

func (i *IndexItem) Type() string {
	return "index_item"
}

func (i *IndexItem) Host() string {
	u, _ := url.Parse(i.URL)
	return u.Host
}

func (i *IndexItem) GetId() string {
	return i.Id
}

func (i *IndexItem) SetId(id string) {
	i.Id = id
}

func (i *IndexItem) GetEscapedName() string {
	r := strings.NewReplacer("-", " ", "_", " ", ".", " ")
	return r.Replace(i.Name)
}

func (i *IndexItem) SetSizeFromHeader() error {
	if i.MIMEType == "directory" {
		return nil
	}
	sem := conn_throttler.Acquire(i.Host())
	defer sem.Release()
	req, err := http.NewRequest("HEAD", i.URL, nil)
	if err != nil {
		return errgo.Mask(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errgo.Mask(err)
	}
	res.Body.Close()
	lengthStr := res.Header.Get("Content-Length")
	if lengthStr == "" {
		log.Println(i, "has no Content-Length header")
		return nil
	}
	length, err := strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		return errgo.Mask(err)
	}
	i.Size = length
	return nil
}

func SearchIndexItemsPerName(from string, typ string, name string) []*IndexItem {
	isRegexp := false
	if strings.ContainsAny(name, "*?+[]{}.") {
		_, err := regexp.Compile(name)
		if err == nil {
			fmt.Println("REGEXP MATCHING")
			isRegexp = true
		} else {
			fmt.Println(name, "err is", err)
		}
	}
	// Fuzzy search query
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"fuzzy_like_this_field": map[string]interface{}{
	// 			"escaped_name": map[string]string{
	// 				"like_text": name,
	// 			},
	// 		},
	// 	},
	// }

	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"fuzzy": map[string]interface{}{
	// 			"escaped_name": name,
	// 		},
	// 	},
	// }

	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"fuzzy": map[string]interface{}{
	// 			"escaped_name": map[string]interface{}{
	// 				"value":     name,
	// 				"fuzziness": 2,
	// 			},
	// 		},
	// 	},
	// }

	// Match with fuzziness
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match": map[string]interface{}{
	// 			"escaped_name": map[string]interface{}{
	// 				"query":     name,
	// 				"fuzziness": "AUTO",
	// 			},
	// 		},
	// 	},
	// }

	var query map[string]interface{}

	if isRegexp {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"regexp": map[string]interface{}{
					"name": name,
				},
			},
		}
	} else {
		if typ != "any" {
			var typRegexp string
			switch typ {
			case "video":
				typRegexp = `.*\.(mp4|flv|mkv|avi)`
			case "audio":
				typRegexp = `.*\.(mp3|m4a|flac|ogg|wav)`
			case "ebook":
				typRegexp = `.*\.(mobi|epub|cbz|djvu)`
			}
			query = map[string]interface{}{
				"query": map[string]interface{}{
					"bool": map[string]interface{}{
						"should": []map[string]interface{}{
							map[string]interface{}{
								"match": map[string]interface{}{
									"escaped_name": map[string]interface{}{
										"query":     name,
										"fuzziness": "auto",
										"type":      "phrase",
									},
								},
							},
							map[string]interface{}{
								"regexp": map[string]interface{}{
									"name": typRegexp,
								},
							},
						},
						"minimum_should_match": 2,
					},
				},
			}
		} else {
			// Full-text query
			query = map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"escaped_name": map[string]interface{}{
							"query":     name,
							"fuzziness": "auto",
							"type":      "phrase",
						},
					},
				},
			}
		}
	}

	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"should": []interface{}{
	// 				map[string]interface{}{
	// 					"match": map[string]interface{}{
	// 						"escaped_name": map[string]interface{}{
	// 							"query":     name,
	// 							"fuzziness": 0.5,
	// 						},
	// 					},
	// 				},
	// 				map[string]interface{}{
	// 					"match": map[string]interface{}{
	// 						"escaped_name": map[string]interface{}{
	// 							"query":     name,
	// 							"fuzziness": 2,
	// 							"type":      "phrase",
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"sort": []interface{}{
	// 		"_score",
	// 	},
	// }

	items := []*IndexItem{}
	var item *IndexItem

	params := map[string]interface{}{
		"size": 10,
	}
	if from != "" {
		params["from"] = from
	}

	log.Printf("Params: %+v\nQuery: %+v\n", params, query)
	res, err := defaultConn.Search(defaultIndex, "index_item", params, query)
	if err != nil {
		fmt.Println("fuzzy search err:", err)
	}

	log.Println("Got", res.Hits.Len(), "hits")

	for _, h := range res.Hits.Hits {
		item = &IndexItem{}

		json.Unmarshal(*h.Source, &item)
		items = append(items, item)
	}

	return items
}

func (i *IndexItem) Index() error {
	i.UpdatedAt = time.Now()
	return Index(i)
}
