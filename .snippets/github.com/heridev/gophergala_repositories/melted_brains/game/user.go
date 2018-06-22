package game

import "math/rand"

type User struct {
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func UserWithRandomName() *User {
	return &User{Name: RandomNames[rand.Intn(len(RandomNames))]}
}

var RandomNames = []string{
	"Washington",
	"Harun al-Rashid",
	"Ashurbanipal",
	"Maria Theresa",
	"Montezuma",
	"Nebuchadnezzar II",
	"Pedro II",
	"Theodora",
	"Dido",
	"Boudicca",
	"Wu Zetian",
	"Harald Bluetooth",
	"Ramesses II",
	"Elizabeth",
	"Haile Selassie",
	"Napoleon",
	"Bismarck",
	"Alexander",
	"Attila",
	"Pachacuti",
	"Gandhi",
	"Gajah Mada",
	"Hiawatha",
	"Oda Nobunaga",
	"Sejong",
	"Pacal",
	"Genghis Khan",
	"Ahmad al-Mansur",
	"William",
	"Suleiman",
	"Darius I",
	"Casimir III",
	"Kamehameha",
	"Maria I",
	"Augustus Caesar",
	"Catherine",
	"Pocatello",
	"Ramkhamhaeng",
	"Askia",
	"Tercio",
	"Gustavus Adolphus",
	"Enrico Dandolo",
	"Shaka",
}
