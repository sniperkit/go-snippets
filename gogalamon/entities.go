package main

import (
	"math/rand"
	"strconv"

	"math"
)

type transform struct {
	x, y   float32
	vx, vy float32
	accel  float32
	speed  float32
}

func (t *transform) adjustV(vx, vy float32) {
	ax := vx - t.vx
	ay := vy - t.vy
	dv := float32(math.Sqrt(float64(ax*ax + ay*ay)))
	if dv > t.accel {
		ax = ax / dv * t.accel
		ay = ay / dv * t.accel
	}
	t.vx += ax
	t.vy += ay
}

func (t *transform) applyV() {
	t.x += t.vx
	t.y += t.vy
}

func NewBullet(x, y, vx, vy float32, t team) {
	var b Bullet
	b.x = x
	b.y = y
	b.vx = vx
	b.vy = vy
	b.timeLeft = framesPerSecond * 3
	b.t = t

	b.renderId = <-NextRenderId
	NewEntity <- &b
}

type Bullet struct {
	transform
	t        team
	renderId int
	timeLeft int
}

func (b *Bullet) update(overworld *Overworld, planets []*Planet) (alive bool) {
	b.applyV()
	b.timeLeft -= 1
	overworld.set(b, b.x, b.y, 8)
	hit := false
	for _, entity := range overworld.query(b, b.x, b.y, 8) {
		if entity, ok := entity.(EntityTeam); ok && entity.team() != b.t {
			if entity, ok := entity.(EntityDamage); ok {
				entity.damage(10, 0)
				hit = true
			}
		}
	}
	return !hit && b.timeLeft > 0
}

func (b *Bullet) RenderInfo() RenderInfo {
	return RenderInfo{
		b.renderId, b.x, b.y, 0, "ball_" + b.t.String(),
	}
}

type Planet struct {
	x, y         float32
	t            team
	rotation     float32
	renderId     int
	set          bool
	allegance    int
	maxAllegance int
	img          string
}

func NewPlanet(x, y float32, img string) {
	var p Planet
	p.x = x
	p.y = y
	p.allegance = framesPerSecond * 10
	p.maxAllegance = p.allegance
	p.t = TeamPirates
	p.img = img

	p.renderId = <-NextRenderId
	NewEntity <- &p

	go NewPirateShip(&p)
}

func (p *Planet) update(overworld *Overworld, planets []*Planet) (alive bool) {
	if !p.set {
		overworld.set(p, p.x, p.y, 512)
		p.set = true
	}

	nextTeam := TeamMax
	for _, entity := range overworld.query(p, p.x, p.y, 512) {
		if entity, ok := entity.(EntityTeam); ok {
			t := entity.team()
			if nextTeam != t {
				if nextTeam == TeamMax {
					nextTeam = t
				} else {
					nextTeam = TeamNone
				}
			}
		}
	}

	if nextTeam != TeamMax && nextTeam != TeamNone {
		if p.t == nextTeam {
			if p.allegance < p.maxAllegance {
				p.allegance += 1
			}
		} else {
			if p.allegance <= 0 {
				p.t = nextTeam
				p.allegance += 1
			} else {
				p.allegance -= 1
			}
		}
	}

	p.rotation += 0.03
	return true
}

func (p *Planet) RenderInfo() RenderInfo {
	return RenderInfo{
		p.renderId, p.x, p.y, p.rotation, p.img,
	}
}

func (p *Planet) planetInfo() PlanetInfo {
	return PlanetInfo{p.x, p.y, p.t.String()}
}

func (p *Planet) Allegance() (float32, string) {
	return float32(p.allegance) / float32(p.maxAllegance), p.t.String()
}

type PlanetInfo struct {
	X, Y float32
	Team string
}

type PlayerShip struct {
	transform
	user           *User
	radius         float32
	health         int
	maxHealth      int
	rotation       float32
	renderId       int
	reloadTime     int
	fullReloadTime int
	t              team

	respawnTime int
	respawning  bool
}

