package recommendation

import (
	"edigophers/user"
	"testing"
)

type RepoMock []user.User

func (m RepoMock) GetUsers() ([]user.User, error) {
	return m, nil
}

func (m RepoMock) GetUser(name string) (*user.User, error) {
	return nil, nil
}

func (m RepoMock) SaveUser(user user.User) error {
	return nil
}

func TestHandleEmptyUserRepository(t *testing.T) {

	var repoMock []user.User
	service := &SimpleRecommender{RepoMock(repoMock)}

	_, err := service.GetRecommendations(&user.User{})

	if err != nil {
		t.Errorf("GetRecommendations failed: %s", err)
	}
}

func TestHandleSimpleRecommendation(t *testing.T) {

	repoMock := make([]user.User, 4)

	userInt := make([]user.Interest, 3)
	userInt[0] = *user.NewInterest("jogging", 9.0)
	userInt[1] = *user.NewInterest("cinema", 6.0)
	userInt[2] = *user.NewInterest("italian", 7.9)

	repoMock[0] = *user.NewUser("Mark", userInt)

	userInt = make([]user.Interest, 3)
	userInt[0] = *user.NewInterest("jogging", 0.0)
	userInt[1] = *user.NewInterest("cinema", 9.0)
	userInt[2] = *user.NewInterest("italian", 0.1)

	repoMock[1] = *user.NewUser("John", userInt)

	userInt = make([]user.Interest, 3)
	userInt[0] = *user.NewInterest("jogging", 0.0)
	userInt[1] = *user.NewInterest("cinema", 0.0)
	userInt[2] = *user.NewInterest("italian", 0.0)

	repoMock[2] = *user.NewUser("Simon", userInt)

	userInt = make([]user.Interest, 3)
	userInt[0] = *user.NewInterest("jogging", 5.0)
	userInt[1] = *user.NewInterest("cinema", 0.0)
	userInt[2] = *user.NewInterest("italian", 0.0)

	repoMock[3] = *user.NewUser("George", userInt)

	userInt = make([]user.Interest, 3)
	userInt[0] = *user.NewInterest("jogging", 8.0)
	userInt[1] = *user.NewInterest("cinema", 0.0)
	userInt[2] = *user.NewInterest("italian", 0.1)

	testCase := user.NewUser("Pepe", userInt)

	service := &SimpleRecommender{RepoMock(repoMock)}

	reco, err := service.GetRecommendations(testCase)

	if err != nil {
		t.Errorf("GetRecommendations failed: %s", err)
	}

	if len(reco) == 0 {
		t.Errorf("Empty recommendations")
		t.FailNow()
	}

	t.Log("best match:", reco[0], " user ")
	if bestMatch := reco[0]; bestMatch.User.Name != repoMock[3].Name {
		t.Errorf("Expected(%q), got (%q)", repoMock[3].Name, bestMatch.User.Name)

	}
}
