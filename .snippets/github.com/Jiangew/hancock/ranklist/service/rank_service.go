package service

import (
	"github.com/go-redis/redis"
	"github.com/jiangew/hancock/ranklist/cache"
	"log"
	"strconv"
	"time"
)

var Rds *redis.Client

const PlayerLevelRankKey = "Rank:PlayerLevel"

type RankService struct {
}

var DefaultRankService = &RankService{}

func levelScoreWithTime(level int, timeStamp int64) float64 {
	if 0 != timeStamp {
		return float64(level) + (1<<14)/(float64(timeStamp)/float64(3600*24))
	} else {
		return float64(level)
	}
}

func (rankService *RankService) GetPlayerByLevelRank(start, count int64) []*cache.PlayerInfo {
	playerInfos := []*cache.PlayerInfo{}

	ids, err := Rds.ZRevRange(PlayerLevelRankKey, start, start+count-1).Result()
	if err != nil {
		log.Println("RankService: GetPlayerByLevelRank: ", err)
		return playerInfos
	}

	for _, idstr := range ids {
		id, err := strconv.Atoi(idstr)
		if nil != err {
			log.Println("RankService: GetPlayerByLevelRank: ", err)
		} else {
			playerInfo := cache.LoadPlayerInfo(id)
			if nil != playerInfos {
				playerInfos = append(playerInfos, playerInfo)
			}
		}
	}

	return playerInfos
}

func (rankService *RankService) SetPlayerLevelRank(playerInfo *cache.PlayerInfo) bool {
	if nil == playerInfo {
		return false
	}

	err := Rds.ZAdd(
		PlayerLevelRankKey,
		redis.Z{
			Score:  levelScoreWithTime(playerInfo.Level, time.Now().Unix()),
			Member: playerInfo.PlayerId,
		}, ).Err()

	if nil != err {
		log.Println("RankService: SetPlayerLevelRank: ", err)
		return false
	}

	return true
}

func (rankService *RankService) AddPlayerExp(playerId, exp int) bool {
	player := cache.GetPlayerInfo(playerId)
	if nil == player {
		return false
	}

	player.Exp += exp

	// 固定经验升级
	if player.Exp >= cache.LevelUpExp {
		player.Level += 1
		player.Exp = player.Exp - cache.LevelUpExp

		rankService.SetPlayerLevelRank(player)
	}

	cache.SetPlayerInfo(player)
	return true
}
