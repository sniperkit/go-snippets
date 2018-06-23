package main

import "fmt"

func Response(u User) UserResponse {
	var ju UserResponse
	ju.ID = u.ID

	var juFollowers []FollowersResponse
	for _, follower := range u.Followers {
		juFollowers = append(juFollowers, FollowersResponse{ID: follower.ID})
	}

	ju.Followers = juFollowers

	return ju
}

type FollowersResponse struct {
	ID int
}

type UserResponse struct {
	ID        int
	Followers []FollowersResponse
}

type User struct {
	ID        int
	Followers []Followers
}

type Followers struct {
	ID int
}

func main() {

	var followers []Followers
	followers = append(followers, Followers{ID: 1})
	followers = append(followers, Followers{ID: 2})
	followers = append(followers, Followers{ID: 3})
	user := User{ID: 2, Followers: followers}
	fmt.Println(user)

	response := Response(user)
	fmt.Println(response)

}