func NewPlayerShip(user *User, t team) {
	var p PlayerShip
	p.x = -12000
	p.y = -12000
	p.user = user
	p.accel = 0.8
	p.radius = 32
	p.maxHealth = 100
	p.health = p.maxHealth
	p.speed = 16
	p.renderId = <-NextRenderId
	p.fullReloadTime = framesPerSecond / 5
	p.t = t
	p.respawning = true
	p.respawnTime = 1

	NewEntity <- &p
}

func (p *PlayerShip) update(overworld *Overworld, planets []*Planet) (alive bool) {
	if p.health <= 0 {
		go NewDebris(p.x, p.y, p.rotation)
		{
			x := p.x
			y := p.y
			go NewSound(x, y, ExplosionSound())
		}
		p.health = p.maxHealth
		p.x = -12000
		p.y = -12000
		p.vx = 0
		p.vy = 0
		p.respawning = true
		p.respawnTime = framesPerSecond * 5
	}
	if p.respawning {
		p.respawnTime--
		if p.respawnTime <= 0 {
			p.respawning = false
			p.x = 0
			p.y = 0
			var respawnPlanets []*Planet
			for _, planet := range planets {
				if planet.t == p.t {
					respawnPlanets = append(respawnPlanets, planet)
				}
			}
			if len(respawnPlanets) > 0 {
				i := int(rand.Int()) % len(respawnPlanets)
				p.x = respawnPlanets[i].x
				p.y = respawnPlanets[i].y
			}
		}

		p.user.health = 0
	} else {
		p.user.health = float32(p.health) / float32(p.maxHealth)
	}

	var dx float32
	var dy float32
	if p.user.Key("a") {
		dx -= p.speed
	}
	if p.user.Key("d") {
		dx += p.speed
	}
	if p.user.Key("w") {
		dy -= p.speed
	}
	if p.user.Key("s") {
		dy += p.speed
	}
	if dy*dx != 0 {
		dx /= 1.41421356237
		dy /= 1.41421356237
	}
	if !p.respawning {
		if dx != 0 || dy != 0 {
			p.rotation = float32(math.Atan2(float64(dx), float64(-1*dy))) /
				(math.Pi * 2) * 360
		}
		p.adjustV(dx, dy)
		p.applyV()

		const userLimit = 10000
		distance := p.x*p.x + p.y*p.y

		if distance > userLimit*userLimit {
			mfactor := userLimit / float32(math.Sqrt(float64(distance)))
			p.x *= mfactor
			p.y *= mfactor
		}

		p.user.viewX = p.x
		p.user.viewY = p.y

	}

	overworld.set(p, p.x, p.y, p.radius)

	//Collision testing code
	// log.Println("____________")
	// log.Println(p.x, p.y)
	//log.Println(overworld.query(nil, p.x, p.y, p.radius))
	// log.Println(overworld.query(nil, p.x, p.y+5, p.radius))

	p.reloadTime += 1
	if p.fullReloadTime < p.reloadTime {
		if p.user.Key("f") {
			r := float64(p.rotation-90) / 180 * math.Pi
			vx := float32(math.Cos(r))*16 + p.vx
			vy := float32(math.Sin(r))*16 + p.vy
			x := float32(math.Cos(r))*25 + p.x
			y := float32(math.Sin(r))*25 + p.y
			r += math.Pi / 2
			dx := float32(math.Cos(r)) * 16
			dy := float32(math.Sin(r)) * 16

			go NewBullet(x+dx, y+dy, vx, vy, p.t)
			go NewBullet(x-dx, y-dy, vx, vy, p.t)
			if p.t == TeamGophers {
				go NewSound(x, y, "laser0")
			} else {
				go NewSound(x, y, "laser2")
			}
		} else if p.health < p.maxHealth {
			p.health += 2
			if p.health > p.maxHealth {
				p.health = p.maxHealth
			}
		}
		p.reloadTime = 0
	}

	return p.user.Connected()
}

