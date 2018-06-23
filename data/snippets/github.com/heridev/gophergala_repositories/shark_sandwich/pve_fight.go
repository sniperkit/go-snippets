package main

import (
	"time"
)

const (
	ENEMY_VARIANCE_LIFE  int64 = 50
	ENEMY_VARIANCE_POWER int64 = 5
	ENEMY_VARIANCE_SPEED int64 = 5
)

type (
	NPCUnit struct {
		IsNPC bool
		BaseStats
	}

	PveFight struct {
		SendEvent chan interface{}
	}
)

func NewPveFight() *PveFight {
	return &PveFight{
		make(chan interface{}),
	}
}

func NewEnemy(h *HeroSheet, name string) *NPCUnit {
	npc := &NPCUnit{
		IsNPC: true,
	}

	npc.genNPCStats(h)
	npc.BaseStats.Name = name
	return npc
}

func (n *NPCUnit) genNPCStats(h *HeroSheet) {
	n.Life = random(h.Life-ENEMY_VARIANCE_LIFE, h.Life+ENEMY_VARIANCE_LIFE)
	n.Power = random(h.Power-ENEMY_VARIANCE_POWER, h.Power+ENEMY_VARIANCE_POWER)
	n.Speed = random(h.Speed-ENEMY_VARIANCE_SPEED, h.Speed+ENEMY_VARIANCE_SPEED)
}

func (f *PveFight) Fight(hero *HeroSheet, npc *NPCUnit) {
	heroLife := hero.Life
	npcLife := npc.Life

	heroWin := make(chan bool, 1)
	npcWin := make(chan bool, 1)

	heroWon := false
	npcWon := false

	for !heroWon && !npcWon {
		select {
		case heroWon = <-heroWin:
			f.SendEvent <- FightEvent{Won: true, EnemyName: npc.BaseStats.Name}
			break
		case npcWon = <-npcWin:
			f.SendEvent <- FightEvent{Won: false, EnemyName: npc.BaseStats.Name}
			break
		case <-time.Tick(time.Duration(hero.Speed)):
			go func() {
				npcLife = npcLife - hero.Power
				if npcLife <= 0 {
					heroWin <- true
				}
			}()
		case <-time.Tick(time.Duration(npc.Speed)):
			go func() {
				heroLife = heroLife - npc.Power
				if heroLife <= 0 {
					npcWin <- true
				}
			}()
		}
	}
}
