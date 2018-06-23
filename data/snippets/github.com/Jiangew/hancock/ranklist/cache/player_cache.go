package cache

import (
	"github.com/go-redis/redis"
	"sync"
	"fmt"
	"github.com/jiangew/hancock/utils"
)

type PlayerInfo struct {
	PlayerId   int `db:"player_id" json:"player_id"`
	PlayerName string `db:"player_name" json:"player_name"`
	Exp        int `db:"exp" json:"exp"`
	Level      int `db:"level" json:"level"`
	Online     bool `db:"online" json:"online"`
}

var Rds *redis.Client
var playerMapMutex sync.Mutex
var playerMap = map[int]*PlayerInfo{}

const LevelUpExp = 100

type setPlayerValueFunc func(ma map[string]string)

func GetPlayerInfo(playerId int) *PlayerInfo {
	playerMapMutex.Lock()
	defer playerMapMutex.Unlock()

	return playerMap[playerId]
}

func SetPlayerInfo(player *PlayerInfo) {
	if nil == player {
		return
	}

	playerMapMutex.Lock()
	defer playerMapMutex.Unlock()

	playerMap[player.PlayerId] = player

	SetPlayerAllValue(player)
}

func setPlayerValue(playerId int, f setPlayerValueFunc) {
	ma := map[string]string{}
	f(ma)

	if len(ma) > 0 {
		err := Rds.HMSet(fmt.Sprintf("playerInfo:%d", playerId), ma).Err()
		if nil != err {
			fmt.Println(err)
		}
	}
}

func SetPlayerAllValue(player *PlayerInfo) {
	f := func(ma map[string]string) {
		for k, v := range utils.Struct2MapString(player) {
			ma[k] = v
		}
	}

	setPlayerValue(player.PlayerId, f)
}

func SetPlayerOnline(playerId int, online bool) {
	f := func(ma map[string]string) {
		ma["online"] = utils.Any(online)
	}

	setPlayerValue(playerId, f)
}

func SetPlayerLevel(playerId int, level int) {
	f := func(ma map[string]string) {
		ma["level"] = utils.Any(level)
	}

	setPlayerValue(playerId, f)
}

func LoadPlayerInfo(playerId int) *PlayerInfo {
	mapobj, err := Rds.HGetAll(fmt.Sprintf("playerInfo:%d", playerId)).Result()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(mapobj) < 1 {
		// redis 中没有缓存，从库中获取并写入 redis
		return nil
	}

	player := &PlayerInfo{}
	if utils.MapString2Struct(mapobj, player) {
		return player
	} else {
		return nil
	}
}
