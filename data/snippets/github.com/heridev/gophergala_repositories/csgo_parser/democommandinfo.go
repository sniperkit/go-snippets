package main

const (
	FDEMO_NORMAL            byte  = 0
	FDEMO_USE_ORIGIN2       byte  = (1 << 0)
	FDEMO_USE_ANGLES2       byte  = (1 << 1)
	FDEMO_NOINTERP          byte  = (1 << 2)
	MAX_SPLITSCREEN_CLIENTS int32 = 2
)

// QAngle setup
type QAngle struct {
	x, y, z float32
}

func (q *QAngle) InitDefault() {
	q.x = 0.0
	q.y = 0.0
	q.z = 0.0
}

func (q *QAngle) Init(x float32, y float32, z float32) {
	q.x = x
	q.y = y
	q.z = z
}

// Vector setup
type Vector struct {
	x, y, z float32
}

func (v *Vector) InitDefault() {
	v.x = 0.0
	v.y = 0.0
	v.z = 0.0
}

func (v *Vector) Init(x float32, y float32, z float32) {
	v.x = x
	v.y = y
	v.z = z
}

// Split type
type Split struct {
	flags byte

	viewOrigin      Vector
	viewAngles      QAngle
	localViewAngles QAngle

	viewOrigin2      Vector
	viewAngles2      QAngle
	localViewAngles2 QAngle
}

// Set default Split values
func (s *Split) Init() {
	s.flags = FDEMO_NORMAL

	s.viewOrigin.InitDefault()
	s.viewAngles.InitDefault()
	s.localViewAngles.InitDefault()

	s.viewOrigin2.InitDefault()
	s.viewAngles2.InitDefault()
	s.localViewAngles.InitDefault()
}

// Copy the src Split to the current one
func (s *Split) Copy(src Split) {
	s.flags = src.flags

	s.viewOrigin = src.viewOrigin
	s.viewAngles = src.viewAngles
	s.localViewAngles = src.localViewAngles

	s.viewOrigin2 = src.viewOrigin2
	s.viewAngles2 = src.viewAngles2
	s.localViewAngles2 = src.localViewAngles2
}

func (s *Split) GetViewOrigin() *Vector {
	if (s.flags & FDEMO_USE_ORIGIN2) >= 1 {
		return &s.viewOrigin2
	}
	return &s.viewOrigin
}

func (s *Split) GetViewAngles() *QAngle {
	if (s.flags & FDEMO_USE_ANGLES2) >= 1 {
		return &s.viewAngles2
	}
	return &s.viewAngles
}

func (s *Split) GetLocalViewAngles() *QAngle {
	if (s.flags & FDEMO_USE_ANGLES2) >= 1 {
		return &s.localViewAngles2
	}
	return &s.localViewAngles
}

func (s *Split) Reset() {
	s.flags = 0
	s.viewOrigin2 = s.viewOrigin
	s.viewAngles2 = s.viewAngles
	s.localViewAngles2 = s.localViewAngles
}

type DemoCmdInfo struct {
	u [MAX_SPLITSCREEN_CLIENTS]Split
}

func (d *DemoCmdInfo) Reset() {
	var i int32 = 0
	for i < MAX_SPLITSCREEN_CLIENTS {
		d.u[i].Reset()
		i++
	}
}
