package recommendation

import (
	"edigophers/user"
	"edigophers/utils"
	"log"
	"math"

	"github.com/kellydunn/golang-geo"
	"github.com/muesli/regommend"
)

//Recommendation contains a user with similiar interests
type Recommendation struct {
	User  user.User
	Score float64
}

//Recommender returns recommendations of users with similiar interests
type Recommender interface {
	GetRecommendations(user *user.User) ([]Recommendation, error)
}

//New creates a new SimpleRecommender instance
func New(userRepo user.Repository) *SimpleRecommender {
	return &SimpleRecommender{uRepo: userRepo}
}

//SimpleRecommender is a simple implementation of the Recommender using regocommend neighbours
type SimpleRecommender struct {
	uRepo user.Repository
}

const (
	minScore       = 0.4  // Min score for neighborhood matching
	maxGeoDistance = 20.0 // km max distance
)

//GetRecommendations is a method for returning recommendations of users with similiar interests
func (sr SimpleRecommender) GetRecommendations(usr user.User) ([]Recommendation, error) {
	userMap, interests, err := prepareData(sr, usr)

	neighbours, err := interests.Neighbors(usr.Name)
	if err != nil {
		return nil, err
	}

	result := make([]Recommendation, 0, 10)
	for _, rec := range neighbours {
		if rec.Key == "" {
			continue
		}

		if rec.Distance >= minScore {
			u, ok := userMap[rec.Key.(string)]
			if !ok {
				log.Printf("[WARN] User map does not contain user with id:(%s)", rec.Key.(string))
			}
			userInterests := usr.Interests.AsMap()
			u.Interests = calculateInterestMatches(userInterests, u.Interests)

			result = append(result, Recommendation{User: u, Score: utils.Round(rec.Distance*100, 0)})
		}
	}

	return result, nil
}

func calculateInterestMatches(srcInterests map[interface{}]float64, targetInterests []user.Interest) []user.Interest {

	for i, interest := range targetInterests {
		srcValue, ok := srcInterests[interest.Name]
		if !ok {
			srcValue = 10
		}

		targetInterests[i].Distance = 10 - math.Abs(srcValue-interest.Rating)
	}

	return targetInterests
}

func prepareData(sr SimpleRecommender, usr user.User) (map[string]user.User, *regommend.RegommendTable, error) {
	users, err := sr.uRepo.GetUsers()
	if err != nil {
		return nil, nil, err
	}

	srcLocation := geo.NewPoint(usr.Location.Latitude, usr.Location.Longitude)

	interests := regommend.Table("interests")
	interests.Flush() //Remove all items from interests table

	userMap := make(map[string]user.User)

	for _, u := range users {
		userMap[u.Name] = u

		targetLocation := geo.NewPoint(u.Location.Latitude, u.Location.Longitude)

		if srcLocation.GreatCircleDistance(targetLocation) <= maxGeoDistance {

			interests.Add(u.Name, u.Interests.AsMap())
		}
	}
	if _, ok := userMap[usr.Name]; !ok {
		interests.Add(usr.Name, usr.Interests.AsMap())
	}

	return userMap, interests, nil
}
