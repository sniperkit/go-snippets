package instagram

type (
	Obj struct {
		Data []struct {
			Attribution string `json:"attribution,omitempty"`
			Caption     struct {
				CreatedTime string `json:"created_time"`
				From        struct {
					FullName       string `json:"full_name"`
					ID             string `json:"id"`
					ProfilePicture string `json:"profile_picture"`
					Username       string `json:"username"`
				} `json:"from"`
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"caption"`
			Comments struct {
				Count uint `json:"count"`
				Data  []struct {
					CreatedTime string `json:"created_time"`
					From        struct {
						FullName       string `json:"full_name"`
						ID             string `json:"id"`
						ProfilePicture string `json:"profile_picture"`
						Username       string `json:"username"`
					} `json:"from"`
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"data"`
			} `json:"comments"`
			CreatedTime string `json:"created_time"`
			Filter      string `json:"filter"`
			ID          string `json:"id"`
			Images      struct {
				LowResolution struct {
					Height uint   `json:"height"`
					URL    string `json:"url"`
					Width  uint   `json:"width"`
				} `json:"low_resolution"`
				StandardResolution struct {
					Height uint   `json:"height"`
					URL    string `json:"url"`
					Width  uint   `json:"width"`
				} `json:"standard_resolution"`
				Thumbnail struct {
					Height uint   `json:"height"`
					URL    string `json:"url"`
					Width  uint   `json:"width"`
				} `json:"thumbnail"`
			} `json:"images"`
			Likes struct {
				Count uint `json:"count"`
				Data  []struct {
					FullName       string `json:"full_name"`
					ID             string `json:"id"`
					ProfilePicture string `json:"profile_picture"`
					Username       string `json:"username"`
				} `json:"data"`
			} `json:"likes"`
			Link     string `json:"link"`
			Location struct {
				Latitude  float32 `json:"latitude"`
				Longitude float32 `json:"longitude"`
			} `json:"location"`
			Tags []string `json:"tags"`
			Type string   `json:"type"`
			User struct {
				Bio            string `json:"bio"`
				FullName       string `json:"full_name"`
				ID             string `json:"id"`
				ProfilePicture string `json:"profile_picture"`
				Username       string `json:"username"`
				Website        string `json:"website"`
			} `json:"user"`
			UserHasLiked bool          `json:"user_has_liked"`
			UsersInPhoto []interface{} `json:"users_in_photo,omitempty"`
		} `json:"data"`
		Meta struct {
			Code uint `json:"code"`
		} `json:"meta"`
		Pagination struct {
			DeprecationWarning string `json:"deprecation_warning"`
			MinTagID           string `json:"min_tag_id"`
			NextMaxID          string `json:"next_max_id"`
			NextMaxTagID       string `json:"next_max_tag_id"`
			NextMinID          string `json:"next_min_id"`
			NextURL            string `json:"next_url"`
		} `json:"pagination"`
	}
)