func (p *PlayerShip) RenderInfo() RenderInfo {
	return RenderInfo{
		p.renderId, p.x, p.y, p.rotation, "ship_" + p.t.String(),
	}
}

func (p *PlayerShip) damage(damage int, teamSource team) {
	p.health -= damage
	{
		x, y := p.x, p.y
		if p.t == TeamGophers {
			go NewSound(x, y, "hit1")
		} else {
			go NewSound(x, y, "hit2")
		}
	}
}

func (p *PlayerShip) team() team {
	return p.t
}

func (p *PlayerShip) shipInfo() shipInfo {
	return shipInfo{p.x, p.y, p.t.String()}
}

type PirateShip struct {
	transform
	radius           float32
	health           int
	respawnRemaining int
	respawning       bool
	maxHealth        int
	rotation         float32
	renderId         int
	reloadTime       int
	fullReloadTime   int
	home             *Planet

	targetFindTime int
	targeting      *PlayerShip
}

func NewPirateShip(home *Planet) {
	var p PirateShip
	p.accel = 0.8
	p.radius = 32
	p.maxHealth = 100
	p.health = p.maxHealth
	p.speed = 16
	p.renderId = <-NextRenderId
	p.fullReloadTime = framesPerSecond / 5
	p.home = home
	p.targetFindTime = rand.Int() % 30
	//random time so they don't all do it at once
	r := rand.Float64() * math.Pi * 2
	p.x = float32(math.Cos(r)) * 12000
	p.y = float32(math.Sin(r)) * 12000

	NewEntity <- &p
}

func (p *PirateShip) update(overworld *Overworld, planets []*Planet) (alive bool) {
	const targetTime = framesPerSecond * 2
	if p.health <= 0 {
		p.health = p.maxHealth
		go NewDebris(p.x, p.y, p.rotation)
		{
			x := p.x
			y := p.y
			go NewSound(x, y, ExplosionSound())
		}

		r := rand.Float64() * math.Pi * 2
		p.x = float32(math.Cos(r)) * 12000
		p.y = float32(math.Sin(r)) * 12000
		p.vx = 0
		p.vy = 0
		p.respawning = true
		p.respawnRemaining = framesPerSecond * 60 * 5
	}

	p.targetFindTime++
	if p.targetFindTime > targetTime {
		p.targeting = nil
		for _, entity := range overworld.query(p, p.x, p.y, 2000) {
			if player, ok := entity.(*PlayerShip); ok {
				if p.targeting == nil {
					p.targeting = player
				} else {
					Dt := (p.targeting.x-p.x)*(p.targeting.x-p.x) + (p.targeting.y-p.y)*(p.targeting.y-p.y)
					Dp := (player.x-p.x)*(player.x-p.x) + (player.y-p.y)*(player.y-p.y)
					if Dp < Dt {
						p.targeting = player
					}
				}
			}
		}
		p.targetFindTime = 0
	}

	var dx float32
	var dy float32
	tx := p.home.x
	ty := p.home.y
	if p.targeting != nil {
		tx = p.targeting.x
		ty = p.targeting.y
	}
	tx -= p.x
	ty -= p.y
	if tx < 5 && tx > -5 {
		tx = 0
	}
	if ty < 5 && ty > -5 {
		ty = 0
	}

	if tx < 0 {
		dx -= p.speed
	}
	if tx > 0 {
		dx += p.speed
	}
	if ty < 0 {
		dy -= p.speed
	}
	if ty > 0 {
		dy += p.speed
	}
	if dy*dx != 0 {
		dx /= 1.41421356237
		dy /= 1.41421356237
	}

	if p.respawning {
		p.respawnRemaining--
		if p.respawnRemaining <= 0 {
			p.respawning = false
		}
		dx = 0
		dy = 0
	}

	{
		if dx != 0 || dy != 0 {
			p.rotation = float32(math.Atan2(float64(dx), float64(-1*dy))) /
				(math.Pi * 2) * 360
		}
		p.adjustV(dx, dy)
		p.applyV()
	}

	overworld.set(p, p.x, p.y, p.radius)

	//Collision testing code
	// log.Println("____________")
	// log.Println(p.x, p.y)
	//log.Println(overworld.query(nil, p.x, p.y, p.radius))
	// log.Println(overworld.query(nil, p.x, p.y+5, p.radius))

	p.reloadTime += 1
	if p.fullReloadTime < p.reloadTime {
		if p.targeting != nil {
			r := float64(p.rotation-90) / 180 * math.Pi
			vx := float32(math.Cos(r))*16 + p.vx
			vy := float32(math.Sin(r))*16 + p.vy
			x := float32(math.Cos(r))*25 + p.x
			y := float32(math.Sin(r))*25 + p.y
			r += math.Pi / 2
			dx := float32(math.Cos(r)) * 16
			dy := float32(math.Sin(r)) * 16

			go NewBullet(x+dx, y+dy, vx, vy, TeamPirates)
			go NewBullet(x-dx, y-dy, vx, vy, TeamPirates)
			go NewSound(x, y, "laser1")
		} else if p.health < p.maxHealth {
			p.health += 2
			if p.health > p.maxHealth {
				p.health = p.maxHealth
			}
		}
		p.reloadTime = 0
	}

	return true
}

