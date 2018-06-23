package intchess

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
)

type User struct {
	Id          int64    `json:"user_id"`
	Username    string   `sql:"type:varchar(60);unique" json:"username"`
	AccessToken string   `sql:"type:varchar(60);" json:"-"`
	CurrentRank int      `json:"current_rank"`
	IsAi        bool     `json:"is_ai"`
	VersesAi    bool     `json:"verses_ai"`
	CreatedAt   NullTime `sql:"type:datetime" json:"created_at"`
	UpdatedAt   NullTime `sql:"type:datetime" json:"updated_at"`
	DeletedAt   NullTime `sql:"type:datetime" json:"-"`
}

func AttemptLogin(propUsername string, propPassword string) *User {
	var propUser User
	err := dbGorm.Where(&User{Username: propUsername}).First(&propUser).Error

	if err == nil {
		if err = bcrypt.CompareHashAndPassword([]byte(propUser.AccessToken), []byte(propPassword)); err == nil {
			//they have passed the login check.
			return &propUser
		}
	}
	return nil
}

func AttemptCreate(propUsername string, propPassword string, isAi bool, versesAi bool) *User {
	var propUser User

	if dbGorm.Where(&User{Username: propUsername}).First(&propUser).RecordNotFound() {
		propUser.Username = propUsername
		hashpass, _ := bcrypt.GenerateFromPassword([]byte(propPassword), 3)
		//if err != nil {
		propUser.AccessToken = string(hashpass)
		propUser.IsAi = isAi
		propUser.VersesAi = isAi
		dbGorm.Create(&propUser)
		return &propUser
		// } else {
		// 	fmt.Println("Error with bcrypt: " + err.Error())
		// 	return nil
		// }
	} else {
		fmt.Println("Username was taken.")
	}
	return nil
}

func (u *User) WonGame(game *ChessGame, opponent *User) (Winner *UserRankChange, Loser *UserRankChange) {
	winnerChange := 1
	if opponent.CurrentRank+1 > winnerChange+u.CurrentRank {
		winnerChange = opponent.CurrentRank + 1 - u.CurrentRank
	}

	Winner = &UserRankChange{
		UserId:      u.Id,
		ChessGameId: game.Id,
		RankChange:  winnerChange,
	}

	Loser = &UserRankChange{
		UserId:      opponent.Id,
		ChessGameId: game.Id,
		RankChange:  -1,
	}

	if err := dbGorm.First(u, u.Id).Error; err != nil { //reload before applying rank change
		fmt.Println("DBError at user:72: " + err.Error())
	}

	if err := dbGorm.First(opponent, opponent.Id).Error; err != nil { //reload before applying rank change
		fmt.Println("DBError at user:76: " + err.Error())
	}

	u.CurrentRank += winnerChange
	opponent.CurrentRank--

	if err := dbGorm.Save(u).Error; err != nil {
		fmt.Println("DBError at user:83: " + err.Error())
	}
	if err := dbGorm.Save(opponent).Error; err != nil {
		fmt.Println("DBError at user:86: " + err.Error())
	}
	if err := dbGorm.Create(Winner).Error; err != nil {
		fmt.Println("DBError at user:89: " + err.Error())
	}
	if err := dbGorm.Create(Loser).Error; err != nil {
		fmt.Println("DBError at user:92: " + err.Error())
	}

	return
}