func (p *PirateShip) RenderInfo() RenderInfo {
	return RenderInfo{
		p.renderId, p.x, p.y, p.rotation, "ship_pirate",
	}
}

func (p *PirateShip) damage(damage int, teamSource team) {
	p.health -= damage
	{
		x, y := p.x, p.y
		go NewSound(x, y, "hit0")
	}
}

func (p *PirateShip) team() team {
	return TeamPirates
}

func (p *PirateShip) shipInfo() shipInfo {
	return shipInfo{p.x, p.y, TeamPirates.String()}
}

type Debris struct {
	transform
	lifespan int
	img      string
	r        float32
	renderId int
	rSpeed   float32
}

func NewDebris(x, y, r float32) {
	for i := 0; i < 4; i++ {
		var d Debris
		d.r = r
		d.x = x
		d.y = y
		d.img = "debris_" + strconv.FormatInt(int64(i), 10)
		d.lifespan = rand.Int()%(framesPerSecond*3) + framesPerSecond*5
		d.renderId = <-NextRenderId
		d.rSpeed = rand.Float32()*10 - 5
		d.vx = rand.Float32()*10 - 5
		d.vy = rand.Float32()*10 - 5
		NewEntity <- &d
	}
}

func (d *Debris) update(overworld *Overworld, planets []*Planet) (alive bool) {
	overworld.set(d, d.x, d.y, 64)
	d.lifespan--
	d.r += d.rSpeed
	d.applyV()
	return d.lifespan > 0
}

func (d *Debris) RenderInfo() RenderInfo {
	return RenderInfo{
		d.renderId, d.x, d.y, d.r, d.img,
	}
}

func NewSound(x, y float32, name string) {
	var s Sound
	s.x = x
	s.y = y
	s.name = name
	s.alive = true

	NewEntity <- &s
}

type Sound struct {
	x, y  float32
	name  string
	alive bool
}

func (s *Sound) update(overworld *Overworld, planets []*Planet) (alive bool) {
	last := s.alive
	overworld.set(s, s.x, s.y, 512)
	s.alive = false

	return last
}

func (s *Sound) RenderInfo() RenderInfo {
	return RenderInfo{
		-10, -12000, -12000, 0, "ball_gopher",
	}
}

func ExplosionSound() string {
	switch rand.Int() % 3 {
	case 0:
		return "explosion0"
	case 1:
		return "explosion1"
	default:
		return "explosion3"
	}
}
